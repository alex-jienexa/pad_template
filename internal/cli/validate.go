package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/alex-jienexa/pad_template/internal/config"
	"github.com/alex-jienexa/pad_template/internal/records"
	"github.com/alex-jienexa/pad_template/internal/yamlmini"
)

var taskIDPattern = regexp.MustCompile(`^TASK-[0-9]{3,}$`)
var reqIDPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]*-REQ-[0-9]{3,}$`)
var adrIDPattern = regexp.MustCompile(`^ADR-[0-9]{4,}$`)

func validateProject(root string) []string {
	var issues []string
	required := []string{"AGENTS.md", "pad.yaml", ".agent/README.md", ".agent/context/PROJECT.md", ".agent/state/current.yaml", ".agent/state/HANDOFF.md", ".agent/tasks", ".agent/requirements", ".agent/decisions", ".agent/protocols/SESSION_START.md", ".agent/protocols/RECOVERY.md", ".agent/schemas/README.md", ".pad/template-base", ".pad/template-base/pad.yaml"}
	for _, name := range required {
		p := filepath.Join(root, filepath.FromSlash(name))
		info, e := os.Stat(p)
		if e != nil {
			issues = append(issues, fmt.Sprintf("missing or inaccessible required path: %s", name))
			continue
		}
		if strings.HasSuffix(name, "tasks") || strings.HasSuffix(name, "requirements") || strings.HasSuffix(name, "decisions") || name == ".pad/template-base" {
			if !info.IsDir() {
				issues = append(issues, fmt.Sprintf("%s must be a directory", name))
			}
		} else if info.IsDir() {
			issues = append(issues, fmt.Sprintf("%s must be a file", name))
		}
	}
	configPath := filepath.Join(root, "pad.yaml")
	cfg, configErr := config.Load(configPath)
	if configErr != nil {
		issues = append(issues, "pad.yaml: "+configErr.Error())
	} else {
		for _, issue := range cfg.Validate() {
			issues = append(issues, "pad.yaml: "+issue)
		}
		if strings.Contains(cfg.Project.Name, "{{") {
			issues = append(issues, "pad.yaml: replace project.name template placeholder")
		}
	}

	statePath := filepath.Join(root, ".agent/state/current.yaml")
	if b, e := os.ReadFile(statePath); e == nil {
		m, e := yamlmini.Parse(string(b))
		if e != nil {
			issues = append(issues, ".agent/state/current.yaml: "+e.Error())
		} else {
			if v := stringValue(m["project_status"]); !contains([]string{"initializing", "active", "paused", "maintenance", "archived"}, v) {
				issues = append(issues, "current.yaml: invalid project_status")
			}
			if v := stringValue(m["health"]); !contains([]string{"healthy", "at_risk", "blocked", "unknown"}, v) {
				issues = append(issues, "current.yaml: invalid health")
			}
			if v := m["updated_at"]; v != nil {
				if _, e := time.Parse(time.RFC3339, stringValue(v)); e != nil {
					issues = append(issues, "current.yaml: updated_at must be an RFC3339 timestamp or null")
				}
			}
			for _, key := range []string{"open_blockers", "next_actions", "recent_changes", "risks"} {
				if _, ok := m[key].([]any); !ok {
					issues = append(issues, "current.yaml: "+key+" must be a YAML list")
				}
			}
		}
	}
	tasks, taskIssues := validateRecords(filepath.Join(root, ".agent/tasks"), "task")
	issues = append(issues, taskIssues...)
	reqs, reqIssues := validateRecords(filepath.Join(root, ".agent/requirements"), "requirement")
	issues = append(issues, reqIssues...)
	decisions, adrIssues := validateRecords(filepath.Join(root, ".agent/decisions"), "adr")
	issues = append(issues, adrIssues...)
	taskIDs := map[string]bool{}
	for _, m := range tasks {
		taskIDs[records.String(m, "id")] = true
	}
	reqIDs := map[string]bool{}
	for _, m := range reqs {
		reqIDs[records.String(m, "id")] = true
	}
	if b, err := os.ReadFile(filepath.Join(root, ".agent", "state", "current.yaml")); err == nil {
		if state, parseErr := yamlmini.Parse(string(b)); parseErr == nil {
			active := stringValue(state["active_task"])
			if active != "" && active != "null" && !taskIDs[active] {
				issues = append(issues, "current.yaml: active_task references unknown task "+active)
			}
		}
	}
	for _, m := range reqs {
		for _, taskID := range records.Strings(m, "related_tasks") {
			if !taskIDs[taskID] {
				issues = append(issues, fmt.Sprintf("%s references unknown task %s", records.String(m, "id"), taskID))
			}
		}
	}
	for _, m := range decisions {
		for _, taskID := range records.Strings(m, "related_tasks") {
			if !taskIDs[taskID] {
				issues = append(issues, fmt.Sprintf("%s references unknown task %s", records.String(m, "id"), taskID))
			}
		}
		for _, reqID := range records.Strings(m, "related_requirements") {
			if !reqIDs[reqID] {
				issues = append(issues, fmt.Sprintf("%s references unknown requirement %s", records.String(m, "id"), reqID))
			}
		}
	}
	for _, m := range tasks {
		id := records.String(m, "id")
		parent := records.String(m, "parent")
		if parent != "" && parent != "null" && !taskIDs[parent] {
			issues = append(issues, fmt.Sprintf("%s references unknown parent task %s", id, parent))
		}
		for _, dep := range records.Strings(m, "depends_on") {
			if !taskIDs[dep] {
				issues = append(issues, fmt.Sprintf("%s depends_on unknown task %s", id, dep))
			}
		}
		for _, req := range records.Strings(m, "requirement_ids") {
			if !reqIDs[req] {
				issues = append(issues, fmt.Sprintf("%s references unknown requirement %s", id, req))
			}
		}
	}
	issues = append(issues, dependencyCycles(tasks)...)
	return issues
}
func validateRecords(dir, kind string) ([]map[string]any, []string) {
	entries, e := os.ReadDir(dir)
	if e != nil {
		return nil, nil
	}
	var recordsOut []map[string]any
	var issues []string
	seen := map[string]bool{}
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".md") || ent.Name() == "README.md" || strings.Contains(strings.ToUpper(ent.Name()), "TEMPLATE") {
			continue
		}
		path := filepath.Join(dir, ent.Name())
		doc, e := records.Read(path)
		if e != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", filepath.ToSlash(path), e))
			continue
		}
		m := doc.Metadata
		id := records.String(m, "id")
		title := records.String(m, "title")
		if id == "" {
			issues = append(issues, ent.Name()+": missing id")
		} else {
			pattern := taskIDPattern
			if kind == "requirement" {
				for _, key := range []string{"created_at", "updated_at", "priority"} {
					if records.String(m, key) == "" {
						issues = append(issues, ent.Name()+": missing required field "+key)
					}
				}
				pattern = reqIDPattern
			}
			if kind == "adr" {
				if _, ok := m["deciders"].([]any); !ok {
					issues = append(issues, ent.Name()+": deciders must be a YAML list")
				}
				pattern = adrIDPattern
			}
			if !pattern.MatchString(id) {
				issues = append(issues, fmt.Sprintf("%s: invalid %s ID %q", ent.Name(), kind, id))
			}
			if seen[id] {
				issues = append(issues, fmt.Sprintf("duplicate %s ID %s", kind, id))
			}
			seen[id] = true
			if !filenameMatchesID(ent.Name(), id) {
				issues = append(issues, fmt.Sprintf("%s: filename must start with metadata ID %s", ent.Name(), id))
			}
		}
		if title == "" {
			issues = append(issues, ent.Name()+": missing title")
		}
		if strings.TrimSpace(doc.Body) == "" {
			issues = append(issues, ent.Name()+": markdown body is empty")
		}
		if kind == "task" {
			for _, key := range []string{"assignee", "created_at", "updated_at"} {
				if records.String(m, key) == "" {
					issues = append(issues, ent.Name()+": missing required field "+key)
				}
			}
			for _, key := range []string{"created_at", "updated_at"} {
				if _, err := time.Parse("2006-01-02", records.String(m, key)); err != nil {
					issues = append(issues, ent.Name()+": "+key+" must use YYYY-MM-DD")
				}
			}
			for _, key := range []string{"depends_on", "requirement_ids", "files_owned", "labels"} {
				if _, ok := m[key].([]any); !ok {
					issues = append(issues, ent.Name()+": "+key+" must be a YAML list")
				}
			}
			if !contains([]string{"backlog", "ready", "in_progress", "review", "blocked", "done", "cancelled"}, records.String(m, "status")) {
				issues = append(issues, ent.Name()+": invalid task status")
			}
			if !contains([]string{"critical", "high", "medium", "low"}, records.String(m, "priority")) {
				issues = append(issues, ent.Name()+": invalid task priority")
			}
		}
		if kind == "requirement" {
			for _, key := range []string{"created_at", "updated_at", "priority"} {
				if records.String(m, key) == "" {
					issues = append(issues, ent.Name()+": missing required field "+key)
				}
			}
			for _, key := range []string{"created_at", "updated_at"} {
				if _, err := time.Parse("2006-01-02", records.String(m, key)); err != nil {
					issues = append(issues, ent.Name()+": "+key+" must use YYYY-MM-DD")
				}
			}
			if !contains([]string{"functional", "non_functional", "constraint"}, records.String(m, "type")) {
				issues = append(issues, ent.Name()+": invalid requirement type")
			}
			if !contains([]string{"proposed", "approved", "implemented", "verified", "deprecated"}, records.String(m, "status")) {
				issues = append(issues, ent.Name()+": invalid requirement status")
			}
		}
		if kind == "adr" {
			if _, ok := m["deciders"].([]any); !ok {
				issues = append(issues, ent.Name()+": deciders must be a YAML list")
			}
			if !contains([]string{"proposed", "accepted", "rejected", "superseded", "deprecated"}, records.String(m, "status")) {
				issues = append(issues, ent.Name()+": invalid ADR status")
			}
			if records.String(m, "date") == "" {
				issues = append(issues, ent.Name()+": missing date")
			} else if _, err := time.Parse("2006-01-02", records.String(m, "date")); err != nil {
				issues = append(issues, ent.Name()+": date must use YYYY-MM-DD")
			}
		}
		recordsOut = append(recordsOut, m)
	}
	return recordsOut, issues
}
func dependencyCycles(tasks []map[string]any) []string {
	graph := map[string][]string{}
	for _, m := range tasks {
		graph[records.String(m, "id")] = records.Strings(m, "depends_on")
	}
	state := map[string]int{}
	var issues []string
	var visit func(string)
	visit = func(id string) {
		if state[id] == 1 {
			issues = append(issues, "task dependency cycle detected at "+id)
			return
		}
		if state[id] == 2 {
			return
		}
		state[id] = 1
		for _, d := range graph[id] {
			if _, ok := graph[d]; ok {
				visit(d)
			}
		}
		state[id] = 2
	}
	for id := range graph {
		visit(id)
	}
	return issues
}
func nestedString(m map[string]any, keys ...string) string {
	var v any = m
	for _, k := range keys {
		x, ok := v.(map[string]any)
		if !ok {
			return ""
		}
		v = x[k]
	}
	return stringValue(v)
}
func stringValue(v any) string { s, _ := v.(string); return s }
func contains(a []string, s string) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}

func nestedMap(m map[string]any, key string) map[string]any {
	v, _ := m[key].(map[string]any)
	if v == nil {
		return map[string]any{}
	}
	return v
}

func filenameMatchesID(filename, id string) bool {
	return filename == id+".md" || strings.HasPrefix(filename, id+"-") && strings.HasSuffix(filename, ".md")
}
