#!/usr/bin/env sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
template="$root/template"

required_files='
AGENTS.md
README.md
pad.yaml
.pad-template.yaml
.agent/README.md
.agent/context/PROJECT.md
.agent/requirements/README.md
.agent/requirements/REQUIREMENT_TEMPLATE.md
.agent/tasks/README.md
.agent/tasks/TASK_TEMPLATE.md
.agent/tasks/TASK-001.md
.agent/decisions/README.md
.agent/decisions/ADR-0001-template.md
.agent/decisions/ADR-0001-pad-baseline.md
.agent/state/current.yaml
.agent/state/HANDOFF.md
.agent/agents/IMPLEMENTER.md
.agent/agents/REVIEWER.md
.agent/agents/INTEGRATOR.md
.agent/protocols/SESSION_START.md
.agent/protocols/RECOVERY.md
.agent/protocols/TASK_EXECUTION.md
.agent/protocols/CHECKPOINT.md
.agent/protocols/COMPLETION.md
.agent/protocols/HANDOFF.md
.agent/protocols/MULTI_AGENT.md
.agent/protocols/THREE_WAY_TEMPLATE_UPDATE.md
.agent/schemas/README.md
.agent/schemas/TASK.md
.agent/schemas/REQUIREMENT.md
.agent/schemas/ADR.md
.agent/schemas/STATE.md
'

printf '%s\n' "$required_files" | while IFS= read -r file; do
  [ -z "$file" ] && continue
  if [ ! -f "$template/$file" ]; then
    echo "Missing template file: $file" >&2
    exit 1
  fi
done

python3 - "$template" <<'PY'
from pathlib import Path
import re, sys
root = Path(sys.argv[1])
checks = {
    ".agent/tasks/TASK-001.md": ["id:", "title:", "status:", "priority:", "assignee:", "created_at:", "updated_at:", "depends_on:", "requirement_ids:"],
    ".agent/decisions/ADR-0001-pad-baseline.md": ["id:", "title:", "status:", "date:", "deciders:"],
    ".agent/requirements/REQUIREMENT_TEMPLATE.md": ["id:", "title:", "type:", "status:", "priority:", "created_at:", "updated_at:"],
}
for name, keys in checks.items():
    content = (root / name).read_text(encoding="utf-8")
    if not content.startswith("---\n"):
        raise SystemExit(f"{name}: missing YAML front matter opener")
    parts = content.split("---", 2)
    if len(parts) < 3:
        raise SystemExit(f"{name}: missing YAML front matter closer")
    front = parts[1]
    for key in keys:
        if not re.search(r"(?m)^" + re.escape(key) + r"\s*", front):
            raise SystemExit(f"{name}: missing front matter key {key}")
print("Template structure and front matter checks passed")
PY
