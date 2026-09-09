# stellaryard-cli — Phase 1 Issues

---

# [Phase 1] Containers start command implementation

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done
>
> **No hand-written HTTP calls to core. Always use the generated API client.**

stellaryard-cli is the terminal client for StellarYard. It talks to stellaryard-core's API and exposes it as scriptable commands. This issue implements the `containers start` command.

---

## Problem

Developers need to start containers from the terminal without opening the dashboard. The `stellaryard containers start` command provides a scriptable way to start Horizon or Soroban RPC containers.

---

## Scope

**In scope:**
- `stellaryard containers start --name <horizon|soroban-rpc>` command
- Default name: "horizon" when --name not provided
- Call core API: POST /containers/{name}/start
- Print confirmation message
- Exit codes: 0 success, 2 core unreachable, 3 core error

**Out of scope:**
- Stop or status commands (separate issues)
- --format json support (separate issue: CL-1.04)
- Log streaming (separate issue)

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `stellaryard containers start --name horizon` calls core API and prints confirmation
2. Default name is "horizon" when --name not provided
3. Prints "Container 'horizon' started" on success
4. Exit code 0 on success
5. Exit code 2 when core is unreachable
6. Exit code 3 when core returns an error
7. Tests verify: success output, default name, exit codes

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/containers.go` — Implement start command

### Architecture constraints
- Exit code must be correct — this is load-bearing for CI
- Use the API client from Phase 0, not raw HTTP calls
- Follow the cobra command pattern already established

### Suggested approach
1. In `internal/cmd/containers.go`, implement `containersStartCmd.RunE`
2. Get `--name` flag value (default "horizon")
3. Call API client `StartContainer(name)`
4. Print confirmation message
5. Return nil on success, classify error for exit code

### Testing approach
- Test prints confirmation on success
- Test uses default name "horizon"
- Test exits with code 2 when core unreachable
- Test exits with code 3 when core error

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `containers start` command exists
- [ ] Default name is "horizon"
- [ ] Prints confirmation on success
- [ ] Exit code 0 on success
- [ ] Exit code 2 when core unreachable
- [ ] Exit code 3 when core error
- [ ] Uses API client (not raw HTTP)
- [ ] Tests verify output and exit codes
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented the command as described
- [ ] I have written tests for output and exit codes
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: implement containers start command`

---

## References

- `ARCHITECTURE.md` → "Command → API Mapping" → containers start/stop/status
- `ARCHITECTURE.md` → "Exit Code Convention"
- `api/openapi.yaml` → POST /containers/{name}/start

---

# [Phase 1] Containers stop command implementation

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done
>
> **No hand-written HTTP calls to core. Always use the generated API client.**

stellaryard-cli needs a stop command. This issue implements `stellaryard containers stop`.

---

## Problem

Developers need to stop containers from the terminal. Stopping an already-stopped container should not be treated as an error.

---

## Scope

**In scope:**
- `stellaryard containers stop --name <name>` command
- Call core API: POST /containers/{name}/stop
- Print confirmation
- Exit codes: 0 success, 2 core unreachable, 3 core error

**Out of scope:**
- --format json support (CL-1.05)
- Force-kill functionality

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `stellaryard containers stop --name horizon` works
2. Prints "Container 'horizon' stopped"
3. Works when container is already stopped (prints confirmation, not error)
4. Correct exit codes

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/containers.go` — Implement stop command

### Architecture constraints
- Stopping an already-stopped container prints confirmation, not error
- Follow same pattern as CL-1.01 (start command)

### Suggested approach
1. Implement `containersStopCmd.RunE`
2. Call API client `StopContainer(name)`
3. Print confirmation
4. Handle exit codes

### Testing approach
- Test prints confirmation on success
- Test handles already-stopped container gracefully
- Test exit codes

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `containers stop` command exists
- [ ] Prints confirmation on success
- [ ] Works when already stopped
- [ ] Exit codes correct
- [ ] Tests verify output and exit codes
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented the command as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: implement containers stop command`

---

## References

- `ARCHITECTURE.md` → "Command → API Mapping"
- Depends on: CL-1.01 (start command pattern)

