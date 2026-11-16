# PAD Template contributor entry point

This repository contains the PAD methodology, canonical starter files under `template/`, and the Go CLI under active development. This file governs work on PAD itself; projects generated from the starter use `template/AGENTS.md` instead.

## Before acting

1. Read `README.md`, `docs/SPECIFICATION.md`, `docs/DECISIONS.md`, and the relevant roadmap section.
2. Inspect Git state (`git status --short --branch`, recent commits, and relevant diffs) before editing.
3. Read the relevant implementation task, architecture notes, and affected template protocols.
4. Preserve approved PAD v1.0 decisions. Any architectural change requires an ADR and explicit project-owner approval.

## Non-negotiable safety rules

- Work locally unless the project owner explicitly authorizes external access. Do not publish, push, create remote repositories, or contact external services as part of ordinary implementation.
- Never run destructive Git commands (`reset --hard`, `clean`, force push, history rewrite) without explicit, current human approval.
- Never overwrite, discard, or silently reformat unrelated user changes.
- Never commit secrets, credentials, `.env` files, private keys, or personal data.
- Do not claim a test, build, deployment, push, or external action succeeded unless it was actually executed and verified.

## Work discipline

- Keep implementation and template content consistent; update the template checks when adding required starter files.
- Use the task's acceptance criteria as the completion contract.
- Run `./scripts/check.sh` after meaningful changes when practical.
- Keep docs concise, actionable, and explicit about what is implemented versus planned.
- Do not create commits using an invented author identity.

## Completion report

Report changed areas, checks actually run and their outcomes, limitations, and the next concrete implementation step. Distinguish local-only work from anything remotely published.
