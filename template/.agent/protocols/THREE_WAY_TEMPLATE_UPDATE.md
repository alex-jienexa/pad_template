# Three-way template update protocol

PAD template updates preserve project customization by comparing three inputs:

- **Base (B):** the exact upstream template version last applied to this project;
- **Upstream (U):** the new template version;
- **Local (L):** the current project files, including user edits.

## Planning

1. Verify the recorded base version and base snapshot exist and are readable.
2. Verify the target is the intended project root and inspect its working tree.
3. Build a file-level plan for additions, modifications, renames, and upstream deletions.
4. For each path, compare B, U, and L. If only U changed, apply the upstream change; if only L changed, preserve L; if both changed identically, accept the common result; if both changed differently, report a conflict.
5. Treat an upstream deletion as a proposed deletion only. Never delete a local file automatically, even if it is unchanged from B, without explicit deletion policy and review.
6. Detect path collisions, type changes, symlinks, and files outside the managed template set.

## Apply

- Produce a preview before modifying files.
- Create a recoverable checkpoint/backup of affected paths.
- Apply non-conflicting changes only; preserve conflicted local files.
- Never use a generic overwrite or recursive delete as a merge strategy.
- Update the stored base snapshot/version only after the apply operation is verified and conflicts are resolved or explicitly deferred.
- Report every conflict and skipped deletion.

## Failure handling

If an operation fails partway through, stop, preserve the partial state, report the exact affected paths, and provide a recovery plan. Do not silently roll back user changes.