---

# [Phase 1] Containers status command with table output

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done
>
> **No hand-written HTTP calls to core. Always use the generated API client.**

stellaryard-cli needs to show container statuses. This issue implements the `containers status` command with table output.

---

## Problem

Developers need to see container status from the terminal. Without a status command, they must use `docker ps` or open the dashboard.

---

## Scope

**In scope:**
- `stellaryard containers status` command
- Call core API: GET /containers
- Table output with Name, State, Health, Started columns
- Exit codes: 0 success, 2 core unreachable, 3 core error

**Out of scope:**
- --format json support (CL-1.06)
- Filtering or sorting

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `stellaryard containers status` shows table of container statuses
2. Columns: Name, State, Health, Started
3. Uses shared formatter package
4. Shows "No containers found" when empty
5. Correct exit codes

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/containers.go` — Implement status command

### Architecture constraints
- Use shared formatter package (not reimplemented)
- Table output must use tabwriter
- Empty state: "No containers found" message

### Suggested approach
1. Implement `containersStatusCmd.RunE`
2. Call API client `ListContainers()`
3. Format as table using formatter package
4. Handle empty list case

### Testing approach
- Test table output has correct columns
- Test empty list shows message
- Test exit codes

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `containers status` command exists
- [ ] Table output with correct columns
- [ ] Uses shared formatter package
- [ ] Shows message for empty list
- [ ] Exit codes correct
- [ ] Tests verify output
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented the command as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: implement containers status command`

---

## References

- `ARCHITECTURE.md` → "Command → API Mapping"
- `internal/formatter/formatter.go` — shared formatter

---

# [Phase 1] Containers start command with --format json

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done
>
> **Every list/show command must support --format json. This is the CLI's core value proposition.**

stellaryard-cli's differentiator is automation-friendliness. This issue adds `--format json` to the start command.

---

## Problem

CI scripts need machine-parseable output. Without --format json, the start command only produces human-readable text that parsers can't handle.

---

## Scope

**In scope:**
- When --format json: output `{"status": "started", "name": "horizon"}`
- When --format table (default): human-readable confirmation
- Error output on stderr also respects --format json

**Out of scope:**
- Other commands' --format json (separate issues)
- NDJSON or streaming formats

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `--format json` produces valid JSON to stdout
2. `--format table` produces human-readable text
3. Error output on stderr is JSON when --format json is active
4. Tests verify both formats

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/containers.go` — Update start command

### Architecture constraints
- --format json on stderr for errors is critical — CI parsers read both streams
- Use shared formatter package

### Suggested approach
1. Check formatStr flag in containersStartCmd
2. If json, marshal response to JSON
3. If table, print human-readable confirmation
4. Errors on stderr also respect format

### Testing approach
- Test --format json produces valid JSON
- Test --format table produces text
- Test error with --format json produces JSON on stderr

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] --format json produces valid JSON
- [ ] --format table produces human-readable text
- [ ] Error output respects --format json
- [ ] Tests verify both formats
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented --format json as described
- [ ] I have written tests for both formats
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: add --format json to containers start command`

---

## References

- `ARCHITECTURE.md` → "Exit Code Convention"
- `AGENTS.md` → rule 4 (--format json required)

---

# [Phase 1] Containers stop command with --format json

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

Same as CL-1.04 but for the stop command.

---

## Problem

CI scripts need machine-parseable output from the stop command.

---

## Scope

**In scope:**
- JSON output: `{"status": "stopped", "name": "horizon"}`
- Table output: human-readable
- Error output respects --format json

**Out of scope:**
- Other commands' --format json

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. --format json produces valid JSON
2. --format table produces human-readable text
3. Errors are JSON on stderr
4. Tests verify both formats

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/containers.go` — Update stop command

### Architecture constraints
- Follow exact same pattern as CL-1.04

### Suggested approach
1. Same as CL-1.04 but for stop command

### Testing approach
- Same as CL-1.04

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] --format json produces valid JSON
- [ ] --format table produces human-readable text
- [ ] Error output respects --format json
- [ ] Tests verify both formats
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented --format json as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: add --format json to containers stop command`

---

## References

