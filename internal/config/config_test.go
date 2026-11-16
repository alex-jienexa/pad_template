package config

import "testing"

const validConfig = `version: "1.0"
project:
  name: Demo
  description: Example
  repository: ""
  default_branch: main
  primary_language: Go
  created_at: "2026-11-16"
agent_policy:
  autonomy: controlled
  require_approval_for:
    - destructive_git_operations
  allow_local_commits: true
  allow_remote_push: false
state:
  current_file: .agent/state/current.yaml
  handoff_file: .agent/state/HANDOFF.md
tasks:
  directory: .agent/tasks
  id_prefix: TASK
  default_assignee: unassigned
template:
  source: local
  version: "1.0.0"
  update_strategy: three_way_merge
  base_directory: .pad/template-base
`

func TestParseAndValidate(t *testing.T) {
	cfg, err := Parse(validConfig)
	if err != nil {
		t.Fatal(err)
	}
	if issues := cfg.Validate(); len(issues) != 0 {
		t.Fatalf("unexpected issues: %v", issues)
	}
	if cfg.AgentPolicy.AllowRemotePush {
		t.Fatal("remote push should be disabled")
	}
}
func TestParseRejectsIncorrectTypes(t *testing.T) {
	bad := `version: "1.0"
project: {}
agent_policy:
  autonomy: controlled
  require_approval_for: []
  allow_local_commits: "yes"
  allow_remote_push: false
`
	if _, err := Parse(bad); err == nil {
		t.Fatal("expected type error")
	}
}
