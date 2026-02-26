# Sync

## The problem with fire-and-forget

All MQSC `START` and `STOP` commands are fire-and-forget -- they return
immediately without waiting for the object to reach its target state.
In practice, tooling that provisions infrastructure needs to wait until
a channel is `RUNNING` or a listener is `STOPPED` before proceeding to
the next step. Writing polling loops by hand is error-prone and
clutters business logic with retry mechanics.

## The sync pattern

The `*Sync` and `Restart*` methods wrap the fire-and-forget commands
with a polling loop that issues `DISPLAY *STATUS` until the object
reaches a stable state or the timeout expires.

## SyncOperation

An integer enum indicating the operation that was performed:

```go
const (
    SyncStarted   SyncOperation = iota  // Object confirmed running
    SyncStopped                         // Object confirmed stopped
    SyncRestarted                       // Stop-then-start completed
)
```

`SyncOperation` implements `fmt.Stringer`, returning `"started"`, `"stopped"`,
or `"restarted"`.

## SyncConfig

A struct controlling the polling behavior:

```go
type SyncConfig struct {
    Timeout      time.Duration  // Max wait before returning TimeoutError (default 30s)
    PollInterval time.Duration  // Duration between status checks (default 1s)
}
```

| Field | Type | Description |
| --- | --- | --- |
| `Timeout` | `time.Duration` | Maximum duration to wait before returning `*TimeoutError` (default: 30s if zero) |
| `PollInterval` | `time.Duration` | Duration between `DISPLAY *STATUS` polls (default: 1s if zero) |

Zero values for `Timeout` and `PollInterval` are replaced with their defaults
(30 seconds and 1 second respectively).

## SyncResult

A struct containing the outcome of a sync operation:

```go
type SyncResult struct {
    Operation      SyncOperation  // What happened: SyncStarted, SyncStopped, or SyncRestarted
    Polls          int            // Number of status polls issued
    ElapsedSeconds float64        // Wall-clock time taken
}
```

| Field | Type | Description |
| --- | --- | --- |
| `Operation` | `SyncOperation` | What happened: `SyncStarted`, `SyncStopped`, or `SyncRestarted` |
| `Polls` | `int` | Number of status polls issued |
| `ElapsedSeconds` | `float64` | Wall-clock seconds from command to confirmation |

## Method signature pattern

All 9 sync methods share the same signature:

```go
func (session *Session) StartChannelSync(
    ctx    context.Context,
    name   string,
    config SyncConfig,
) (SyncResult, error)
```

Every sync method takes `context.Context` as its first parameter,
following standard Go conventions for I/O operations. The `config`
parameter is always required; pass a zero-value `SyncConfig{}` for
defaults.

## Basic usage

```go
ctx := context.Background()

// Start a channel and wait until it is RUNNING
result, err := session.StartChannelSync(ctx, "TO.PARTNER", mqrestadmin.SyncConfig{})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Channel running after %d poll(s), %.1fs\n", result.Polls, result.ElapsedSeconds)

// Stop a listener and wait until it is STOPPED
result, err = session.StopListenerSync(ctx, "TCP.LISTENER", mqrestadmin.SyncConfig{})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result.Operation) // "stopped"
```

Passing a zero-value `SyncConfig{}` uses the defaults (30-second timeout,
1-second poll interval).

## Custom timeout and poll interval

Pass a `SyncConfig` with non-zero fields to override the defaults:

```go
// Aggressive polling for fast local development
fast := mqrestadmin.SyncConfig{
    Timeout:      10 * time.Second,
    PollInterval: 250 * time.Millisecond,
}
result, err := session.StartServiceSync(ctx, "MY.SVC", fast)

// Patient polling for remote queue managers
patient := mqrestadmin.SyncConfig{
    Timeout:      120 * time.Second,
    PollInterval: 5 * time.Second,
}
result, err = session.StartChannelSync(ctx, "REMOTE.CHL", patient)
```

## Restart convenience

The `Restart*` methods perform a synchronous stop followed by a
synchronous start. Each phase gets the full timeout independently --
worst case is 2x the configured timeout.

The returned `SyncResult` reports **total** polls and **total** elapsed
time across both phases:

```go
ctx := context.Background()

result, err := session.RestartChannel(ctx, "TO.PARTNER", mqrestadmin.SyncConfig{})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Restarted in %.1fs (%d total polls)\n", result.ElapsedSeconds, result.Polls)
```

## Timeout handling