- Depends on: CL-1.04 (start command --format json pattern)

---

# [Phase 1] Containers status command with --format json

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

Same as CL-1.04 but for the status command with array output.

---

## Problem

CI scripts need machine-parseable status output.

---

## Scope

**In scope:**
- JSON output: array of container status objects
- Table output: tabular format
- Empty list: `[]` for JSON, "No containers found" for table

**Out of scope:**
- Filtering or sorting

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. --format json produces valid JSON array
2. --format table produces tabular output
3. Empty list returns `[]`, not `null`
4. Tests verify both formats and empty state

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/containers.go` — Update status command

### Architecture constraints
- Empty list must return `[]`, not `null`
- Use shared formatter package

### Suggested approach
1. Check formatStr flag in containersStatusCmd
2. JSON: marshal array to stdout
3. Table: use formatter.WriteTable()

### Testing approach
- Test --format json with running containers
- Test --format json with empty list returns []
- Test --format table output

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] --format json produces valid JSON array
- [ ] --format table produces tabular output
- [ ] Empty list returns `[]` for JSON
- [ ] Empty list shows message for table
- [ ] Tests verify all cases
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented --format json as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: add --format json to containers status command`

---

## References

- `AGENTS.md` → rule 4 (--format json required on every list/show command)

---

# [Phase 1] Logs command implementation (non-follow mode)

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli needs to display container logs. This issue implements non-follow mode (fetch recent logs and exit).

---

## Problem

Developers need to see container logs from the terminal. Non-follow mode is useful for scripts that need a snapshot of recent logs.

---

## Scope

**In scope:**
- `stellaryard logs <container-name>` command (no --follow)
- Fetch recent logs from core
- Print log lines to stdout
- Exit after printing
- Exit codes: 0 success, 1 missing argument, 2 core unreachable, 3 core error

**Out of scope:**
- --follow mode (CL-1.08)
- --format json (CL-1.10)
- Signal handling (CL-1.09)

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `stellaryard logs horizon` prints recent log lines
2. Exits after printing (non-follow is default)
3. Exit code 1 when no container name provided
4. Exit code 2 when core unreachable
5. Exit code 3 when core error

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/logs.go` — Implement logs command

### Architecture constraints
- Container name is required (cobra.ExactArgs(1))
- Non-follow mode is the default

### Suggested approach
1. Implement `logsCmd.RunE`
2. Validate container name argument
3. Fetch recent logs from core
4. Print each line to stdout

### Testing approach
- Test prints log lines
- Test exits with code 1 when no argument
- Test exits with code 2 when core unreachable

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `logs` command exists
- [ ] Requires container name argument
- [ ] Prints log lines in non-follow mode
- [ ] Exits after printing
- [ ] Exit codes correct
- [ ] Tests verify output and exit codes
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented the command as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: implement logs command in non-follow mode`

---

## References

- `ARCHITECTURE.md` → "Command → API Mapping" → logs

---

# [Phase 1] Logs command --follow mode with WebSocket streaming

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Medium
> **Points:** 150
> **Estimated time:** 2-4 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli needs real-time log streaming. This issue adds the --follow flag that connects to core's WebSocket log endpoint.

---

## Problem

Developers need to watch container logs in real time from the terminal. Without --follow, they must repeatedly run the command.

---

## Scope

**In scope:**
- `--follow` flag (boolean, default false)
- When --follow: connect to core WS endpoint
- Stream logs in real time
- Keep running until Ctrl+C

**Out of scope:**
- Signal handling for clean Ctrl+C (CL-1.09)
- --format json (CL-1.10)
- Reconnection logic (core-side handles this)

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `stellaryard logs horizon --follow` connects to WS endpoint
2. Prints each log line as it arrives
3. Keeps running until Ctrl+C
4. Non-follow mode still works
5. Tests verify WS connection and output

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/logs.go` — Add --follow flag

### Architecture constraints
- WebSocket URL: ws://localhost:8080/api/v1/containers/{name}/logs
- Use gorilla/websocket or nhooyr.io/websocket for WS client
- Channel must be closed when connection drops

### Suggested approach
1. Add --follow flag to logsCmd
2. When follow=true, connect to WS endpoint
3. Read messages from WS, print to stdout
4. Handle WS errors gracefully

### Testing approach
- Test --follow connects to WS
- Test non-follow still works
- Test WS error is handled

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] --follow flag exists
- [ ] --follow connects to WS endpoint
- [ ] Streams logs in real time
- [ ] Non-follow mode still works
- [ ] WS errors handled gracefully
- [ ] Tests verify WS connection
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented --follow as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: add --follow flag to logs command with WS streaming`

