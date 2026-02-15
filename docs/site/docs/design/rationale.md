# Design Rationale

--8<-- "design/rationale.md"

## Go-specific design choices

### Zero external dependencies

`mqrestadmin` has zero external runtime dependencies. The entire package is
built on the Go standard library:

- `net/http` for HTTP communication
- `encoding/json` for JSON serialization
- `crypto/tls` for TLS and mutual TLS
- `embed` for mapping data resources

This keeps the dependency tree empty and avoids version conflicts in
downstream projects.

### Single flat package

The entire public API lives in one package: `mqrestadmin`. There are no
sub-packages to import. This follows the Go convention of small, focused
packages and keeps the import path simple:

```go
import "github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
```

### Functional options

Session configuration uses the functional options pattern rather than a
configuration struct or builder:

```go
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.LTPAAuth{Username: "mqadmin", Password: "mqadmin"},
    mqrestadmin.WithVerifyTLS(false),
    mqrestadmin.WithTimeout(30 * time.Second),
)
```

This provides a clean API with sensible defaults, optional parameters
without nil checks, and backward-compatible extensibility.

### context.Context integration

All I/O methods accept `context.Context` as their first parameter:

```go
results, err := session.DisplayQueue(ctx, "MY.QUEUE")
```

This enables cancellation, deadline propagation, and tracing integration
without inventing a custom mechanism.

### Errors as values

Go errors follow the standard `error` interface with typed error structs
for classification:

```go
results, err := session.DisplayQueue(ctx, "MY.QUEUE")
if err != nil {
    var cmdErr *mqrestadmin.CommandError
    if errors.As(err, &cmdErr) {
        fmt.Printf("status: %d, payload: %v\n", cmdErr.StatusCode, cmdErr.Payload)
    }
    return err
}
```

There is no exception hierarchy. Errors are values that can be inspected
with `errors.As()` and `errors.Is()`.

### Method naming conventions

Command methods use `PascalCase` following Go export conventions. The
pattern is `<Verb><Qualifier>`:

| MQSC command | Go method |
| --- | --- |
| `DISPLAY QUEUE` | `DisplayQueue()` |
| `DEFINE QLOCAL` | `DefineQlocal()` |
| `DELETE CHANNEL` | `DeleteChannel()` |
| `ALTER QMGR` | `AlterQmgr()` |

### Return shapes

**DISPLAY commands** return `([]map[string]any, error)`. An empty slice
means no objects matched -- this is not an error. The caller can range
over the result without nil checks.

**Queue manager singletons** (`DisplayQmgr`, `DisplayQmstatus`, etc.)
return `(map[string]any, error)`. A nil map means no result was returned.

**Non-DISPLAY commands** (`Define`, `Delete`, `Alter`, etc.) return
`error` only. A nil error means success.
