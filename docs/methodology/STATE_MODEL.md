# PAD state model

PAD deliberately separates durable history from volatile current state.

| Record | Purpose | Update frequency |
|---|---|---|
| Git commits | Authoritative history of committed implementation changes | At logical checkpoints |
| Requirements | Stable obligations and acceptance criteria | When scope/behavior changes |
| ADRs | Historical rationale for material decisions | When a decision is made or superseded |
| Task records | Work ownership, progress, dependencies, completion | During task lifecycle |
| `current.yaml` | Compact current snapshot for tools | When material state changes |
| `HANDOFF.md` | Human-readable continuation instructions | At meaningful handoffs |

## Consistency rules

- The task's status and the active task in `current.yaml` must not contradict each other.
- `HANDOFF.md` must identify the next concrete action, not merely summarize the past.
- `last_verified_commit` may contain only a commit whose relevant checks were actually run and recorded.
- `last_checkpoint` is not evidence of a remote backup.
- A blocker must identify its impact and the person or role that can resolve it.
- Do not duplicate the same status narrative in several files; keep one canonical location for each kind of information.

## Suggested update transaction

1. Finish or pause a meaningful unit of work.
2. Update the task's status and progress log.
3. Update `current.yaml` from the verified repository state.
4. Replace `HANDOFF.md` with the next actionable continuation note.
5. Review the combined diff for contradictions and accidental sensitive data.
6. Create a logical checkpoint if authorized.
