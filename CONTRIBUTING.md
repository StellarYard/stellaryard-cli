# Contributing to StellarYard CLI

Thank you for your interest in contributing to StellarYard CLI! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful, constructive, and professional. We're building tools for the Stellar ecosystem together.

## Getting Started

### Prerequisites

- Go 1.22+
- stellaryard-core running locally (for integration testing)

### Setup

```bash
# Clone the repo
git clone https://github.com/StellarYard/stellaryard-cli.git
cd stellaryard-cli

# Install dependencies
go mod tidy

# Build
go build -o stellaryard ./cmd/stellaryard

# Run
./stellaryard --help
```

### Running Tests

```bash
go test ./...
go vet ./...
```

## How to Contribute

### Finding Issues

1. Check the [open issues](https://github.com/StellarYard/stellaryard-cli/issues) for tasks labeled `ready`
2. Issues labeled `good-first-issue` are ideal for first-time contributors
3. Read the issue description carefully — each issue includes acceptance criteria and implementation guidelines

### Submitting Changes

1. **Fork** the repository
2. **Create a branch** from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```
3. **Make your changes** following the coding standards below
4. **Write or update tests** for your changes
5. **Update ROADMAP.md** as part of your PR — this is required
6. **Commit** with a descriptive message:
   ```bash
   git commit -m "feat: add containers start command"
   ```
7. **Push** your branch:
   ```bash
   git push origin feat/your-feature-name
   ```
8. **Open a Pull Request** against `main`

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): description

[optional body]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring without behavior change
- `chore`: Maintenance tasks

Examples:
```
feat(cmd): add containers start command with --format json
fix(formatter): handle empty list output consistently
docs: update README with installation instructions
test(client): add mock HTTP server for API client tests
```

### Pull Request Guidelines

- **One logical change per PR** — don't bundle unrelated changes
- **Include tests** — PRs without test coverage will be sent back
- **Update ROADMAP.md** — every PR must update the roadmap to reflect what was done
- **Keep PRs small** — ideally under 500 lines of diff
- **Describe what and why** — not just what changed, but why

## Coding Standards

### Go Style

- Follow [Effective Go](https://go.dev/doc/effective_go) conventions
- Use `gofmt` and `go vet` — no manual formatting
- Error messages should be lowercase, no punctuation
- Use `cobra` conventions for command definitions
- Exit codes: `0` (success), `1` (argument error), `2` (core unreachable), `3` (app error)

### Exit Code Convention

| Code | Meaning | Example |
|------|---------|---------|
| 0 | Success | Command completed normally |
| 1 | Argument error | Invalid flag, missing required arg |
| 2 | Core unreachable | Cannot connect to stellaryard-core |
| 3 | Application error | Core returned an error |

### Output Formatting

- All list/show commands must support `--format table` (default) and `--format json`
- Use the shared `formatter` package — never reimplement output formatting per command
- Table output uses `tabwriter` for alignment
- JSON output uses `json.Encoder` with standard library only

### Project Structure

```
cmd/stellaryard/     — Entry point
internal/cmd/        — Cobra command definitions
internal/client/     — API client (generated from OpenAPI spec)
internal/formatter/  — Output formatting (table/JSON)
```

### Architecture Rules (Non-Negotiable)

1. **API client is generated from core's openapi.yaml** — never hand-write client methods
2. **Core's OpenAPI spec is the single source of truth** — never fake backend behavior
3. **No cross-repo modifications** — note dependencies in PR descriptions
4. **ROADMAP.md must be updated in every PR**

## Reporting Issues

- Use GitHub Issues for bug reports and feature requests
- Include the command you ran and the output you got
- Include your Go version and OS

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
