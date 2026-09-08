# stellaryard-cli

A scriptable terminal client for **StellarYard**. It talks to [stellaryard-core](../stellaryard-core)'s REST/WebSocket API — the same contract the [dashboard](../stellaryard-dashboard) consumes — but exposes it as composable CLI commands instead of a browser UI.

This is not "the dashboard but text." Its differentiator is **automation-friendliness** — predictable exit codes, machine-parseable `--format json` output, and scriptability for CI pipelines.

## What StellarYard Does

StellarYard is a local development environment manager for the Stellar network. The CLI lets you manage containers, accounts, contracts, and ledger data entirely from the terminal.

## V1 Commands

### Container Management
```bash
stellaryard containers start --name horizon
stellaryard containers stop --name soroban-rpc
stellaryard containers status
```

### Account Management
```bash
stellaryard accounts create --label "my-account"
stellaryard accounts list --format table|json
stellaryard accounts show <publicKey>
```

### Ledger Inspection
```bash
stellaryard ledger snapshot
stellaryard ledger tx list --format table|json --limit N
```

### Contract Operations
```bash
stellaryard contracts deploy <wasm-path>
stellaryard contracts invoke <contractId> <method> [args...]
```

### Log Streaming
```bash
stellaryard logs horizon --follow
stellaryard logs soroban-rpc
```

## Exit Code Convention

Designed for CI/automation — every code is intentional and distinguishable:

| Code | Meaning | Example |
|------|---------|---------|
| `0` | Success | Command completed successfully |
| `1` | Command/argument error | Bad args, invalid contract ID |
| `2` | Core unreachable | Connection refused, timeout |
| `3` | Core application error | Deploy failed, contract error |

CI scripts can branch on exit codes to retry/wait vs. fail hard.

## Machine-Readable Output

Every `list` and `show` command supports `--format json` for scripting:

```bash
# Human-readable (default)
stellaryard accounts list

# Machine-parseable for CI/scripts
stellaryard accounts list --format json
```

Error output also supports `--format json` on stderr for structured error handling.

## Tech Stack

- **Language**: Go (1.22+)
- **CLI Framework**: `cobra` (standard for Go CLIs, good subcommand ergonomics)
- **Output Formatting**: `--format table|json` via internal formatter package; table output via stdlib `tabwriter`
- **API Client**: Thin Go client generated from `stellaryard-core`'s `openapi.yaml` — same source of truth the dashboard's TypeScript client uses

## Why a Generated Client?

Hand-written HTTP calls in two repos against one evolving API leads to silent contract drift. The CLI's API client is generated from the same OpenAPI spec as the dashboard's TypeScript client. If core's schema changes, both clients are regenerated — not manually patched.

## System Architecture

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

| Command | Core Endpoint |
|---------|---------------|
| `containers start/stop/status` | `/containers/*` |
| `accounts create/list/show` | `/accounts*` |
| `ledger snapshot` | `/ledger/snapshot` |
| `ledger tx list` | `/ledger/transactions` |
| `contracts deploy` | `/contracts/deploy` |
| `contracts invoke` | `/contracts/{id}/invoke` |
| `logs <container> --follow` | WS `/containers/{name}/logs` |

## Roadmap

Full roadmap: [`ROADMAP.md`](./ROADMAP.md) — **must be updated with every contribution** (see agent rules below).

**Blocking dependency**: `stellaryard-core`'s `/api/openapi.yaml` must be merged before any item below starts.