---

## References

- `ARCHITECTURE.md` → "Command → API Mapping" → logs
- `api/openapi.yaml` → WS /containers/{name}/logs

---

# [Phase 1] Logs --follow signal handling (clean Ctrl+C)

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Small
> **Points:** 150
> **Estimated time:** 2-4 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli's logs --follow must handle Ctrl+C gracefully. This issue ensures the WebSocket connection is closed cleanly.

---

## Problem

Without clean signal handling, Ctrl+C kills the process and leaves a dangling WebSocket connection on core. This leaks resources.

---

## Scope

**In scope:**
- Catch SIGINT during --follow mode
- Send WebSocket close frame
- Wait for server to acknowledge
- Exit with code 0

**Out of scope:**
- Other signals (SIGTERM, SIGQUIT)
- Reconnection logic

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. Ctrl+C during --follow triggers clean WS close
2. Process exits with code 0
3. No dangling WS connection on core
4. Tests verify clean close behavior

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/logs.go` — Add signal handling

### Architecture constraints
- Use signal.NotifyContext for clean shutdown
- Send WS close frame, don't just os.Exit
- Wait for server to acknowledge close

### Suggested approach
1. Use signal.NotifyContext to catch SIGINT
2. On SIGINT: send ws.CloseMessage
3. Wait for close response with timeout
4. Exit with code 0

### Testing approach
- Test SIGINT triggers clean WS close
- Test process exits with code 0

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] SIGINT triggers clean WS close
- [ ] Process exits with code 0
- [ ] No dangling WS connection
- [ ] Tests verify clean close
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented signal handling as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: handle Ctrl+C gracefully in logs --follow mode`

---

## References

- `ARCHITECTURE.md` → "Known Failure Modes" → WS signal handling
- Depends on: CL-1.08 (--follow mode)

---

# [Phase 1] Logs command --format json support

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Small
> **Points:** 150
> **Estimated time:** 2-4 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli's logs command needs --format json for CI automation.

---

## Problem

CI scripts need machine-parseable log output. Without --format json, logs are only human-readable.

---

## Scope

**In scope:**
- Non-follow + JSON: log lines as JSON strings in array
- Follow + JSON: each line as `{"line": "...", "timestamp": "..."}`
- Error output respects --format json

**Out of scope:**
- Log filtering or search
- NDJSON format (follow mode uses NDJSON naturally)

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. Non-follow + JSON: output is valid JSON array
2. Follow + JSON: each line is JSON object
3. Errors on stderr are JSON
4. Tests verify both modes

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/logs.go` — Add --format json support

### Architecture constraints
- Follow mode JSON should be NDJSON (one JSON object per line)
- Non-follow mode wraps lines in JSON array

### Suggested approach
1. Check formatStr flag
2. Non-follow JSON: wrap lines in array
3. Follow JSON: each line is `{"line": "...", "timestamp": "..."}`
4. Errors on stderr respect format

### Testing approach
- Test non-follow JSON output
- Test follow JSON output
- Test error output

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] --format json works in non-follow mode
- [ ] --format json works in follow mode (NDJSON)
- [ ] Error output respects --format json
- [ ] Tests verify both modes
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented --format json as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: add --format json to logs command`

---

## References

- `AGENTS.md` → rule 4 (--format json required)

---

# [Phase 1] API client method for container start

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done
>
> **No hand-written HTTP calls. Always use the generated API client.**

stellaryard-cli's API client needs a StartContainer method. This is used by the containers start command.

---

## Problem

The CLI commands need a typed API client method to call core's start endpoint. Without it, commands would need raw HTTP calls (violates AGENTS.md rule 1).

