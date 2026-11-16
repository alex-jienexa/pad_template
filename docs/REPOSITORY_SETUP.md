# Repository setup checklist

This document is a local preparation checklist. It does not imply that a remote repository has been created or configured.

## Before publication

- [ ] Create the GitHub repository under the intended owner.
- [ ] Confirm repository name and visibility.
- [ ] Push the reviewed local history and verify the remote default branch.
- [ ] Enable GitHub Template Repository mode for the canonical starter repository.
- [ ] Configure branch protection/rulesets for `main` (pull request review and required CI checks appropriate to the project).
- [ ] Confirm CI checks on an actual GitHub Actions run.
- [ ] Add release notes and version tags only after local validation.
- [ ] Confirm license, copyright holder, security contact, and contribution policy.

## Local repository preparation

```bash
git init -b main
git add .
git diff --cached --check
# Review the staged diff before committing.
git status --short
```

Set `user.name` and `user.email` explicitly in your own Git environment before creating the first commit. Do not use an invented identity. Add a remote and push only when you are ready to publish.

## Branch protection suggestions

- Require pull requests for changes to `main`.
- Require CI checks after the workflow has been verified on GitHub.
- Require conversation resolution.
- Restrict force pushes and branch deletion.
- Keep administrator bypass intentional and limited.
