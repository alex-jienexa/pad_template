# PAD Template

**Persistent Agent Development** — an open, agent-agnostic method and starter structure for projects developed with AI agents over multiple sessions.

> **Status:** PAD v1.0 methodology, starter template, and local CLI implementation are available for review. The CLI supports initialization, metadata validation, task/requirement/ADR management, agent-profile inspection, checkpoints, status inspection, and local three-way template updates. This is a local working repository, not a published release.

## What PAD provides

- durable project context stored in version control;
- requirements and tasks represented as Markdown with YAML metadata;
- a single authoritative handoff and machine-readable project state;
- recovery procedures for sessions where the agent has lost context;
- explicit Git safety rules and checkpoint practices;
- protocols for sequential and parallel agents;
- controlled template updates based on three-way merge.

PAD is a methodology and toolkit, not an AI model, agent runtime, or guarantee of correctness. The human project owner remains responsible for decisions, access, and review.

## Repository map

- `template/` — canonical starter files copied into user projects.
- `cmd/pad/`, `internal/` — Go CLI implementation.
- `docs/SPECIFICATION.md` — approved PAD v1.0 baseline.
- `docs/methodology/` — operating model, state, task design, Git, security, and agent assignment guidance.
- `docs/ARCHITECTURE.md` — component boundaries and data ownership.
- `docs/CLI.md` — implemented CLI commands and current limitations.
- `.github/` — CI and contribution workflows.

## Quick start

The canonical starter structure is in `template/`. Build the CLI with `go build -o ./bin/pad ./cmd/pad`. Preview initialization with `pad init --dry-run ./my-project`, then create a configured project with `pad init --name "My Project" --description "Project purpose" ./my-project`. The CLI validates the supported PAD YAML subset, manages tasks/requirements/decisions, creates checkpoints, and compares projects against the template embedded in the CLI. It never fetches remote content.

## Development

Requires Go 1.23+ and Git.

```bash
go test ./...
go vet ./...
go build ./cmd/pad
./scripts/check.sh
```

## License

MIT. See [LICENSE](LICENSE). Projects created from the starter should choose and include their own license deliberately; the PAD tool's license does not automatically determine the generated project's license.
