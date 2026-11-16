package records

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAndParseRoundTripEmptyLists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "TASK-001.md")
	meta := map[string]any{"id": "TASK-001", "title": "Example task", "status": "backlog", "priority": "medium", "assignee": "unassigned", "created_at": "2026-11-16", "updated_at": "2026-11-16", "depends_on": []string{}, "requirement_ids": []string{}, "labels": []string{}}
	if err := Write(path, meta, "# Example\n"); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := ParseMarkdown(b)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"depends_on", "requirement_ids", "labels"} {
		v, ok := doc.Metadata[key].([]any)
		if !ok || len(v) != 0 {
			t.Fatalf("%s should be an empty list, got %#v", key, doc.Metadata[key])
		}
	}
}
func TestParseMarkdownRequiresClosingDelimiter(t *testing.T) {
	if _, err := ParseMarkdown([]byte("---\nid: x\n# no close")); err == nil {
		t.Fatal("expected missing delimiter error")
	}
}
