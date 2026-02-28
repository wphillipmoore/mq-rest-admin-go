// Environment provisioner example.
//
// Defines a complete set of queues, channels, and remote queue
// definitions across two queue managers, then verifies connectivity.
// Includes teardown.
//
// Usage:
//
//	go run ./examples/cmd/provisionenv
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

	qm1, err := mqrestadmin.NewSession(
		envOr("MQ_REST_BASE_URL", "https://localhost:9463/ibmmq/rest/v2"),
		"QM1",
		mqrestadmin.BasicAuth{
			Username: envOr("MQ_ADMIN_USER", "mqadmin"),
			Password: envOr("MQ_ADMIN_PASSWORD", "mqadmin"),
		},
		mqrestadmin.WithVerifyTLS(false),
	)
	if err != nil {
		log.Fatal(err)
	}

	qm2, err := mqrestadmin.NewSession(
		envOr("MQ_REST_BASE_URL_QM2", "https://localhost:9464/ibmmq/rest/v2"),
		"QM2",
		mqrestadmin.BasicAuth{
			Username: envOr("MQ_ADMIN_USER", "mqadmin"),
			Password: envOr("MQ_ADMIN_PASSWORD", "mqadmin"),
		},
		mqrestadmin.WithVerifyTLS(false),
	)
	if err != nil {
		log.Fatal(err)
	}

	examples.PrintProvision(ctx, qm1, qm2)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
