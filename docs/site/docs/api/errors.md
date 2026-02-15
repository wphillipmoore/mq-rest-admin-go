# Errors

## Overview

All error types are concrete structs that implement the `error` interface. They
are designed for use with Go's `errors.As()` for targeted error handling. Unlike
Java's sealed exception hierarchy, Go errors are matched structurally using
type assertions or `errors.As()`.

```text
*TransportError   -- Network/connection failures
*ResponseError    -- Malformed JSON, unexpected structure
*AuthError        -- Authentication/authorization failures
*CommandError     -- MQSC command returned error codes
*TimeoutError     -- Polling timeout exceeded
*MappingError     -- Attribute mapping failures (separate concern)
```

## TransportError

Returned when the HTTP request fails at the network level -- connection refused,
DNS resolution failure, TLS handshake error, etc. Wraps the underlying error
via `Unwrap()`.

```go
type TransportError struct {
    URL string  // The URL that was being accessed
    Err error   // The underlying error
}
```

| Field | Type | Description |
| --- | --- | --- |
| `URL` | `string` | The URL that was being accessed |
| `Err` | `error` | The underlying network error |

`TransportError` implements `Unwrap() error`, so you can inspect the root cause
with `errors.Unwrap()` or match nested errors with `errors.Is()`.

```go
queues, err := session.DisplayQueue(ctx, "*")
if err != nil {
    var transportErr *mqrestadmin.TransportError
    if errors.As(err, &transportErr) {
        fmt.Println("Cannot reach:", transportErr.URL)
        fmt.Println("Cause:", transportErr.Err)
    }
}
```

## ResponseError

Returned when the HTTP request succeeds but the response cannot be parsed --
invalid JSON, missing expected fields, unexpected response structure.

```go
type ResponseError struct {
    ResponseText string  // Raw response body
    StatusCode   int     // HTTP status code
}
```

| Field | Type | Description |
| --- | --- | --- |
| `ResponseText` | `string` | Raw response body |
| `StatusCode` | `int` | HTTP status code |

```go
queues, err := session.DisplayQueue(ctx, "*")
if err != nil {
    var responseErr *mqrestadmin.ResponseError
    if errors.As(err, &responseErr) {
        fmt.Printf("Bad response (HTTP %d): %s\n",
            responseErr.StatusCode, responseErr.ResponseText)
    }
}
```

## AuthError

Returned when authentication or authorization fails -- invalid credentials,
expired tokens, insufficient permissions (HTTP 401/403), or LTPA login failure.

```go
type AuthError struct {
    URL        string  // The URL that was being accessed
    StatusCode int     // HTTP status code
}
```

| Field | Type | Description |
| --- | --- | --- |
| `URL` | `string` | The URL that was being accessed |
| `StatusCode` | `int` | HTTP status code (401 or 403) |

```go
qmgr, err := session.DisplayQmgr(ctx)
if err != nil {
    var authErr *mqrestadmin.AuthError
    if errors.As(err, &authErr) {
        fmt.Printf("Auth failed: HTTP %d for %s\n",
            authErr.StatusCode, authErr.URL)
    }
}
```

## CommandError

Returned when the MQSC command returns a non-zero completion or reason code.
This is the most commonly encountered error -- it indicates the command was
delivered to MQ but the queue manager rejected it.

```go
type CommandError struct {
    Payload    map[string]any  // Full response payload
    StatusCode int             // HTTP status code
}
```

| Field | Type | Description |
| --- | --- | --- |
| `Payload` | `map[string]any` | Full response payload including completion and reason codes |
| `StatusCode` | `int` | HTTP status code |

```go
err := session.DefineQlocal(ctx, "MY.QUEUE",
    mqrestadmin.WithRequestParameters(map[string]any{}))
if err != nil {
    var cmdErr *mqrestadmin.CommandError
    if errors.As(err, &cmdErr) {
        fmt.Println("Command failed:", cmdErr.Error())
        fmt.Println("HTTP status:", cmdErr.StatusCode)
        fmt.Println("Response:", cmdErr.Payload)
    }
}
```

## TimeoutError

Returned when a synchronous polling operation exceeds its configured timeout
duration. Only produced by [sync methods](sync.md).

```go
type TimeoutError struct {
    Name           string         // Resource name being polled
    Operation      SyncOperation  // Operation being performed
    ElapsedSeconds float64        // Elapsed time in seconds
}
```

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Resource name being polled |
| `Operation` | `SyncOperation` | The sync operation (`SyncStarted`, `SyncStopped`) |
| `ElapsedSeconds` | `float64` | Elapsed time in seconds |

```go
result, err := session.StartChannelSync(ctx, "TO.PARTNER", mqrestadmin.SyncConfig{
    Timeout: 5 * time.Second,
})
if err != nil {
    var timeoutErr *mqrestadmin.TimeoutError
    if errors.As(err, &timeoutErr) {
        fmt.Printf("Timed out %s %s after %.1fs\n",
            timeoutErr.Operation, timeoutErr.Name, timeoutErr.ElapsedSeconds)
    }
}
```

## MappingError

Returned when attribute mapping fails in strict mode. Separate from the
transport/command error types. Contains the list of `MappingIssue` instances
that caused the failure.

```go
type MappingError struct {
    Issues []MappingIssue
}
```

| Field | Type | Description |
| --- | --- | --- |
| `Issues` | `[]MappingIssue` | List of attribute translation failures |

See [Mapping](mapping.md) for details on `MappingIssue`, `MappingDirection`,
and `MappingReason`.

```go
queues, err := session.DisplayQueue(ctx, "*",
    mqrestadmin.WithRequestParameters(map[string]any{
        "invalid_attr": "value",
    }),
)
if err != nil {
    var mappingErr *mqrestadmin.MappingError
    if errors.As(err, &mappingErr) {
        for _, issue := range mappingErr.Issues {
            fmt.Printf("%s %s: %s\n",
                issue.Direction, issue.Reason, issue.AttributeName)
        }
    }
}
```

## Error handling patterns

Use `errors.As()` for targeted recovery, or check the error interface for
broad handling:

```go
err := session.DefineQlocal(ctx, "MY.QUEUE",
    mqrestadmin.WithRequestParameters(map[string]any{"max_queue_depth": 50000}))
if err != nil {
    var cmdErr *mqrestadmin.CommandError
    var authErr *mqrestadmin.AuthError
    var transportErr *mqrestadmin.TransportError
    var mappingErr *mqrestadmin.MappingError

    switch {
    case errors.As(err, &cmdErr):
        // MQSC command failed -- check reason code in payload
        fmt.Println("Command failed:", cmdErr.Error())
    case errors.As(err, &authErr):
        // Credentials rejected
        fmt.Printf("Not authorized: HTTP %d\n", authErr.StatusCode)
    case errors.As(err, &transportErr):
        // Network error
        fmt.Println("Connection failed to", transportErr.URL)
    case errors.As(err, &mappingErr):
        // Attribute name/value not recognized
        fmt.Println("Mapping failed:", mappingErr.Error())
    default:
        // Catch-all
        fmt.Println("Unexpected error:", err)
    }
}
```
