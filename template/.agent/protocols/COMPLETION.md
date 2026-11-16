# Task completion protocol

A task may move to `done` only when:

- every mandatory acceptance criterion is met or has a documented, owner-approved exception;
- relevant automated and manual checks are recorded with exact outcomes;
- no known blocker remains hidden;
- documentation and state reflect the actual result;
- review is complete where required.

The completion report includes task ID, summary, changed files, tests/checks, commit reference if any, unresolved risks, and next action. If any mandatory criterion is unmet, use `review` or `blocked`, not `done`.