---

## Scope

**In scope:**
- `StartContainer(ctx, name) error` method
- HTTP POST to /api/v1/containers/{name}/start
- Return nil on success
- Return classified errors

**Out of scope:**
- Other API client methods (separate issues)
- Client generation from OpenAPI (Phase 0)

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `StartContainer()` method exists
2. Makes HTTP POST to correct endpoint
3. Returns nil on 200
4. Returns error on failure (classified as transient/permanent)
5. Tests verify success and error cases

---

## Implementation Guidelines

### Key files to modify or create
- `internal/client/client.go` — Add StartContainer method

### Architecture constraints
- Error classification matters for exit code determination
- Eventually generated from openapi.yaml — hand-written is fine for V1

### Suggested approach
1. Implement `StartContainer(ctx, name) error`
2. Create HTTP POST request
3. Check response status code
4. Return nil on 200, error on failure

### Testing approach
- Test returns nil on 200
- Test returns error on 404
- Test returns error on connection refused

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `StartContainer()` method exists
- [ ] Makes HTTP POST to correct endpoint
- [ ] Returns nil on success
- [ ] Returns classified errors on failure
- [ ] Tests verify all cases
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented the method as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: add StartContainer method to API client`

---

## References

- `api/openapi.yaml` → POST /containers/{name}/start

---

# [Phase 1] API client method for container stop

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

Same pattern as CL-1.11 but for the stop endpoint.

---

## Problem

The CLI needs a typed client method for the stop endpoint.

---

## Scope

**In scope:**
- `StopContainer(ctx, name) error` method
- HTTP POST to /api/v1/containers/{name}/stop

**Out of scope:**
- Other client methods

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `StopContainer()` method exists
2. Returns nil on success (including already-stopped)
3. Returns classified errors on failure

---

## Implementation Guidelines

### Key files to modify or create
- `internal/client/client.go` — Add StopContainer method

### Architecture constraints
- Follow same pattern as CL-1.11
- Already-stopped container returns nil, not error

### Suggested approach
1. Implement `StopContainer(ctx, name) error`
2. Same pattern as StartContainer

### Testing approach
- Test returns nil on 200
- Test returns error on connection failure

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `StopContainer()` method exists
- [ ] Returns nil on success (including already-stopped)
- [ ] Returns classified errors
- [ ] Tests verify all cases
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented the method as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: add StopContainer method to API client`

---

## References

- Depends on: CL-1.11 (StartContainer pattern)

---

# [Phase 1] API client method for container list status

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

Same pattern as CL-1.11 but for the list endpoint.

---

## Problem

The CLI needs a typed client method for listing container statuses.

---

## Scope

**In scope:**
- `ListContainers(ctx) ([]ContainerStatus, error)` method
- HTTP GET to /api/v1/containers

**Out of scope:**
- Filtering or pagination

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `ListContainers()` method exists
2. Returns typed slice
3. Returns empty slice (not nil) for empty list
4. Returns error on failure

---

## Implementation Guidelines

### Key files to modify or create
- `internal/client/client.go` — Add ListContainers method

### Architecture constraints
- Return empty slice, not nil, for empty lists
- Nil-safety: check for nil before iterating

### Suggested approach
1. Implement `ListContainers(ctx) ([]ContainerStatus, error)`
2. Parse JSON response
3. Return empty slice for empty array

