# Local MQ Container

--8<-- "development/local-mq-container.md"

## Go-specific notes

### Running integration tests

```bash
# Start MQ and seed configuration
scripts/dev/mq_start.sh
scripts/dev/mq_seed.sh

# Run integration tests
MQ_REST_ADMIN_RUN_INTEGRATION=1 go test -v ./...

# Stop MQ when done
scripts/dev/mq_stop.sh
```

### Environment variables

| Variable | Default | Description |
| --- | --- | --- |
| `MQ_REST_ADMIN_RUN_INTEGRATION` | (unset) | Set to `1` to enable integration tests |

### Gateway routing example

```go
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM2",
    mqrestadmin.LTPAAuth{Username: "mqadmin", Password: "mqadmin"},
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
