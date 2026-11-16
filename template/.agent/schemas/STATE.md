# State snapshot schema (v1)

`state/current.yaml` is a compact navigation snapshot, not a full event log.

| Key | Type | Meaning |
|---|---|---|
| `schema_version` | string | State schema version |
| `project_status` | enum | `initializing`, `active`, `paused`, `maintenance`, `archived` |
| `health` | enum | `healthy`, `at_risk`, `blocked`, `unknown` |
| `updated_at` | timestamp or null | ISO 8601 with timezone |
| `current_focus` | string | Current project-level focus |
| `active_task` | task ID or null | Primary task being executed |
| `active_branch` | string or null | Current Git branch |
| `last_checkpoint` | string or null | Local checkpoint identifier |
| `last_verified_commit` | string or null | Commit with recorded verification |
| `open_blockers` | list | Current blockers |
| `next_actions` | list[string] | Ordered concrete actions |
| `recent_changes` | list | Short verified changes |
| `risks` | list | Known risks |

Keep detailed history in Git and task progress logs. Do not put credentials or sensitive personal information in this file.
