package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	padtemplate "github.com/alex-jienexa/pad_template/internal/template"
	"github.com/alex-jienexa/pad_template/internal/yamlmini"
)

// runAdopt records the embedded template as the common base for a project created
// outside `pad init` (for example, from a GitHub Template Repository).
func runAdopt(args []string, out, errw io.Writer) int {
	fs := flag.NewFlagSet("adopt", flag.ContinueOnError)
	fs.SetOutput(errw)
	if e := fs.Parse(args); e != nil {
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	base := filepath.Join(root, ".pad", "template-base")
	if _, e = os.Lstat(base); e == nil {
		fmt.Fprintf(errw, "adopt: template baseline already exists: %s\n", base)
		return 1
	} else if !os.IsNotExist(e) {
		fmt.Fprintln(errw, e)
		return 1
	}
	currentPath := filepath.Join(root, ".pad-template.yaml")
	current, e := os.ReadFile(currentPath)
	if e != nil {
		fmt.Fprintf(errw, "adopt: cannot read .pad-template.yaml: %v\n", e)
		return 1
	}
	embedded, e := padtemplate.Read(".pad-template.yaml")
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	currentVersion, e := manifestVersion(current)
	if e != nil {
		fmt.Fprintf(errw, "adopt: invalid project template manifest: %v\n", e)
		return 1
	}
	embeddedVersion, e := manifestVersion(embedded)
	if e != nil {
		fmt.Fprintf(errw, "adopt: invalid embedded template manifest: %v\n", e)
		return 1
	}
	if currentVersion != embeddedVersion {
		fmt.Fprintf(errw, "adopt: project template version %q does not match this CLI's embedded version %q; refusing to guess a merge base\n", currentVersion, embeddedVersion)
		return 1
	}
	// Refuse symlinked target paths before creating any baseline files.
	if e = projectPathSafe(root); e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	if e = padtemplate.Snapshot(root); e != nil {
		fmt.Fprintf(errw, "adopt: %v\n", e)
		return 1
	}
	fmt.Fprintf(out, "Adopted embedded template version %s as the local merge baseline. Project files were not modified.\n", currentVersion)
	return 0
}
func manifestVersion(data []byte) (string, error) {
	m, e := yamlmini.Parse(string(data))
	if e != nil {
		return "", e
	}
	section, ok := m["template"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("missing template mapping")
	}
	v, ok := section["version"].(string)
	if !ok || v == "" {
		return "", fmt.Errorf("missing template.version string")
	}
	return v, nil
}
func projectPathSafe(root string) error {
	for _, p := range []string{".pad", filepath.Join(".pad", "template-base")} {
		full := filepath.Join(root, p)
		info, e := os.Lstat(full)
		if e == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to adopt through symlink: %s", full)
		}
		if e != nil && !os.IsNotExist(e) {
			return e
		}
	}
	return nil
}
