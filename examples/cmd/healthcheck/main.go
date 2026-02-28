// Health check example.
//
// Connects to one or more queue managers and checks QMGR status,
// command server availability, and listener state.
//
// Usage:
//
//	go run ./examples/cmd/healthcheck
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

	sessions := []*mqrestadmin.Session{mustSession(
		envOr("MQ_REST_BASE_URL", "https://localhost:9463/ibmmq/rest/v2"),
		envOr("MQ_QMGR_NAME", "QM1"),
	)}

	if qm2URL := os.Getenv("MQ_REST_BASE_URL_QM2"); qm2URL != "" {
		sessions = append(sessions, mustSession(qm2URL, "QM2"))
	}

	examples.PrintHealthCheck(ctx, sessions)
}

func mustSession(baseURL, qmgrName string) *mqrestadmin.Session {
	session, err := mqrestadmin.NewSession(
		baseURL, qmgrName,
		mqrestadmin.BasicAuth{
			Username: envOr("MQ_ADMIN_USER", "mqadmin"),
			Password: envOr("MQ_ADMIN_PASSWORD", "mqadmin"),
		},
		mqrestadmin.WithVerifyTLS(false),
	)
	if err != nil {
		log.Fatal(err)
	}
	return session
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
