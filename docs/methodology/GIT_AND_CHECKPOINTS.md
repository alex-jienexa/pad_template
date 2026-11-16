# Git and checkpoint policy

Git is the authoritative technical history for committed source changes. PAD documentation supplies intent, constraints, and navigation; it is not a substitute for inspecting the repository.

## Safe default

Start with read-only inspection. Make scoped edits. Review the diff. Run checks. Synchronize task/state/handoff. Create a logical commit only when policy permits. Never assume the working tree is clean.

## Protected operations

The following require explicit, current human approval: `git reset --hard`, `git clean`, force push, rewriting shared history, deleting branches or tags that may contain work, destructive bulk file operations, and overwriting user changes. Approval for one operation does not authorize future repetitions.

## Logical commits

A commit should represent one coherent, reviewable change. Do not commit every trivial edit, and do not combine unrelated changes merely to reduce commit count. Inspect staged content before committing. Exclude credentials, `.env`, private keys, tokens, and personal data.

## Checkpoint versus backup

A local commit is a checkpoint, not a remote backup. A remote backup is confirmed only after the push succeeds and the remote ref is verified. If remote access is unavailable or the user has prohibited network access, report the checkpoint as local only.

## Near context or execution limits

Stop starting new scope. Update task progress, current state, and the single handoff; record uncommitted changes and exact verification status. If authorized and safe, create a logical local checkpoint. Do not push unless remote access is permitted and explicitly within scope.
