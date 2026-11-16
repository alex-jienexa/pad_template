# PAD Architecture

## Components

1. **Methodology:** operating principles and protocols under `.agent/protocols/`.
2. **Project memory:** requirements, decisions, context, task records, state, and handoff under `.agent/`.
3. **Template:** canonical starter project under `template/`.
4. **CLI:** Go executable under `cmd/pad` with private packages in `internal/`.
5. **Git/GitHub:** local version history and optional remote collaboration/backup.

## Data ownership

| Data | Canonical location | Notes |
|---|---|---|
| Project configuration | `pad.yaml` | Human-edited YAML |
| Agent entry point | `AGENTS.md` | Short and stable |
| Requirements | `.agent/requirements/*.md` | YAML front matter + Markdown |
| Tasks | `.agent/tasks/*.md` | YAML front matter + Markdown |
| Decisions | `.agent/decisions/*.md` | Append-only decision history |
| Current state | `.agent/state/current.yaml` | Machine-readable snapshot |
| Handoff | `.agent/state/HANDOFF.md` | One current handoff |
| Agent roles | `.agent/agents/*.md` | Responsibilities and boundaries |
| Protocols | `.agent/protocols/*.md` | Operational procedures |
| Upstream template base | `.pad/template-base/` | Baseline for future three-way merge; managed by CLI |

## CLI package boundaries

- `cmd/pad`: argument parsing and process entry.
- `internal/cli`: command routing and user-facing output.
- `internal/config`: typed PAD configuration parsing and validation.
- `internal/yamlmini`: the restricted, data-only YAML subset accepted by PAD metadata.
- `internal/records`: Markdown front matter parsing and deterministic record writing.
- `internal/project`: project discovery and safe relative-path checks.
- `internal/template`: embedded template, initialization, and merge-base snapshots.
- `internal/cli`: task, requirement, ADR, profile, checkpoint, status, validation, and update commands. Git is invoked locally through fixed argument lists; no remote operation is implemented.

## Safety invariants

- Never overwrite a non-empty target during initialization without explicit `--force` and a preview.
- Never delete user files during template updates by default.
- Always calculate update plans before applying changes.
- Preserve a recoverable base version for three-way merge.
- Do not run shell commands by interpolating untrusted user-controlled strings.
- Report external actions only after verifying their result.

## Template distribution

The same canonical `template/` tree is intended to populate the GitHub Template Repository. The template repository and CLI source may be published separately after repository boundaries and release automation are finalized. GitHub template-generated repositories begin with unrelated history, so template updates are handled by PAD's base-aware merge rather than GitHub pull requests against the template repository.
