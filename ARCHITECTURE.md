# Architecture: stellaryard-cli

## Tech Stack

- **Language**: Go (1.22+) — same as core; also lets this repo import a shared Go client library for core's API rather than reimplementing HTTP calls independently (see below)
- **CLI framework**: `cobra` (standard for Go CLIs, well-understood, good subcommand ergonomics matching the command tree in the PRD)
- **Output formatting**: `--format table|json` via a small internal formatter package; table output via `tabwriter` (stdlib, no extra dependency needed for this)
- **API client**: a thin Go client package generated from `stellaryard-core`'s `openapi.yaml` (same source of truth dashboard's TypeScript client is generated from) — this is the mechanism that prevents CLI and dashboard from silently drifting apart in what they assume core's API looks like

## Why This Stack

Go + Cobra is the default choice for CLIs in this ecosystem and keeps this repo in the same language as core, which matters for Wave contributors who may move between core and CLI issues. The critical decision here isn't the framework choice — it's generating the API client from the same OpenAPI spec dashboard uses, instead of hand-writing HTTP calls. Hand-written clients in two repos against one evolving API is how you get silent contract drift; a shared generation source turns a "did the API change?" question into a build-time regeneration step.

## System Overview

```
┌───────────────────────────┐
│ stellaryard-cli             │
│                             │
│  ┌───────────┐             │
│  │ cobra      │             │
│  │ commands   │             │
│  │ - containers│            │
│  │ - accounts │             │
│  │ - ledger   │             │
│  │ - contracts│             │
│  │ - logs     │             │
│  └─────┬─────┘             │
│        │                    │
│  ┌─────▼─────┐             │
│  │ Formatter  │             │
│  │ (table/json)│            │
│  └─────┬─────┘             │
│        │                    │
│  ┌─────▼─────┐             │
│  │ core API   │             │
│  │ client     │             │
│  │ (generated)│             │
│  └─────┬─────┘             │
└────────┼───────────────────┘
         │ REST + WS
         ▼
  stellaryard-core (separate repo)
```

## Command → API Mapping

| Command | Core endpoint |
|---|---|
| `containers start/stop/status` | `/containers/*` |
| `accounts create/list/show` | `/accounts*` |
| `ledger snapshot` | `/ledger/snapshot` |
| `ledger tx list` | `/ledger/transactions` |
| `contracts deploy` | `/contracts/deploy` |
| `contracts invoke` | `/contracts/{id}/invoke` |
| `logs <container> --follow` | WS `/containers/{name}/logs` |

## Exit Code Convention

This is a specific, deliberate design point because it's what makes the CLI genuinely useful in CI rather than just "dashboard with worse ergonomics":

- `0` — success
- `1` — command-level error (bad args, invalid contract ID, etc.)
- `2` — core unreachable (connection refused, timeout) — distinguishable so a CI script can retry/wait rather than treat it as a hard failure
- `3` — core reported an application error (e.g., deploy failed) — response body error detail printed to stderr, machine-parseable with `--format json`

## Known Failure Modes

- **Exit code misclassification**: The boundary between exit 2 (core unreachable) and exit 3 (core application error) depends on whether core returned *any* HTTP response. A TCP connection reset is exit 2; a 500 Internal Server Error is exit 3. A timeout could be either depending on whether the TCP handshake completed. This boundary must be tested explicitly, not assumed.
- **Generated client type mismatches**: If core's actual JSON response doesn't match the generated Go types (e.g., a field is `null` where the type expects a string), the CLI will panic or return garbage. The generated client needs nil-safety wrappers, not raw generated code.
- **`--format json` on stderr**: Error output must also be JSON when `--format json` is active. A human-readable error string on stderr while the user expects JSON on stdout will break CI parsers that read both streams.
- **WS signal handling**: `logs --follow` opens a persistent WebSocket. Ctrl+C must trigger a clean WS close handshake, not just os.Exit. A dangling WS connection leaks resources on core's side.
- **Large list output**: No streaming or paginated JSON output. A ledger with thousands of transactions produces unbounded stdout that can break pipes.

## What's Overengineered for V1

- **Four exit codes**: The 2/3 distinction (core unreachable vs core error) is useful for CI retry logic, but V1 has no documented CI examples. Two codes (success/failure) would suffice until real automation use cases emerge.
- **Separate formatter package**: For ~7 commands, a shared formatter package with interfaces is forward-looking but adds indirection. A flat helper function set would be simpler.

## Non-Goals (v1)

- No TUI/interactive mode.
- No independent key generation or signing logic — delegates to core's `/accounts` and `Signer` interface entirely, same as dashboard.
- No mainnet commands until core supports them.
