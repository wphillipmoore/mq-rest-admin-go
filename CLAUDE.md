# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Auto-memory policy

**Do NOT use MEMORY.md.** Claude Code's auto-memory feature stores behavioral
rules outside of version control, making them invisible to code review,
inconsistent across repos, and unreliable across sessions. All behavioral rules,
conventions, and workflow instructions belong in managed, version-controlled
documentation (CLAUDE.md, AGENTS.md, skills, or docs/).

If you identify a pattern, convention, or rule worth preserving:

1. **Stop.** Do not write to MEMORY.md.
2. **Discuss with the user** what you want to capture and why.
3. **Together, decide** the correct managed location (CLAUDE.md, a skill file,
   standards docs, or a new issue to track the gap).

This policy exists because MEMORY.md is per-directory and per-machine — it
creates divergent agent behavior across the multi-repo environment this project
operates in. Consistency requires all guidance to live in shared, reviewable
documentation.

## Documentation Strategy

This repository uses two complementary approaches for AI agent guidance:

- **AGENTS.md**: Generic AI agent instructions using include directives to force documentation indexing. Contains canonical standards references, shared skills loading, and user override support.
- **CLAUDE.md** (this file): Claude Code-specific guidance with prescriptive commands, architecture details, and development workflows optimized for `/init`.

<!-- include: docs/standards-and-conventions.md -->
<!-- include: docs/repository-standards.md -->

## Project Overview

This is a Go port of `pymqrest`, providing a Go wrapper for the IBM MQ administrative REST API. The project provides typed Go functions for every MQSC command exposed by the `runCommandJSON` REST endpoint, with automatic attribute name translation between Go `snake_case` and native MQSC parameter names.

**Project name**: mq-rest-admin-go

**Status**: Pre-alpha (initial setup)

**Module path**: `github.com/wphillipmoore/mq-rest-admin-go`

**Package name**: `mqrestadmin`

**Canonical Standards**: This repository follows standards at <https://github.com/wphillipmoore/standards-and-conventions> (local path: `../standards-and-conventions` if available)

## Development Commands

### Environment Setup

- **Go**: 1.25+ (CI tests 1.25 and 1.26; go.mod declares 1.26)
- **golangci-lint**: `brew install golangci-lint` (not in `tools.go` per project recommendation)
- **Dev tools** (pinned in `tools.go`): `go install golang.org/x/vuln/cmd/govulncheck && go install github.com/vladopajic/go-test-coverage/v2 && go install github.com/fzipp/gocyclo/cmd/gocyclo`
- **Git hooks**: `git config core.hooksPath scripts/git-hooks` (required before committing)

### Build

```bash
go build ./...          # Compile all packages
go vet ./...            # Static analysis
```

### Validation

```bash
scripts/dev/validate_local.sh   # Canonical validation (runs all checks below)
go vet ./...                    # Static analysis
golangci-lint run ./...         # Lint checks
gocyclo -over 15 ./mqrestadmin/ # Cyclomatic complexity gate
go test -race -count=1 ./...   # Unit tests with race detection
govulncheck ./...               # Vulnerability scanning
```

### Testing

```bash
go test ./...                                   # Unit tests
go test -race -count=1 ./...                    # Unit tests with race detection
go test -run TestFunctionName ./mqrestadmin/...  # Run a single test
go test -race -count=1 -tags=integration ./...  # Integration tests (needs MQ env)
go test -coverprofile=coverage.out ./... && go-test-coverage --config .testcoverage.yml  # Coverage gate
go tool cover -html=coverage.out                # View coverage in browser
```

- **Framework**: stdlib `testing` package
- **Coverage**: Target 100% line coverage (file, package, and total), enforced in
  CI via `go-test-coverage` with `.testcoverage.yml`. Structurally untestable
  lines (e.g., `json.Marshal` on `map[string]any`, embedded JSON parse errors) are
  annotated with `// coverage-ignore -- <reason>` on the **preceding line** (the
  `{` line) and excluded from measurement.
- **Integration tests**: Require `MQ_REST_ADMIN_GO_RUN_INTEGRATION=1` and a running
  MQ container. CI uses the `wphillipmoore/mq-rest-admin-dev-environment` action.

## Architecture

Direct port of `pymqrest`'s architecture, adapted to Go idioms. Uses the
`mq-rest-admin` Java port as a secondary reference.

### Technology Stack

- **HTTP client**: `net/http` (stdlib, zero runtime dependencies)
- **JSON library**: `encoding/json` (stdlib)
- **TLS/mTLS**: `crypto/tls` (stdlib)
- **Mapping data**: `go:embed` (stdlib)
- **Runtime dependencies**: 0 (stdlib only)

### API Surface

- **Style**: Method-per-command mirroring pymqrest (`DisplayQueue()`, `DefineQlocal()`, etc.)
- **Parameters/results**: `map[string]any` for MQ attributes (dynamic attribute sets)
- **Fixed-schema types**: Go structs for `TransportResponse`, `SyncConfig`, `SyncResult`, `EnsureResult`, `MappingIssue`, credential types
- **Credentials**: Interface with unexported method (closed set: `BasicAuth`, `LTPAAuth`, `CertificateAuth`)
- **Session**: Single `Session` struct with functional options constructor
- **Context**: All I/O methods accept `context.Context` as first parameter

