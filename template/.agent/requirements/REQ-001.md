---
id: PAD-REQ-001
title: "Preserve durable project context"
type: functional
status: approved
priority: must
created_at: "2026-11-16"
updated_at: "2026-11-16"
verification: inspection
---

# Statement

The project shall store its essential working context, current state, and handoff instructions in version-controlled files so a new agent session can resume without access to prior conversation history.

## Rationale

Conversation context may be unavailable, truncated, or tied to a particular tool. Repository records provide a reviewable and portable continuity layer.

## Acceptance criteria

- [ ] The agent entry point identifies the required read order.
- [ ] Current state and handoff have one canonical location each.
- [ ] Recovery instructions explain how to reconcile notes with Git and source files.
- [ ] No secret or sensitive personal data is required to use the mechanism.
