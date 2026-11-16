package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	padtemplate "github.com/alex-jienexa/pad_template/internal/template"
)

const Version = "1.0.0-dev"

// Run executes the PAD CLI and returns a process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stdout)
		return 0
	}
	switch args[0] {
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	case "version", "--version":
		fmt.Fprintf(stdout, "pad %s\n", Version)
		return 0
	case "init":
		return runInit(args[1:], stdout, stderr)
	case "adopt":
		return runAdopt(args[1:], stdout, stderr)
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "status":
		return runStatus(args[1:], stdout, stderr)
	case "task":
		return runTask(args[1:], stdout, stderr)
	case "requirement", "req":
		return runRequirement(args[1:], stdout, stderr)
	case "decision", "adr":
		return runDecision(args[1:], stdout, stderr)
	case "agent":
		return runAgent(args[1:], stdout, stderr)
	case "checkpoint":
		return runCheckpoint(args[1:], stdout, stderr)
	case "update":
		return runUpdate(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		printUsage(stderr)
		return 2
	}
}
func printUsage(w io.Writer) {
	fmt.Fprintln(w, `PAD — Persistent Agent Development

Usage:
  pad <command> [options]

Project lifecycle:
  init [path]             Initialize a project from the embedded template
  adopt [path]            Record the baseline for a matching template-generated project
  validate [path]         Validate configuration, records, and references
  status [path]           Show project state, task counts, and Git status
  update [path]           Safely update managed files using a three-way baseline
  checkpoint              Write a durable session checkpoint

Work management:
  task list|show|add|set  Manage project tasks
  requirement list|show|add|set  Manage project requirements
  decision list|show|add|set     Manage architecture decisions
  agent list|show                Inspect agent role profiles

Other:
  version                 Print CLI version
  help                    Show this help

Run "pad <command> --help" for command options. The CLI does not access the network.`)
}
func normalizeInitArgs(args []string) []string {
	var options, positional []string
	valueFlags := map[string]bool{"--name": true, "--description": true, "--repository": true, "--language": true}
	boolFlags := map[string]bool{"--force": true, "--dry-run": true}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if boolFlags[arg] || strings.HasPrefix(arg, "--name=") || strings.HasPrefix(arg, "--description=") || strings.HasPrefix(arg, "--repository=") || strings.HasPrefix(arg, "--language=") {
			options = append(options, arg)
			continue
		}
		if valueFlags[arg] {
			options = append(options, arg)
			if i+1 < len(args) {
				i++
				options = append(options, args[i])
			}
			continue
		}
		positional = append(positional, arg)
	}
	return append(options, positional...)
}

