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

type decisionRecord struct{ path, id, title, status, date string }

func findDecisions(root string) ([]decisionRecord, error) {
	dir := filepath.Join(root, ".agent", "decisions")
	entries, e := os.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	var a []decisionRecord
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".md") || ent.Name() == "README.md" || strings.Contains(strings.ToUpper(ent.Name()), "TEMPLATE") {
			continue
		}
		p := filepath.Join(dir, ent.Name())
		d, e := records.Read(p)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", ent.Name(), e)
		}
		a = append(a, decisionRecord{p, records.String(d.Metadata, "id"), records.String(d.Metadata, "title"), records.String(d.Metadata, "status"), records.String(d.Metadata, "date")})
	}
	sort.Slice(a, func(i, j int) bool { return a[i].id < a[j].id })
	return a, nil
}
func runDecision(args []string, out, errw io.Writer) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(out, "Usage: pad decision <list|show|add|set> [options]")
		return 0
	}
	switch args[0] {
	case "list":
		return decisionList(args[1:], out, errw)
	case "show":
		return decisionShow(args[1:], out, errw)
	case "add":
		return decisionAdd(args[1:], out, errw)
	case "set":
		return decisionSet(args[1:], out, errw)
	default:
		fmt.Fprintf(errw, "decision: unknown subcommand %q\n", args[0])
		return 2
	}
}
func decisionList(args []string, out, errw io.Writer) int {
	root, e := discoverArg(args)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	a, e := findDecisions(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	fmt.Fprintf(out, "%-12s %-12s %-12s %s\n", "ID", "STATUS", "DATE", "TITLE")
	for _, d := range a {
		fmt.Fprintf(out, "%-12s %-12s %-12s %s\n", d.id, d.status, d.date, d.title)
	}
	fmt.Fprintf(out, "%d decision(s)\n", len(a))
	return 0
}
func decisionShow(args []string, out, errw io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(errw, "decision show requires an ID")
		return 2
	}
	id := args[0]
	root, e := discoverArg(args[1:])
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	a, e := findDecisions(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	for _, d := range a {
		if d.id == id {
			b, e := os.ReadFile(d.path)
			if e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			_, _ = out.Write(b)
			return 0
		}
	}
	fmt.Fprintf(errw, "decision %s not found\n", id)
	return 1
}
func decisionAdd(args []string, out, errw io.Writer) int {
	fs := flag.NewFlagSet("decision add", flag.ContinueOnError)
	fs.SetOutput(errw)
	title := fs.String("title", "", "decision title (required)")
	decider := fs.String("decider", "project-owner", "accountable decision-maker")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	if strings.TrimSpace(*title) == "" {
		fmt.Fprintln(errw, "decision add: --title is required")
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	a, e := findDecisions(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	max := 0
	for _, d := range a {
		var n int
		_, _ = fmt.Sscanf(d.id, "ADR-%d", &n)
		if n > max {
			max = n
		}
	}
	id := fmt.Sprintf("ADR-%04d", max+1)
	meta := map[string]any{"id": id, "title": *title, "status": "proposed", "date": records.Date(), "deciders": []string{*decider}, "supersedes": nil, "superseded_by": nil, "related_requirements": []string{}, "related_tasks": []string{}}
	body := "# Context\n\nDescribe the problem or decision context.\n\n## Decision drivers\n\n- Driver\n\n## Considered options\n\n1. Option A\n2. Option B\n\n## Decision\n\nState the selected option and rationale.\n\n## Consequences\n\n### Positive\n\n- To be documented\n\n### Negative / trade-offs\n\n- To be documented\n\n## Follow-up actions\n\n- [ ] Record implementation tasks\n"
	path := filepath.Join(root, ".agent", "decisions", id+".md")
	if e := records.Write(path, meta, body); e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	fmt.Fprintf(out, "Created %s: %s\n", id, path)
	return 0
}
func decisionSet(args []string, out, errw io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(errw, "decision set requires an ID")
		return 2
	}
	id := args[0]
	fs := flag.NewFlagSet("decision set", flag.ContinueOnError)
	fs.SetOutput(errw)
	status := fs.String("status", "", "proposed|accepted|rejected|superseded|deprecated")
	if e := fs.Parse(args[1:]); e != nil {
		return 2
	}
	if !contains([]string{"proposed", "accepted", "rejected", "superseded", "deprecated"}, *status) {
		fmt.Fprintln(errw, "decision set: provide a valid --status")
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	a, e := findDecisions(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	transitions := map[string]map[string]bool{"proposed": {"accepted": true, "rejected": true}, "accepted": {"superseded": true, "deprecated": true}, "rejected": {"proposed": true}, "superseded": {"deprecated": true}, "deprecated": {"proposed": true}}
	for _, d := range a {
		if d.id != id {
			continue
		}
		doc, e := records.Read(d.path)
		if e != nil {
			fmt.Fprintln(errw, e)
			return 1
		}
		old := records.String(doc.Metadata, "status")
		if old != *status && !transitions[old][*status] {
			fmt.Fprintf(errw, "invalid decision transition: %s -> %s\n", old, *status)
			return 1
		}
		doc.Metadata["status"] = *status
		if e := records.Write(d.path, doc.Metadata, doc.Body); e != nil {
			fmt.Fprintln(errw, e)
			return 1
		}
		fmt.Fprintf(out, "Updated %s: %s -> %s\n", id, old, *status)
		return 0
	}
	fmt.Fprintf(errw, "decision %s not found\n", id)
	return 1
}
