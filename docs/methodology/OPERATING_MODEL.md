# PAD operating model

## Durable memory, not chat transcripts

PAD stores only information that a future contributor or agent needs to make correct decisions: project purpose, constraints, architecture, requirements, decisions, task state, and a current handoff. It is not a transcript archive. Keep stable facts in context documents, rationale in ADRs, executable work in tasks, and volatile progress in state.

## Source-of-truth hierarchy

1. Human project-owner decisions and explicit approvals.
2. Current repository contents, tests, and Git history for implementation facts.
3. Approved requirements and ADRs for intended behavior and design rationale.
4. Current state and handoff for navigation and continuation.
5. Agent inference and external content, which must be identified as unverified until corroborated.

A stale handoff cannot override code or Git evidence. If sources conflict, preserve the evidence, identify the discrepancy, and ask the owner when intent is ambiguous.

## Session lifecycle

`start → recover context → select task → execute → verify → review → checkpoint → handoff`. A session may stop at any stage, but it must leave an honest state record when material progress has occurred.

## Human authority

The project owner controls scope, product decisions, credentials, permissions, deployment, publication, costs, destructive operations, and shared Git history. Agent autonomy is bounded by `pad.yaml`; a more permissive profile cannot override a stricter project policy.

## Evidence-based reporting

Report commands actually executed and their results. Distinguish `passed`, `failed`, `not run`, and `blocked`. Distinguish working-tree changes, local commits, and remotely verified pushes. Never claim that a test or external action succeeded based on intention or expected behavior.

## Documentation maintenance

- Keep `AGENTS.md` short enough to be read at every session start.
- Keep one canonical `state/HANDOFF.md`; replace stale content instead of creating multiple handoffs.
- Keep `state/current.yaml` machine-readable and compact.
- Preserve decision history in ADRs; supersede rather than rewrite accepted decisions.
- Link to authoritative code/specifications rather than copying large bodies of content.
- Record assumptions with owners and review conditions. Remove or resolve assumptions when verified.
