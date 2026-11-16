# PAD v1.0 implementation roadmap

## Foundation and methodology

- [x] Repository skeleton, Go module, CI workflow, contribution/security policies.
- [x] MIT license selected and added.
- [x] PAD v1.0 methodology specification and project template baseline.
- [x] Core agent protocols: session start, recovery, task execution, checkpoint, completion/handoff, multi-agent coordination, template update.
- [x] Requirement, task, ADR, state metadata contracts and starter records.
- [x] Local template quality gate and metadata examples.

## CLI

- [x] `pad init`: embedded template initialization, safe preview, collision handling, placeholder substitution, and pristine baseline capture.
- [x] `pad adopt`: safe baseline adoption for projects generated directly from the matching GitHub Template Repository.
- [x] `pad validate`: supported YAML parsing, typed configuration validation, front matter, IDs, enums, dates, references, and task dependency-cycle checks.
- [x] `pad status`: project/task summary and read-only Git status.
- [x] `pad task`: list/show/create/update with checked status transitions.
- [x] `pad requirement`: list/show/create/update with checked status transitions.
- [x] `pad decision`: list/show/create/update ADR records with checked status transitions.
- [x] `pad agent`: inspect installed role profiles.
- [x] `pad checkpoint`: timestamped checkpoint, synchronized state and handoff; optional scoped local commit with explicit opt-in.
- [x] `pad update`: local line-level three-way merge, conflict artifacts, manual deletion review, and explicit conflict-baseline acceptance.
- [x] CLI unit, integration, and workflow tests.

## Release preparation (owner-operated)

- [ ] Run the supported-platform build/test matrix on target environments.
- [ ] Review the YAML subset limitation and decide whether to adopt a full YAML library in a later release.
- [ ] Prepare migration notes and release changelog.
- [ ] Verify license, repository metadata, and GitHub Template Repository settings.
- [ ] Create/publish the GitHub repository and release artifacts manually.

The CLI is local-only. It does not fetch a remote template, create a GitHub repository, configure remotes, push, or publish releases.
