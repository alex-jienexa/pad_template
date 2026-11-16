# Task execution protocol

1. Read the complete task and all linked requirements and decisions.
2. Confirm the task is `ready` or obtain owner approval to start it from another state.
3. Inspect the repository and identify relevant code, tests, and conventions.
4. Restate the objective and acceptance criteria internally; ask if they conflict or are not testable.
5. Work only within the approved scope and declared file ownership. Keep unrelated user changes untouched.
6. Implement the smallest coherent solution. Add or update tests alongside behavior.
7. Run focused checks, then broader checks when practical. Record exact commands and outcomes; distinguish not-run from failed.
8. Review the diff for accidental changes, secrets, generated artifacts, and scope creep.
9. Update task status and progress log. A task enters `review` when implementation is ready for independent review; it enters `done` only after its completion contract is satisfied.
10. Synchronize `state/current.yaml` and `state/HANDOFF.md` when the next action or project state materially changes.