### Transport Layer

- `Transport` interface — enables mock-based testing
- `HTTPTransport` implementation using `net/http`
- Supports TLS/mTLS via `crypto/tls`, timeouts via `time.Duration`

### Attribute Mapping

- Direct port of pymqrest's 3-layer pipeline: key map, value map, key-value map
- Two directions: request and response
- Strict/permissive modes with `MappingIssue` tracking
- Mapping data embedded via `go:embed` from JSON file
- Override mechanism with merge/replace modes

### Error Types

Go typed errors (not exception hierarchy), used with `errors.As()`:

```text
TransportError   — network/connection failures
ResponseError    — malformed JSON, unexpected structure
AuthError        — authentication/authorization failures
CommandError     — MQSC command returned error codes
TimeoutError     — polling timeout exceeded
MappingError     — attribute translation failures (strict mode)
```

### Package Structure

Single flat package under `mqrestadmin/`. Key files:

- `session.go` — `Session` struct, `NewSession`, functional options, core `mqscCommand` dispatch
- `session_commands.go` — All MQSC command methods (DISPLAY, DEFINE, ALTER, DELETE, START, STOP, etc.)
- `session_ensure.go` / `session_sync.go` — Idempotent ensure and synchronous polling methods
- `auth.go` — `Credentials` interface (closed set via unexported method), `BasicAuth`, `LTPAAuth`, `CertificateAuth`
- `transport.go` — `Transport` interface, `HTTPTransport`
- `mapping.go` + `mapping_data.go` + `mapping_data.json` — 3-layer attribute translation pipeline
- `errors.go` — All typed error types

### Testing Patterns

Tests use a `mockTransport` (in `mock_test.go`) that records calls and returns
pre-configured responses. Helper constructors:

- `newTestSession(transport)` — mapping disabled, for simple assertions
- `newTestSessionWithMapping(transport)` — mapping enabled
- `newTestSessionWithClock(transport, clock)` — mock clock for sync/polling tests
- `generateSelfSignedCert(t)` — creates test TLS certificates (in `testhelpers_test.go`)

## Branching and PR Workflow

- **Protected branches**: `main`, `develop`, `release/*` — no direct commits (enforced by pre-commit hook)
- **Branch naming**: `feature/*`, `bugfix/*`, or `hotfix/*` only
- **Feature/bugfix PRs** target `develop` with squash merge: `gh pr merge --auto --squash --delete-branch`
- **Release PRs** target `main` with regular merge: `gh pr merge --auto --merge --delete-branch`
- **Pre-flight**: Always check branch with `git status -sb` before modifying files. If on `develop`, create a `feature/*` branch first.
- **Go version gate**: PRs to `develop` must have a `go.mod` version different from `main`

## Commit and PR Scripts

**NEVER use raw `git commit`** — always use `scripts/dev/commit.sh`.
**NEVER use raw `gh pr create`** — always use `scripts/dev/submit-pr.sh`.

### Committing

```bash
scripts/dev/commit.sh --type feat --scope mapping --message "add reverse lookup" --agent claude
scripts/dev/commit.sh --type fix --message "correct off-by-one in polling" --agent claude
scripts/dev/commit.sh --type docs --message "update README examples" --body "Expanded usage section" --agent claude
```

- `--type` (required): `feat|fix|docs|style|refactor|test|chore|ci|build`
- `--message` (required): commit description
- `--agent` (required): `claude` or `codex` — resolves the correct `Co-Authored-By` identity
- `--scope` (optional): conventional commit scope
- `--body` (optional): detailed commit body

### Submitting PRs

```bash
scripts/dev/submit-pr.sh --issue 42 --summary "Add reverse lookup for attribute mapping"
scripts/dev/submit-pr.sh --issue 42 --linkage Ref --summary "Update docs" --docs-only
scripts/dev/submit-pr.sh --issue 42 --summary "Fix polling bug" --notes "Tested with MQ 9.4"
```

- `--issue` (required): GitHub issue number (just the number)
- `--summary` (required): one-line PR summary
- `--linkage` (optional, default: `Fixes`): `Fixes|Closes|Resolves|Ref`
- `--title` (optional): PR title (default: most recent commit subject)
- `--notes` (optional): additional notes
- `--docs-only` (optional): applies docs-only testing exception
- `--dry-run` (optional): print generated PR without executing

## Key References

**Canonical Standards**: <https://github.com/wphillipmoore/standards-and-conventions>

- Local path (preferred): `../standards-and-conventions`
- Load all skills from: `<standards-repo-path>/skills/**/SKILL.md`

**Reference implementation**: `../pymqrest` (Python version)

**Java port reference**: `../mq-rest-admin` (Java version)

**External Documentation**:

- IBM MQ 9.4 administrative REST API
- MQSC command reference

**User Overrides**: `~/AGENTS.md` (optional, applied if present and readable)
