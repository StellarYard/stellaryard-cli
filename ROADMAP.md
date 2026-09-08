# ROADMAP: stellaryard-cli

> **This file must be updated with every contribution.** Before opening a PR: mark completed items done, add newly-surfaced work, or note if your change invalidates an assumption below. See `AGENTS.md` rule 6.

Status legend: `[ ]` not started · `[~]` in progress · `[x]` done

## Blocking dependency

- [ ] **`stellaryard-core`'s `/api/openapi.yaml` must be merged before any item below starts.** See `stellaryard-core/ROADMAP.md` Phase 0.

## Phase 0 — Foundation

- [ ] Scaffold Go module, `cobra` command tree skeleton, CI (lint + test on PR)
- [ ] Generate API client from `stellaryard-core`'s `openapi.yaml`; set up regeneration step
- [ ] Formatter package (`table` / `json` output modes) — build once, shared by every command, not reimplemented per-command
- [ ] Exit code handling wired at the top level (root command error handling), so individual subcommands don't each reinvent exit-code logic

## Phase 1 — Containers

- [ ] `containers start|stop|status`
- [ ] `logs <container> [--follow]` (WS streaming)

## Phase 2 — Accounts

- [ ] `accounts create [--label]`
- [ ] `accounts list [--format table|json]`
- [ ] `accounts show <publicKey>`

## Phase 3 — Ledger

- [ ] `ledger snapshot`
- [ ] `ledger tx list [--format table|json] [--limit N]` (depends on core's transaction detail shape being finalized — see `stellaryard-core/ROADMAP.md` Phase 3 open question)

## Phase 4 — Contracts

- [ ] `contracts deploy <wasm-path>`
- [ ] `contracts invoke <contractId> <method> [args...]`

## Phase 5 — Hardening (required for "100% ready")

- [ ] Consistent error message formatting across all commands (stderr, `--format json` error shape) — audit for consistency, not just presence
- [ ] Integration test suite that runs commands against a real local `stellaryard-core` instance (not just mocked client) — container/network commands especially need this given how much can go wrong at that boundary
- [ ] Shell completion (bash/zsh) — not scoped in any phase above yet, needs explicit decision on whether it's in v1 scope

## Explicitly deferred

- [ ] TUI/interactive mode — not v1, don't start without an explicit scope change in `PRD.md`
- [ ] Mainnet commands — blocked on `stellaryard-core` shipping mainnet support

## What would break

- **Exit code misclassification**: If exit 2/3 boundary is wrong, CI scripts either fail hard on retryable errors or retry hard failures. Must be tested explicitly.
- **Generated client drift**: Core ships endpoint change without updating `openapi.yaml` → CLI compiles against stale types → panics or wrong output at runtime.
- **`--format json` not applied to errors**: Human-readable stderr while JSON mode is active breaks CI parsers.
- **Unbounded list output**: Thousands of transactions produce unbounded stdout that breaks pipes.
- **Dangling WS on Ctrl+C**: `logs --follow` without clean close leaks connections on core.
- **Core startup race**: First command after `docker-compose up` may fail with exit 2.

## Edge cases not yet addressed

- Empty list output shape — must be documented and consistent (`[]` vs `{"accounts": []}`)
- WASM file validation before deploy — no local check, cryptic core errors for bad files
- Windows backslash paths — `contracts deploy .\path\to\file.wasm` needs normalization
- Concurrent CLI instances — two `accounts create` commands race on Friendbot funding
- Large JSON output — no streaming or paginated output for big lists
- `--format json` with zero results — shape must be documented

## What's overengineered for V1

- Four exit codes (2/3 distinction) — useful but no documented CI use cases yet
- Separate formatter package — flat helpers would suffice for ~7 commands
- Shell completion — no confirmed demand, adds complexity

## Open questions blocking full readiness

- Transaction detail shape — inherited unknown from `stellaryard-core/ROADMAP.md` Phase 3
- Shell completion — in scope or not, undecided
- Nil-safety wrappers for generated client responses — needed to prevent panics
- Empty output shape contract — must be defined and enforced across all commands
