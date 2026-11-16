package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const handoffCheckpointStart = "<!-- PAD:CHECKPOINT:START -->"
const handoffCheckpointEnd = "<!-- PAD:CHECKPOINT:END -->"

func runCheckpoint(args []string, out, errw io.Writer) int {
	fs := flag.NewFlagSet("checkpoint", flag.ContinueOnError)
	fs.SetOutput(errw)
	message := fs.String("message", "", "checkpoint summary (required)")
	commit := fs.Bool("commit", false, "create a local commit containing only checkpoint/state/handoff files")
	commitMessage := fs.String("commit-message", "", "optional local commit message")
	if e := fs.Parse(args); e != nil {
		return 2
	}
	if strings.TrimSpace(*message) == "" {
		fmt.Fprintln(errw, "checkpoint: --message is required")
		return 2
	}
	if strings.ContainsAny(*message, "\r\n") {
		fmt.Fprintln(errw, "checkpoint: --message must be a single line")
		return 2
	}
	root, e := discoverArg(fs.Args())
	if e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	if *commit {
		if _, e := runGit(root, "rev-parse", "--is-inside-work-tree"); e != nil {
			fmt.Fprintf(errw, "checkpoint: --commit requires an existing Git repository: %v\n", e)
			return 1
		}
		if _, e := runGit(root, "diff", "--cached", "--quiet"); e != nil {
			fmt.Fprintln(errw, "checkpoint: refusing to commit because the index already contains staged changes; commit or unstage them first")
			return 1
		}
	}
	now := time.Now()
	stamp := now.Format("20060102-150405.000000000")
	rel := filepath.Join(".agent", "state", "checkpoints", stamp+".md")
	stateRel := filepath.Join(".agent", "state", "current.yaml")
	handoffRel := filepath.Join(".agent", "state", "HANDOFF.md")
	for _, p := range []string{rel, stateRel, handoffRel} {
		if e := rejectSymlinkPath(root, p); e != nil {
			fmt.Fprintln(errw, e)
			return 1
		}
	}
	full := filepath.Join(root, rel)
	if e = os.MkdirAll(filepath.Dir(full), 0o755); e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	body := fmt.Sprintf("# Checkpoint %s\n\n- Created: %s\n- Summary: %s\n\n## Git status at checkpoint\n\n```text\n%s\n```\n\n## Resume instructions\n\n1. Read `AGENTS.md` and `.agent/protocols/RECOVERY.md`.\n2. Reconcile this checkpoint with the current Git state.\n3. Verify uncommitted changes before continuing.\n", stamp, now.Format(time.RFC3339), *message, gitStatus(root))
	if e = os.WriteFile(full, []byte(body), 0o644); e != nil {
		fmt.Fprintln(errw, e)
		return 1
	}
	state := filepath.Join(root, stateRel)
	b, e := os.ReadFile(state)
	if e != nil {
		fmt.Fprintf(errw, "checkpoint created, but cannot update current.yaml: %v\n", e)
		return 1
	}
	s := string(b)
	s = replaceYAMLScalar(s, "last_checkpoint", stamp)
	s = replaceYAMLScalar(s, "updated_at", now.Format(time.RFC3339))
	if e = os.WriteFile(state, []byte(s), 0o644); e != nil {
		fmt.Fprintf(errw, "checkpoint created, but cannot update current.yaml: %v\n", e)
		return 1
	}
	handoff := filepath.Join(root, handoffRel)
	hb, e := os.ReadFile(handoff)
	if e != nil {
		fmt.Fprintf(errw, "checkpoint created, but cannot update HANDOFF.md: %v\n", e)
		return 1
	}
	section := fmt.Sprintf("%s\n\n**Latest checkpoint:** [%s](checkpoints/%s.md)\n\n**Summary:** %s\n\n%s", handoffCheckpointStart, stamp, stamp, *message, handoffCheckpointEnd)
	hs := string(hb)
	if strings.Contains(hs, handoffCheckpointStart) && strings.Contains(hs, handoffCheckpointEnd) {
		start := strings.Index(hs, handoffCheckpointStart)
		end := strings.Index(hs, handoffCheckpointEnd) + len(handoffCheckpointEnd)
		hs = hs[:start] + section + hs[end:]
	} else {
		hs = strings.TrimRight(hs, "\n") + "\n\n" + section + "\n"
	}
	if e = os.WriteFile(handoff, []byte(hs), 0o644); e != nil {
		fmt.Fprintf(errw, "checkpoint created, but cannot update HANDOFF.md: %v\n", e)
		return 1
	}
	fmt.Fprintf(out, "Created checkpoint: %s\nUpdated current.yaml and HANDOFF.md\n", filepath.ToSlash(rel))
	if *commit {
		msg := *commitMessage
		if msg == "" {
			msg = "chore(pad): checkpoint " + stamp
		}
		addArgs := []string{"add", "--", filepath.ToSlash(rel), filepath.ToSlash(stateRel), filepath.ToSlash(handoffRel)}
		if _, e := runGit(root, addArgs...); e != nil {
			fmt.Fprintf(errw, "checkpoint files created, but git add failed: %v\n", e)
			return 1
		}
		if _, e := runGit(root, "commit", "-m", msg); e != nil {
			fmt.Fprintf(errw, "checkpoint files staged, but git commit failed: %v\n", e)
			return 1
		}
		hash, _ := runGit(root, "rev-parse", "--short", "HEAD")
		fmt.Fprintf(out, "Created local checkpoint commit %s\n", strings.TrimSpace(hash))
	}
	return 0
}
func replaceYAMLScalar(src, key, value string) string {
	re := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(key) + `:\s*.*$`)
	replacement := key + ": " + value
	if re.MatchString(src) {
		return re.ReplaceAllString(src, replacement)
	}
	return src + "\n" + replacement + "\n"
}
func gitStatus(root string) string {
	b, e := runGit(root, "status", "--short", "--branch")
	if e != nil {
		return "Git status unavailable: " + e.Error()
	}
	s := strings.TrimSpace(b)
	if s == "" {
		return "(clean)"
	}
	return s
}
func runGit(root string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	b, e := cmd.CombinedOutput()
	if e != nil {
		return string(b), fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), e, strings.TrimSpace(string(b)))
	}
	return string(b), nil
}
