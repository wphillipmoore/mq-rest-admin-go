# Sync

## Overview

The sync methods provide 9 synchronous start/stop/restart operations on
`Session`. These methods wrap fire-and-forget `START` and `STOP` commands with
a polling loop that waits until the object reaches its target state or the
timeout expires.

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

## Available sync methods

| Method | Operation | Object type |
| --- | --- | --- |
| `StartChannelSync()` | Start | Channel |
| `StopChannelSync()` | Stop | Channel |
| `RestartChannel()` | Restart | Channel |
| `StartListenerSync()` | Start | Listener |
| `StopListenerSync()` | Stop | Listener |
| `RestartListener()` | Restart | Listener |
| `StartServiceSync()` | Start | Service |
| `StopServiceSync()` | Stop | Service |
| `RestartService()` | Restart | Service |

## Usage

```go
ctx := context.Background()

// Start a channel and wait for RUNNING status
result, err := session.StartChannelSync(ctx, "TO.PARTNER", mqrestadmin.SyncConfig{})
if err != nil {
    var timeoutErr *mqrestadmin.TimeoutError
    if errors.As(err, &timeoutErr) {
        fmt.Printf("Timed out after %.1fs\n", timeoutErr.ElapsedSeconds)
    }
    log.Fatal(err)
}
fmt.Printf("Channel %s after %d polls (%.1fs)\n",
    result.Operation, result.Polls, result.ElapsedSeconds)

// Custom polling configuration
result, err = session.StopChannelSync(ctx, "TO.PARTNER", mqrestadmin.SyncConfig{
    Timeout:      10 * time.Second,
    PollInterval: 500 * time.Millisecond,
})

// Restart (stop then start, polling at each step)
result, err = session.RestartChannel(ctx, "TO.PARTNER", mqrestadmin.SyncConfig{})
// result.Operation == SyncRestarted
// result.Polls includes polls from both stop and start phases
```

See [Sync Methods](../sync-methods.md) for the full conceptual overview,
polling behavior, and status detection logic.
