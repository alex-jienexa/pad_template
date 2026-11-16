# PAD project layout

A project initialized from the template has this logical structure:

```text
project/
├── AGENTS.md                 # short mandatory agent entry point
├── pad.yaml                  # project identity and autonomy policy
├── .agent/
│   ├── README.md             # memory index and read order
│   ├── context/              # stable project/domain/architecture facts
│   ├── requirements/         # durable product and quality requirements
│   ├── tasks/                # bounded work items and acceptance criteria
│   ├── decisions/            # append-only ADR history
│   ├── state/
│   │   ├── current.yaml      # compact machine-readable snapshot
│   │   └── HANDOFF.md        # single canonical continuation note
│   ├── agents/               # role profiles and authority boundaries
│   ├── protocols/            # repeatable operating procedures
│   └── schemas/              # metadata contracts and enum values
└── .pad/
    └── template-base/        # CLI-managed base for three-way updates
```

## Separation of concerns

- `AGENTS.md` tells an agent where to start and what it must never do.
- `pad.yaml` configures the project and its autonomy boundary.
- `.agent/context/` stores stable facts, not session-by-session logs.
- `.agent/requirements/` records what must be true.
- `.agent/tasks/` records work to make it true.
- `.agent/decisions/` records why material choices were made.
- `.agent/state/` records what is true now and what should happen next.
- `.agent/agents/` defines specialized roles.
- `.agent/protocols/` defines repeatable workflows.
- `.pad/template-base/` is owned by the PAD updater and supports safe three-way merges; do not edit it manually.

Keep project source code in its normal locations. PAD is a context and workflow layer, not a requirement to move application code into `.agent/`.
