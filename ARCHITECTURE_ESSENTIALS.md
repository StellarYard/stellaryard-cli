# Architecture Essentials: stellaryard-cli

> Quick-reference outline. Full detail: ARCHITECTURE.md.

## Stack
- Go 1.22+, `cobra` for CLI, stdlib `tabwriter` for table output
- API client generated from `stellaryard-core`'s `openapi.yaml` — same source dashboard's TS client uses

## Non-negotiable rule
**No hand-written HTTP calls to core, no independent domain types.** Always go through the generated client. This is what keeps CLI and dashboard from silently disagreeing about core's API shape.

## Repo role
Scriptable terminal client for core's API. Differentiator from dashboard is automation-friendliness (exit codes, `--format json`), not just "text instead of UI" — keep that in mind when scoping new commands.

## Command tree
`containers` / `accounts` / `ledger` / `contracts` / `logs` — see full mapping table in ARCHITECTURE.md

## Exit codes (memorize, this is load-bearing for CI use)
- `0` success
- `1` command/arg error
- `2` core unreachable
- `3` core-side application error

## Explicit v1 non-goals
- No TUI/interactive mode
- No independent signing/key logic — delegates to core entirely
- No mainnet commands until core exposes them

## What would break
- Exit code misclassification (exit 2 vs 3) → CI scripts retry hard failures or fail on retryable ones
- Generated client drift from core's `openapi.yaml` → CLI compiles but returns wrong data or panics
- `--format json` ignored on error paths → CI parsers fail on human-readable stderr
- Unbounded list output → pipes break, memory exhaustion
- Dangling WS on Ctrl+C → leaked connections on core
- Core startup race → first command after `docker-compose up` fails with exit 2

## Edge cases missing
- Empty list output shape (JSON array vs empty object) — must be documented and consistent
- WASM file not validated locally before deploy → cryptic core errors
- Windows backslash paths in `contracts deploy` → needs normalization
- Concurrent CLI instances racing on Friendbot funding
- `--format json` with no data → must return documented empty shape

## What's overengineered
- Four exit codes (2/3 distinction) for V1 with no documented CI use cases
- Separate formatter package for ~7 commands (flat helpers would suffice)
- Shell completion (no confirmed demand, adds complexity)

## When in doubt
If a new command needs a core endpoint that doesn't exist yet, that's a core repo issue first. Don't fake it client-side or call raw HTTP outside the generated client.
