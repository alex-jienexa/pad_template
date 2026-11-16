# Checkpoint protocol

A checkpoint is a recoverable point in project history, not merely a note that work happened.

## Before checkpointing

- Inspect `git status --short` and the complete diff.
- Run relevant tests and checks; record any failures or omissions.
- Update task metadata, `state/current.yaml`, and the single `state/HANDOFF.md`.
- Exclude secrets, credentials, `.env`, private keys, personal data, and unrelated files.
- Confirm the proposed commit contains only the intended logical unit.

## Commit policy

- Prefer a small number of meaningful commits over a commit for every trivial edit.
- Follow the repository's commit message convention, if defined.
- Creating a local commit is allowed only when `agent_policy.allow_local_commits` permits it and the task's authority allows it.
- Never amend, rebase, reset, clean, force-push, or rewrite shared history without explicit approval for that exact operation.

## Backup semantics

A local commit is not an off-device backup. A remote backup exists only after a push has completed and the resulting remote state has been verified. If network access or credentials are unavailable, report that the checkpoint is local only.
