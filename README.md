# mq-rest-admin-go

Go wrapper for the IBM MQ administrative REST API.

`mqrestadmin` provides typed Go methods for every MQSC command exposed
by the IBM MQ 9.4 `runCommandJSON` REST endpoint. Attribute names are
automatically translated between Go `snake_case` and native MQSC
parameter names, so you work with idiomatic Go identifiers throughout.

## Table of Contents

- [Installation](#installation)
- [Quick start](#quick-start)
- [API overview](#api-overview)
- [Documentation](#documentation)
- [Development](#development)
- [License](#license)

## Installation

```bash
go get github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin
```

Requires Go 1.25+. Zero external runtime dependencies.

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

func main() {
    session, err := mqrestadmin.NewSession(
        "https://localhost:9443/ibmmq/rest/v2",
        "QM1",
        mqrestadmin.LTPAAuth{Username: "mqadmin", Password: "mqadmin"},
        mqrestadmin.WithTimeout(30*time.Second),
        mqrestadmin.WithVerifyTLS(false),
    )
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    // Query the queue manager
    qmgr, err := session.DisplayQmgr(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Println(qmgr["queue_manager_name"])

    // List all local queues
    queues, err := session.DisplayQlocal(ctx, "*")
    if err != nil {
        panic(err)
    }
    for _, q := range queues {
        fmt.Println(q["queue_name"], q["current_queue_depth"])
    }

    // Idempotent object management
    result, err := session.EnsureQlocal(ctx, "APP.REQUESTS", map[string]any{
        "max_queue_depth": "50000",
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(result.Action) // created, updated, or unchanged
}
```

## API overview

### Session

`NewSession` creates a session that manages authentication, connection
settings, and attribute mapping. All command methods are called on the
session.

```go
session, err := mqrestadmin.NewSession(
    "https://host:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.LTPAAuth{Username: "user", Password: "pass"},
    mqrestadmin.WithMapAttributes(true),   // snake_case <-> MQSC (default)
    mqrestadmin.WithMappingStrict(true),    // error on unknown attributes (default)
    mqrestadmin.WithVerifyTLS(true),        // TLS verification (default)
    mqrestadmin.WithTimeout(30*time.Second), // HTTP timeout (default)
)
```

### Commands

Over 140 methods cover the MQSC command set:

| Verb | Methods | Returns | Example |
| --- | --- | --- | --- |
| `Display*` | 44 | `([]map[string]any, error)` | `session.DisplayQlocal(ctx, "*")` |
| `Define*` | 19 | `error` | `session.DefineQlocal(ctx, "Q1", params)` |
| `Alter*` | 17 | `error` | `session.AlterQlocal(ctx, "Q1", params)` |
| `Delete*` | 16 | `error` | `session.DeleteQlocal(ctx, "Q1")` |
| Other | 48 | `error` | `StartChannel`, `StopListener`, `ClearQlocal`, ... |

All methods accept `context.Context` as the first parameter. Display
commands accept optional `CommandOption` functions for request/response
parameter filtering.

### Ensure methods

Idempotent `Ensure*` methods implement a declarative upsert pattern
for 15 object types (queues, channels, topics, listeners, and more):

- **Define** when the object does not exist
- **Alter** only the attributes that differ
- **No-op** when all specified attributes already match

Returns an `EnsureResult` whose `Action` is `EnsureCreated`,
`EnsureUpdated`, or `EnsureUnchanged`.

### Attribute mapping

When `WithMapAttributes(true)` (the default), attribute names and
values are translated automatically:

| Direction | From | To | Example |
| --- | --- | --- | --- |
| Request | `max_queue_depth` | `MAXDEPTH` | snake_case to MQSC |
| Response | `MAXDEPTH` | `max_queue_depth` | MQSC to snake_case |

Disable per-session (`WithMapAttributes(false)`) or per-call for raw
MQSC parameter access.

### Authentication

Three credential types are supported:

- `CertificateAuth` — mutual TLS client certificates
- `LTPAAuth` — LTPA token login (automatic at session creation)
- `BasicAuth` — HTTP Basic authentication

## Documentation

Full documentation will be published at a later date.

## Development

### Prerequisites

- **Go** 1.25+
- **golangci-lint**: `brew install golangci-lint`
- **Dev tools** (govulncheck, go-test-coverage, gocyclo) are pinned in
  `tools.go` and can be installed from the module:

```bash
go install golang.org/x/vuln/cmd/govulncheck
go install github.com/vladopajic/go-test-coverage/v2
go install github.com/fzipp/gocyclo/cmd/gocyclo
```

Ensure `$(go env GOPATH)/bin` is on your `PATH`.

### Validation

```bash
scripts/dev/validate_local.sh   # runs all checks
go vet ./...                    # static analysis
golangci-lint run ./...         # lint
go test -race -count=1 ./...   # unit tests
govulncheck ./...               # vulnerability scan
```

## License

GPL-3.0-or-later. See `LICENSE`.
