# Architecture

## Component overview

`mqrestadmin` is structured as a single flat Go package. The main
components are:

- **`Session`**: The main entry point. A single struct that owns connection
  details, authentication, mapping configuration, diagnostic state, and all
  ~144 command methods plus 15 ensure methods and 9 sync methods. Created
  via `NewSession` with functional options.
- **Command methods**: Exported methods on `Session` (e.g. `DisplayQlocal()`,
  `DefineQlocal()`, `DeleteChannel()`). Each method is a thin wrapper that
  calls the internal `mqscCommand()` dispatcher with the correct verb and
  qualifier.
- **`attributeMapper`**: Internal type that handles bidirectional attribute
  translation using mapping data loaded from an embedded JSON resource. See
  the [mapping pipeline](mapping-pipeline.md) for details.
- **Error types**: Typed error structs (`TransportError`, `CommandError`, etc.)
  used with `errors.As()`. All errors are values, not an exception hierarchy.

## Request lifecycle

When you call a command method like `session.DisplayQlocal(ctx, "*")`:

1. The method delegates to the internal `mqscCommand()` dispatcher with the
   verb (`DISPLAY`), qualifier (`QLOCAL`), and name.
2. If mapping is enabled, request attributes are translated from `snake_case`
   to MQSC parameter names.
3. The dispatcher builds a `runCommandJSON` payload and sends it via the
   `Transport` interface.
4. The response JSON is parsed and validated.
5. If mapping is enabled, response attributes are translated from MQSC back
   to `snake_case`.
6. The session retains diagnostic state from the most recent command.

```go
session.DisplayQlocal(ctx, "MY.QUEUE")

session.LastCommandPayload    // the JSON sent to MQ
session.LastResponsePayload   // the parsed JSON response
session.LastHTTPStatus        // HTTP status code
session.LastResponseText      // raw response body
```

## Transport abstraction

The `Transport` interface abstracts HTTP communication:

```go
type Transport interface {
    PostJSON(ctx context.Context, url string, payload map[string]any,
        headers map[string]string, timeout time.Duration, verifyTLS bool,
    ) (*TransportResponse, error)
}
```

The default `HTTPTransport` uses `net/http`. Custom implementations can be
injected via `WithTransport()` for testing or specialized HTTP handling.

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

See [ensure methods](ensure-methods.md) for details on the idempotent
create-or-update pipeline.

## Sync pipeline

See [sync methods](sync-methods.md) for details on the synchronous
polling pipeline.
