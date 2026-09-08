# CLAUDE.md — stellaryard-cli

This project's agent instructions live in [`AGENTS.md`](./AGENTS.md). Read it in full before making changes — it is the source of truth and this file does not duplicate it.

Claude-specific notes:

- Before implementing a new subcommand, verify the backing core API endpoint actually exists (check `stellaryard-core`'s current `openapi.yaml`, not just its roadmap).
- Don't forget `--format json` support — it's easy to ship table-only output and consider a command "done." It isn't, per `AGENTS.md` rule 4.
- Update `ROADMAP.md` as part of the same change, not as an afterthought — see `AGENTS.md` rule 6.
- **Exit code misclassification is the #1 CLI bug.** Always test the boundary between exit 2 (core unreachable) and exit 3 (core error). A timeout during TCP handshake is exit 2; a timeout after receiving a partial response is exit 3. This distinction matters for CI retry logic.
- **Nil-check generated client responses.** The generated Go types don't guarantee non-nil fields. A `null` JSON field will panic if accessed without a nil check. Always wrap generated client calls with nil guards.
- **`--format json` applies to errors too.** When the user passes `--format json`, stderr error output must also be JSON. Don't print human-readable error strings on stderr when JSON mode is active.
- **Empty output shapes must be consistent.** If `accounts list` returns `[]` (JSON array), then `accounts list --format json` with zero results must also return `[]`, not `null` or `{"accounts": []}`. Pick one shape and enforce it everywhere.
