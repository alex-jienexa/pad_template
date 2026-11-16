package project

import (
	"path/filepath"
	"testing"
)

func TestSafeRelativePath(t *testing.T) {
	for _, input := range []string{"../outside", "a/../../outside", "/absolute", "."} {
		if _, err := SafeRelativePath(input); err == nil {
			t.Errorf("SafeRelativePath(%q) expected error", input)
		}
	}
	got, err := SafeRelativePath(".agent/tasks/TASK-001.md")
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.Join(".agent", "tasks", "TASK-001.md") {
		t.Fatalf("unexpected safe path: %q", got)
	}
}
