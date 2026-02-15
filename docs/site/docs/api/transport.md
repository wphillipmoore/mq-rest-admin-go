# Transport

## Overview

The transport layer abstracts HTTP communication from the session logic. The
session builds `runCommandJSON` payloads and delegates HTTP delivery to a
transport implementation. This separation enables testing the entire command
pipeline without an MQ server by injecting a mock transport.

## Transport interface

The `Transport` interface defines a single method for posting JSON payloads:

```go
type Transport interface {
    PostJSON(
        ctx       context.Context,
        url       string,
        payload   map[string]any,
        headers   map[string]string,
        timeout   time.Duration,
        verifyTLS bool,
    ) (*TransportResponse, error)
}
```

| Parameter | Type | Description |
| --- | --- | --- |
| `ctx` | `context.Context` | Request context for cancellation and deadlines |
| `url` | `string` | Fully-qualified endpoint URL |
| `payload` | `map[string]any` | The `runCommandJSON` request body |
| `headers` | `map[string]string` | Authentication, CSRF token, and optional gateway headers |
| `timeout` | `time.Duration` | Per-request timeout duration |
| `verifyTLS` | `bool` | Whether to verify server certificates |

Returns `*TransportResponse` on success or a `*TransportError` on network
failures.

## TransportResponse

A struct containing the HTTP response data:

```go
type TransportResponse struct {
    StatusCode int               // HTTP status code
    Body       string            // Response body (empty string if no body)
    Headers    map[string]string // Response headers
}
```

| Field | Type | Description |
| --- | --- | --- |
| `StatusCode` | `int` | HTTP status code |
| `Body` | `string` | Response body text |
| `Headers` | `map[string]string` | Response headers (first value per key) |

## HTTPTransport

The default `Transport` implementation using `net/http` from the Go standard
library (zero external dependencies):

```go
type HTTPTransport struct {
    TLSConfig *tls.Config  // optional TLS configuration for mTLS or custom CAs
}
```

```go
// Default -- verifies TLS certificates
transport := &mqrestadmin.HTTPTransport{}

// Custom TLS configuration for mTLS
transport := &mqrestadmin.HTTPTransport{
    TLSConfig: &tls.Config{
        Certificates: []tls.Certificate{cert},
    },
}
```

`HTTPTransport` handles:

- HTTPS connections with configurable `tls.Config`
- Automatic TLS certificate verification (or disabled via `verifyTLS=false`)
- Request timeouts via `time.Duration`
- JSON serialization/deserialization with `encoding/json`
- Custom HTTP headers
- Context-aware requests with cancellation support

When `CertificateAuth` credentials are provided and no custom transport is set,
`NewSession` automatically creates an `HTTPTransport` with the client
certificate loaded into `TLSConfig`.

## Injecting a custom transport

Use `WithTransport` to provide a custom `Transport` implementation for testing
or specialized HTTP handling:

```go
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
    mqrestadmin.WithTransport(mockTransport),
)
```

## Mock transport for testing

Because `Transport` is an interface with a single method, it is straightforward
to create mock implementations for unit tests:

```go
type mockTransport struct {
    response *mqrestadmin.TransportResponse
    err      error
}

func (m *mockTransport) PostJSON(
    ctx context.Context, url string, payload map[string]any,
    headers map[string]string, timeout time.Duration, verifyTLS bool,
) (*mqrestadmin.TransportResponse, error) {
    return m.response, m.err
}

// Use in tests
mock := &mockTransport{
    response: &mqrestadmin.TransportResponse{
        StatusCode: 200,
        Body:       `{"commandResponse":[]}`,
        Headers:    map[string]string{},
    },
}

session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
    mqrestadmin.WithTransport(mock),
)
```

This pattern is used extensively in the library's own test suite to verify
command payload construction, response parsing, and error handling without
network access.
