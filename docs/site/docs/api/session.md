# Session

## Overview

The main entry point for interacting with an IBM MQ queue manager's
administrative REST API. A `Session` encapsulates connection details,
authentication, attribute mapping configuration, and diagnostic state. It
provides ~144 command methods covering all MQSC verbs and qualifiers, plus 15
idempotent ensure methods and 9 synchronous sync methods.

The Go implementation uses functional options for construction, following
standard Go library conventions.

## Creating a session

Use `NewSession` with functional options:

```go
import "github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"

session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
)
```

`NewSession` validates configuration, initializes the attribute mapper, and
(for LTPA credentials) performs the login request immediately. Errors in
configuration are returned at construction time.

## Signature

```go
func NewSession(
    restBaseURL string,
    qmgrName    string,
    credentials Credentials,
    opts        ...Option,
) (*Session, error)
```

| Parameter | Description |
| --- | --- |
| `restBaseURL` | Base URL of the MQ REST API (e.g. `https://host:9443/ibmmq/rest/v2`) |
| `qmgrName` | Target queue manager name |
| `credentials` | Authentication credentials (`BasicAuth`, `LTPAAuth`, or `CertificateAuth`) |
| `opts` | Zero or more functional options |

## Functional options

| Option | Type | Description |
| --- | --- | --- |
| `WithTransport(Transport)` | `Transport` | Custom transport implementation (default: `HTTPTransport`) |
| `WithGatewayQmgr(string)` | `string` | Gateway queue manager for remote routing |
| `WithVerifyTLS(bool)` | `bool` | Verify server TLS certificates (default: `true`) |
| `WithTimeout(time.Duration)` | `time.Duration` | HTTP request timeout (default: 30s) |
| `WithMapAttributes(bool)` | `bool` | Enable/disable attribute mapping (default: `true`) |
| `WithMappingStrict(bool)` | `bool` | Strict or permissive mapping mode (default: `true`) |
| `WithCSRFToken(*string)` | `*string` | Custom CSRF token value; `nil` omits the header |
| `WithMappingOverrides(map[string]any, MappingOverrideMode)` | `map[string]any` | Custom mapping overrides with merge or replace mode |

### Minimal example

```go
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
)
if err != nil {
    log.Fatal(err)
}
```

### Full example

```go
session, err := mqrestadmin.NewSession(
    "https://mq-server.example.com:9443/ibmmq/rest/v2",
    "QM2",
    mqrestadmin.BasicAuth{Username: "mqadmin", Password: "mqadmin"},
    mqrestadmin.WithGatewayQmgr("QM1"),
    mqrestadmin.WithMapAttributes(true),
    mqrestadmin.WithMappingStrict(false),
    mqrestadmin.WithVerifyTLS(true),
    mqrestadmin.WithTimeout(30*time.Second),
    mqrestadmin.WithMappingOverrides(overrides, mqrestadmin.MappingOverrideMerge),
)
if err != nil {
    log.Fatal(err)
}
```

## Command methods

The session provides ~144 command methods, one for each MQSC verb + qualifier
combination. See [Commands](commands.md) for the full list.

```go
ctx := context.Background()

// DISPLAY commands return a slice of maps
queues, err := session.DisplayQueue(ctx, "APP.*")

// Queue manager singletons return a single map
qmgr, err := session.DisplayQmgr(ctx)

// Non-DISPLAY commands return only an error
err = session.DefineQlocal(ctx, "MY.QUEUE",
    mqrestadmin.WithRequestParameters(map[string]any{"max_queue_depth": 50000}))
err = session.DeleteQueue(ctx, "MY.QUEUE")
```

## Ensure methods

The session provides 15 ensure methods for declarative object management. Each
method implements an idempotent upsert: DEFINE if the object does not exist,
ALTER only the attributes that differ, or no-op if already correct.

```go
result, err := session.EnsureQlocal(ctx, "MY.QUEUE",
    map[string]any{"max_queue_depth": 50000})
// result.Action is EnsureCreated, EnsureUpdated, or EnsureUnchanged
```

See [Ensure](ensure.md) for details.

## Diagnostic fields

The session retains the most recent request and response for inspection. These
are exported struct fields, updated after every command:

```go
session.DisplayQueue(ctx, "MY.QUEUE")

fmt.Println(session.LastCommandPayload)    // the JSON sent to MQ
fmt.Println(session.LastResponsePayload)   // the parsed JSON response
fmt.Println(session.LastHTTPStatus)        // HTTP status code
fmt.Println(session.LastResponseText)      // raw response body
```

### Exported fields and accessors

| Field / Method | Type | Description |
| --- | --- | --- |
| `LastHTTPStatus` | `int` | HTTP status code from last command |
| `LastResponseText` | `string` | Raw response body from last command |
| `LastResponsePayload` | `map[string]any` | Parsed response from last command |
| `LastCommandPayload` | `map[string]any` | Command payload sent for last command |
| `QmgrName()` | `string` | Queue manager name |
| `GatewayQmgr()` | `string` | Gateway queue manager (or empty string) |
