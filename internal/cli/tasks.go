package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alex-jienexa/pad_template/internal/project"
	"github.com/alex-jienexa/pad_template/internal/records"
)

var taskStatuses = map[string]bool{"backlog": true, "ready": true, "in_progress": true, "review": true, "blocked": true, "done": true, "cancelled": true}
var taskTransitions = map[string]map[string]bool{
	"backlog": {"ready": true, "blocked": true, "cancelled": true}, "ready": {"in_progress": true, "blocked": true, "cancelled": true},
	"in_progress": {"review": true, "blocked": true, "ready": true, "cancelled": true}, "review": {"done": true, "in_progress": true, "blocked": true},
	"blocked": {"ready": true, "in_progress": true, "cancelled": true}, "done": {"ready": true}, "cancelled": {"backlog": true},
}

func runTask(args []string, out, errw io.Writer) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(out, "Usage: pad task <list|show|add|set> [options]\n  list [--status STATUS] [path]\n  show ID [path]\n  add --title TITLE [--priority medium] [--assignee NAME] [path]\n  set ID --status STATUS [path]")
		return 0
	}
	switch args[0] {
	case "list":
		return taskList(args[1:], out, errw)
	case "show":
		return taskShow(args[1:], out, errw)
	case "add":
		return taskAdd(args[1:], out, errw)
	case "set":
		return taskSet(args[1:], out, errw)
	default:
		fmt.Fprintf(errw, "task: unknown subcommand %q\n", args[0])
		return 2
	}
}

type taskRecord struct{ path, id, title, status, priority, assignee string }

