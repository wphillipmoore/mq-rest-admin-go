# Go Port Plan: IBM MQ REST Admin Client

## Table of Contents

- [Overview](#overview)
- [Source Analysis](#source-analysis)
- [Go Port Architecture](#go-port-architecture)
- [Implementation Phases](#implementation-phases)
- [CI/CD Gates (per standards)](#cicd-gates-per-standards)
- [Key Differences from Python/Java](#key-differences-from-pythonjava)
- [Resolved Decisions](#resolved-decisions)
- [Open Questions](#open-questions)

## Overview

Port the pymqrest Python library to Go, using the mq-rest-admin Java port as a
secondary reference. The library wraps the IBM MQ 9.4 administrative REST API,
providing typed Go functions for every MQSC command exposed by the
`runCommandJSON` REST endpoint.

## Source Analysis

### pymqrest (Python original)

- 10 source modules, ~4,500 lines
- `MQRESTSession` class composed from 3 mixins: commands, ensure, sync
- 115+ MQSC command methods, 15 ensure methods, 9 sync methods
- 3-layer attribute mapping pipeline (key-value map, key map, value map)
- `requests` library for HTTP; pluggable transport protocol
- 3 auth types: BasicAuth, LTPAAuth, CertificateAuth
- Sealed exception hierarchy (5 concrete types)

### mq-rest-admin (Java port)

- 32 source files, ~4,057 lines
- Single `MqRestSession` class (no mixin decomposition)
- Builder pattern for session construction
- Gson for JSON; JDK `HttpClient` for HTTP
- Sealed interfaces for credentials and exceptions
- Records for DTOs (TransportResponse, SyncResult, EnsureResult, etc.)
- 100% test coverage enforced

### Key decisions from the Java port

1. Kept all command methods as explicit methods (no code generation)
2. Used `Map<String, Object>` for dynamic MQSC attributes
3. Loaded mapping data from a JSON resource file
4. Single runtime dependency (Gson)

## Go Port Architecture

### Naming

- **Repository**: `mq-rest-admin-go`
- **Module path**: `github.com/wphillipmoore/mq-rest-admin-go`
- **Package name**: `mqrest` (short, lowercase, single word per Go conventions)

Follows the family naming pattern established by `mq-rest-admin` (Java), with
a `-go` suffix to distinguish the language.

### Project Structure

```text
mq-rest-admin-go/
├── docs/
│   ├── decisions/
│   └── standards-and-conventions.md
├── mqrestadmin/                    # Main package
│   ├── session.go             # Session type, builder/options, core dispatch
│   ├── session_commands.go    # MQSC command methods (display, define, alter, delete, etc.)
│   ├── session_ensure.go      # Idempotent ensure methods
│   ├── session_sync.go        # Synchronous polling methods
│   ├── auth.go                # Credential types (BasicAuth, LTPAAuth, CertificateAuth)
│   ├── transport.go           # Transport interface + default net/http implementation
│   ├── mapping.go             # Attribute mapper (3-layer pipeline)
│   ├── mapping_data.go        # Embedded mapping data (via go:embed)
│   ├── mapping_data.json      # Mapping definitions (reused from Java port)
│   ├── errors.go              # Error types
│   ├── ensure.go              # EnsureResult, EnsureAction
│   ├── sync.go                # SyncConfig, SyncResult, SyncOperation
│   └── doc.go                 # Package documentation
├── mqrestadmin/mqresttest/         # Optional test-helper subpackage
│   └── mock_transport.go      # Mock transport for consumer tests
├── examples/
│   └── basic/
│       └── main.go
├── go.mod
├── go.sum
├── CLAUDE.md
├── AGENTS.md
├── README.md
└── .github/
    └── workflows/
```

### Design Decisions

#### 1. Single package, multiple files

Go favours flat package structures. All public types live in `mqrestadmin`. Split
across files by concern (session, commands, ensure, sync, mapping, auth,
errors, transport) mirroring the Python module structure without introducing
sub-packages.

#### 2. Functional options for session construction

Go does not have builders or default parameter values. Use the functional
options pattern:

```go
session, err := mqrestadmin.NewSession(
    "https://host:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.WithBasicAuth("user", "pass"),
    mqrestadmin.WithTimeout(30 * time.Second),
    mqrestadmin.WithGatewayQmgr("GATEWAY"),
    mqrestadmin.WithMappingStrict(true),
)
```

#### 3. Credential types as an interface

```go
type Credentials interface {
    applyAuth(req *http.Request, session *Session) error
    // unexported method restricts implementations to this package
}

type BasicAuth struct {
    Username string
    Password string
}

type LTPAAuth struct {
    Username string
    Password string
}

type CertificateAuth struct {
    CertPath string
    KeyPath  string  // optional; empty if combined PEM
}
```

The unexported `applyAuth` method makes this a closed set (Go's equivalent of
sealed types), preventing external implementations while keeping the types
exported and constructable.

#### 4. Error types (not a hierarchy)

Go uses error values, not exception hierarchies. Define sentinel-style typed
errors:

```go
type TransportError struct {
    URL string
    Err error
}

type ResponseError struct {
    ResponseText string
}

type AuthError struct {
    URL        string
    StatusCode int
}

type CommandError struct {
    Payload    map[string]any
    StatusCode int
}

type TimeoutError struct {
    Name            string
    Operation       SyncOperation
    ElapsedSeconds  float64
}

type MappingError struct {
    Issues []MappingIssue
}
```

All implement `error` interface. Callers use `errors.As()` for type-specific
handling:

```go
var cmdErr *mqrestadmin.CommandError
if errors.As(err, &cmdErr) {
    fmt.Println(cmdErr.Payload)
}
```

#### 5. Return types

- **DISPLAY commands**: `([]map[string]any, error)` for multi-row,
  `(map[string]any, error)` for single-row (e.g., `DisplayQmgr`)
- **DEFINE/ALTER/DELETE/other void commands**: `error`
- **Ensure methods**: `(EnsureResult, error)`
- **Sync methods**: `(SyncResult, error)`

#### 6. Attribute mapping via embedded JSON

Use `go:embed` to embed `mapping_data.json` (reused from the Java port with
no format changes):

```go
//go:embed mapping_data.json
var mappingDataJSON []byte
```

Parse at init time with `encoding/json`. The 3-layer pipeline (key-value map,
key map, value map) ports directly.

#### 7. Transport interface

```go
type Transport interface {
    PostJSON(ctx context.Context, url string, payload map[string]any,
        headers map[string]string, timeout time.Duration, verifyTLS bool,
    ) (*TransportResponse, error)
}

type TransportResponse struct {
    StatusCode int
    Body       string
    Headers    map[string]string
}
```

Default implementation uses `net/http` (stdlib). No external HTTP dependency.

#### 8. Context support

Go convention: accept `context.Context` as first parameter on methods that do
I/O. All command methods should accept a context:

```go
func (session *Session) DisplayQlocal(ctx context.Context, name string,
    opts ...CommandOption) ([]map[string]any, error)
```

This is a Go-idiomatic addition not present in the Python or Java ports.

#### 9. Zero external runtime dependencies

- HTTP: `net/http` (stdlib)
- JSON: `encoding/json` (stdlib)
- TLS/mTLS: `crypto/tls` (stdlib)
- Embedding: `embed` (stdlib)

Test dependencies only: `testify` (or stdlib `testing` only, per preference).

### Type Mapping: Python/Java to Go

| Concept | Python | Java | Go |
| ------------- | ---------------------- | ------------------------ | --------------------------------- |
| Session | `MQRESTSession` class | `MqRestSession` class | `Session` struct |
| Construction | `__init__` kwargs | Builder pattern | Functional options |
| Auth types | Union type | Sealed interface | Interface with unexported method |
| Exceptions | Exception hierarchy | Sealed unchecked hierarchy | Typed error structs |
| DTOs | `@dataclass` | Records | Structs |
| Enums | `enum.Enum` | Java enums | Typed constants (`iota`) |
| Nullable | `None` / `Optional` | `@Nullable` | Zero values / pointers |
| Dynamic attrs | `dict[str, object]` | `Map<String, Object>` | `map[string]any` |
| Transport | Protocol (duck typing) | Interface | Interface |
| JSON | `json` stdlib | Gson | `encoding/json` stdlib |
| HTTP | `requests` | JDK `HttpClient` | `net/http` stdlib |

### Enum Types

```go
type EnsureAction int
const (
    EnsureCreated   EnsureAction = iota
    EnsureUpdated
    EnsureUnchanged
)

type SyncOperation int
const (
    SyncStarted  SyncOperation = iota
    SyncStopped
    SyncRestarted
)

type MappingDirection int
const (
    MappingRequest  MappingDirection = iota
    MappingResponse
)

type MappingReason int
const (
    MappingUnknownKey MappingReason = iota
    MappingUnknownValue
    MappingUnknownQualifier
)

type MappingOverrideMode int
const (
    MappingOverrideMerge   MappingOverrideMode = iota
    MappingOverrideReplace
)
```

## Implementation Phases

### Phase 1: Foundation

1. **Project scaffolding** -- `go.mod`, directory structure, CLAUDE.md,
   AGENTS.md, `docs/standards-and-conventions.md`
2. **Error types** (`errors.go`) -- all 6 error structs with `Error()` methods
3. **Transport interface and default implementation** (`transport.go`) --
   `Transport` interface, `HTTPTransport` struct using `net/http`
4. **Auth types** (`auth.go`) -- `Credentials` interface, `BasicAuth`,
   `LTPAAuth`, `CertificateAuth`
5. **Session core** (`session.go`) -- `Session` struct, `NewSession` with
   functional options, `mqscCommand` base dispatcher, LTPA login flow,
   header construction, response parsing, error detection

### Phase 2: Attribute Mapping

1. **Mapping data** (`mapping_data.json`, `mapping_data.go`) -- embed and
   parse the JSON mapping definitions
2. **Attribute mapper** (`mapping.go`) -- 3-layer pipeline, strict/permissive
   modes, request and response mapping, `MappingIssue` collection

### Phase 3: Command Methods

1. **DISPLAY commands** (`session_commands.go`) -- 44 display methods
2. **DEFINE commands** -- 19 define methods
3. **ALTER commands** -- 17 alter methods
4. **DELETE commands** -- 16 delete methods
5. **Other commands** -- START, STOP, CLEAR, RESET, SUSPEND, RESUME, etc.

### Phase 4: Ensure and Sync

1. **Ensure methods** (`session_ensure.go`) -- 15 idempotent ensure methods
   returning `EnsureResult`
2. **Sync methods** (`session_sync.go`) -- 9 polling methods returning
   `SyncResult`, `SyncConfig`

### Phase 5: Quality and Documentation

1. **Unit tests** -- mock transport, test all command methods, mapping,
   ensure logic, sync polling, auth, error handling
2. **Integration tests** -- against live MQ (gated by env var)
3. **Examples** -- basic usage, ensure, sync, auth types
4. **CI/CD** -- GitHub Actions: `go vet`, `golangci-lint`, `govulncheck`,
   `go test -race -cover`, coverage enforcement
5. **Documentation** -- README, package doc comments, ADRs for key decisions

## CI/CD Gates (per standards)

### Hard Gates (merge-blocking)

- `test-and-validate`: `go vet`, `golangci-lint`, `go test -race -count=1 ./...`
- `integration-tests`: gated by `MQREST_RUN_INTEGRATION=1`
- `dependency-audit`: `govulncheck ./...`

### Coverage

- Enforce no coverage decline (track via CI artifact or Codecov)
- Target: 100% line coverage (matching Java port standard)

## Key Differences from Python/Java

| Aspect | Go approach |
| -------------------- | ----------------------------------------------------------------- |
| No inheritance | Composition via embedding not needed; single Session struct |
| No generics needed | `map[string]any` for dynamic attributes |
| Context propagation | `context.Context` on all I/O methods |
| Error handling | Return `error`, caller uses `errors.As` |
| No default params | Functional options pattern |
| No sealed types | Interface with unexported method |
| Concurrent safety | Document thread-safety; `sync.Mutex` if Session holds mutable state |
| Testability | Interface-based transport; no mocking framework needed |
| Zero deps | Everything from stdlib |

## Resolved Decisions

1. **Repo name**: `mq-rest-admin-go`

## Open Questions

1. **Test framework**: stdlib `testing` only, or allow `testify`?
2. **Minimum Go version**: Go 1.22+ (for `go:embed`, generics if needed)?
3. **Code generation**: hand-code all 130+ methods (like Java), or generate
   from mapping data?
4. **Context on all methods**: confirmed as the Go-idiomatic approach?
