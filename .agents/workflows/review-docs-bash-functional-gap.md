---
description: Workflow for documenting technical deviations, workarounds, and sandbox limitations in docs/bash/functional_gap.md
---

// turbo-all

This workflow guides the AI agent on how to capture and document technical rationales for Bash features that cannot be implemented 1:1 with upstream behavior, ensuring all deviations are justified and tracked in `docs/bash/functional_gap.md`.

## Goal
Provide a transparent record of the simulator's technical limitations and the strategic decisions behind specific Bash builtins workarounds.

## Phase 1: Audit Alignment
Before starting, review `docs/bash/tasks.md` for the target command.
- Identify flags/features marked as `[-]` (skipped) or `[x]` (implemented but potentially via workaround).
- If a flag in the task list has a brief note like *[Not Simulation-friendly]*, it likely requires a detailed entry here.

## Phase 2: Technical Root Cause Analysis
Investigate why the Bash feature deviates from upstream:
- **Sandbox Barriers**: Does it require OS-level control unavailable in the browser?
- **Execution Model**: Does it rely on process suspension or complex job control?

## Phase 3: Workaround Documentation
If a workaround exists:
- Locate the internal implementation in `internal/commands/` or `internal/shell/`.
- Verify the specific logic that replaces the upstream behavior.

## Phase 4: functional_gap.md Update
Update the `docs/bash/functional_gap.md` file using the established status codes:

- `[x]` **Workaround**: For features implemented using simulator-specific logic. 
  - **Requirement**: Path to local Go code + Rationale.
- `[-]` **Unsupported**: For features explicitly excluded.
  - **Requirement**: Path to upstream C code + Detailed rationale (e.g., Process management limitation).
- `[ ]` **Pending**: For identified gaps without a clear resolution path.

### Example Entry

```markdown
### `suspend`

- `[-]` Process Suspension (Unsupported): `internal/commands/suspend/suspend.go`
  > Rationale: WebAssembly processes in the browser cannot be suspended/resumed by the shell in the same way as OS-level processes.
```

## Phase 5: Verification
Ensure that every entry in the gap map has a corresponding reference (even if just a placeholder) in `docs/bash/tasks.md` to maintain a single source of truth for progress.