| Phase | Status | Scope |
|-------|--------|-------|
| **0 — Foundation** | Not started | Go module scaffold, cobra command tree, generate API client from core's OpenAPI, formatter package, exit code wiring |
| **1 — Containers** | Not started | `containers start|stop|status`, `logs --follow` |
| **2 — Accounts** | Not started | `accounts create/list/show` |
| **3 — Ledger** | Not started | `ledger snapshot`, `ledger tx list` (blocked on core's Phase 3 XDR decision) |
| **4 — Contracts** | Not started | `contracts deploy`, `contracts invoke` |
| **5 — Hardening** | Not started | Error formatting audit, integration tests against real core, shell completion (scope TBD) |

### Open Questions
- Transaction detail shape — inherited unknown from `stellaryard-core` Phase 3
- Shell completion — in scope or not, undecided

### Explicitly Deferred (not v1)
- TUI/interactive mode — don't start without explicit scope change in `PRD.md`
- Mainnet commands — blocked on core shipping mainnet support

## Known Failure Modes

1. **Exit code misclassification**: The most dangerous CLI bug. If a core-side error (exit 3) is misclassified as a command error (exit 1), CI scripts will treat retryable failures as hard failures. If a core-unreachable error (exit 2) is misclassified as an application error (exit 3), CI scripts will fail hard instead of retrying. The boundary between exit 2 and exit 3 is determined by whether core returned *any* HTTP response — a timeout is exit 2, a 500 is exit 3.
2. **Generated client drift**: The API client is generated from core's `openapi.yaml`. If core ships an endpoint change without updating the spec, the CLI will compile against stale types and produce wrong output or panic at runtime. No build-time drift detection exists yet.
3. **`--format json` on stderr for errors**: Error output must also be machine-parseable when `--format json` is set. If a command prints human-readable error text to stderr while `--format json` is active, CI parsers will fail. This is easy to miss because errors are tested manually, not in automated output-format tests.
4. **Large output breaking pipes**: `stellaryard ledger tx list --format json` with thousands of transactions will produce unbounded JSON output. Piping to `jq` or writing to a file may hit OS pipe buffer limits or disk issues. No streaming/paginated output mode exists.
5. **WS log streaming signal handling**: `stellaryard logs horizon --follow` opens a WebSocket. If the user presses Ctrl+C, the WS connection must be cleanly closed. Ungraceful termination leaves a dangling connection on core's side.
6. **Core connection state at command start**: The CLI doesn't check if core is reachable before executing. A slow core startup means the first command after `docker-compose up` may fail with exit 2 even though core is starting. No `--wait` or retry flag exists.

## Edge Cases Not Yet Addressed

- Empty account list: `stellaryard accounts list` with no accounts should show a clear message, not an empty table or JSON `[]` with no context.
- WASM file validation: `stellaryard contracts deploy` doesn't validate the WASM file locally before sending to core. A corrupted or non-WASM file gets sent over the network and fails with a cryptic core error.
- Windows path handling: `contracts deploy ./path/to/contract.wasm` — backslash paths on Windows need normalization.
- Concurrent CLI instances: Two `stellaryard accounts create` commands running simultaneously can race on Friendbot funding.
- `--format json` with empty results: What does `stellaryard accounts list --format json` return when there are no accounts? An empty array `[]`? An object with an empty `accounts` field? This shape must be consistent and documented.
- Signal handling for `--follow` mode: Ctrl+C must close the WS connection cleanly, not just kill the process.

## What's Overengineered for V1

- **Four distinct exit codes**: The distinction between exit 2 (core unreachable) and exit 3 (core application error) is useful for CI, but V1 has a single consumer (this CLI) and no documented CI integration examples. Two exit codes (0/1) might suffice until real CI use cases emerge. The four-code convention is forward-looking but adds testing burden now.
- **Separate formatter package**: For V1 with ~7 commands, a simple `formatTable()`/`formatJSON()` helper function set might suffice. A full package with interfaces adds indirection that's hard to justify until there are 15+ commands.
- **Shell completion**: Explicitly undecided in the roadmap. Adding it in V1 adds cobra complexity for a feature that has no confirmed user demand.

## V1 Non-Goals

- ❌ TUI / interactive mode (plausible future feature, not v1)
- ❌ Independent key generation or signing logic — delegates to core entirely
- ❌ Mainnet commands until core exposes mainnet endpoints

## Contributing

- New subcommands are addable as isolated Wave issues without touching the underlying core client library.
- **If a new command needs a core endpoint that doesn't exist yet** — that's a core repo issue first. Don't fake it client-side or call raw HTTP outside the generated client.
- The CLI must remain scriptable and composable — predictable exit codes, machine-parseable output.
- **Every PR must update `ROADMAP.md`** — mark completed items, add new work, or note invalidated assumptions. An unupdated roadmap is an incomplete PR.

### Agent Instructions

This repo includes instructions for AI coding agents:
- [`AGENTS.md`](./AGENTS.md) — Non-negotiable rules (no hand-written HTTP, exit codes are load-bearing, `--format json` required on every list/show command, no TUI in v1, roadmap updates required)
- [`CLAUDE.md`](./CLAUDE.md) — Claude-specific notes (verify core endpoints exist before building, don't forget `--format json`, update ROADMAP)

## License

See [LICENSE](./LICENSE).
