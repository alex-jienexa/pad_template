# Getting started with PAD

## 1. Initialize

Use the GitHub **Use this template** workflow or the local CLI. Preview initialization with `pad init --dry-run ./path/to/project`, then run `pad init ./path/to/project`. The CLI embeds the starter and refuses to overwrite colliding files unless `--force` is explicitly supplied. Review the complete plan before using `--force`. Do not initialize over an existing project without reviewing the preview. A template-generated repository starts with its own Git history; PAD's update mechanism is designed to compare the recorded template base with upstream and local files.

## 2. Configure

Edit `pad.yaml`: set project name, description, repository, primary language/toolchain, default branch, and agent policy. Replace every `{{...}}` placeholder and remove examples that do not apply. Do not leave contradictory placeholder values in a live project.

## 3. Establish context

Read `.agent/README.md`. Complete `.agent/context/PROJECT.md` with verified facts future sessions need: domain definitions, constraints, architecture, important commands, and authoritative references. Avoid copying entire source files or chat transcripts into memory documents.

## 4. Establish requirements and tasks

Write requirements as stable, testable obligations. Decompose implementation into bounded tasks with explicit dependencies, ownership, acceptance criteria, and verification. Start by adapting `TASK-001.md` or replacing it with the project's first real task. Replace the sample ADR if the project does not want to adopt the PAD baseline as an explicit decision.

## 5. Start a session

Follow `.agent/protocols/SESSION_START.md`. Inspect Git state first. Confirm the task, scope, permissions, and acceptance criteria before editing.

## 6. Execute and verify

Follow `.agent/protocols/TASK_EXECUTION.md`. Add tests alongside behavior. Report exact commands and outcomes; do not label unrun checks as passed.

## 7. Finish and hand off

Follow `.agent/protocols/COMPLETION.md` and `.agent/protocols/HANDOFF.md`. Update the task, `state/current.yaml`, and the single `state/HANDOFF.md`. Create a logical checkpoint only when authorized. Push to a remote only when network access is permitted and the task/policy allows it; verify the remote result before claiming backup.

## 8. Update the template

Follow `.agent/protocols/THREE_WAY_TEMPLATE_UPDATE.md`. The updater must preserve local customizations, preview changes, retain a recoverable base, and surface conflicts. Upstream deletions require manual review by default.
