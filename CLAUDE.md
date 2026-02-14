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

**Package name**: `mqrest`

**Canonical Standards**: This repository follows standards at https://github.com/wphillipmoore/standards-and-conventions (local path: `../standards-and-conventions` if available)

## Development Commands

### Environment Setup

- **Go**: 1.22+ (install via `brew install go` or https://go.dev/dl/)
- **golangci-lint**: `brew install golangci-lint`
- **govulncheck**: `go install golang.org/x/vuln/cmd/govulncheck@latest`

### Build

```bash
go build ./...          # Compile all packages
go vet ./...            # Static analysis
```

### Validation

```bash
go vet ./...                    # Static analysis
golangci-lint run ./...         # Lint checks
go test -race -count=1 ./...   # Unit tests with race detection
govulncheck ./...               # Vulnerability scanning
```

### Testing

```bash
go test ./...                                   # Unit tests
go test -race -count=1 ./...                    # Unit tests with race detection
go test -race -count=1 -tags=integration ./...  # Integration tests
go test -coverprofile=coverage.out ./...        # Coverage report
go tool cover -html=coverage.out                # View coverage in browser
```

- **Framework**: stdlib `testing` package
- **Coverage**: Target 100% line coverage (matching Java port standard)

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

```
TransportError   — network/connection failures
ResponseError    — malformed JSON, unexpected structure
AuthError        — authentication/authorization failures
CommandError     — MQSC command returned error codes
TimeoutError     — polling timeout exceeded
MappingError     — attribute translation failures (strict mode)
```

### Package Structure

```
mqrest/                        # Single flat package
    session.go                 # Session struct, NewSession, functional options, core dispatch
    session_commands.go        # MQSC command methods (display, define, alter, delete, etc.)
    session_ensure.go          # Idempotent ensure methods
    session_sync.go            # Synchronous polling methods
    auth.go                    # Credentials interface, BasicAuth, LTPAAuth, CertificateAuth
    transport.go               # Transport interface, HTTPTransport, TransportResponse
    mapping.go                 # AttributeMapper, 3-layer pipeline
    mapping_data.go            # go:embed for mapping-data.json
    mapping_data.json          # Mapping definitions (shared with Java port)
    errors.go                  # All error types
    ensure.go                  # EnsureResult, EnsureAction
    sync.go                    # SyncConfig, SyncResult, SyncOperation
    doc.go                     # Package documentation
```

## Key References

**Canonical Standards**: https://github.com/wphillipmoore/standards-and-conventions
- Local path (preferred): `../standards-and-conventions`
- Load all skills from: `<standards-repo-path>/skills/**/SKILL.md`

**Reference implementation**: `../pymqrest` (Python version)

**Java port reference**: `../mq-rest-admin` (Java version)

**External Documentation**:
- IBM MQ 9.4 administrative REST API
- MQSC command reference

**User Overrides**: `~/AGENTS.md` (optional, applied if present and readable)
