# Requirements

Requirements describe stable product behavior or quality constraints. Keep them implementation-independent where practical.

## Naming and lifecycle

- File name: `PAD-REQ-NNN-short-name.md`.
- IDs are stable and never reused.
- Lifecycle: `proposed`, `approved`, `implemented`, `verified`, `deprecated`.
- Each requirement has observable acceptance criteria and a verification method.
- Link tasks that implement or verify the requirement; do not mark a requirement verified solely because code exists.

Copy `REQUIREMENT_TEMPLATE.md` for each new requirement. Requirement IDs use the `PAD-REQ-NNN` form; priorities are `must`, `should`, `could`, or `wont`. The template is intentionally not a requirement record.