func findTasks(root string) ([]taskRecord, error) {
	dir := filepath.Join(root, ".agent", "tasks")
	entries, e := os.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	var a []taskRecord
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".md") || ent.Name() == "README.md" || strings.Contains(strings.ToUpper(ent.Name()), "TEMPLATE") {
			continue
		}
		p := filepath.Join(dir, ent.Name())
		d, e := records.Read(p)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", ent.Name(), e)
		}
		a = append(a, taskRecord{p, records.String(d.Metadata, "id"), records.String(d.Metadata, "title"), records.String(d.Metadata, "status"), records.String(d.Metadata, "priority"), records.String(d.Metadata, "assignee")})
	}
	sort.Slice(a, func(i, j int) bool { return a[i].id < a[j].id })
	return a, nil
}
func taskList(args []string, out, errw io.Writer) int {
	fs := flag.NewFlagSet("task list", flag.ContinueOnError)
	fs.SetOutput(errw)
	status := fs.String("status", "", "filter by status")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	tasks, e := findTasks(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	n := 0
	fmt.Fprintf(out, "%-14s %-13s %-9s %-18s %s\n", "ID", "STATUS", "PRIORITY", "ASSIGNEE", "TITLE")
	for _, t := range tasks {
		if *status != "" && t.status != *status {
			continue
		}
		fmt.Fprintf(out, "%-14s %-13s %-9s %-18s %s\n", t.id, t.status, t.priority, t.assignee, t.title)
		n++
	}
	fmt.Fprintf(out, "%d task(s)\n", n)
	return 0
}
func taskShow(args []string, out, errw io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(errw, "task show requires an ID")
		return 2
	}
	id := args[0]
	root, e := discoverArg(args[1:])
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	tasks, e := findTasks(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	for _, t := range tasks {
		if t.id == id {
			b, e := os.ReadFile(t.path)
			if e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			_, _ = out.Write(b)
			return 0
		}
	}
	fmt.Fprintf(errw, "task %s not found\n", id)
	return 1
}
func taskAdd(args []string, out, errw io.Writer) int {
	fs := flag.NewFlagSet("task add", flag.ContinueOnError)
	fs.SetOutput(errw)
	title := fs.String("title", "", "task title (required)")
	priority := fs.String("priority", "medium", "critical|high|medium|low")
	assignee := fs.String("assignee", "unassigned", "task owner")
	estimate := fs.String("estimate", "", "optional estimate")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	if strings.TrimSpace(*title) == "" {
		fmt.Fprintln(errw, "task add: --title is required")
		return 2
	}
	if !map[string]bool{"critical": true, "high": true, "medium": true, "low": true}[*priority] {
		fmt.Fprintln(errw, "task add: invalid priority")
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	tasks, e := findTasks(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	max := 0
	for _, t := range tasks {
		var n int
		_, _ = fmt.Sscanf(t.id, "TASK-%d", &n)
		if n > max {
			max = n
		}
	}
	id := fmt.Sprintf("TASK-%03d", max+1)
	meta := map[string]any{"id": id, "title": *title, "status": "backlog", "priority": *priority, "assignee": *assignee, "created_at": records.Date(), "updated_at": records.Date(), "estimate": nil, "parent": nil, "depends_on": []string{}, "requirement_ids": []string{}, "files_owned": []string{}, "labels": []string{}}
	if *estimate != "" {
		meta["estimate"] = *estimate
	}
	body := "# " + *title + "\n\n## Objective\n\nDescribe the intended outcome.\n\n## Scope\n\n- Define included work\n\n## Acceptance criteria\n\n- [ ] Observable completion condition\n\n## Verification\n\nDescribe how the result will be checked.\n"
	path := filepath.Join(root, ".agent", "tasks", id+".md")
	if e := records.Write(path, meta, body); e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	fmt.Fprintf(out, "Created %s: %s\n", id, path)
	return 0
}
func taskSet(args []string, out, errw io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(errw, "task set requires an ID")
		return 2
	}
	id := args[0]
	fs := flag.NewFlagSet("task set", flag.ContinueOnError)
	fs.SetOutput(errw)
	status := fs.String("status", "", "new status")
	priority := fs.String("priority", "", "new priority")
	assignee := fs.String("assignee", "", "new assignee")
	reason := fs.String("reason", "", "required explanation when reopening a completed/cancelled task")
	if e := fs.Parse(args[1:]); e != nil {
		return 2
	}
	if *status == "" && *priority == "" && *assignee == "" {
		fmt.Fprintln(errw, "task set: provide at least one of --status, --priority, or --assignee")
		return 2
	}
	if *status != "" && !taskStatuses[*status] {
		fmt.Fprintln(errw, "task set: invalid status")
		return 2
	}
	if *priority != "" && !map[string]bool{"critical": true, "high": true, "medium": true, "low": true}[*priority] {
		fmt.Fprintln(errw, "task set: invalid priority")
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	tasks, e := findTasks(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	for _, t := range tasks {
		if t.id != id {
			continue
		}
		d, e := records.Read(t.path)
		if e != nil {
			fmt.Fprintln(errw, e)
			return 1
		}
		old := records.String(d.Metadata, "status")
		if *status != "" && old != *status && !taskTransitions[old][*status] {
			fmt.Fprintf(errw, "invalid task transition: %s -> %s\n", old, *status)
			return 1
		}
		if old == "done" && *status == "ready" || old == "cancelled" && *status == "backlog" {
			if strings.TrimSpace(*reason) == "" {
				fmt.Fprintln(errw, "task set: --reason is required when reopening a done or cancelled task")
				return 2
			}
			d.Body = strings.TrimRight(d.Body, "\n") + fmt.Sprintf("\n\n## Status history\n\n- %s — reopened from %s to %s. Reason: %s\n", records.Date(), old, *status, *reason)
		}
		changes := []string{}
		if *status != "" {
			d.Metadata["status"] = *status
			changes = append(changes, "status="+*status)
		}
		if *priority != "" {
			d.Metadata["priority"] = *priority
			changes = append(changes, "priority="+*priority)
		}
		if *assignee != "" {
			d.Metadata["assignee"] = *assignee
			changes = append(changes, "assignee="+*assignee)
		}
		d.Metadata["updated_at"] = records.Date()
		if e := records.Write(t.path, d.Metadata, d.Body); e != nil {
			fmt.Fprintln(errw, e)
			return 1
		}
		fmt.Fprintf(out, "Updated %s: %s\n", id, strings.Join(changes, ", "))
		return 0
	}
	fmt.Fprintf(errw, "task %s not found\n", id)
	return 1
}

func discoverArg(args []string) (string, error) {
	if len(args) > 1 {
		return "", fmt.Errorf("too many path arguments")
	}
	p := "."
	if len(args) == 1 {
		p = args[0]
	}
	return project.Discover(p)
}
