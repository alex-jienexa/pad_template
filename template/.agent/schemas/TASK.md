# Task metadata schema (v1)

The canonical machine-readable contract is `task.schema.yaml`. Required YAML front matter keys are `id`, `title`, `status`, `priority`, `created_at`, and `updated_at`.

- `id`: stable identifier using `TASK-NNN`.
- `status`: `backlog`, `ready`, `in_progress`, `review`, `blocked`, `done`, or `cancelled`.
- `priority`: `critical`, `high`, `medium`, or `low`.
- `assignee`: person, agent role, or `unassigned`.
- `depends_on`: list of task IDs.
- `requirement_ids`: list of requirement IDs.
- `labels`: list of strings.
- `estimate`: string or null.
- `due_date`: date or null.
- `files_owned`: optional list of paths defining write ownership.

Dates use `YYYY-MM-DD`. The Markdown body should include an objective, context, scope, acceptance criteria, verification, and progress log. The CLI validator must not silently mutate records.