### Testing approach
- Test returns slice with correct data
- Test returns empty slice for empty array
- Test returns error on failure

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `ListContainers()` method exists
- [ ] Returns typed slice
- [ ] Returns empty slice (not nil) for empty list
- [ ] Returns error on failure
- [ ] Tests verify all cases
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented the method as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: add ListContainers method to API client`

---

## References

- Depends on: CL-1.11 (client pattern)

---

# [Phase 1] API client method for WebSocket log streaming

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Medium
> **Points:** 150
> **Estimated time:** 2-4 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli's logs --follow needs a WebSocket client method. This issue adds the client-side WS connection.

---

## Problem

The logs command needs to connect to core's WebSocket endpoint for real-time streaming.

---

## Scope

**In scope:**
- `StreamLogs(ctx, name) (<-chan string, error)` method
- Connect to ws://localhost:8080/api/v1/containers/{name}/logs
- Return channel that receives log lines
- Close channel when WS closes

**Out of scope:**
- Reconnection logic (core handles this)
- Signal handling (CL-1.09)
- --format json (CL-1.10)

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `StreamLogs()` returns a channel
2. Channel receives log lines from WS
3. Channel closes when connection drops
4. Tests verify with mock WS server

---

## Implementation Guidelines

### Key files to modify or create
- `internal/client/client.go` — Add StreamLogs method

### Architecture constraints
- Use gorilla/websocket (already a dependency in core)
- Channel must close when WS closes — otherwise consumer hangs

### Suggested approach
1. Connect to WS endpoint
2. Read messages in goroutine
3. Send lines to channel
4. Close channel when connection closes

### Testing approach
- Test returns a channel
- Test receives lines from mock WS server
- Test channel closes on connection drop

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `StreamLogs()` method exists
- [ ] Returns channel that receives log lines
- [ ] Channel closes when WS closes
- [ ] Tests verify with mock WS
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented the method as described
- [ ] I have written tests
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: add StreamLogs method to API client`

---

## References

- `api/openapi.yaml` → WS /containers/{name}/logs

---

# [Phase 1] CLI exit code classification for container commands

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Small
> **Points:** 150
> **Estimated time:** 2-4 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli's exit codes are load-bearing for CI. This issue implements proper classification for all container commands.

---

## Problem

Without proper exit code classification, CI scripts can't distinguish between retryable (core unreachable) and permanent (core error) failures. This is the most critical CLI feature for automation.

---

## Scope

**In scope:**
- Exit 2: core unreachable (connection refused, timeout)
- Exit 3: core returned error (4xx, 5xx)
- Exit 1: bad arguments
- `classifyError(err) int` helper
- Apply to all container commands

**Out of scope:**
- Exit codes for non-container commands (Phase 2-4)
- Exit code 0 (success — already works)

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `classifyError()` helper exists
2. Connection errors → exit 2
3. HTTP errors (>= 400) → exit 3
4. Argument errors → exit 1
5. All container commands use this classification
6. Tests verify both sides of the 2/3 boundary

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/root.go` or new `exitcode.go` — classifyError helper
- `internal/cmd/containers.go` — Apply to all commands

### Architecture constraints
- The boundary between exit 2 and 3 depends on whether core returned ANY HTTP response
- Test both sides of the boundary explicitly

### Suggested approach
1. Implement `classifyError(err error) int`
2. Connection errors → 2, HTTP errors → 3, argument errors → 1
3. Apply to all container command RunE functions

### Testing approach
- Test connection refused → exit 2
- Test HTTP 500 → exit 3
- Test HTTP 404 → exit 3
- Test missing argument → exit 1

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `classifyError()` helper exists
- [ ] Connection errors → exit 2
- [ ] HTTP errors → exit 3
- [ ] Argument errors → exit 1
- [ ] All container commands use classification
- [ ] Tests verify both sides of 2/3 boundary
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented exit code classification as described
- [ ] I have written tests for both sides of the boundary
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: classify exit codes for container commands`

---

## References

- `ARCHITECTURE.md` → "Exit Code Convention"
- `ARCHITECTURE.md` → "Known Failure Modes" → Exit code misclassification

---

# [Phase 1] Container command --help text and usage documentation

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli's --help text is the first thing new users see. This issue adds comprehensive help to all container commands.

---

## Problem

Without good --help text, users must read the README to understand command options. --help should be self-sufficient.

---

## Scope

**In scope:**
- Update Use, Short, Long, Example fields for all container commands
- Document --name flag for start/stop
- Document --follow flag for logs
- Include example invocations

