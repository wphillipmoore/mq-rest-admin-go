// Dead letter queue inspector example.
//
// Checks the dead letter queue configuration, reports depth and
// capacity, and suggests actions when messages are present.
//
// Usage:
//
//	go run ./examples/cmd/dlqinspector
package main

import (
	"context"
	"log"
	"os"

	"github.com/wphillipmoore/mq-rest-admin-go/examples"
	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

func main() {
	ctx := context.Background()

	session, err := mqrestadmin.NewSession(
		envOr("MQ_REST_BASE_URL", "https://localhost:9463/ibmmq/rest/v2"),
		envOr("MQ_QMGR_NAME", "QM1"),
		mqrestadmin.BasicAuth{
			Username: envOr("MQ_ADMIN_USER", "mqadmin"),
			Password: envOr("MQ_ADMIN_PASSWORD", "mqadmin"),
		},
		mqrestadmin.WithVerifyTLS(false),
	)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := examples.PrintDLQInspection(ctx, session); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
