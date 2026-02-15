# Ensure

## Overview

The ensure methods provide 15 idempotent ensure operations on `Session`. These
methods implement a declarative upsert pattern: DEFINE if the object does not
exist, ALTER only attributes that differ, or no-op if the object already matches
the desired state.

## EnsureAction

An integer enum indicating the action taken by an ensure method:

```go
const (
    EnsureCreated   EnsureAction = iota  // Object did not exist; DEFINE was issued
    EnsureUpdated                        // Object existed but attributes differed; ALTER was issued
    EnsureUnchanged                      // Object already matched the desired state
)
```

`EnsureAction` implements `fmt.Stringer`, returning `"created"`, `"updated"`,
or `"unchanged"`.

## EnsureResult

A struct containing the action taken and the list of attribute names that
triggered the change (if any):

```go
type EnsureResult struct {
    Action  EnsureAction  // What happened: EnsureCreated, EnsureUpdated, or EnsureUnchanged
    Changed []string      // Attribute names that triggered an ALTER (in the caller's namespace)
}
```

| Field | Type | Description |
| --- | --- | --- |
| `Action` | `EnsureAction` | What happened: `EnsureCreated`, `EnsureUpdated`, or `EnsureUnchanged` |
| `Changed` | `[]string` | Attribute names that triggered an ALTER (in the caller's namespace) |

## Method signature patterns

Named-object ensure methods share this signature:

```go
func (session *Session) EnsureQlocal(
    ctx               context.Context,
    name              string,
    requestParameters map[string]any,
) (EnsureResult, error)
```

The queue manager ensure method omits the name parameter:

```go
func (session *Session) EnsureQmgr(
    ctx               context.Context,
    requestParameters map[string]any,
) (EnsureResult, error)
```

## Available ensure methods

| Method | Object type |
| --- | --- |
| `EnsureQmgr()` | Queue manager (always exists; never returns `EnsureCreated`) |
| `EnsureQlocal()` | Local queue |
| `EnsureQremote()` | Remote queue |
| `EnsureQalias()` | Alias queue |
| `EnsureQmodel()` | Model queue |
| `EnsureChannel()` | Channel |
| `EnsureAuthinfo()` | Authentication information object |
| `EnsureListener()` | Listener |
| `EnsureNamelist()` | Namelist |
| `EnsureProcess()` | Process |
| `EnsureService()` | Service |
| `EnsureTopic()` | Topic |
| `EnsureSub()` | Subscription |
| `EnsureStgclass()` | Storage class |
| `EnsureComminfo()` | Communication information object |
| `EnsureCfstruct()` | CF structure |

## Usage

```go
ctx := context.Background()

result, err := session.EnsureQlocal(ctx, "MY.QUEUE",
    map[string]any{
        "max_queue_depth": 50000,
        "description":     "App queue",
    })
if err != nil {
    log.Fatal(err)
}

switch result.Action {
case mqrestadmin.EnsureCreated:
    fmt.Println("Queue created")
case mqrestadmin.EnsureUpdated:
    fmt.Println("Changed:", result.Changed)
case mqrestadmin.EnsureUnchanged:
    fmt.Println("Already correct")
}
```

## How comparison works

The ensure methods compare desired attributes against the current state using
case-insensitive string comparison after trimming whitespace. Only attributes
that differ are included in the ALTER command. Attributes not specified in the
request are left unchanged.

See [Ensure Methods](../ensure-methods.md) for the full conceptual overview
and comparison logic.
