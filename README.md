<p align="center">
  <img src="https://img.shields.io/badge/stellar-yard%20cli-blue?style=for-the-badge&logo=stellar&logoColor=white" alt="StellarYard CLI"/>
</p>

<h1 align="center">stellaryard-cli</h1>

<p align="center">
  A scriptable terminal client for StellarYard. Manage containers, accounts, contracts, and ledger data from the command line.
</p>

<p align="center">
  <a href="https://github.com/StellarYard/stellaryard-cli/blob/main/LICENSE"><img src="https://img.shields.io/github/license/StellarYard/stellaryard-cli?style=flat-square" alt="License"/></a>
  <a href="https://github.com/StellarYard/stellaryard-cli/actions"><img src="https://img.shields.io/github/actions/workflow/status/StellarYard/stellaryard-cli/ci.yml?style=flat-square&label=CI" alt="CI"/></a>
  <a href="https://github.com/StellarYard/stellaryard-cli/issues"><img src="https://img.shields.io/github/issues/StellarYard/stellaryard-cli?style=flat-square" alt="Issues"/></a>
  <a href="https://www.drips.network/wave/stellar"><img src="https://img.shields.io/badge/Drips%20Wave-Stellar-7c3aed?style=flat-square" alt="Drips Wave"/></a>
</p>

---

## What is StellarYard CLI?

StellarYard CLI is a terminal client that talks to [stellaryard-core](../stellaryard-core)'s API. It's designed for **automation** — predictable exit codes, machine-parseable `--format json` output, and scriptability for CI pipelines.

## Quick Start

```bash
# Install
git clone https://github.com/StellarYard/stellaryard-cli.git
cd stellaryard-cli
go build -o stellaryard ./cmd/stellaryard

# Start core first (separate terminal)
cd ../stellaryard-core && docker-compose up -d && go run cmd/server/main.go

# Use the CLI
./stellaryard containers status
./stellaryard accounts create --label "test-account"
./stellaryard contracts deploy ./my-contract.wasm
```

## Commands

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

## Exit Codes

Designed for CI/automation:

| Code | Meaning | Example |
|------|---------|---------|
| `0` | Success | Command completed |
| `1` | Argument error | Bad flags, missing args |
| `2` | Core unreachable | Connection refused |
| `3` | Core app error | Deploy failed |

## Architecture

```
┌───────────────────────────┐
│ stellaryard-cli             │
│  ┌───────────┐             │
│  │ cobra      │             │
│  │ commands   │             │
│  └─────┬─────┘             │
│  ┌─────▼─────┐             │
│  │ Formatter  │             │
│  │ (table/json)│            │
│  └─────┬─────┘             │
│  ┌─────▼─────┐             │
│  │ core API   │             │
│  │ client     │             │
│  │ (generated)│             │
│  └─────┬─────┘             │
└────────┼───────────────────┘
         │ REST + WS
         ▼
  stellaryard-core
```

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.22 |
| CLI Framework | [cobra](https://github.com/spf13/cobra) |
| Output | table (tabwriter) / JSON |
| API Client | Generated from core's OpenAPI spec |

## Roadmap

| Phase | Scope | Status |
|-------|-------|--------|
| 0 — Foundation | Scaffold, command tree, API client generation | Not started |
| 1 — Containers | start/stop/status, logs --follow | Not started |
| 2 — Accounts | create/list/show | Not started |
| 3 — Ledger | snapshot, tx list | Not started |
| 4 — Contracts | deploy, invoke | Not started |
| 5 — Hardening | Error handling, integration tests, shell completion | Not started |

Full roadmap: [`ROADMAP.md`](./ROADMAP.md)

## Contributing

We welcome contributions! See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

- Check [open issues](https://github.com/StellarYard/stellaryard-cli/issues) for `ready` tasks
- Issues labeled `good-first-issue` are ideal for first-time contributors
- Every PR must update `ROADMAP.md`

## Maintainers

| Name | GitHub | Contact |
|------|--------|---------|
| Adejumo-2 | [@Adejumo-2](https://github.com/Adejumo-2) | [Telegram](https://t.me/Adejumo-2) |

## Community

- [GitHub Discussions](https://github.com/StellarYard/stellaryard-cli/discussions)
- [Drips Wave — Stellar](https://www.drips.network/wave/stellar)

## License

[Apache 2.0](./LICENSE)

---

<p align="center">
  Built for the <a href="https://www.drips.network/wave/stellar">Stellar Wave Program</a>
</p>
