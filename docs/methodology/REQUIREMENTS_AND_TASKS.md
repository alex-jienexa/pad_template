# Requirements and task design

## Requirements

A requirement describes a durable capability, constraint, or quality attribute. It should answer what must be true and why, without prematurely prescribing implementation. Use one stable ID per requirement. Acceptance criteria must be observable and testable; verification states how evidence will be collected. Requirements move through `proposed → approved → implemented → verified`; `deprecated` records a deliberate retirement.

Avoid vague words such as fast, intuitive, robust, or secure unless they are quantified or tied to a verification method. Separate independent obligations when they can be accepted independently.

## Tasks

A task is an executable unit of work, not a second copy of a requirement. It should have one objective, bounded scope, dependencies, file ownership when useful, acceptance criteria, verification steps, and risks. Tasks may implement multiple requirements, and a requirement may require multiple tasks. Keep these links explicit.

Task lifecycle: `backlog → ready → in_progress → review → done`. `blocked` can be entered from any non-terminal state; `cancelled` records intentionally abandoned work. A task is not complete because an agent says it is complete: required checks and review conditions must be satisfied or explicitly waived by the owner.

## Sizing and decomposition

Split work when separate parts have independent acceptance, different owners, distinct risk, or can be reviewed separately. Do not split merely to create many tiny tasks. Prefer logical deliverables that can be checkpointed and handed off safely.

## Traceability

Use task metadata to link requirements and dependencies. Use ADR metadata to link decisions. Keep IDs stable and never recycle them. If a requirement changes after implementation begins, update the requirement first and review the impact on linked tasks.
