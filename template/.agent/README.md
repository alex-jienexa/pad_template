# Project memory index

This directory is the durable context store for AI-assisted work. Keep it concise, factual, and current. Prefer links to source files over duplicated implementation details.

## Read order

1. `../AGENTS.md` — mandatory agent rules.
2. `../pad.yaml` — project configuration and autonomy policy.
3. `state/current.yaml` — current machine-readable state.
4. `state/HANDOFF.md` — next actions and session continuation.
5. `context/PROJECT.md` — purpose, scope, domain terms, and constraints.
6. The assigned task in `tasks/`, its requirements in `requirements/`, and relevant ADRs in `decisions/`.
7. Relevant protocols in `protocols/` and role definitions in `agents/`.
8. `schemas/` — metadata contracts and allowed values.

## Directory guide

- `context/`: stable project facts and architecture summaries.
- `requirements/`: durable product and quality requirements.
- `tasks/`: executable work items.
- `decisions/`: append-only architecture/product decision records.
- `state/`: current snapshot and single handoff document.
- `agents/`: role boundaries for specialized or parallel agents.
- `protocols/`: repeatable operating procedures.
- `schemas/`: metadata contracts and enums.

## Maintenance rules

- Facts should be traceable to code, a requirement, a decision, or a human statement.
- Distinguish confirmed facts from assumptions and open questions.
- Do not store secrets or sensitive personal information.
- Update rather than append repetitive status prose; preserve historical decisions in ADRs and Git.
- Keep exactly one canonical `state/HANDOFF.md`.
- If a note conflicts with the repository, inspect the code and Git history, then correct the note. Never blindly trust stale memory.
