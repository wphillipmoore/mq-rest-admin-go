# Examples

The `examples/` directory contains practical scripts that demonstrate common
MQ administration tasks using `mqrestadmin`. Each example is self-contained
and can be run against the local Docker environment.

## Prerequisites

Start the multi-queue-manager Docker environment and seed both queue managers:

```bash
./scripts/dev/mq_start.sh
./scripts/dev/mq_seed.sh
```

This starts two queue managers (`QM1` on port 9443, `QM2` on port 9444) on a
shared Docker network. See [local MQ container](development/local-mq-container.md) for details.

## Health check

Connect to one or more queue managers and check:

- Queue manager attributes via `DisplayQmgr()`
- Running status via `DisplayQmstatus()`
- Listener definitions via `DisplayListener()`

```go
ctx := context.Background()

qmgr, err := session.DisplayQmgr(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Queue manager:", qmgr["queue_manager_name"])

status, err := session.DisplayQmstatus(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Status:", status["channel_initiator_status"])

listeners, err := session.DisplayListener(ctx, "*")
if err != nil {
    log.Fatal(err)
}
for _, listener := range listeners {
    fmt.Printf("Listener: %s port=%v\n",
        listener["listener_name"], listener["port"])
}
```

## Queue depth monitor

Display all local queues with their current depth and flag queues
approaching capacity:

```go
ctx := context.Background()

queues, err := session.DisplayQueue(ctx, "*")
if err != nil {
    log.Fatal(err)
}

for _, queue := range queues {
    depth, _ := strconv.Atoi(fmt.Sprint(queue["current_queue_depth"]))
    maxDepth, _ := strconv.Atoi(fmt.Sprint(queue["max_queue_depth"]))
    pct := 0
    if maxDepth > 0 {
        pct = depth * 100 / maxDepth
    }
    flag := ""
    if pct > 80 {
        flag = " *** HIGH ***"
    }
    fmt.Printf("%-40s %5d / %5d (%d%%)%s\n",
        queue["queue_name"], depth, maxDepth, pct, flag)
}
```

## Channel status report

Cross-reference channel definitions with live channel status:

```go
ctx := context.Background()

channels, err := session.DisplayChannel(ctx, "*")
if err != nil {
    log.Fatal(err)
}

statuses, err := session.DisplayChstatus(ctx, "*")
if err != nil {
    log.Fatal(err)
}

running := make(map[string]bool)
for _, s := range statuses {
    running[fmt.Sprint(s["channel_name"])] = true
}

for _, ch := range channels {
    name := fmt.Sprint(ch["channel_name"])
    state := "INACTIVE"
    if running[name] {
        state = "RUNNING"
    }
    fmt.Println(name + ": " + state)
}
```

## Environment provisioner

Demonstrate bulk provisioning across two queue managers using ensure
methods:

```go
ctx := context.Background()

// Ensure application queues exist on QM1
_, err := session.EnsureQlocal(ctx, "APP.REQUESTS", map[string]any{
    "max_queue_depth":     "50000",
    "default_persistence": "persistent",
})
if err != nil {
    log.Fatal(err)
}

_, err = session.EnsureQlocal(ctx, "APP.RESPONSES", map[string]any{
    "max_queue_depth":     "50000",
    "default_persistence": "persistent",
})
if err != nil {
    log.Fatal(err)
}

// Ensure listeners are running
config := mqrestadmin.SyncConfig{Timeout: 60 * time.Second}
_, err = session.StartListenerSync(ctx, "TCP.LISTENER", config)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Environment provisioned")
```

## Dead letter queue inspector

Inspect the dead letter queue configuration:

```go
ctx := context.Background()

qmgr, err := session.DisplayQmgr(ctx)
if err != nil {
    log.Fatal(err)
}

dlqName, ok := qmgr["dead_letter_q_name"].(string)
if ok && dlqName != "" {
    dlq, err := session.DisplayQueue(ctx, dlqName)
    if err != nil {
        log.Fatal(err)
    }
    if len(dlq) > 0 {
        fmt.Printf("DLQ: %s depth=%v max=%v\n",
            dlqName, dlq[0]["current_queue_depth"], dlq[0]["max_queue_depth"])
    }
} else {
    fmt.Println("No dead letter queue configured")
}
```