**Out of scope:**
- Man pages or external documentation
- Auto-generated help from code comments

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `stellaryard containers --help` shows all subcommands
2. Each subcommand has Short description, Long description, Example
3. All flags are documented
4. Help text is clear and complete

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/containers.go` — Update cobra command fields
- `internal/cmd/logs.go` — Update logs command fields

### Architecture constraints
- Cobra auto-generates help from Use, Short, Long, Example fields
- Include examples users can copy-paste

### Suggested approach
1. Update Use, Short, Long, Example for each command
2. Document all flags with descriptions
3. Include realistic examples

### Testing approach
- Run each command with --help
- Verify output is clear and complete

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] All container commands have help text
- [ ] Short description is one line
- [ ] Long description is detailed paragraph
- [ ] Example invocations are included
- [ ] All flags are documented
- [ ] Help is clear and complete
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have updated help text as described
- [ ] I have verified --help output for all commands
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`docs: add comprehensive help text to container commands`

---

## References

- `README.md` — existing documentation to align with

---

# [Phase 1] Container command integration tests

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Medium
> **Points:** 150
> **Estimated time:** 4-8 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli's container commands need comprehensive tests. This issue adds integration tests using mocked API client.

---

## Problem

Without integration tests, commands may have wrong output, wrong exit codes, or missing --format json support. Tests catch these before they reach users.

---

## Scope

**In scope:**
- Create `internal/cmd/containers_test.go`
- Mock API client responses
- Test each command: start, stop, status, logs
- Test --format json output
- Test --help output
- Test exit codes

**Out of scope:**
- Tests against real core instance (Phase 5)
- Formatter tests (separate issue)

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `containers_test.go` exists with comprehensive tests
2. Mock API client returns controlled responses
3. Table-driven tests for each command
4. Tests verify stdout, stderr, exit codes
5. Tests verify --format json produces valid JSON
6. All tests pass

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/containers_test.go` — New test file

### Architecture constraints
- Mock the API client, not the HTTP layer
- Test user-visible behavior, not implementation

### Suggested approach
1. Create mock API client for tests
2. Write table-driven tests for each command
3. Test stdout, stderr, exit codes

### Testing approach
- Table-driven tests with subtests
- Mock API client with controlled responses
- Verify JSON output is valid JSON

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `containers_test.go` exists
- [ ] Mock API client is implemented
- [ ] Tests for each command
- [ ] Tests verify output, exit codes
- [ ] Tests verify --format json
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented comprehensive tests as described
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`test: add integration tests for container commands`

---

## References

- `ARCHITECTURE.md` → "Testing expectations"

---

# [Phase 1] Formatter package tests

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli's formatter package is used by every command. This issue adds tests to prevent regressions.

---

## Problem

Without formatter tests, changes to output formatting can break all commands silently.

---

## Scope

**In scope:**
- Create `internal/formatter/formatter_test.go`
- Test table output: columns, alignment, empty state
- Test JSON output: valid JSON, correct structure
- Test error output: correct shape on stderr

**Out of scope:**
- Integration tests (separate issue)
- Performance testing

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. `formatter_test.go` exists
2. Tests cover table and JSON modes
3. Tests cover empty input scenarios
4. Tests cover error output
5. All tests pass

---

## Implementation Guidelines

### Key files to modify or create
- `internal/formatter/formatter_test.go` — New test file

### Architecture constraints
- Test empty arrays return `[]` in JSON, not `null`
- Formatter is used by every command — tests prevent regressions

### Suggested approach
1. Write tests for Writer.Write() with both modes
2. Write tests for Writer.WriteTable() with various inputs
3. Write tests for Writer.WriteError()

### Testing approach
- Table-driven tests
- Test with various data types
- Test empty input

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] `formatter_test.go` exists
- [ ] Tests cover table output
- [ ] Tests cover JSON output
- [ ] Tests cover error output
- [ ] Tests cover empty input
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented tests as described
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`test: add formatter package tests for table and JSON output`

---

## References

- `internal/formatter/formatter.go` — the code being tested

---

# [Phase 1] Root command error handling integration

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Small
> **Points:** 150
> **Estimated time:** 2-4 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli's root command must handle errors from subcommands and exit with the correct code.

---

## Problem

Without centralized error handling, each subcommand would need to call os.Exit directly, leading to inconsistent behavior.

---

## Scope

**In scope:**
- Root command's Execute() handles errors from subcommands
- Classify errors using classifyError() from CL-1.15
- Print error to stderr
- Exit with classified code
- Respect --format json for stderr output

