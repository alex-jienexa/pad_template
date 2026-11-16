package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestInitAcceptsFlagAfterTarget(t *testing.T) {
	target := filepath.Join(t.TempDir(), "project")
	var stdout, stderr bytes.Buffer
	code := Run([]string{"init", target, "--dry-run"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("Run exit code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Dry run") {
		t.Fatalf("expected dry-run message, got %q", stdout.String())
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("dry run created target: %v", err)
	}
}

func TestUnknownCommandReturnsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"unknown"}, &stdout, &stderr); code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("expected unknown-command error, got %q", stderr.String())
	}
}

func runForTest(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errout bytes.Buffer
	code := Run(args, &out, &errout)
	return code, out.String(), errout.String()
}

func TestFullProjectWorkflow(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatalf("init: %s", err)
	}
	if code, out, err := runForTest(t, "validate", root); code != 0 {
		t.Fatalf("validate: %s %s", out, err)
	}
	if code, out, err := runForTest(t, "task", "add", "--title", "Implement parser", root); code != 0 || !strings.Contains(out, "TASK-002") {
		t.Fatalf("task add: code=%d out=%s err=%s", code, out, err)
	}
	if code, out, err := runForTest(t, "task", "set", "TASK-002", "--status", "ready", root); code != 0 || !strings.Contains(out, "status=ready") {
		t.Fatalf("task set: code=%d out=%s err=%s", code, out, err)
	}
	if code, out, err := runForTest(t, "requirement", "add", "--title", "Validate all task references", root); code != 0 || !strings.Contains(out, "PAD-REQ-002") {
		t.Fatalf("requirement add: code=%d out=%s err=%s", code, out, err)
	}
	if code, out, err := runForTest(t, "checkpoint", "--message", "Workflow smoke test", root); code != 0 || !strings.Contains(out, "Created checkpoint") {
		t.Fatalf("checkpoint: code=%d out=%s err=%s", code, out, err)
	}
	if code, out, err := runForTest(t, "update", "--dry-run", root); code != 0 || !strings.Contains(out, "Plan:") {
		t.Fatalf("update dry run: code=%d out=%s err=%s", code, out, err)
	}
	if code, out, err := runForTest(t, "validate", root); code != 0 {
		t.Fatalf("final validate: %s %s", out, err)
	}
}

func TestTaskTransitionRejectsInvalidJump(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	code, _, err := runForTest(t, "task", "set", "TASK-001", "--status", "done", root)
	if code == 0 || !strings.Contains(err, "invalid task transition") {
		t.Fatalf("expected transition error, got code=%d err=%s", code, err)
	}
}

func TestMergeFileBytesCleanAndConflict(t *testing.T) {
	merged, clean, err := mergeFileBytes([]byte("local\nkeep\nbase-line\n"), []byte("base\nkeep\nbase-line\n"), []byte("base\nkeep\nupstream-line\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !clean {
		t.Fatalf("expected clean merge: %s", merged)
	}
	if !strings.Contains(string(merged), "local") || !strings.Contains(string(merged), "upstream-line") {
		t.Fatalf("merged result lost independent changes: %s", merged)
	}
	merged, clean, err = mergeFileBytes([]byte("local\n"), []byte("base\n"), []byte("upstream\n"))
	if err != nil {
		t.Fatal(err)
	}
	if clean {
		t.Fatalf("expected conflict, got %s", merged)
	}
	if !strings.Contains(string(merged), "<<<<<<<") {
		t.Fatalf("expected conflict markers, got %s", merged)
	}
}

func TestUpdatePerformsCleanThreeWayMerge(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	upstream, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	base := bytes.Replace(upstream, []byte("# PAD Template"), []byte("# PAD Template (old baseline)"), 1)
	current := append(append([]byte{}, base...), []byte("\nLocal project note.\n")...)
	basePath := filepath.Join(root, ".pad", "template-base", "README.md")
	if err = os.WriteFile(basePath, base, 0644); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(root, "README.md"), current, 0644); err != nil {
		t.Fatal(err)
	}
	code, out, errout := runForTest(t, "update", root)
	if code != 0 {
		t.Fatalf("update failed: %s %s", out, errout)
	}
	merged, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(merged, []byte("# {{PROJECT_NAME}}")) || !bytes.Contains(merged, []byte("Local project note.")) {
		t.Fatalf("clean merge lost upstream or local edits: %s", merged)
	}
	backups, err := filepath.Glob(filepath.Join(root, ".pad", "update-backups", "*", "local", "README.md"))
	if err != nil || len(backups) != 1 {
		t.Fatalf("expected a pre-update local backup, got %v (%v)", backups, err)
	}
}

