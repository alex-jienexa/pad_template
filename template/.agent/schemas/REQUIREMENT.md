# Requirement metadata schema (v1)

The canonical machine-readable contract is `requirement.schema.yaml`. Required YAML front matter keys are `id`, `title`, `type`, `status`, `created_at`, and `updated_at`.

- `id`: stable identifier using `PAD-REQ-NNN`.
- `type`: `functional`, `non_functional`, or `constraint`.
- `status`: `proposed`, `approved`, `implemented`, `verified`, or `deprecated`.
- `priority`: `must`, `should`, `could`, or `wont`.
- `verification`: `test`, `inspection`, `analysis`, `demonstration`, or `manual`.
- Optional ownership and traceability: `owner`, `related_requirements`, `related_tasks`.

Dates use `YYYY-MM-DD`. The Markdown body must include a statement, rationale, acceptance criteria, and verification method.
