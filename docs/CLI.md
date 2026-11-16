# PAD CLI reference

The PAD CLI is a local project-management utility. It does not fetch templates, call GitHub APIs, configure remotes, push commits, or publish releases. The embedded template is the CLI's update source.

## Build and install

Requires Go 1.23+ and Git. Git is used for read-only status and local `git merge-file` operations.

```bash
go test ./...
go vet ./...
go build -o ./bin/pad ./cmd/pad
# Optional: install into the Go bin directory
go install ./cmd/pad
```

## Commands

### `pad init [path]`

Initializes a project from the embedded template and saves its pristine merge base under `.pad/template-base/`.

```bash
pad init --dry-run ./my-project
pad init --name "Example Project" --description "Project summary" ./my-project
```

- `--dry-run` previews the file plan without creating directories or files.
- `--name`, `--description`, `--repository`, and `--language` fill project placeholders. The name defaults to the target directory name; other optional values have safe defaults.
- The pristine baseline retains the original template placeholders; project-specific substitutions do not contaminate it.
- Existing template paths cause the operation to stop before writing unless `--force` is supplied.
- `--force` overwrites colliding template files and resets an existing merge baseline. It is intentionally destructive to managed customizations; it does not delete unrelated files.
- Symbolic links on managed paths and the `.pad` baseline path are rejected.
- Does not run `git init`, configure remotes, or make commits.

### `pad adopt [path]`

Records the CLI's embedded template as the merge base for a project created directly from a matching GitHub Template Repository rather than with `pad init`. It requires `.pad-template.yaml` to have the same `template.version` as the embedded template, refuses to overwrite an existing baseline, and does not modify project files. Use it only when the project was generated from that exact template version.

### `pad validate [path]`

Discovers `pad.yaml` in the given directory or a parent and validates:

- required project files and directories;
- configuration and state YAML syntax and key types;
- task, requirement, and ADR front matter;
- required metadata, allowed status/priority/type values, and ID patterns;
- record filename/ID consistency for tasks, requirements, and ADRs;
- task-to-requirement and requirement-to-task references;
- task dependencies and dependency cycles.

PAD metadata uses a deliberately restricted YAML subset: nested mappings, scalar lists, quoted/plain strings, booleans, numbers, and null. Anchors, tags, block scalars, complex keys, and YAML object constructors are unsupported and rejected. This keeps the CLI dependency-free and limits the accepted metadata surface; it is not a general-purpose YAML implementation.

### `pad status [path]`

Displays project status, task counts, selected fields from `.agent/state/current.yaml`, and `git status --short --branch`. Git status is read-only. A project without a Git repository is still reportable.

### `pad task`

```bash
pad task list [--status ready] [path]
pad task show TASK-001 [path]
pad task add --title "Implement parser" [--priority medium] [--assignee unassigned] [--estimate "2d"] [path]
pad task set TASK-001 --status in_progress [--priority high] [--assignee agent-a] [path]
```

- IDs are allocated as `TASK-NNN` using the highest existing numeric ID.
- New tasks start in `backlog`.
- Status transitions are checked. Allowed statuses: `backlog`, `ready`, `in_progress`, `review`, `blocked`, `done`, `cancelled`.
- Priority values: `critical`, `high`, `medium`, `low`.
- `set` can change status, priority, and/or assignee; at least one field is required.
- Reopening `done` or `cancelled` tasks requires `--reason`; the reason is appended to the task's status history.
- Task creation includes an editable objective, scope, acceptance criteria, and verification outline.

### `pad requirement` (alias: `pad req`)

```bash
pad requirement list [path]
pad requirement show PAD-REQ-001 [path]
pad requirement add --title "Validate task references" [--type functional] [--priority should] [path]
pad requirement set PAD-REQ-001 --status approved [path]
```

- IDs are allocated as `PAD-REQ-NNN`.
- Types: `functional`, `non_functional`, `constraint`.
- Priorities: `must`, `should`, `could`, `wont`.
- Statuses: `proposed`, `approved`, `implemented`, `verified`, `deprecated`.
- Status transitions are checked.

### `pad decision` (alias: `pad adr`)

```bash
pad decision list [path]
pad decision show ADR-0001 [path]
pad decision add --title "Choose storage engine" [--decider project-owner] [path]
pad decision set ADR-0002 --status accepted [path]
```

Decision records use `ADR-NNNN` IDs and the statuses `proposed`, `accepted`, `rejected`, `superseded`, and `deprecated`. Status transitions are checked. New records include context, decision drivers, options, consequences, and follow-up sections.

### `pad agent`

```bash
pad agent list [path]
pad agent show IMPLEMENTER [path]
```

Lists or displays the role profiles stored in `.agent/agents/`. This command does not launch an agent runtime.

### `pad checkpoint --message "..." [--commit] [--commit-message "..."] [path]`

Writes a timestamped checkpoint under `.agent/state/checkpoints/`, including the current Git status and recovery instructions. Updates `last_checkpoint` and `updated_at` in `.agent/state/current.yaml`, and refreshes a marked latest-checkpoint section in `HANDOFF.md` without replacing the rest of the handoff. By default it does not stage or commit files. With explicit `--commit`, it creates a local commit containing only the generated checkpoint, `current.yaml`, and `HANDOFF.md`; it refuses to run if the Git index already contains staged changes. It never stages the rest of the working tree and never pushes.

### `pad update [path]`

Plans and applies a local three-way update from the CLI's embedded template.

```bash
pad update --dry-run ./my-project
pad update ./my-project
# After manually resolving conflict files:
pad update --accept-resolved ./my-project
```

For each managed file, PAD compares:

- **base** — the exact version saved during `pad init` or the last successful update;
- **local** — the current project file;
- **upstream** — the template embedded in this CLI build.

Behavior:

- local equals base: upstream changes are applied;
- upstream equals base: local-only edits are retained;
- local equals upstream: the file is already current;
- both changed: Git's local `merge-file` performs a line-level three-way merge when possible;
- unresolved conflicts: the local file is preserved; upstream candidate and merge output (with conflict markers) are written under `.pad/update-conflicts/` for review. After editing the project file to resolve the conflict, `--accept-resolved` advances its merge base to the saved upstream candidate;
- files removed from the upstream template are reported for manual review and never deleted automatically;
- files deleted locally are preserved as local deletions.

Before applying changes, the updater saves copies of affected local files and merge bases under `.pad/update-backups/`; these temporary backup and conflict directories are ignored by the generated project's `.gitignore`. A successful update advances the baseline for files updated or cleanly merged. Conflict files retain their old base until manually resolved. The command never contacts a remote service. Run `pad validate` and review `git diff` after an update.

## Exit codes

- `0` — command completed successfully (a conflict report from `update` is a successful safe completion; inspect the report).
- `1` — project, validation, or execution error.
- `2` — invalid command, arguments, or usage.

## Safety boundaries

- No remote Git operations are implemented.
- No Git operation is automatic. `checkpoint --commit` is an explicit opt-in and commits only the generated checkpoint/state/handoff files; no reset, clean, or push is performed.
- No user file is automatically deleted by template updates.
- `--force` on `init` is destructive to colliding template paths; review `--dry-run` first.
- Keep `.pad/template-base/` in version control: it is required for correct future three-way updates.
