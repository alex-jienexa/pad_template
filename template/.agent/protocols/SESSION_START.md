# Session start protocol

Follow this sequence at the beginning of every substantial session.

1. **Read entry points:** `AGENTS.md`, `pad.yaml`, and `.agent/README.md`.
2. **Inspect repository reality:** run `git status --short --branch`, `git log -5 --oneline`, and inspect relevant diffs. If there is no Git repository or no commit yet, state that explicitly.
3. **Reconcile state:** read `state/current.yaml` and `state/HANDOFF.md`; treat them as guidance, not proof. Correct stale statements after checking repository evidence.
4. **Load task context:** read the assigned task, linked requirements, relevant ADRs, and only the protocols/agent profiles needed for the work.
5. **Confirm execution contract:** identify objective, included/excluded scope, files owned, dependencies, acceptance criteria, and required verification.
6. **Check authorization:** identify operations requiring human approval. Ask before crossing the configured autonomy boundary.
7. **Begin with a small inspection or plan** when uncertainty or impact is high; otherwise proceed within the task's approved scope.

Do not infer that a remote backup exists from a local commit. Do not infer that a previous agent completed work merely because a task says `done`; verify the relevant artifacts and checks.