func TestUpdatePreservesLocalFileOnConflict(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	base := []byte("# Baseline title\nshared body\n")
	local := []byte("# Local title\nshared body\n")
	if err := os.WriteFile(filepath.Join(root, ".pad", "template-base", "README.md"), base, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), local, 0644); err != nil {
		t.Fatal(err)
	}
	code, out, errout := runForTest(t, "update", root)
	if code != 0 {
		t.Fatalf("update failed: %s %s", out, errout)
	}
	actual, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, local) {
		t.Fatalf("local file was overwritten: %s", actual)
	}
	if _, err = os.Stat(filepath.Join(root, ".pad", "update-conflicts", "README.md.upstream")); err != nil {
		t.Fatalf("upstream conflict candidate missing: %v", err)
	}
	if !strings.Contains(out, "conflict") {
		t.Fatalf("expected conflict in output: %s", out)
	}
	mergeCandidate, err := os.ReadFile(filepath.Join(root, ".pad", "update-conflicts", "README.md.merge"))
	if err != nil || !bytes.Contains(mergeCandidate, []byte("<<<<<<<")) {
		t.Fatalf("conflict merge candidate missing markers: %s %v", mergeCandidate, err)
	}
	resolved := []byte("# Resolved title\nshared body\n")
	if err := os.WriteFile(filepath.Join(root, "README.md"), resolved, 0644); err != nil {
		t.Fatal(err)
	}
	if code, out, err := runForTest(t, "update", "--accept-resolved", root); code != 0 || !strings.Contains(out, "Accepted 1") {
		t.Fatalf("accept resolved: %d %s %s", code, out, err)
	}
	if code, out, err := runForTest(t, "update", "--dry-run", root); code != 0 || strings.Contains(out, "conflict  ") {
		t.Fatalf("resolved file still conflicts: %d %s %s", code, out, err)
	}
}

