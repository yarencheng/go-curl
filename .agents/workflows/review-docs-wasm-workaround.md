---
description: Workflow for documenting technical deviations, workarounds, and sandbox limitations in docs/functional_gap.md
---

// turbo-all

This workflow guides the AI agent on how to capture and document technical rationales for `curl` features that cannot be implemented 1:1 with upstream behavior, ensuring all deviations are justified and tracked in `docs/functional_gap.md`.

## Goal
Provide a transparent record of the implementation's technical limitations and the strategic decisions behind specific `curl` feature workarounds.

## Phase 1: Audit Alignment
Before starting, review `docs/tasks.md` for the target flag or feature.
- Identify flags/features marked as `[-]` (skipped) or `[x]` (implemented but potentially via workaround).
- If a flag in the task list has a brief note like *[Not Sandbox-friendly]*, it likely requires a detailed entry here.

## Phase 2: Technical Root Cause Analysis
Investigate why the `curl` feature deviates from upstream:
- **Sandbox Barriers**: Does it require OS-level socket control or low-level network access unavailable in the browser/WASM?
- **Security Policy**: Does it conflict with browser CORS or security headers?
- **Filesystem Model**: Does it rely on persistent disk storage not available in the in-memory VFS?

## Phase 3: Workaround Documentation
If a workaround exists:
- Locate the internal implementation in `internal/commands/curl/`.
- Verify the specific logic that replaces the upstream behavior.

## Phase 4: functional_gap.md Update
Update the `docs/functional_gap.md` file using the established status codes:

- `[x]` **Workaround**: For features implemented using simulator-specific logic. 
  - **Requirement**: Path to local Go code + Rationale.
- `[-]` **Unsupported**: For features explicitly excluded.
  - **Requirement**: Path to upstream C code in `third_party/curl/src/` + Detailed rationale.
- `[ ]` **Pending**: For identified gaps without a clear resolution path.

### Example Entry

```markdown
### SSL/TLS
- `[-]` **Custom CA Bundles**: Browser environments typically manage their own certificate trust; `-k` (insecure) may be simulated.
  > Rationale: The browser's `fetch` API does not allow overriding the certificate validation logic programmatically for individual requests.
```

## Phase 5: Verification
Ensure that every entry in the gap map has a corresponding reference in `docs/tasks.md` to maintain a single source of truth for progress.
