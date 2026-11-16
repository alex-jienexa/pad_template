---
id: TASK-001
title: Establish project context and development workflow
status: ready
priority: high
assignee: project-owner
created_at: "{{DATE_YYYY_MM_DD}}"
updated_at: "{{DATE_YYYY_MM_DD}}"
estimate: null
parent: null
depends_on: []
requirement_ids: []
files_owned:
  - .agent/context/PROJECT.md
  - .agent/state/current.yaml
  - .agent/state/HANDOFF.md
labels:
  - bootstrap
---

# Establish project context and development workflow

## Objective

Replace the starter context with verified project facts and make the project reproducible for a fresh contributor or agent session.

## Scope

- Complete `.agent/context/PROJECT.md`.
- Verify setup, build, and test commands.
- Record the initial project state and next task.

## Acceptance criteria

- [ ] Project purpose, scope, constraints, and domain language are documented.
- [ ] Setup/build/test commands have been executed or explicitly marked unavailable.
- [ ] Open questions and assumptions have named owners.
- [ ] `current.yaml` and `HANDOFF.md` agree with the repository state.

## Verification

A fresh session can follow `AGENTS.md` and identify the next concrete action without relying on chat history.

## Progress log

- YYYY-MM-DD — Starter task created by the PAD template.
