# Examples

Runnable example programs demonstrate common MQ administration tasks using
`mqrestadmin`. Each example has a core function in the
[`examples/`](https://github.com/wphillipmoore/mq-rest-admin-go/tree/main/examples)
package and a standalone `main.go` entry point in
[`examples/cmd/`](https://github.com/wphillipmoore/mq-rest-admin-go/tree/main/examples/cmd).

## Prerequisites

Start the multi-queue-manager Docker environment and seed both queue managers:

```bash
./scripts/dev/mq_start.sh
./scripts/dev/mq_seed.sh
```

This starts two queue managers (`QM1` on port 9463, `QM2` on port 9464) on a
shared Docker network. See [local MQ container](development/local-mq-container.md) for details.

## Environment variables

| Variable               | Default                                | Description                   |
|------------------------|----------------------------------------|-------------------------------|
| `MQ_REST_BASE_URL`     | `https://localhost:9463/ibmmq/rest/v2` | QM1 REST endpoint             |
| `MQ_REST_BASE_URL_QM2` | `https://localhost:9464/ibmmq/rest/v2` | QM2 REST endpoint             |
| `MQ_QMGR_NAME`         | `QM1`                                  | Queue manager name            |
| `MQ_ADMIN_USER`        | `mqadmin`                              | Admin username                |
| `MQ_ADMIN_PASSWORD`    | `mqadmin`                              | Admin password                |
| `DEPTH_THRESHOLD_PCT`  | `80`                                   | Queue depth warning threshold |

## Health check

Connects to one or more queue managers and checks QMGR status,
command server availability, and listener state. Produces a pass/fail
summary for each queue manager.

```bash
go run ./examples/cmd/healthcheck
```

See [`examples/healthcheck.go`](https://github.com/wphillipmoore/mq-rest-admin-go/blob/main/examples/healthcheck.go).

## Queue depth monitor

Displays local queues with their current depth, flags queues
approaching capacity, and sorts by depth percentage.

```bash
go run ./examples/cmd/depthmonitor
```

See [`examples/depthmonitor.go`](https://github.com/wphillipmoore/mq-rest-admin-go/blob/main/examples/depthmonitor.go).

## Channel status report

Displays channel definitions alongside live channel status, identifies
channels that are defined but not running, and shows connection details.

```bash
go run ./examples/cmd/channelstatus
```

See [`examples/channelstatus.go`](https://github.com/wphillipmoore/mq-rest-admin-go/blob/main/examples/channelstatus.go).

## Environment provisioner

Defines a complete set of queues, channels, and remote queue definitions
across two queue managers, then verifies connectivity. Includes teardown.

```bash
go run ./examples/cmd/provisionenv
```

See [`examples/provisionenv.go`](https://github.com/wphillipmoore/mq-rest-admin-go/blob/main/examples/provisionenv.go).

## Dead letter queue inspector

Checks the dead letter queue configuration, reports depth and capacity,
and suggests actions when messages are present.

```bash
go run ./examples/cmd/dlqinspector
```

See [`examples/dlqinspector.go`](https://github.com/wphillipmoore/mq-rest-admin-go/blob/main/examples/dlqinspector.go).

## Queue status and connection handles

Demonstrates `DISPLAY QSTATUS TYPE(HANDLE)` and `DISPLAY CONN TYPE(HANDLE)`
queries, showing how `mqrestadmin` flattens nested object response
structures into uniform flat maps.

```bash
go run ./examples/cmd/queuestatus
```

See [`examples/queuestatus.go`](https://github.com/wphillipmoore/mq-rest-admin-go/blob/main/examples/queuestatus.go).
