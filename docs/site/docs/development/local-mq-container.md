# Local MQ Container

A containerized IBM MQ environment provides two queue managers for
development and integration testing.

## Prerequisites

- Docker Desktop or compatible Docker Engine.
- IBM MQ container image access (license acceptance required).
- The `mq-dev-environment` repository cloned as a sibling directory
  (`../mq-dev-environment`), or set `MQ_DEV_ENV_PATH` to its location.

## Configuration

The Docker Compose file in the `mq-dev-environment` repository runs two
queue managers on a shared network (`mq-dev-net`):

| Setting | QM1 | QM2 |
| --- | --- | --- |
| Queue manager | `QM1` | `QM2` |
| MQ listener port | `1414` | `1415` |
| REST API port | `9443` | `9444` |
| Container name | `mq-dev-qm1` | `mq-dev-qm2` |

Both queue managers share the same credentials:

| Setting | Value |
| --- | --- |
| Admin credentials | `mqadmin` / `mqadmin` |
| Read-only credentials | `mqreader` / `mqreader` |
| QM1 REST base URL | `https://localhost:9443/ibmmq/rest/v2` |
| QM2 REST base URL | `https://localhost:9444/ibmmq/rest/v2` |

## Quick start

Start both queue managers:

```bash
./scripts/dev/mq_start.sh
```

Seed deterministic test objects on both QMs (all prefixed with `DEV.`):

```bash
./scripts/dev/mq_seed.sh
```

Verify REST-based MQSC responses on both QMs:

```bash
./scripts/dev/mq_verify.sh
```

## Seed objects

QM1 receives the full set of test objects (queues, channels, topics,
namelists, listeners, processes) plus cross-QM objects for communicating
with QM2. QM2 receives a smaller set of objects plus the reciprocal
cross-QM definitions.

The seed scripts are maintained in the `mq-dev-environment` repository
at `seed/base-qm1.mqsc` and `seed/base-qm2.mqsc`. Both use `REPLACE`
so they can be re-run at any time without side effects.

## Lifecycle scripts

| Script | Purpose |
| --- | --- |
| `scripts/dev/mq_start.sh` | Start both queue managers and wait for REST readiness |
| `scripts/dev/mq_seed.sh` | Seed deterministic test objects on both QMs |
| `scripts/dev/mq_verify.sh` | Verify REST-based MQSC responses on both QMs |
| `scripts/dev/mq_stop.sh` | Stop both queue managers |
| `scripts/dev/mq_reset.sh` | Reset to clean state (removes data volumes) |

## Running integration tests

```bash
# Start MQ and seed configuration
scripts/dev/mq_start.sh
scripts/dev/mq_seed.sh

# Run integration tests
MQ_REST_ADMIN_RUN_INTEGRATION=1 go test -v ./...

# Stop MQ when done
scripts/dev/mq_stop.sh
```

## Environment variables

| Variable | Default | Description |
| --- | --- | --- |
| `MQ_REST_BASE_URL` | `https://localhost:9443/ibmmq/rest/v2` | QM1 REST API base URL |
| `MQ_REST_BASE_URL_QM2` | `https://localhost:9444/ibmmq/rest/v2` | QM2 REST API base URL |
| `MQ_ADMIN_USER` | `mqadmin` | Admin username |
| `MQ_ADMIN_PASSWORD` | `mqadmin` | Admin password |
| `MQ_IMAGE` | `icr.io/ibm-messaging/mq:latest` | Container image |
| `MQ_DEV_ENV_PATH` | `../mq-dev-environment` | Path to mq-dev-environment project |
| `MQ_REST_ADMIN_RUN_INTEGRATION` | (unset) | Set to `1` to enable integration tests |

## Gateway routing

The two-QM local setup supports gateway routing out of the box. The seed
scripts create QM aliases and sender/receiver channels so each queue manager
can route MQSC commands to the other.

### curl example

Query QM2's queue manager attributes through QM1's REST API:

```bash
curl -k -u mqadmin:mqadmin \
  -H "Content-Type: application/json" \
  -H "ibm-mq-rest-csrf-token: local" \
  -H "ibm-mq-rest-gateway-qmgr: QM1" \
  -d '{"type": "runCommandJSON", "command": "DISPLAY", "qualifier": "QMGR"}' \
  https://localhost:9443/ibmmq/rest/v2/admin/action/qmgr/QM2/mqsc
```

### mqrestadmin example

```go
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM2",
    mqrestadmin.WithBasicAuth("mqadmin", "mqadmin"),
    mqrestadmin.WithGatewayQmgr("QM1"),
    mqrestadmin.WithVerifyTLS(false),
)
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()
qmgr, err := session.DisplayQmgr(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println(qmgr) // QM2's attributes, routed through QM1
```

## Reset workflow

To return to a completely clean state (removes both data volumes):

```bash
./scripts/dev/mq_reset.sh
```

## Troubleshooting

If the REST API is not reachable, ensure the embedded web server is
binding to all interfaces:

```bash
docker compose -f ../mq-dev-environment/config/docker-compose.yml exec -T qm1 \
    setmqweb properties -k httpHost -v "*"
```

Then restart the containers and retry the verification workflow.
