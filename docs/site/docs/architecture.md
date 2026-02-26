# Architecture

## Component overview

--8<-- "architecture/component-overview.md"

In the Go implementation, the core components map to these types:

- **`Session`**: The main entry point. A single struct that owns connection
  details, authentication, mapping configuration, diagnostic state, and all
  ~144 command methods plus 16 ensure methods and 9 sync methods. Created
  via `NewSession` with functional options.
- **Command methods**: Exported methods on `Session` (e.g. `DisplayQueue()`,
  `DefineQlocal()`, `DeleteChannel()`). Each method is a thin wrapper that
  calls the internal `mqscCommand()` dispatcher with the correct verb and
  qualifier.
- **`attributeMapper`**: Internal type that handles bidirectional attribute
  translation using mapping data loaded from an embedded JSON resource. See
  the [mapping pipeline](mapping-pipeline.md) for details.
- **Error types**: Typed error structs (`TransportError`, `CommandError`, etc.)
  used with `errors.As()`. All errors are values, not an exception hierarchy.

## Request lifecycle

--8<-- "architecture/request-lifecycle.md"

In Go, the command dispatcher is the internal `mqscCommand()` method on
`Session`. Every exported command method (e.g. `DisplayQueue()`,
`DefineQlocal()`) delegates to it with the appropriate verb and qualifier.

The session retains diagnostic state from the most recent command for
inspection:

```go
session.DisplayQueue(ctx, "MY.QUEUE")

session.LastCommandPayload    // the JSON sent to MQ
session.LastResponsePayload   // the parsed JSON response
session.LastHTTPStatus        // HTTP status code
session.LastResponseText      // raw response body
```

## Transport abstraction

--8<-- "architecture/transport-abstraction.md"

In Go, the transport is defined by the `Transport` interface:

```go
type Transport interface {
    PostJSON(ctx context.Context, url string, payload map[string]any,
        headers map[string]string, timeout time.Duration, verifyTLS bool,
    ) (*TransportResponse, error)
}
```

The default `HTTPTransport` uses `net/http`. Custom implementations can be
injected via `WithTransport()` for testing or specialized HTTP handling.

For testing, inject a mock transport:

```go
type mockTransport struct{}

func (m *mockTransport) PostJSON(ctx context.Context, url string,
    payload map[string]any, headers map[string]string,
    timeout time.Duration, verifyTLS bool,
) (*TransportResponse, error) {
    return &mqrestadmin.TransportResponse{
        StatusCode: 200,
        Body:       responseJSON,
        Headers:    map[string]string{},
    }, nil
}

session, _ := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2", "QM1",
    mqrestadmin.LTPAAuth{Username: "admin", Password: "pass"},
    mqrestadmin.WithTransport(&mockTransport{}),
)
```

This makes the entire command pipeline testable without an MQ server.

## Single-endpoint design

--8<-- "architecture/single-endpoint-design.md"

In Go, this means every command method on `Session` ultimately calls the
same `PostJSON()` method on the transport with the same URL pattern. The only
variation is the JSON payload content.

## Gateway routing

--8<-- "architecture/gateway-routing.md"

In Go, configure gateway routing via a functional option:

```go
session, err := mqrestadmin.NewSession(
    "https://qm1-host:9443/ibmmq/rest/v2",
    "QM2",                                      // target (remote) queue manager
    mqrestadmin.LTPAAuth{Username: "mqadmin", Password: "mqadmin"},
    mqrestadmin.WithGatewayQmgr("QM1"),         // local gateway queue manager
)
```

## Context support

All I/O methods accept `context.Context` as their first parameter, enabling
cancellation, timeouts, and tracing integration.

## Zero dependencies

The package uses only the Go standard library:

- `net/http` for HTTP
- `encoding/json` for JSON
- `crypto/tls` for TLS/mTLS
- `embed` for mapping data

## Ensure pipeline

See [ensure](api/ensure.md) for details on the idempotent
create-or-update pipeline.

## Sync pipeline

See [sync](api/sync.md) for details on the synchronous
polling pipeline.
