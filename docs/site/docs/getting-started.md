# Getting Started

## Prerequisites

- **Go**: 1.25 or later
- **IBM MQ**: A running queue manager with the administrative REST API enabled

## Installation

```bash
go get github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin
```

## Creating a session

All interaction with IBM MQ goes through a `Session`. You need the
REST API base URL, queue manager name, and credentials:

```go
package main

import (
    "context"
    "time"

    "github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

func main() {
    session, err := mqrestadmin.NewSession(
        "https://localhost:9443/ibmmq/rest/v2",
        "QM1",
        mqrestadmin.LTPAAuth{Username: "mqadmin", Password: "mqadmin"},
        mqrestadmin.WithTimeout(30*time.Second),
        mqrestadmin.WithVerifyTLS(false), // for local development only
    )
    if err != nil {
        panic(err)
    }
    _ = session
}
```

## Running a command

Every MQSC command has a corresponding method on the session. Method names
follow the pattern `VerbQualifier` in PascalCase:

```go
ctx := context.Background()

// DISPLAY QUEUE — returns a slice of maps
queues, err := session.DisplayQlocal(ctx, "*")
if err != nil {
    panic(err)
}

for _, queue := range queues {
    fmt.Println(queue["queue_name"], queue["current_queue_depth"])
}
```

## Attribute mapping

By default, attribute names are translated between `snake_case` and MQSC
parameter names:

```go
// Request: snake_case → MQSC
err := session.DefineQlocal(ctx, "APP.REQUESTS", map[string]any{
    "max_queue_depth":      "50000",
    "default_persistence":  "persistent",
})

// Response: MQSC → snake_case
queues, _ := session.DisplayQlocal(ctx, "APP.REQUESTS")
fmt.Println(queues[0]["max_queue_depth"])      // "50000"
fmt.Println(queues[0]["default_persistence"])  // "persistent"
```

Disable mapping per-session with `WithMapAttributes(false)` to work with
raw MQSC parameter names.

See [mapping pipeline](mapping-pipeline.md) for a detailed explanation of
how mapping works.

## Idempotent operations

Ensure methods implement a declarative upsert pattern — define if the object
does not exist, alter only the attributes that differ, or no-op if everything
matches:

```go
result, err := session.EnsureQlocal(ctx, "APP.REQUESTS", map[string]any{
    "max_queue_depth": "50000",
})
if err != nil {
    panic(err)
}
fmt.Println(result.Action) // created, updated, or unchanged
```

See [ensure methods](ensure-methods.md) for details.

## Error handling

All error types are designed for use with `errors.As()`:

```go
var cmdErr *mqrestadmin.CommandError
if errors.As(err, &cmdErr) {
    fmt.Println("MQSC command failed:", cmdErr.Payload)
}
```

## What's next

- [Architecture](architecture.md) — how the components fit together
- [Mapping Pipeline](mapping-pipeline.md) — how attribute translation works
- [Ensure Methods](ensure-methods.md) — idempotent object management
- [Sync Methods](sync-methods.md) — synchronous polling operations