func runInit(args []string, out, errw io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(errw)
	force := fs.Bool("force", false, "overwrite colliding starter files after reviewing the plan")
	dry := fs.Bool("dry-run", false, "show the plan without writing files")
	name := fs.String("name", "", "project name (defaults to target directory name)")
	description := fs.String("description", "Describe the project purpose and scope.", "project description")
	repository := fs.String("repository", "", "project repository URL (optional)")
	language := fs.String("language", "Not specified", "primary implementation language")
	if e := fs.Parse(normalizeInitArgs(args)); e != nil {
		return 2
	}
	if fs.NArg() > 1 {
		fmt.Fprintln(errw, "init accepts at most one target path")
		return 2
	}
	target := "."
	if fs.NArg() == 1 {
		target = fs.Arg(0)
	}
	plan, e := padtemplate.Init(target, *force, *dry)
	for _, it := range plan {
		action := "create"
		if it.Exists && it.WillWrite {
			action = "overwrite"
		} else if it.Exists {
			action = "preserve"
		}
		fmt.Fprintf(out, "%-10s %s\n", action, it.Path)
	}
	if e != nil {
		fmt.Fprintf(errw, "init: %v\n", e)
		return 1
	}
	if *dry {
		fmt.Fprintf(out, "Dry run: %d template files planned; no files written.\n", len(plan))
	} else {
		abs, _ := filepath.Abs(target)
		projectName := *name
		if projectName == "" {
			projectName = filepath.Base(abs)
		}
		if e := substituteProjectPlaceholders(abs, projectName, *description, *repository, *language); e != nil {
			fmt.Fprintf(errw, "init: placeholder substitution failed: %v\n", e)
			return 1
		}
		fmt.Fprintf(out, "Initialized PAD project %q in %s (%d template files plus merge baseline).\n", projectName, target, len(plan))
	}
	return 0
}
func substituteProjectPlaceholders(root, name, description, repository, language string) error {
	for label, value := range map[string]string{"name": name, "description": description, "repository": repository, "language": language} {
		if strings.ContainsAny(value, "\r\n") {
			return fmt.Errorf("project %s must be a single line", label)
		}
	}
	raw := map[string]string{"{{PROJECT_NAME}}": name, "{{PROJECT_DESCRIPTION}}": description, "{{REPOSITORY_URL}}": repository, "{{PRIMARY_LANGUAGE}}": language, "{{DATE_YYYY_MM_DD}}": time.Now().Format("2006-01-02")}
	paths, err := padtemplate.Files()
	if err != nil {
		return err
	}
	for _, rel := range paths {
		path := filepath.Join(root, filepath.FromSlash(rel))
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		replacements := raw
		if filepath.ToSlash(rel) == "pad.yaml" {
			replacements = make(map[string]string, len(raw))
			for token, value := range raw {
				replacements[token] = strings.ReplaceAll(strings.ReplaceAll(value, `\`, `\\`), `"`, `\"`)
			}
		}
		updated := string(data)
		for token, value := range replacements {
			updated = strings.ReplaceAll(updated, token, value)
		}
		if updated != string(data) {
			if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}

func runValidate(args []string, out, errw io.Writer) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(errw)
	if e := fs.Parse(args); e != nil {
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintf(errw, "validate: %v\n", e)
		return 1
	}
	issues := validateProject(root)
	if len(issues) > 0 {
		for _, x := range issues {
			fmt.Fprintf(errw, "ERROR %s\n", x)
		}
		fmt.Fprintf(errw, "Validation failed: %d issue(s).\n", len(issues))
		return 1
	}
	fmt.Fprintf(out, "PAD project validation passed: %s\n", root)
	fmt.Fprintln(out, "Validated configuration, state, task/requirement/ADR metadata, identifiers, references, and dependency cycles.")
	return 0
}
func runStatus(args []string, out, errw io.Writer) int {
	root, e := discoverArg(args)
	if e != nil {
		fmt.Fprintf(errw, "status: %v\n", e)
		return 1
	}
	fmt.Fprintf(out, "PAD project: %s\n", root)
	tasks, e := findTasks(root)
	if e == nil {
		counts := map[string]int{}
		for _, t := range tasks {
			counts[t.status]++
		}
		fmt.Fprintf(out, "Tasks: %d total", len(tasks))
		for _, s := range []string{"backlog", "ready", "in_progress", "review", "blocked", "done", "cancelled"} {
			if counts[s] > 0 {
				fmt.Fprintf(out, " · %s=%d", s, counts[s])
			}
		}
		fmt.Fprintln(out)
	}
	state := filepath.Join(root, ".agent", "state", "current.yaml")
	if b, e := os.ReadFile(state); e == nil {
		for _, line := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(line, "project_status:") || strings.HasPrefix(line, "health:") || strings.HasPrefix(line, "current_focus:") || strings.HasPrefix(line, "active_task:") {
				fmt.Fprintln(out, line)
			}
		}
	}
	cmd := exec.Command("git", "-C", root, "status", "--short", "--branch")
	b, e := cmd.Output()
	if e != nil {
		fmt.Fprintf(out, "Git status unavailable: %v\n", e)
	} else {
		fmt.Fprintln(out, "Git:")
		s := strings.TrimSpace(string(b))
		if s == "" {
			fmt.Fprintln(out, "  (clean)")
		} else {
			for _, line := range strings.Split(s, "\n") {
				fmt.Fprintf(out, "  %s\n", line)
			}
		}
	}
	return 0
}
