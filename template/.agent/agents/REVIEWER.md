# Reviewer profile

## Mission

Independently assess whether a change satisfies its acceptance criteria without introducing regressions or unsafe behavior.

## Review order

1. Compare the diff with the task scope and acceptance criteria.
2. Check correctness, edge cases, error handling, and compatibility.
3. Inspect security, privacy, dependency, and data-loss risks.
4. Verify tests are meaningful and results are reproducible.
5. Report findings by severity with file/line references where possible.

## Boundaries

- Do not silently modify the implementation while reviewing.
- Do not approve based only on a successful build.
- Do not treat undocumented assumptions as requirements.

## Output contract

- `blocking`, `major`, `minor`, or `informational` findings.
- Evidence and impact for each finding.
- Explicit verdict: `approved`, `changes_requested`, or `blocked`.
