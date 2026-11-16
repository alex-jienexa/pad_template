package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alex-jienexa/pad_template/internal/project"
	padtemplate "github.com/alex-jienexa/pad_template/internal/template"
)

type updateItem struct{ path, action string }

func runUpdate(args []string, out, errw io.Writer) int {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.SetOutput(errw)
	dry := fs.Bool("dry-run", false, "show update plan without writing files")
	acceptResolved := fs.Bool("accept-resolved", false, "advance merge bases for manually resolved conflicts")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	if *acceptResolved {
		if *dry {
			fmt.Fprintln(errw, "update: --accept-resolved cannot be combined with --dry-run")
			return 2
		}
		n, err := acceptResolvedUpdates(root)
		if err != nil {
			fmt.Fprintln(errw, err)
			return 1
		}
		fmt.Fprintf(out, "Accepted %d manually resolved file(s); merge bases advanced to embedded upstream versions.\n", n)
		return 0
	}
	items, conflicts, e := planUpdate(root)
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	for _, it := range items {
		fmt.Fprintf(out, "%-18s %s\n", it.action, it.path)
	}
	fmt.Fprintf(out, "Plan: %d file(s), %d conflict(s)\n", len(items), conflicts)
	if *dry {
		return 0
	}
	backup, backupErr := createUpdateBackup(root, items)
	if backupErr != nil {
		fmt.Fprintf(errw, "update: cannot create safety backup: %v\n", backupErr)
		return 1
	}
	if backup != "" {
		fmt.Fprintf(out, "Safety backup: %s\n", filepath.ToSlash(backup))
	}
	applied := 0
	for _, it := range items {
		if it.action == "update" || it.action == "add" || it.action == "merge" {
			data, e := padtemplate.Read(it.path)
			if it.action == "merge" {
				data, e = mergeCurrentWithUpstream(root, it.path)
			}
			if e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			if e = rejectSymlinkPath(root, it.path); e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			dest := filepath.Join(root, filepath.FromSlash(it.path))
			if e = os.MkdirAll(filepath.Dir(dest), 0o755); e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			if e = os.WriteFile(dest, data, 0o644); e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			if e = writeBase(root, it.path, data); e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			applied++
		} else if it.action == "unchanged" {
			data, e := padtemplate.Read(it.path)
			if e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			if e = writeBase(root, it.path, data); e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
		} else if it.action == "conflict" {
			data, e := padtemplate.Read(it.path)
			if e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			conflictRel := filepath.ToSlash(filepath.Join(".pad", "update-conflicts", filepath.FromSlash(it.path)+".upstream"))
			if e = rejectSymlinkPath(root, conflictRel); e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			conflict := filepath.Join(root, filepath.FromSlash(conflictRel))
			if e = os.MkdirAll(filepath.Dir(conflict), 0o755); e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			if e = os.WriteFile(conflict, data, 0o644); e != nil {
				fmt.Fprintln(errw, e)
				return 1
			}
			if merged, me := mergeCurrentWithUpstream(root, it.path); me == nil {
				mergeRel := filepath.ToSlash(filepath.Join(".pad", "update-conflicts", filepath.FromSlash(it.path)+".merge"))
				if e = rejectSymlinkPath(root, mergeRel); e == nil {
					mergePath := filepath.Join(root, filepath.FromSlash(mergeRel))
					_ = os.MkdirAll(filepath.Dir(mergePath), 0o755)
					_ = os.WriteFile(mergePath, merged, 0o644)
				}
			}
		} else if it.action == "baseline-missing" {
			fmt.Fprintf(out, "Baseline missing for %s; local file and baseline were preserved.\n", it.path)
		}
	}
	fmt.Fprintf(out, "Applied %d automatic update(s).", applied)
	if conflicts > 0 {
		fmt.Fprintf(out, " %d conflict(s) need manual resolution; local files were preserved.", conflicts)
	}
	fmt.Fprintln(out)
	return 0
}
func createUpdateBackup(root string, items []updateItem) (string, error) {
	stamp := time.Now().Format("20060102-150405.000000000")
	relRoot := filepath.Join(".pad", "update-backups", stamp)
	needsBackup := false
	for _, it := range items {
		if it.action == "update" || it.action == "merge" || it.action == "unchanged" {
			needsBackup = true
			break
		}
	}
	if !needsBackup {
		return "", nil
	}
	if err := rejectSymlinkPath(root, relRoot); err != nil {
		return "", err
	}
	for _, it := range items {
		if it.action != "update" && it.action != "merge" && it.action != "unchanged" {
			continue
		}
		for _, source := range []struct{ rel, category string }{{it.path, "local"}, {filepath.ToSlash(filepath.Join(".pad", "template-base", filepath.FromSlash(it.path))), "base"}} {
			data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(source.rel)))
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return "", err
			}
			destRel := filepath.ToSlash(filepath.Join(relRoot, source.category, filepath.FromSlash(it.path)))
			if err := rejectSymlinkPath(root, destRel); err != nil {
				return "", err
			}
			dest := filepath.Join(root, filepath.FromSlash(destRel))
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return "", err
			}
			if err := os.WriteFile(dest, data, 0o644); err != nil {
				return "", err
			}
		}
	}
	readme := "# PAD template update backup\n\nThis directory contains pre-update copies of project files and merge bases. It was created automatically before `pad update`. Keep it until the update has been reviewed and validated.\n"
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(filepath.Join(relRoot, "README.md"))), []byte(readme), 0o644); err != nil {
		return "", err
	}
	return relRoot, nil
}

