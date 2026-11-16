# Changelog

All notable changes to PAD are documented here. This project is under active development; entries describe local development state until a release is published.

## Unreleased — PAD v1.0 implementation

- Added the approved PAD v1.0 methodology, architecture, state model, agent protocols, and project template.
- Added MIT licensing, contribution/security policies, GitHub issue/PR templates, and local quality checks.
- Implemented `pad init` with placeholder substitution, dry-run planning, collision protection, and a pristine local merge baseline; added `pad adopt` for matching GitHub Template Repository projects.
- Implemented typed configuration parsing and validation for PAD's restricted, data-only YAML subset.
- Implemented `pad validate` for configuration/state, task/requirement/ADR metadata, identifiers, dates, references, and dependency cycles.
- Implemented `pad status`, task/requirement/decision CRUD workflows, and agent-profile inspection.
- Implemented timestamped checkpoints with synchronized state and handoff; optional local commits are explicit and limited to checkpoint metadata files.
- Implemented local three-way template updates with line-level merge, safety backups, conflict artifacts, manual deletion review, and explicit merge-base acceptance after conflict resolution.
- Added unit, integration, workflow, conflict-resolution, and Git-safety tests.
- Added a GitHub Actions build/test matrix for Linux, Windows, and macOS. Hosted CI has not been run from this local environment.

## Known v1.0 boundaries

- No remote template fetching, GitHub API access, repository creation, remote configuration, push, or release publishing.
- YAML support is intentionally limited; it is not a general-purpose YAML 1.2 implementation.
- Agent commands inspect role profiles; they do not launch or orchestrate AI runtimes.
- Template deletion and semantic migrations remain manual-review operations.
