# ADR metadata schema (v1)

Required YAML front matter keys: `id`, `title`, `status`, `date`, and `deciders`.

- `id`: stable identifier using `ADR-NNNN`.
- `status`: `proposed`, `accepted`, `rejected`, `superseded`, or `deprecated`.
- `date`: `YYYY-MM-DD`.
- `deciders`: list of accountable decision-makers.
- Optional links: `supersedes`, `superseded_by`, `related_requirements`, `related_tasks`.

The Markdown body must include context, decision drivers, considered options, decision, consequences, and follow-up actions.
