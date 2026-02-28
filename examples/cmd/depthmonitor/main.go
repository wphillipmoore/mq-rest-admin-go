// Queue depth monitor example.
//
// Displays local queues with their current depth, flags queues
// approaching capacity, and sorts by depth percentage.
//
// Usage:
//
//	go run ./examples/cmd/depthmonitor
package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/wphillipmoore/mq-rest-admin-go/examples"
	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

func main() {
	ctx := context.Background()

	threshold := 80.0
	if s := os.Getenv("DEPTH_THRESHOLD_PCT"); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			threshold = v
		}
	}

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

	if _, err := examples.PrintQueueDepths(ctx, session, threshold); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
