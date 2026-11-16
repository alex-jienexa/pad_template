# PAD metadata schemas

The YAML files in this directory define the v1 metadata contracts for tasks and requirements. `TASK.md`, `REQUIREMENT.md`, `ADR.md`, and `STATE.md` explain the human-readable conventions.

The current CLI performs structural checks only. Full YAML parsing, enum/type validation, ID/reference validation, and cross-file consistency checks are planned; do not assume they are already enforced. Validation must never silently rewrite user-authored records.

## Canonical values

- Task status: `backlog`, `ready`, `in_progress`, `review`, `blocked`, `done`, `cancelled`.
- Requirement status: `proposed`, `approved`, `implemented`, `verified`, `deprecated`.
- ADR status: `proposed`, `accepted`, `rejected`, `superseded`, `deprecated`.
- Task priority: `critical`, `high`, `medium`, `low`.
- Requirement priority: `must`, `should`, `could`, `wont`.
- Agent autonomy: `read_only`, `controlled`, `delegated`.
- Project status: `initializing`, `active`, `paused`, `maintenance`, `archived`.
- Health: `healthy`, `at_risk`, `blocked`, `unknown`.

## Identifier conventions

- Tasks: `TASK-NNN`, matching `task.schema.yaml`.
- Requirements: `PAD-REQ-NNN`, matching `requirement.schema.yaml`.
- ADRs: `ADR-NNNN`.

Timestamps use ISO 8601; date-only metadata uses `YYYY-MM-DD`. Lists should be empty lists rather than omitted when a schema defines a list. IDs are stable and must not be reused.