func TestInitSubstitutesProjectPlaceholdersButKeepsPristineBaseline(t *testing.T) {
	root := filepath.Join(t.TempDir(), "named-project")
	code, _, errMsg := runForTest(t, "init", "--name", "Example Project", "--description", "A demo.", root)
	if code != 0 {
		t.Fatal(errMsg)
	}
	config, err := os.ReadFile(filepath.Join(root, "pad.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(config), `name: "Example Project"`) || strings.Contains(string(config), "{{PROJECT_NAME}}") {
		t.Fatalf("project placeholders not substituted: %s", config)
	}
	baseline, err := os.ReadFile(filepath.Join(root, ".pad", "template-base", "pad.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(baseline), "{{PROJECT_NAME}}") {
		t.Fatalf("baseline should retain canonical template placeholders: %s", baseline)
	}
}

func TestDecisionAndAgentCommands(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	if code, out, err := runForTest(t, "decision", "add", "--title", "Choose storage engine", root); code != 0 || !strings.Contains(out, "ADR-0002") {
		t.Fatalf("decision add: %d %s %s", code, out, err)
	}
	if code, out, err := runForTest(t, "decision", "set", "ADR-0002", "--status", "accepted", root); code != 0 || !strings.Contains(out, "proposed -> accepted") {
		t.Fatalf("decision set: %d %s %s", code, out, err)
	}
	if code, out, err := runForTest(t, "agent", "list", root); code != 0 || !strings.Contains(out, "IMPLEMENTER") {
		t.Fatalf("agent list: %d %s %s", code, out, err)
	}
	if code, out, err := runForTest(t, "validate", root); code != 0 {
		t.Fatalf("validate after decision: %s %s", out, err)
	}
}

func TestInitRejectsSymlinkedBaselineBeforeWriting(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, ".pad")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	code, _, err := runForTest(t, "init", root)
	if code == 0 || !strings.Contains(err, "symlink") {
		t.Fatalf("expected symlink rejection, code=%d err=%s", code, err)
	}
	if _, err := os.Stat(filepath.Join(root, "pad.yaml")); !os.IsNotExist(err) {
		t.Fatalf("init wrote files before rejecting baseline symlink: %v", err)
	}
}

func TestInitAcceptsFlagsAfterTarget(t *testing.T) {
	root := filepath.Join(t.TempDir(), "after-flags")
	code, _, err := runForTest(t, "init", root, "--name", "After Flags")
	if code != 0 {
		t.Fatalf("init with trailing options failed: %s", err)
	}
	b, e := os.ReadFile(filepath.Join(root, "pad.yaml"))
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(string(b), `name: "After Flags"`) {
		t.Fatalf("name not applied: %s", b)
	}
}

func TestCheckpointCommitIsScopedAndLocal(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	git := func(args ...string) string {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		b, e := cmd.CombinedOutput()
		if e != nil {
			t.Fatalf("git %v: %s (%v)", args, b, e)
		}
		return string(b)
	}
	git("init")
	git("config", "user.name", "PAD Test")
	git("config", "user.email", "pad-test@example.invalid")
	if err := os.WriteFile(filepath.Join(root, "notes.txt"), []byte("keep uncommitted\n"), 0644); err != nil {
		t.Fatal(err)
	}
	code, out, err := runForTest(t, "checkpoint", "--message", "Scoped test", "--commit", root)
	if code != 0 {
		t.Fatalf("checkpoint commit: %s %s", out, err)
	}
	names := strings.Fields(git("show", "--pretty=format:", "--name-only", "HEAD"))
	sort.Strings(names)
	expected := []string{".agent/state/HANDOFF.md", ".agent/state/current.yaml"}
	entries, _ := filepath.Glob(filepath.Join(root, ".agent", "state", "checkpoints", "*.md"))
	if len(entries) != 1 {
		t.Fatalf("expected one checkpoint record, got %d", len(entries))
	}
	expected = append(expected, filepath.ToSlash(strings.TrimPrefix(entries[0], root+string(filepath.Separator))))
	sort.Strings(expected)
	if strings.Join(names, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("checkpoint commit included unexpected files: %v", names)
	}
	if _, err := os.Stat(filepath.Join(root, "notes.txt")); err != nil {
		t.Fatal(err)
	}
	if status := git("status", "--short"); !strings.Contains(status, "notes.txt") {
		t.Fatalf("unrelated file should remain uncommitted: %s", status)
	}
	handoff, _ := os.ReadFile(filepath.Join(root, ".agent", "state", "HANDOFF.md"))
	if !bytes.Contains(handoff, []byte("Latest checkpoint")) {
		t.Fatalf("handoff was not synchronized: %s", handoff)
	}
}

func TestCheckpointCommitRefusesExistingStagedChanges(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	git := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if b, e := cmd.CombinedOutput(); e != nil {
			t.Fatalf("git %v: %s (%v)", args, b, e)
		}
	}
	git("init")
	git("config", "user.name", "PAD Test")
	git("config", "user.email", "pad-test@example.invalid")
	if err := os.WriteFile(filepath.Join(root, "staged.txt"), []byte("staged"), 0644); err != nil {
		t.Fatal(err)
	}
	git("add", "staged.txt")
	code, _, err := runForTest(t, "checkpoint", "--message", "Should not commit", "--commit", root)
	if code == 0 || !strings.Contains(err, "index already contains staged changes") {
		t.Fatalf("expected staged-change refusal, got %d %s", code, err)
	}
}

func TestTaskReopeningRequiresReasonAndRecordsHistory(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	if code, _, err := runForTest(t, "task", "set", "TASK-001", "--status", "in_progress", root); code != 0 {
		t.Fatal(err)
	}
	if code, _, err := runForTest(t, "task", "set", "TASK-001", "--status", "review", root); code != 0 {
		t.Fatal(err)
	}
	if code, _, err := runForTest(t, "task", "set", "TASK-001", "--status", "done", root); code != 0 {
		t.Fatal(err)
	}
	if code, _, err := runForTest(t, "task", "set", "TASK-001", "--status", "ready", root); code == 0 || !strings.Contains(err, "--reason is required") {
		t.Fatalf("expected missing reason error: %d %s", code, err)
	}
	if code, _, err := runForTest(t, "task", "set", "TASK-001", "--status", "ready", "--reason", "Acceptance criteria changed", root); code != 0 {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, ".agent", "tasks", "TASK-001.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(b, []byte("Acceptance criteria changed")) {
		t.Fatalf("reopening reason not recorded: %s", b)
	}
}