func planUpdate(root string) ([]updateItem, int, error) {
	paths, e := padtemplate.Files()
	if e != nil {
		return nil, 0, e
	}
	sort.Strings(paths)
	var items []updateItem
	conflicts := 0
	for _, name := range paths {
		safe, e := project.SafeRelativePath(name)
		if e != nil {
			return nil, 0, e
		}
		if e := rejectSymlinkPath(root, safe); e != nil {
			return nil, 0, e
		}
		baseRel := filepath.ToSlash(filepath.Join(".pad", "template-base", filepath.FromSlash(safe)))
		if e := rejectSymlinkPath(root, baseRel); e != nil {
			return nil, 0, e
		}
		curPath := filepath.Join(root, filepath.FromSlash(safe))
		basePath := filepath.Join(root, ".pad", "template-base", filepath.FromSlash(safe))
		current, ce := os.ReadFile(curPath)
		base, be := os.ReadFile(basePath)
		upstream, ue := padtemplate.Read(name)
		if ue != nil {
			return nil, 0, ue
		}
		if os.IsNotExist(ce) {
			if os.IsNotExist(be) {
				items = append(items, updateItem{name, "add"})
			} else {
				items = append(items, updateItem{name, "preserve-deletion"})
			}
			continue
		}
		if ce != nil {
			return nil, 0, ce
		}
		if os.IsNotExist(be) {
			if bytes.Equal(current, upstream) {
				items = append(items, updateItem{name, "unchanged"})
			} else {
				items = append(items, updateItem{name, "baseline-missing"})
				conflicts++
			}
			continue
		}
		if be != nil {
			return nil, 0, be
		}
		switch {
		case bytes.Equal(current, base) && !bytes.Equal(upstream, base):
			items = append(items, updateItem{name, "update"})
		case bytes.Equal(current, upstream):
			items = append(items, updateItem{name, "unchanged"})
		case bytes.Equal(upstream, base):
			items = append(items, updateItem{name, "local-only"})
		default:
			if _, clean, mergeErr := mergeFileBytes(current, base, upstream); mergeErr == nil && clean {
				items = append(items, updateItem{name, "merge"})
			} else {
				items = append(items, updateItem{name, "conflict"})
				conflicts++
			}
		}
	}
	// Surface files removed from the new template, but never delete them automatically.
	canonical := make(map[string]bool, len(paths))
	for _, name := range paths {
		canonical[filepath.ToSlash(name)] = true
	}
	baseRoot := filepath.Join(root, ".pad", "template-base")
	walkErr := filepath.WalkDir(baseRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				return nil
			}
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(baseRoot, path)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(rel)
		if !canonical[name] {
			items = append(items, updateItem{name, "removed-upstream"})
		}
		return nil
	})
	if walkErr != nil && !os.IsNotExist(walkErr) {
		return nil, conflicts, walkErr
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].path == items[j].path {
			return items[i].action < items[j].action
		}
		return items[i].path < items[j].path
	})
	return items, conflicts, nil
}
func writeBase(root, name string, data []byte) error {
	dest := filepath.Join(root, ".pad", "template-base", filepath.FromSlash(name))
	if e := os.MkdirAll(filepath.Dir(dest), 0o755); e != nil {
		return e
	}
	return os.WriteFile(dest, data, 0o644)
}

