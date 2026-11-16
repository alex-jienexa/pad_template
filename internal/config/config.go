// Package config loads and validates the typed PAD project configuration.
package config

import (
	"fmt"
	"os"

	"github.com/alex-jienexa/pad_template/internal/yamlmini"
)

type Config struct {
	Version     string
	Project     Project
	AgentPolicy AgentPolicy
	State       State
	Tasks       Tasks
	Template    Template
}
type Project struct{ Name, Description, Repository, DefaultBranch, PrimaryLanguage, CreatedAt string }
type AgentPolicy struct {
	Autonomy                           string
	RequireApprovalFor                 []string
	AllowLocalCommits, AllowRemotePush bool
}
type State struct{ CurrentFile, HandoffFile string }
type Tasks struct{ Directory, IDPrefix, DefaultAssignee string }
type Template struct{ Source, Version, UpdateStrategy, BaseDirectory string }

func Load(path string) (Config, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return Config{}, e
	}
	return Parse(string(b))
}
func Parse(src string) (Config, error) {
	m, e := yamlmini.Parse(src)
	if e != nil {
		return Config{}, e
	}
	var c Config
	c.Version, e = str(m, "version")
	if e != nil {
		return c, e
	}
	p, e := mapping(m, "project")
	if e != nil {
		return c, e
	}
	c.Project.Name, e = str(p, "name")
	if e != nil {
		return c, e
	}
	c.Project.Description, e = str(p, "description")
	if e != nil {
		return c, e
	}
	c.Project.Repository, e = str(p, "repository")
	if e != nil {
		return c, e
	}
	c.Project.DefaultBranch, e = str(p, "default_branch")
	if e != nil {
		return c, e
	}
	c.Project.PrimaryLanguage, e = str(p, "primary_language")
	if e != nil {
		return c, e
	}
	c.Project.CreatedAt, e = str(p, "created_at")
	if e != nil {
		return c, e
	}
	a, e := mapping(m, "agent_policy")
	if e != nil {
		return c, e
	}
	c.AgentPolicy.Autonomy, e = str(a, "autonomy")
	if e != nil {
		return c, e
	}
	c.AgentPolicy.RequireApprovalFor, e = strings(a, "require_approval_for")
	if e != nil {
		return c, e
	}
	c.AgentPolicy.AllowLocalCommits, e = boolean(a, "allow_local_commits")
	if e != nil {
		return c, e
	}
	c.AgentPolicy.AllowRemotePush, e = boolean(a, "allow_remote_push")
	if e != nil {
		return c, e
	}
	st, e := mapping(m, "state")
	if e != nil {
		return c, e
	}
	c.State.CurrentFile, e = str(st, "current_file")
	if e != nil {
		return c, e
	}
	c.State.HandoffFile, e = str(st, "handoff_file")
	if e != nil {
		return c, e
	}
	t, e := mapping(m, "tasks")
	if e != nil {
		return c, e
	}
	c.Tasks.Directory, e = str(t, "directory")
	if e != nil {
		return c, e
	}
	c.Tasks.IDPrefix, e = str(t, "id_prefix")
	if e != nil {
		return c, e
	}
	c.Tasks.DefaultAssignee, e = str(t, "default_assignee")
	if e != nil {
		return c, e
	}
	tp, e := mapping(m, "template")
	if e != nil {
		return c, e
	}
	c.Template.Source, e = str(tp, "source")
	if e != nil {
		return c, e
	}
	c.Template.Version, e = str(tp, "version")
	if e != nil {
		return c, e
	}
	c.Template.UpdateStrategy, e = str(tp, "update_strategy")
	if e != nil {
		return c, e
	}
	c.Template.BaseDirectory, e = str(tp, "base_directory")
	if e != nil {
		return c, e
	}
	return c, nil
}
func (c Config) Validate() []string {
	var out []string
	if c.Version == "" {
		out = append(out, "version is required")
	}
	if c.Project.Name == "" {
		out = append(out, "project.name is required")
	}
	if c.Project.DefaultBranch == "" {
		out = append(out, "project.default_branch is required")
	}
	if c.AgentPolicy.Autonomy != "read_only" && c.AgentPolicy.Autonomy != "controlled" && c.AgentPolicy.Autonomy != "delegated" {
		out = append(out, "agent_policy.autonomy must be read_only, controlled, or delegated")
	}
	if c.Template.UpdateStrategy != "three_way_merge" {
		out = append(out, "template.update_strategy must be three_way_merge")
	}
	if c.Tasks.Directory == "" || c.Tasks.IDPrefix == "" {
		out = append(out, "tasks.directory and tasks.id_prefix are required")
	}
	if c.Template.Source == "" || c.Template.Version == "" || c.Template.BaseDirectory != ".pad/template-base" {
		out = append(out, "template.source, template.version, and template.base_directory (.pad/template-base) are required")
	}
	if c.Project.Description == "" || c.Project.PrimaryLanguage == "" || c.Project.CreatedAt == "" {
		out = append(out, "project.description, project.primary_language, and project.created_at are required")
	}
	if c.State.CurrentFile == "" || c.State.HandoffFile == "" {
		out = append(out, "state.current_file and state.handoff_file are required")
	}
	return out
}
func mapping(m map[string]any, key string) (map[string]any, error) {
	v, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("missing mapping %q", key)
	}
	x, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s must be a mapping", key)
	}
	return x, nil
}
func str(m map[string]any, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("missing field %q", key)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", key)
	}
	return s, nil
}
func boolean(m map[string]any, key string) (bool, error) {
	v, ok := m[key]
	if !ok {
		return false, fmt.Errorf("missing field %q", key)
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("%s must be a boolean", key)
	}
	return b, nil
}
func strings(m map[string]any, key string) ([]string, error) {
	v, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("missing field %q", key)
	}
	a, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("%s must be a list", key)
	}
	out := make([]string, 0, len(a))
	for _, item := range a {
		s, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("%s items must be strings", key)
		}
		out = append(out, s)
	}
	return out, nil
}
