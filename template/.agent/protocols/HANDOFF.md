# Handoff writing guide

Keep `.agent/state/HANDOFF.md` as the only current handoff. It should let a fresh agent resume in a few minutes.

Include:
- verified branch, commit, and working-tree status;
- completed work and evidence;
- active task and exact stopping point;
- decisions made and relevant ADRs;
- blockers, assumptions, and unresolved questions;
- exact next action and useful commands;
- checks run and checks still needed;
- whether work is committed and whether it is pushed remotely.

Avoid long narrative, duplicated task descriptions, speculative status, and sensitive data. Update the existing handoff instead of creating dated copies; Git history preserves prior versions.
