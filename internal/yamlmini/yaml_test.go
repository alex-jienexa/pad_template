package yamlmini

import "testing"

func TestParseNestedMappingsAndSequences(t *testing.T) {
	got, err := Parse("version: '1.0'\nagent_policy:\n  autonomy: controlled\n  allow_remote_push: false\n  require_approval_for:\n    - production\n    - secrets\nstate:\n  updated_at: null\n")
	if err != nil {
		t.Fatal(err)
	}
	policy := got["agent_policy"].(map[string]any)
	if policy["autonomy"] != "controlled" || policy["allow_remote_push"] != false {
		t.Fatalf("unexpected policy: %#v", policy)
	}
	list := policy["require_approval_for"].([]any)
	if len(list) != 2 || list[1] != "secrets" {
		t.Fatalf("unexpected list: %#v", list)
	}
	state := got["state"].(map[string]any)
	if state["updated_at"] != nil {
		t.Fatalf("expected null, got %#v", state["updated_at"])
	}
}

func TestParseRejectsDuplicateKeysAndUnsupportedFeatures(t *testing.T) {
	for _, src := range []string{"a: 1\na: 2\n", "value: &anchor hello\n", "bad line\n"} {
		if _, err := Parse(src); err == nil {
			t.Errorf("Parse(%q) expected error", src)
		}
	}
}
