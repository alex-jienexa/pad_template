package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func runAgent(args []string, out, errw io.Writer) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(out, "Usage: pad agent <list|show> [profile] [path]")
		return 0
	}
	rootArgs := args[1:]
	root, e := discoverArgForAgent(args[0], rootArgs)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	dir := filepath.Join(root, ".agent", "agents")
	switch args[0] {
	case "list":
		entries, e := os.ReadDir(dir)
		if e != nil {
			fmt.Fprintln(errw, e)
			return 1
		}
		var names []string
		for _, ent := range entries {
			if !ent.IsDir() && strings.HasSuffix(ent.Name(), ".md") && ent.Name() != "README.md" {
				names = append(names, strings.TrimSuffix(ent.Name(), ".md"))
			}
		}
		sort.Strings(names)
		fmt.Fprintf(out, "Available agent profiles (%d):\n", len(names))
		for _, n := range names {
			fmt.Fprintf(out, "- %s\n", n)
		}
		return 0
	case "show":
		if len(rootArgs) < 1 {
			fmt.Fprintln(errw, "agent show requires a profile name")
			return 2
		}
		name := rootArgs[0]
		if filepath.Base(name) != name {
			fmt.Fprintln(errw, "invalid profile name")
			return 2
		}
		p := filepath.Join(dir, name)
		if !strings.HasSuffix(p, ".md") {
			p += ".md"
		}
		b, e := os.ReadFile(p)
		if e != nil {
			fmt.Fprintln(errw, e)
			return 1
		}
		_, _ = out.Write(b)
		return 0
	default:
		fmt.Fprintf(errw, "agent: unknown subcommand %q\n", args[0])
		return 2
	}
}
func discoverArgForAgent(command string, args []string) (string, error) {
	if command == "show" && len(args) > 0 {
		args = args[1:]
	}
	return discoverArg(args)
}
