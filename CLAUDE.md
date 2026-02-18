# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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
- **golangci-lint**: `brew install golangci-lint`
- **govulncheck**: `go install golang.org/x/vuln/cmd/govulncheck@latest`
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
