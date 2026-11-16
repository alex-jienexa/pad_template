package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitDryRunDoesNotCreateTarget(t *testing.T) {
	target := filepath.Join(t.TempDir(), "new-project")
	plan, err := Init(target, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) == 0 {
		t.Fatal("expected a non-empty plan")
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("dry run created target: stat error = %v", err)
	}
}

func TestInitCopiesHiddenFiles(t *testing.T) {
	target := filepath.Join(t.TempDir(), "project")
	if _, err := Init(target, false, false); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{".gitignore", ".pad-template.yaml", ".agent/protocols/SESSION_START.md"} {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(name))); err != nil {
			t.Errorf("expected embedded file %s: %v", name, err)
		}
	}
}

func TestInitPreservesCollidingFiles(t *testing.T) {
	target := t.TempDir()
	path := filepath.Join(target, "AGENTS.md")
	original := []byte("user content\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Init(target, false, false); err == nil {
		t.Fatal("expected collision error")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatalf("existing file was modified: %q", got)
	}
}

func TestInitRejectsSymlink(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "project")
	outside := filepath.Join(parent, "outside")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(target, ".agent")); err != nil {
		t.Skipf("symlink not supported: %v", err)
	}
	if _, err := Init(target, true, false); err == nil {
		t.Fatal("expected symlink rejection")
	}
}
