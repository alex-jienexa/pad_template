# Context recovery protocol

Use when chat history, agent memory, or the previous handoff is missing, truncated, or contradictory.

## Recovery order

1. Preserve the current working tree. Do not reset, clean, stash, or overwrite files as a recovery shortcut.
2. Inspect `git status --short --branch`, `git log -10 --oneline --decorate`, current branch, and relevant diffs.
3. Inspect recent commits and file timestamps only as supporting evidence; timestamps are not authoritative history.
4. Read `pad.yaml`, `AGENTS.md`, `state/current.yaml`, and `state/HANDOFF.md`.
5. Inspect active tasks and their linked requirements/ADRs. Search the repository for TODO/FIXME markers only when relevant to the task.
6. Compare documented claims with code, tests, and Git. Mark discrepancies as stale or uncertain; do not silently invent missing context.
7. Reconstruct a concise working summary: confirmed facts, current changes, completed checks, blockers, and next action.
8. Update state/handoff only after the reconstruction is supported by evidence.

## If recovery is incomplete

- Stop before risky or irreversible operations.
- Ask the project owner one focused question for the missing decision or fact.
- Record the unresolved uncertainty and its impact.