func mergeCurrentWithUpstream(root, name string) ([]byte, error) {
	current, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	if err != nil {
		return nil, err
	}
	base, err := os.ReadFile(filepath.Join(root, ".pad", "template-base", filepath.FromSlash(name)))
	if err != nil {
		return nil, err
	}
	upstream, err := padtemplate.Read(name)
	if err != nil {
		return nil, err
	}
	merged, _, err := mergeFileBytes(current, base, upstream)
	if err != nil {
		return nil, err
	}
	return merged, nil
}

// mergeFileBytes delegates the diff3 operation to Git's local merge-file utility.
// It never modifies source files; all inputs are temporary files.
func mergeFileBytes(current, base, upstream []byte) ([]byte, bool, error) {
	paths := make([]string, 0, 3)
	for _, data := range [][]byte{current, base, upstream} {
		f, err := os.CreateTemp("", "pad-merge-*")
		if err != nil {
			return nil, false, err
		}
		paths = append(paths, f.Name())
		if _, err = f.Write(data); err != nil {
			f.Close()
			return nil, false, err
		}
		if err = f.Close(); err != nil {
			return nil, false, err
		}
	}
	defer func() {
		for _, p := range paths {
			_ = os.Remove(p)
		}
	}()
	cmd := exec.Command("git", "merge-file", "-p", paths[0], paths[1], paths[2])
	output, err := cmd.Output()
	if err == nil {
		return output, true, nil
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 1 {
		return output, false, nil
	}
	return nil, false, fmt.Errorf("git merge-file failed: %w", err)
}

func rejectSymlinkPath(root, relative string) error {
	safe, err := project.SafeRelativePath(relative)
	if err != nil {
		return err
	}
	current := root
	for _, part := range strings.Split(filepath.Clean(safe), string(filepath.Separator)) {
		current = filepath.Join(current, part)
		info, e := os.Lstat(current)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return e
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to update through symlink: %s", current)
		}
	}
	return nil
}

// acceptResolvedUpdates advances only the merge base for conflict candidates that
// the user explicitly reviewed and resolved in the corresponding project files.
func acceptResolvedUpdates(root string) (int, error) {
	if err := rejectSymlinkPath(root, filepath.Join(".pad", "update-conflicts")); err != nil {
		return 0, err
	}
	dir := filepath.Join(root, ".pad", "update-conflicts")
	var accepted int
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".upstream") {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		conflictRel := filepath.ToSlash(filepath.Join(".pad", "update-conflicts", rel))
		if err := rejectSymlinkPath(root, conflictRel); err != nil {
			return err
		}
		name := strings.TrimSuffix(filepath.ToSlash(rel), ".upstream")
		if err := rejectSymlinkPath(root, name); err != nil {
			return err
		}
		current := filepath.Join(root, filepath.FromSlash(name))
		if _, err := os.Stat(current); err != nil {
			return fmt.Errorf("cannot accept %s: resolved project file is missing", name)
		}
		baseRel := filepath.ToSlash(filepath.Join(".pad", "template-base", filepath.FromSlash(name)))
		if err := rejectSymlinkPath(root, baseRel); err != nil {
			return err
		}
		candidate, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := writeBase(root, name, candidate); err != nil {
			return err
		}
		accepted++
		return nil
	})
	if os.IsNotExist(err) {
		return 0, fmt.Errorf("no update conflicts found")
	}
	return accepted, err
}
