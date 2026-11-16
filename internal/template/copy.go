package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	padtemplate "github.com/alex-jienexa/pad_template"
	"github.com/alex-jienexa/pad_template/internal/project"
)

// PlanEntry describes one file from the embedded starter and its destination state.
type PlanEntry struct {
	Path      string
	Exists    bool
	WillWrite bool
}

// Plan returns the files that initialization would create or replace.
func Plan(target string, force bool) ([]PlanEntry, error) {
	if info, err := os.Lstat(target); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("refusing to initialize through symlink target: %s", target)
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	// Never reset an existing merge base implicitly. Reinitialization must be explicit.
	baselinePath := filepath.Join(target, ".pad", "template-base")
	if _, statErr := os.Lstat(baselinePath); statErr == nil && !force {
		return nil, fmt.Errorf("template baseline already exists at %s; use --force only if you intend to replace managed files and reset the baseline", baselinePath)
	} else if statErr != nil && !os.IsNotExist(statErr) {
		return nil, statErr
	}
	// Preflight all merge-baseline paths before any project file is written.
	paths, err := Files()
	if err != nil {
		return nil, err
	}
	for _, name := range paths {
		if err := rejectSymlinkComponents(target, filepath.Join(".pad", "template-base", filepath.FromSlash(name))); err != nil {
			return nil, err
		}
	}
	var plan []PlanEntry
	err = fs.WalkDir(padtemplate.TemplateFS, "template", func(source string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel := source[len("template/"):]
		safe, err := project.SafeRelativePath(rel)
		if err != nil {
			return err
		}
		destination := filepath.Join(target, safe)
		if err := rejectSymlinkComponents(target, safe); err != nil {
			return err
		}
		info, statErr := os.Lstat(destination)
		exists := statErr == nil
		if statErr != nil && !os.IsNotExist(statErr) {
			return statErr
		}
		if exists && !info.Mode().IsRegular() {
			return fmt.Errorf("template file path is occupied by a non-regular file: %s", destination)
		}
		plan = append(plan, PlanEntry{Path: filepath.ToSlash(safe), Exists: exists, WillWrite: !exists || force})
		return nil
	})
	return plan, err
}

// rejectSymlinkComponents prevents writes from following links out of the project root.
func rejectSymlinkComponents(root, relative string) error {
	current := root
	parts := strings.Split(filepath.Clean(relative), string(filepath.Separator))
	for _, part := range parts {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to write through symlink: %s", current)
		}
	}
	return nil
}

// Init copies the embedded starter into target. Existing files are preserved unless force is true.
// A dry run never creates the target directory or writes any file.
func Init(target string, force, dryRun bool) ([]PlanEntry, error) {
	abs, err := filepath.Abs(target)
	if err != nil {
		return nil, err
	}
	plan, err := Plan(abs, force)
	if err != nil {
		return nil, err
	}
	for _, item := range plan {
		if item.Exists && !force {
			return plan, fmt.Errorf("target contains existing template path %q; use --force only after reviewing the plan", item.Path)
		}
	}
	if dryRun {
		return plan, nil
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return plan, err
	}
	for _, item := range plan {
		if !item.WillWrite {
			continue
		}
		data, err := Read(item.Path)
		if err != nil {
			return plan, err
		}
		destination := filepath.Join(abs, filepath.FromSlash(item.Path))
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return plan, err
		}
		if err := os.WriteFile(destination, data, 0o644); err != nil {
			return plan, err
		}
	}
	if err := Snapshot(abs); err != nil {
		return plan, fmt.Errorf("write template baseline: %w", err)
	}
	return plan, nil
}

// Files returns canonical template paths in stable order.
func Files() ([]string, error) {
	var paths []string
	err := fs.WalkDir(padtemplate.TemplateFS, "template", func(source string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		paths = append(paths, filepath.ToSlash(source[len("template/"):]))
		return nil
	})
	return paths, err
}

// Read returns the embedded canonical file for a relative template path.
func Read(name string) ([]byte, error) {
	safe, err := project.SafeRelativePath(name)
	if err != nil {
		return nil, err
	}
	return padtemplate.TemplateFS.ReadFile("template/" + filepath.ToSlash(safe))
}

// Snapshot stores the exact starter version used for a project. It is the merge base
// for future local template updates and is intentionally kept outside managed paths.
func Snapshot(root string) error {
	paths, err := Files()
	if err != nil {
		return err
	}
	for _, rel := range []string{".pad", filepath.Join(".pad", "template-base")} {
		p := filepath.Join(root, rel)
		info, statErr := os.Lstat(p)
		if statErr == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to write template baseline through symlink: %s", p)
		}
		if statErr != nil && !os.IsNotExist(statErr) {
			return statErr
		}
	}
	for _, name := range paths {
		rel := filepath.Join(".pad", "template-base", filepath.FromSlash(name))
		if err := rejectSymlinkComponents(root, rel); err != nil {
			return err
		}
		data, err := Read(name)
		if err != nil {
			return err
		}
		dest := filepath.Join(root, ".pad", "template-base", filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
