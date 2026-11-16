# Tasks

Tasks are bounded execution units. A task must be small enough to review and have a clear completion contract.

## Naming and lifecycle

- File name: `TASK-NNN-short-name.md`; IDs use the `TASK-NNN` form.
- IDs are stable and never reused.
- Lifecycle: `backlog`, `ready`, `in_progress`, `blocked`, `review`, `done`, `cancelled`.
- `blocked` can be entered from any non-terminal state; record the blocker and owner.
- A task is `done` only when acceptance criteria and required checks pass, or an explicitly approved exception is documented.
- Reopening a completed task requires a reason in the task history.

Copy `TASK_TEMPLATE.md` for new work. Use `created_at`, `updated_at`, `depends_on`, and `requirement_ids` as defined in `schemas/task.schema.yaml`. `TASK-001.md` is an example of a first project bootstrap task; adapt or replace it.
