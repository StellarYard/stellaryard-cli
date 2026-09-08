# AGENTS.md — stellaryard-cli

Instructions for AI coding agents working in this repository. Read this before making changes.

## What this repo is

Scriptable terminal client for StellarYard, consuming `stellaryard-core`'s REST/WS API. Its differentiator from `stellaryard-dashboard` is automation-friendliness — exit codes, `--format json` — not just "text instead of UI." See `PRD.md` and `ARCHITECTURE.md` for full context; `ARCHITECTURE_ESSENTIALS.md` for a fast-reference outline.

## Non-negotiable rules

1. **No hand-written HTTP calls to core.** Always use the generated API client (from `stellaryard-core`'s `openapi.yaml`). If the client doesn't support what you need, that's a signal core's spec is missing something — flag it, don't route around the client with raw `net/http` calls.
2. **No independent signing or key-generation logic.** Delegates entirely to core's `/accounts` endpoints and `Signer` interface, same as dashboard.
3. **Exit code convention is load-bearing — do not casually add new exit codes or repurpose existing ones.** `0` success, `1` command/arg error, `2` core unreachable, `3` core-side application error. CI scripts depend on this being stable; changing it is a breaking change to every consumer of this CLI, not a small tweak.
4. **Every list/show command must support `--format json`, not just the default table format.** A new command that only implements table output is incomplete — this is the CLI's core value proposition over the dashboard, not an optional nice-to-have.
5. **No TUI/interactive mode in v1.** If a task description implies building interactive/`bubbletea`-style UI, that's scope creep beyond what's currently defined — flag it rather than building it.
6. **Update `ROADMAP.md` in every contribution.** Before opening a PR: mark your item done, add newly-surfaced work, or note if your change invalidates a roadmap assumption. Treat an unupdated `ROADMAP.md` as an incomplete PR.

## Before starting any task

- Read `ARCHITECTURE_ESSENTIALS.md` first.
- Check `ROADMAP.md` for whether your task is already scoped and what phase it belongs to.
- Confirm the core API endpoint you need actually exists in `stellaryard-core`'s current `openapi.yaml` before building a command around it.

## What you cannot do from this repo

- Modify `stellaryard-core` or `stellaryard-dashboard`. If your task needs a core API change, note it in your PR description and in this repo's `ROADMAP.md` as a blocked/dependent item.

## Known pitfalls

- **Exit code boundaries are fragile.** The line between exit 2 (core unreachable) and exit 3 (core error) depends on whether core returned *any* HTTP response. A TCP reset is exit 2; a 500 is exit 3. A timeout could be either. Test both paths explicitly, don't assume.
- **Generated client types can panic on nil.** If core returns a field as `null` where the generated Go type expects a string, you'll get a nil pointer dereference. Always nil-check generated client responses before accessing fields.
- **`--format json` must apply to errors too.** When `--format json` is active, stderr error output must also be JSON. Human-readable error text on stderr while the user expects JSON will break CI parsers.
- **WS connections need clean close on signal.** If you're implementing `logs --follow`, ensure Ctrl+C triggers a WS close handshake, not just `os.Exit`. Dangling connections leak on core's side.
- **Empty output needs a documented shape.** `accounts list --format json` with zero accounts must return a consistent, documented empty shape (e.g., `{"accounts": []}`), not `null` or `[]`.