When the timeout expires, a `*TimeoutError` is returned with
diagnostic fields:

```go
result, err := session.StartChannelSync(ctx, "BROKEN.CHL", mqrestadmin.SyncConfig{
    Timeout:      15 * time.Second,
    PollInterval: 1 * time.Second,
})
if err != nil {
    var timeoutErr *mqrestadmin.TimeoutError
    if errors.As(err, &timeoutErr) {
        fmt.Printf("Name: %s\n", timeoutErr.Name)               // "BROKEN.CHL"
        fmt.Printf("Operation: %s\n", timeoutErr.Operation)      // "started"
        fmt.Printf("Elapsed: %.1fs\n", timeoutErr.ElapsedSeconds) // 15.0
    }
}
```

`TimeoutError` can be matched with `errors.As`, following standard Go
error-handling conventions.

## Available methods

| Method | Operation | START/STOP qualifier | Status qualifier |
| --- | --- | --- | --- |
| `StartChannelSync()` | Start | `CHANNEL` | `CHSTATUS` |
| `StopChannelSync()` | Stop | `CHANNEL` | `CHSTATUS` |
| `RestartChannel()` | Restart | `CHANNEL` | `CHSTATUS` |
| `StartListenerSync()` | Start | `LISTENER` | `LSSTATUS` |
| `StopListenerSync()` | Stop | `LISTENER` | `LSSTATUS` |
| `RestartListener()` | Restart | `LISTENER` | `LSSTATUS` |
| `StartServiceSync()` | Start | `SERVICE` | `SVSTATUS` |
| `StopServiceSync()` | Stop | `SERVICE` | `SVSTATUS` |
| `RestartService()` | Restart | `SERVICE` | `SVSTATUS` |

## Status detection

The polling loop checks the `STATUS` attribute in the `DISPLAY *STATUS`
response. The target values are:

- **Start**: `RUNNING`
- **Stop**: `STOPPED` or `INACTIVE`

### Channel stop edge case

When a channel stops, its `CHSTATUS` record may disappear entirely
(the `DISPLAY CHSTATUS` response returns no rows). The channel sync
methods treat an empty status result as successfully stopped. Listener
and service status records are always present, so empty results are not
treated as stopped for those object types.

## Attribute mapping

The sync methods call the internal MQSC command layer, so they participate
in the same [mapping pipeline](../mapping-pipeline.md) as all other
command methods. The status key is checked using both the mapped
`snake_case` name and the raw MQSC name, so polling works correctly
regardless of whether mapping is enabled or disabled.

## Provisioning example

The sync methods pair naturally with the
[ensure methods](ensure.md) for end-to-end provisioning:

```go
ctx := context.Background()
config := mqrestadmin.SyncConfig{Timeout: 60 * time.Second}

// Ensure listeners exist for application and admin traffic
_, err := session.EnsureListener(ctx, "APP.LISTENER", map[string]any{
    "transport_type": "TCP",
    "port":           1415,
    "start_mode":     "MQSVC_CONTROL_Q_MGR",
})
if err != nil {
    log.Fatal(err)
}

_, err = session.EnsureListener(ctx, "ADMIN.LISTENER", map[string]any{
    "transport_type": "TCP",
    "port":           1416,
    "start_mode":     "MQSVC_CONTROL_Q_MGR",
})
if err != nil {
    log.Fatal(err)
}

// Start them synchronously
_, err = session.StartListenerSync(ctx, "APP.LISTENER", config)
if err != nil {
    log.Fatal(err)
}
_, err = session.StartListenerSync(ctx, "ADMIN.LISTENER", config)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Listeners ready")
```

## Rolling restart example

Restart all listeners with error handling -- useful when a queue
manager serves multiple TCP ports for different client populations:

```go
ctx := context.Background()
listeners := []string{"APP.LISTENER", "ADMIN.LISTENER", "PARTNER.LISTENER"}
config := mqrestadmin.SyncConfig{
    Timeout:      30 * time.Second,
    PollInterval: 2 * time.Second,
}

for _, name := range listeners {
    result, err := session.RestartListener(ctx, name, config)
    if err != nil {
        var timeoutErr *mqrestadmin.TimeoutError
        if errors.As(err, &timeoutErr) {
            fmt.Printf("%s: timed out after %.1fs\n", name, timeoutErr.ElapsedSeconds)
            continue
        }
        log.Fatalf("%s: %v", name, err)
    }
    fmt.Printf("%s: restarted in %.1fs\n", name, result.ElapsedSeconds)
}
```
