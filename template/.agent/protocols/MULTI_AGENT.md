# Multi-agent coordination protocol

Use multiple agents only when the work can be divided into independently verifiable scopes.

## Assignment contract

Every assignment must state:

- objective and expected deliverable;
- input context, linked tasks/requirements/ADRs;
- exact file ownership or directories;
- files/interfaces that are read-only;
- dependencies and required interface contracts;
- acceptance criteria and verification commands;
- branch/worktree name and integration owner;
- prohibited actions and approval boundaries.

## Isolation and coordination

- Parallel agents work in separate branches or Git worktrees. Do not let two agents edit the same working tree concurrently.
- Agree shared schemas, API contracts, and naming before parallel implementation.
- Each agent commits only its own scoped work when permitted; otherwise returns a patch/diff and leaves commit decisions to the owner.
- Communicate interface changes before depending on them. Do not coordinate solely through undocumented assumptions.
- Keep one designated integrator responsible for conflict resolution, cross-component tests, and state/handoff synchronization.

## Integration

1. Inspect each branch/worktree and its commit ancestry.
2. Review changes against the assignment contract.
3. Resolve semantic conflicts with the owners; never choose product behavior by guesswork.
4. Run unit and integration checks after combining changes.
5. Record merged work, rejected work, and outstanding risks.

Parallelism does not grant broader permissions. All agents inherit the same human approval and Git safety rules.