func TestInitDoesNotResetExistingBaselineWithoutForce(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	base := filepath.Join(root, ".pad", "template-base", "pad.yaml")
	original, err := os.ReadFile(base)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(base, []byte("custom baseline"), 0644); err != nil {
		t.Fatal(err)
	}
	code, _, errMsg := runForTest(t, "init", root)
	if code == 0 || !strings.Contains(errMsg, "template baseline already exists") {
		t.Fatalf("expected baseline collision refusal, got %d %s", code, errMsg)
	}
	after, err := os.ReadFile(base)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, []byte("custom baseline")) {
		t.Fatalf("baseline was overwritten: %s", after)
	}
	if code, _, errMsg = runForTest(t, "init", "--force", root); code != 0 {
		t.Fatalf("explicit force should permit reinitialization: %s", errMsg)
	}
	restored, err := os.ReadFile(base)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(restored, original) {
		t.Fatal("force did not reset baseline to embedded template")
	}
}

func TestValidateDetectsTaskDependencyCycles(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	if code, _, err := runForTest(t, "task", "add", "--title", "Second task", root); code != 0 {
		t.Fatal(err)
	}
	first := filepath.Join(root, ".agent", "tasks", "TASK-001.md")
	second := filepath.Join(root, ".agent", "tasks", "TASK-002.md")
	for path, dep := range map[string]string{first: "TASK-002", second: "TASK-001"} {
		b, e := os.ReadFile(path)
		if e != nil {
			t.Fatal(e)
		}
		updated := strings.Replace(string(b), "depends_on: []", "depends_on:\n  - "+dep, 1)
		if updated == string(b) {
			t.Fatalf("depends_on field not found in %s", path)
		}
		if e = os.WriteFile(path, []byte(updated), 0644); e != nil {
			t.Fatal(e)
		}
	}
	code, _, err := runForTest(t, "validate", root)
	if code == 0 || !strings.Contains(err, "dependency cycle") {
		t.Fatalf("expected dependency cycle error, got %d %s", code, err)
	}
}

func TestAdoptRecordsMatchingGitHubTemplateBaseline(t *testing.T) {
	root := filepath.Join(t.TempDir(), "template-generated")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(root, ".pad", "template-base")); err != nil {
		t.Fatal(err)
	}
	// Simulate a repository created directly from the GitHub template rather than pad init.
	if code, _, _ := runForTest(t, "validate", root); code == 0 {
		t.Fatal("validation should require a merge baseline")
	}
	code, out, err := runForTest(t, "adopt", root)
	if code != 0 || !strings.Contains(out, "Adopted embedded template version") {
		t.Fatalf("adopt: %d %s %s", code, out, err)
	}
	if code, out, err = runForTest(t, "validate", root); code != 0 {
		t.Fatalf("validate after adopt: %s %s", out, err)
	}
	if code, _, err = runForTest(t, "adopt", root); code == 0 || !strings.Contains(err, "baseline already exists") {
		t.Fatalf("adopt should not overwrite baseline: %d %s", code, err)
	}
}

func TestAdoptRejectsTemplateVersionMismatch(t *testing.T) {
	root := filepath.Join(t.TempDir(), "template-generated")
	if code, _, err := runForTest(t, "init", root); code != 0 {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(root, ".pad", "template-base")); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(root, ".pad-template.yaml")
	b, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	b = bytes.Replace(b, []byte(`version: "1.0.0"`), []byte(`version: "0.9.0"`), 1)
	if err = os.WriteFile(manifest, b, 0644); err != nil {
		t.Fatal(err)
	}
	code, _, errMsg := runForTest(t, "adopt", root)
	if code == 0 || !strings.Contains(errMsg, "does not match") {
		t.Fatalf("expected version mismatch rejection, got %d %s", code, errMsg)
	}
}
