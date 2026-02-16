# Declarative Object Management

## The problem with ALTER

Every `Alter*()` call sends an `ALTER` command to the queue manager,
even when every specified attribute already matches the current state.
MQ updates `ALTDATE` and `ALTTIME` on every `ALTER`, regardless of
whether any values actually changed. This makes `ALTER` unsuitable for
declarative configuration management where idempotency matters -- running
the same configuration twice should not corrupt audit timestamps.

## The ensure pattern

The `Ensure*()` methods implement a declarative upsert pattern:

1. **DEFINE** the object when it does not exist.
2. **ALTER** only the attributes that differ from the current state.
3. **Do nothing** when all specified attributes already match,
   preserving `ALTDATE` and `ALTTIME`.

Each call returns an `EnsureResult` indicating what action was taken:

```go
// EnsureAction constants:
//   EnsureCreated   -- Object did not exist, was defined
//   EnsureUpdated   -- Object existed, attributes were altered
//   EnsureUnchanged -- Object existed, no changes needed

// EnsureResult struct:
//   Action  EnsureAction  -- the action taken
//   Changed []string      -- attribute names that triggered ALTER
```

## Basic usage

```go
ctx := context.Background()

// First call -- queue does not exist yet
result, err := session.EnsureQlocal(ctx, "APP.REQUEST.Q", map[string]any{
    "max_queue_depth": 50000,
    "description":     "Application request queue",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result.Action) // "created"

// Second call -- same attributes, nothing to change
result, err = session.EnsureQlocal(ctx, "APP.REQUEST.Q", map[string]any{
    "max_queue_depth": 50000,
    "description":     "Application request queue",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result.Action) // "unchanged"

// Third call -- description changed, only that attribute is altered
result, err = session.EnsureQlocal(ctx, "APP.REQUEST.Q", map[string]any{
    "max_queue_depth": 50000,
    "description":     "Updated request queue",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result.Action)  // "updated"
fmt.Println(result.Changed) // ["description"]
```

## Comparison logic

The ensure methods compare only the attributes the caller passes in
`requestParameters` against the current state returned by `DISPLAY`.
Attributes not specified by the caller are ignored.

Comparison is:

- **Case-insensitive** -- `"ENABLED"` matches `"enabled"`.
- **Type-normalizing** -- integer `5000` matches string `"5000"`.
- **Whitespace-trimming** -- `" YES "` matches `"YES"`.

An attribute present in `requestParameters` but absent from the
`DISPLAY` response is treated as changed and included in the `ALTER`.

## Selective ALTER

When an update is needed, only the changed attributes are sent in the
`ALTER` command. Attributes that already match are excluded from the
request. This minimizes the scope of each `ALTER` to the strict delta.

## Available methods

Each method targets a specific MQ object type with the correct
MQSC qualifier triple (DISPLAY / DEFINE / ALTER):

| Method | Object type | DISPLAY | DEFINE | ALTER |
| --- | --- | --- | --- | --- |
| `EnsureQmgr()` | Queue manager | `QMGR` | -- | `QMGR` |
| `EnsureQlocal()` | Local queue | `QUEUE` | `QLOCAL` | `QLOCAL` |
| `EnsureQremote()` | Remote queue | `QUEUE` | `QREMOTE` | `QREMOTE` |
| `EnsureQalias()` | Alias queue | `QUEUE` | `QALIAS` | `QALIAS` |
| `EnsureQmodel()` | Model queue | `QUEUE` | `QMODEL` | `QMODEL` |
| `EnsureChannel()` | Channel | `CHANNEL` | `CHANNEL` | `CHANNEL` |
| `EnsureAuthinfo()` | Auth info | `AUTHINFO` | `AUTHINFO` | `AUTHINFO` |
| `EnsureListener()` | Listener | `LISTENER` | `LISTENER` | `LISTENER` |
| `EnsureNamelist()` | Namelist | `NAMELIST` | `NAMELIST` | `NAMELIST` |
| `EnsureProcess()` | Process | `PROCESS` | `PROCESS` | `PROCESS` |
| `EnsureService()` | Service | `SERVICE` | `SERVICE` | `SERVICE` |
| `EnsureTopic()` | Topic | `TOPIC` | `TOPIC` | `TOPIC` |
| `EnsureSub()` | Subscription | `SUB` | `SUB` | `SUB` |
| `EnsureStgclass()` | Storage class | `STGCLASS` | `STGCLASS` | `STGCLASS` |
| `EnsureComminfo()` | Comm info | `COMMINFO` | `COMMINFO` | `COMMINFO` |
| `EnsureCfstruct()` | CF structure | `CFSTRUCT` | `CFSTRUCT` | `CFSTRUCT` |

Most methods share the same signature:

```go
func (session *Session) EnsureQlocal(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error)
```

All ensure methods take `context.Context` as their first parameter,
following standard Go conventions for I/O operations.

`responseParameters` is not exposed -- the ensure logic always requests
`["all"]` internally so it can compare the full current state.

### Queue manager (singleton)

`EnsureQmgr()` has no `name` parameter because the queue manager is a
singleton that always exists. It can only return `EnsureUpdated` or
`EnsureUnchanged` (never `EnsureCreated`):

```go
func (session *Session) EnsureQmgr(ctx context.Context, requestParameters map[string]any) (EnsureResult, error)
```

This makes it ideal for asserting queue manager-level settings such as
statistics, monitoring, events, and logging attributes without
corrupting `ALTDATE`/`ALTTIME` on every run.

## Attribute mapping

The ensure methods participate in the same
[mapping pipeline](mapping-pipeline.md) as all other command methods.
Pass `snake_case` attribute names in `requestParameters` and the
mapping layer translates them to MQSC names for the DISPLAY, DEFINE,
and ALTER commands automatically.

## Configuration management example

The ensure pattern is designed for programs that declare desired state:

```go
func configureQueueManager(ctx context.Context, session *mqrestadmin.Session) error {
    // Ensure queue manager settings
    result, err := session.EnsureQmgr(ctx, map[string]any{
        "queue_statistics":   "on",
        "channel_statistics": "on",
        "queue_monitoring":   "medium",
        "channel_monitoring": "medium",
    })
    if err != nil {
        return fmt.Errorf("ensure qmgr: %w", err)
    }
    fmt.Printf("Queue manager: %s\n", result.Action)

    // Ensure application queues
    queues := map[string]map[string]any{
        "APP.REQUEST.Q": {"max_queue_depth": 50000, "default_persistence": "yes"},
        "APP.REPLY.Q":   {"max_queue_depth": 10000, "default_persistence": "no"},
        "APP.DLQ":       {"max_queue_depth": 100000, "default_persistence": "yes"},
    }

    for name, attrs := range queues {
        result, err := session.EnsureQlocal(ctx, name, attrs)
        if err != nil {
            return fmt.Errorf("ensure %s: %w", name, err)
        }
        fmt.Printf("%s: %s\n", name, result.Action)
    }

    return nil
}
```

Running this function repeatedly produces no side effects when the configuration
is already correct. Only genuine changes trigger `ALTER` commands, keeping
`ALTDATE`/`ALTTIME` accurate.
