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
queues, err := session.DisplayQueue(ctx, "*")
if err != nil {
    panic(err)
}

for _, queue := range queues {
    fmt.Println(queue["queue_name"], queue["current_queue_depth"])
}
```

```go
// DISPLAY QMGR — returns a single map or nil
qmgr, err := session.DisplayQmgr(ctx)
if err != nil {
    panic(err)
}
if qmgr != nil {
    fmt.Println(qmgr["queue_manager_name"])
}
```

## Attribute mapping

By default, the session maps between developer-friendly `snake_case` names
and MQSC parameter names. This applies to both request and response attributes:

```go
// With mapping enabled (default)
queues, _ := session.DisplayQueue(ctx, "MY.QUEUE",
    map[string]any{"response_parameters": []string{"current_queue_depth", "max_queue_depth"}},
)
// Returns: [{"queue_name": "MY.QUEUE", "current_queue_depth": 0, "max_queue_depth": 5000}]

// With mapping disabled
queues, _ = session.DisplayQueue(ctx, "MY.QUEUE",
    map[string]any{"response_parameters": []string{"CURDEPTH", "MAXDEPTH"}},
)
// Returns: [{"queue": "MY.QUEUE", "curdepth": 0, "maxdepth": 5000}]
```

Mapping can be disabled at the session level:

```go
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2", "QM1",
    mqrestadmin.LTPAAuth{Username: "mqadmin", Password: "mqadmin"},
    mqrestadmin.WithMapAttributes(false),
)
```

See [mapping pipeline](mapping-pipeline.md) for a detailed explanation of how
mapping works.

## Strict vs lenient mapping

By default, mapping runs in lenient mode. Unknown attribute names or values
pass through unchanged. In strict mode, unknown attributes return an error:

```go
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2", "QM1",
    mqrestadmin.LTPAAuth{Username: "mqadmin", Password: "mqadmin"},
    mqrestadmin.WithStrictMapping(true),
)
```

## Custom mapping overrides

Sites with existing naming conventions can override individual entries in the
built-in mapping tables without replacing them entirely. Pass override data
when creating the session:

```go
overrideData := map[string]any{
    "qualifiers": map[string]any{
        "queue": map[string]any{
            "response_key_map": map[string]any{
                "CURDEPTH": "queue_depth",      // override built-in mapping
                "MAXDEPTH": "queue_max_depth",   // override built-in mapping
            },
        },
    },
}

session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2", "QM1",
    mqrestadmin.LTPAAuth{Username: "mqadmin", Password: "mqadmin"},
    mqrestadmin.WithMappingOverride(overrideData, mqrestadmin.MergeOverride),
)

queues, _ := session.DisplayQueue(ctx, "MY.QUEUE")
// Returns: [{"queue_depth": 0, "queue_max_depth": 5000, ...}]
```

Overrides are **sparse** — you only specify the entries you want to change. All
other mappings in the qualifier continue to work as normal. In the example above,
only `CURDEPTH` and `MAXDEPTH` are remapped; every other queue attribute keeps
its default `snake_case` name.

Overrides support all five sub-maps per qualifier: `request_key_map`,
`request_value_map`, `request_key_value_map`, `response_key_map`, and
`response_value_map`. See [mapping pipeline](mapping-pipeline.md) for details
on how each sub-map is used.

## Gateway queue manager

The MQ REST API is available on all supported IBM MQ platforms (Linux, AIX,
Windows, z/OS, and IBM i). mqrestadmin is developed and tested against the
**Linux** implementation only.

In enterprise environments, a **gateway queue manager** can route MQSC
commands to remote queue managers via MQ channels — the same mechanism used
by `runmqsc -w` and the MQ Console.

To use a gateway, pass `WithGatewayQmgr` when creating the session. The
base URL and queue manager name specify the **target** (remote) queue manager,
while `WithGatewayQmgr` names the **local** queue manager whose REST API
routes the command:

```go
// Route commands to QM2 through QM1's REST API
session, err := mqrestadmin.NewSession(
    "https://qm1-host:9443/ibmmq/rest/v2",
    "QM2",                                     // target queue manager
    mqrestadmin.LTPAAuth{Username: "mqadmin", Password: "mqadmin"},
    mqrestadmin.WithGatewayQmgr("QM1"),        // local gateway queue manager
    mqrestadmin.WithVerifyTLS(false),
)

qmgr, _ := session.DisplayQmgr(ctx)
// Returns QM2's queue manager attributes, routed through QM1
```

Prerequisites:

- The gateway queue manager must have a running REST API.
- MQ channels must be configured between the gateway and target queue managers.
- A QM alias (QREMOTE with empty RNAME) must map the target QM name to the
  correct transmission queue on the gateway.

## Error handling

`DISPLAY` commands return an empty slice when no objects match. Queue manager
display methods return `nil` when no match is found. Non-display commands
return a `*CommandError` on failure:

```go
// Empty slice — no error
result, err := session.DisplayQueue(ctx, "NONEXISTENT.*")
// result == []

// Define returns error on failure
var cmdErr *mqrestadmin.CommandError
err = session.DefineQlocal(ctx, "MY.QUEUE", nil)
if errors.As(err, &cmdErr) {
    fmt.Println(cmdErr.Error())
    fmt.Println("HTTP status:", cmdErr.StatusCode)
    fmt.Println(cmdErr.Payload) // full MQ response payload
}
```

## Diagnostic state

The session retains the most recent request and response for inspection:

```go
session.DisplayQueue(ctx, "MY.QUEUE")

fmt.Println(session.LastCommandPayload)    // the JSON sent to MQ
fmt.Println(session.LastResponsePayload)   // the parsed JSON response
fmt.Println(session.LastHTTPStatus)        // HTTP status code
fmt.Println(session.LastResponseText)      // raw response body
```