**Out of scope:**
- Subcommand-specific error handling
- Error logging or monitoring

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. Execute() handles errors from subcommands
2. Errors are classified with classifyError()
3. Errors printed to stderr
4. Process exits with classified code
5. --format json produces JSON on stderr
6. Tests verify all exit codes

---

## Implementation Guidelines

### Key files to modify or create
- `internal/cmd/root.go` — Update Execute() function

### Architecture constraints
- Subcommands should NOT call os.Exit directly
- Root command handles all exit logic

### Suggested approach
1. Update Execute() to handle errors
2. Classify with classifyError()
3. Print to stderr (respect --format json)
4. Call os.Exit with classified code

### Testing approach
- Test exit code 0 on success
- Test exit code 1 on argument error
- Test exit code 2 on connection error
- Test exit code 3 on application error
- Test --format json on stderr

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] Execute() handles subcommand errors
- [ ] Errors classified correctly
- [ ] Errors printed to stderr
- [ ] --format json respected on stderr
- [ ] Tests verify all exit codes
- [ ] All tests pass
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have implemented error handling as described
- [ ] I have written tests for all exit codes
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`feat: wire root command error handling with exit code classification`

---

## References

- `ARCHITECTURE.md` → "Exit Code Convention"
- Depends on: CL-1.15 (classifyError helper)

---

# [Phase 1] README update for container commands

> **Repository:** stellaryard-cli
> **Phase:** Phase 1 — Container Commands
> **Complexity:** Trivial
> **Points:** 100
> **Estimated time:** 1-2 hours

---

## Repo Context

> ⚠️ **Before starting this issue, read the following files in the repo for full context:**
> - `AGENTS.md` — Non-negotiable rules for all contributors
> - `ARCHITECTURE_ESSENTIALS.md` — Quick reference for architecture decisions
> - `ARCHITECTURE.md` — Full architecture document (read if essentials don't cover your question)
> - `ROADMAP.md` — Check that your task is scoped and update it when done

stellaryard-cli's README is the primary documentation. This issue updates it with complete container command documentation.

---

## Problem

Without README documentation, users must guess command options and output formats.

---

## Scope

**In scope:**
- Document all container commands with examples
- Document exit code convention
- Document --format json usage
- Document --follow flag for logs
- Include CI/automation examples

**Out of scope:**
- API documentation (core's responsibility)
- Architecture documentation (already exists)

---

## What "Done" Looks Like

> A contributor should be able to read this section and know exactly when to stop.

1. README documents all container commands
2. Exit code table is included
3. --format json examples are included
4. CI integration example is included
5. README renders correctly on GitHub

---

## Implementation Guidelines

### Key files to modify or create
- `README.md` — Update with container command docs

### Architecture constraints
- Include real command examples users can copy-paste
- Document the --format json contract for CI users

### Suggested approach
1. Update V1 Commands section
2. Add container command docs with examples
3. Add exit code table
4. Add --format json examples
5. Add CI example

### Testing approach
- Verify README renders on GitHub
- Verify all documented commands work

---

## Acceptance Criteria

> PRs that don't meet ALL of these criteria will be sent back for revision.

- [ ] All container commands documented
- [ ] Exit code table included
- [ ] --format json examples included
- [ ] CI example included
- [ ] README renders correctly
- [ ] All documented commands work
- [ ] `ROADMAP.md` is updated
- [ ] No violations of `AGENTS.md` rules

---

## Checklist

> Tick each item in your PR description to confirm completion.

- [ ] I have read `AGENTS.md`, `ARCHITECTURE_ESSENTIALS.md`, and `ROADMAP.md`
- [ ] I understand the non-negotiable rules for this repo
- [ ] I have updated README as described
- [ ] I have verified all documented commands work
- [ ] All existing tests still pass
- [ ] I have updated `ROADMAP.md`
- [ ] I have included a clear commit message
- [ ] I have not violated any rules in `AGENTS.md`

---

## Example Commit Message

`docs: update README with container command documentation`

---

## References

- `README.md` — the file to update
- `ARCHITECTURE.md` → "Exit Code Convention"
