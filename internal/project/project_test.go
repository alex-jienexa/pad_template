package project

import (
	"path/filepath"
	"testing"
)

func TestSafeRelativePath(t *testing.T) {
	for _, input := range []string{
		"../outside",
		"a/../../outside",
		"/absolute",
		".",
		"C:/absolute",
		"C:relative",
		`C:\absolute`,
		`\\server\share`,
	} {
		if _, err := SafeRelativePath(input); err == nil {
			t.Errorf("SafeRelativePath(%q) expected error", input)
		}
	}
	got, err := SafeRelativePath(".agent/tasks/TASK-001.md")
	if err != nil {
		t.Fatal(err)
	}
	windowsPath, err := SafeRelativePath(
		`.agent\state\current.yaml`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if windowsPath != filepath.Join(
		".agent",
		"state",
		"current.yaml",
	) {
		t.Fatalf("unexpected Windows path: %q", windowsPath)
	}
	if got != filepath.Join(".agent", "tasks", "TASK-001.md") {
		t.Fatalf("unexpected safe path: %q", got)
	}
}
