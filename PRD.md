# PRD: stellaryard-cli

## What We're Building

`stellaryard-cli` is a terminal client for StellarYard. It talks to `stellaryard-core`'s REST/WS API — same contract the dashboard consumes — but exposes it as scriptable commands instead of a browser UI. It's for developers who want StellarYard's orchestration (containers, accounts, contracts) inside scripts, CI, or a terminal-first workflow, without opening a browser.

Like dashboard, this repo holds no independent business logic. Unlike dashboard, it has one additional responsibility dashboard doesn't: **it must be scriptable and composable** (predictable exit codes, machine-parseable output modes), because its primary differentiator from the dashboard isn't "worse UI," it's "usable in automation." If this repo just becomes "the dashboard's commands but text," it hasn't earned being a separate client.

## Who It's For

- Developers who prefer terminal workflows over a browser dashboard for day-to-day local dev.
- CI/automation: e.g., a test pipeline that needs to spin up a local network, fund accounts, deploy a contract, run tests, tear down — scriptably, without a browser.
- Wave contributors adding individual subcommands as isolated issues (e.g., "add `stellaryard tx list --format json`").

## What The Product Actually Needs To Do

### V1

1. **Container commands**: `stellaryard containers start|stop|status [--name horizon|soroban-rpc]`
2. **Account commands**: `stellaryard accounts create [--label NAME]`, `stellaryard accounts list [--format table|json]`, `stellaryard accounts show <publicKey>`
3. **Ledger commands**: `stellaryard ledger snapshot`, `stellaryard ledger tx list [--format table|json] [--limit N]`
4. **Contract commands**: `stellaryard contracts deploy <wasm-path>`, `stellaryard contracts invoke <contractId> <method> [args...]`
5. **Log streaming**: `stellaryard logs <container-name> [--follow]`, consuming core's WS endpoint.
6. **Machine-readable output**: every list/show command supports `--format json` for scripting; default is human-readable table.
7. **Predictable exit codes**: 0 success, non-zero on failure, distinguishable failure classes (e.g., core unreachable vs. command error) so CI scripts can branch on them.

### Explicitly NOT in V1

- No interactive/TUI mode (that's a plausible future repo or feature, not v1 — don't scope-creep an issue into building a `bubbletea`-style TUI when a flat command was asked for).
- No local key generation independent of core — CLI calls core's `/accounts` endpoint like dashboard does, it doesn't implement its own signer.
- No mainnet commands until core exposes mainnet endpoints.

## Success Criteria

- A developer can replace their current manual CLI/curl workflow entirely with `stellaryard` subcommands.
- A CI pipeline can use the CLI non-interactively with `--format json` and rely on exit codes, with no human in the loop.
- New subcommands are addable as isolated Wave issues without touching the underlying core client library.

## What Would Break

- **Exit code misclassification**: If a core timeout (should be exit 2) is misclassified as an application error (exit 3), CI scripts will fail hard instead of retrying. This is the highest-risk bug because it silently breaks automation.
- **Generated client drift**: If core ships an endpoint change without updating `openapi.yaml`, the CLI compiles against stale types and panics or returns wrong data at runtime. No build-time drift detection exists.
- **`--format json` inconsistency on errors**: If error output is human-readable while `--format json` is active, CI parsers fail. Every error path must respect the active format.
- **Unbounded output in list commands**: `stellaryard ledger tx list --format json` with thousands of transactions produces unbounded output that can break pipes or exhaust memory.
- **Dangling WebSocket on Ctrl+C**: `logs --follow` must cleanly close the WS connection on signal. Ungraceful termination leaves orphaned connections on core.
- **Core startup race**: The first command after `docker-compose up` may fail with exit 2 if core hasn't started yet. No retry/wait mechanism exists.

## Edge Cases Not Yet Addressed

- Empty list output: `accounts list` with no accounts — empty table vs empty JSON array? Shape must be consistent and documented.
- WASM file validation: `contracts deploy` doesn't validate the file locally before sending to core. Cryptic core errors for non-WASM files.
- Windows path handling: Backslash paths in `contracts deploy` need normalization.
- Concurrent CLI instances: Two `accounts create` commands can race on Friendbot funding.
- `--format json` with no data: What shape does `accounts list --format json` return when there are zero accounts? Must be documented.
- Large transaction pages: No streaming/paginated JSON output for lists with thousands of entries.

## What's Overengineered

- **Four exit codes**: Useful for CI, but V1 has no documented CI integration examples. Two codes (0/1) might suffice until real use cases emerge. The four-code convention adds testing burden now for a feature without confirmed demand.
- **Separate formatter package**: For ~7 commands, a simple helper function set might suffice. A full package with interfaces adds indirection.
- **Shell completion**: No confirmed user demand. Adds cobra complexity for V1.
