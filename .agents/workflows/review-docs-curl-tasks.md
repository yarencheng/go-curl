---
description: Workflow guide for AI agents to audit upstream parity and maintain docs/tasks.md
---

// turbo-all

This workflow instructs the AI agent on how to systematically review and identify missing parity features directly from the upstream `curl` repository, maintaining alignment checklists in `docs/tasks.md`.

## Goal
Achieve high-fidelity functional parity with `curl` by identifying implementation gaps and tracking them systematically.

## Phase 1: Context Discovery
Examine the tracking file located at `docs/tasks.md`.
- Review existing entries for the target flag or category to avoid redundant research.

## Phase 2: Upstream Deep-Dive
Investigate the original implementation in the upstream curl repository:
- **Flag Parsing & Logic:** `third_party/curl/src/tool_getparam.c`
- **Help/Usage:** `third_party/curl/src/tool_listhelp.c`, `third_party/curl/src/tool_help.c`

**Action:** Use `grep_search` to scan for:
- `long_options`: To find all available flags.
- `Parameter Error`: To find flag-specific error handling and validation logic.
- `CASE(`: In `tool_getparam.c` to identify handling for specific flag characters.
- Core operation flags (e.g., `-X`, `-d`, `-H`).

## Phase 3: Local Implementation Audit
Search the local codebase (`internal/`, `cmd/`) for the Go implementation of the `curl` command.
- Identify which flags are already handled.
- Check `internal/commands/curl/curl.go` for the `Execute` logic.
- Check for existing tests to confirm implementation status.

## Phase 4: Gap Mapping & tasks.md Update
Log every identified flag and feature in a strict hierarchical checkbox list in `docs/tasks.md`.

**Review Checkbox Syntax:**
- `[x]` : Feature is fully implemented and verified. **Requirement:** Link to the local Go implementation.
- `[ ]` : Feature is incomplete or missing. **Requirement:** Link to the specific file and line range in `third_party/curl/src/*` where this is handled upstream.
- `[-]` : Deliberately skipped. **Requirement:** State a brief rationale (e.g., *[Not Sandbox-friendly] Relies on OS network stack features*).

### Example Output for `docs/tasks.md`

```markdown
## Parity: curl Flags

- [x] Flag `-X` (custom request): `internal/commands/curl/curl.go`
- [ ] Flag `-d` (post data): `third_party/curl/src/tool_getparam.c:L...`
- [-] Flag `--interface`: Not sandbox-friendly (requires bind to specific local IP).
```

## Phase 5: Actionable Planning
After finalizing the audit, propose the next steps (e.g., "Implement `-H` flag for headers") and suggest a TDD approach.
