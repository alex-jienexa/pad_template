package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alex-jienexa/pad_template/internal/records"
)

func runRequirement(args []string, out, errw io.Writer) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(out, "Usage: pad requirement <list|show|add|set> [options]")
		return 0
	}
	switch args[0] {
	case "list":
		return requirementList(args[1:], out, errw)
	case "show":
		return requirementShow(args[1:], out, errw)
	case "add":
		return requirementAdd(args[1:], out, errw)
	case "set":
		return requirementSet(args[1:], out, errw)
	default:
		fmt.Fprintf(errw, "requirement: unknown subcommand %q\n", args[0])
		return 2
	}
}

type reqRecord struct{ path, id, title, typ, status, priority string }

func findRequirements(root string) ([]reqRecord, error) {
	dir := filepath.Join(root, ".agent", "requirements")
	entries, e := os.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	var a []reqRecord
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".md") || ent.Name() == "README.md" || strings.Contains(strings.ToUpper(ent.Name()), "TEMPLATE") {
			continue
		}
		p := filepath.Join(dir, ent.Name())
		d, e := records.Read(p)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", ent.Name(), e)
		}
		a = append(a, reqRecord{p, records.String(d.Metadata, "id"), records.String(d.Metadata, "title"), records.String(d.Metadata, "type"), records.String(d.Metadata, "status"), records.String(d.Metadata, "priority")})
	}
	sort.Slice(a, func(i, j int) bool { return a[i].id < a[j].id })
	return a, nil
}
func requirementList(args []string, out, errw io.Writer) int {
	root, e := discoverArg(args)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	a, e := findRequirements(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	fmt.Fprintf(out, "%-20s %-16s %-12s %-8s %s\n", "ID", "STATUS", "TYPE", "PRIORITY", "TITLE")
	for _, r := range a {
		fmt.Fprintf(out, "%-20s %-16s %-12s %-8s %s\n", r.id, r.status, r.typ, r.priority, r.title)
	}
	fmt.Fprintf(out, "%d requirement(s)\n", len(a))
	return 0
}
func requirementShow(args []string, out, errw io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(errw, "requirement show requires an ID")
		return 2
	}
	id := args[0]
	root, e := discoverArg(args[1:])
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	a, e := findRequirements(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	for _, r := range a {
		if r.id == id {
			b, e := os.ReadFile(r.path)
			if e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			_, _ = out.Write(b)
			return 0
		}
	}
	fmt.Fprintf(errw, "requirement %s not found\n", id)
	return 1
}
func requirementAdd(args []string, out, errw io.Writer) int {
	fs := flag.NewFlagSet("requirement add", flag.ContinueOnError)
	fs.SetOutput(errw)
	title := fs.String("title", "", "requirement title (required)")
	typ := fs.String("type", "functional", "functional|non_functional|constraint")
	priority := fs.String("priority", "should", "must|should|could|wont")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	if strings.TrimSpace(*title) == "" {
		fmt.Fprintln(errw, "requirement add: --title is required")
		return 2
	}
	if !contains([]string{"functional", "non_functional", "constraint"}, *typ) || !contains([]string{"must", "should", "could", "wont"}, *priority) {
		fmt.Fprintln(errw, "requirement add: invalid type or priority")
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	a, e := findRequirements(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	max := 0
	for _, r := range a {
		var n int
		_, _ = fmt.Sscanf(r.id, "PAD-REQ-%d", &n)
		if n > max {
			max = n
		}
	}
	id := fmt.Sprintf("PAD-REQ-%03d", max+1)
	date := records.Date()
	meta := map[string]any{"id": id, "title": *title, "type": *typ, "status": "proposed", "priority": *priority, "created_at": date, "updated_at": date, "verification": "inspection", "related_requirements": []string{}, "related_tasks": []string{}, "tags": []string{}}
	body := "# Statement\n\nDescribe the requirement in verifiable terms.\n\n## Rationale\n\nExplain why this requirement is needed.\n\n## Acceptance criteria\n\n- [ ] Observable verification condition\n"
	path := filepath.Join(root, ".agent", "requirements", id+".md")
	if e := records.Write(path, meta, body); e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	fmt.Fprintf(out, "Created %s: %s\n", id, path)
	return 0
}

func requirementSet(args []string, out, errw io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(errw, "requirement set requires an ID")
		return 2
	}
	id := args[0]
	fs := flag.NewFlagSet("requirement set", flag.ContinueOnError)
	fs.SetOutput(errw)
	status := fs.String("status", "", "proposed|approved|implemented|verified|deprecated")
	if e := fs.Parse(args[1:]); e != nil {
		return 2
	}
	if !contains([]string{"proposed", "approved", "implemented", "verified", "deprecated"}, *status) {
		fmt.Fprintln(errw, "requirement set: provide a valid --status")
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	items, e := findRequirements(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	transitions := map[string]map[string]bool{"proposed": {"approved": true, "deprecated": true}, "approved": {"implemented": true, "deprecated": true}, "implemented": {"verified": true, "approved": true, "deprecated": true}, "verified": {"deprecated": true}, "deprecated": {"proposed": true}}
	for _, r := range items {
		if r.id != id {
			continue
		}
		d, e := records.Read(r.path)
		if e != nil {
			fmt.Fprintln(errw, e)
			return 1
		}
		old := records.String(d.Metadata, "status")
		if old != *status && !transitions[old][*status] {
			fmt.Fprintf(errw, "invalid requirement transition: %s -> %s\n", old, *status)
			return 1
		}
		d.Metadata["status"] = *status
		d.Metadata["updated_at"] = records.Date()
		if e := records.Write(r.path, d.Metadata, d.Body); e != nil {
			fmt.Fprintln(errw, e)
			return 1
		}
		fmt.Fprintf(out, "Updated %s: %s -> %s\n", id, old, *status)
		return 0
	}
	fmt.Fprintf(errw, "requirement %s not found\n", id)
	return 1
}
