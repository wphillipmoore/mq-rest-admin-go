package examples_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/wphillipmoore/mq-rest-admin-go/examples"
	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

func TestMain(m *testing.M) {
	if os.Getenv("MQ_REST_ADMIN_RUN_INTEGRATION") != "1" {
		fmt.Println("skipping examples integration tests (MQ_REST_ADMIN_RUN_INTEGRATION != 1)")
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func qm1Session(t *testing.T) *mqrestadmin.Session {
	t.Helper()
	session, err := mqrestadmin.NewSession(
		envOr("MQ_REST_BASE_URL", "https://localhost:9463/ibmmq/rest/v2"),
		"QM1",
		mqrestadmin.BasicAuth{
			Username: envOr("MQ_ADMIN_USER", "mqadmin"),
			Password: envOr("MQ_ADMIN_PASSWORD", "mqadmin"),
		},
		mqrestadmin.WithVerifyTLS(false),
	)
	if err != nil {
		t.Fatalf("create QM1 session: %v", err)
	}
	return session
}

func qm2Session(t *testing.T) *mqrestadmin.Session {
	t.Helper()
	session, err := mqrestadmin.NewSession(
		envOr("MQ_REST_BASE_URL_QM2", "https://localhost:9464/ibmmq/rest/v2"),
		"QM2",
		mqrestadmin.BasicAuth{
			Username: envOr("MQ_ADMIN_USER", "mqadmin"),
			Password: envOr("MQ_ADMIN_PASSWORD", "mqadmin"),
		},
		mqrestadmin.WithVerifyTLS(false),
	)
	if err != nil {
		t.Fatalf("create QM2 session: %v", err)
	}
	return session
}

// ---------------------------------------------------------------------------
// Health check
// ---------------------------------------------------------------------------

func TestHealthCheckQM1(t *testing.T) {
	ctx := context.Background()
	result := examples.CheckHealth(ctx, qm1Session(t))

	if !result.Reachable {
		t.Fatal("QM1 should be reachable")
	}
	if !result.Passed {
		t.Fatal("QM1 health check should pass")
	}
	if result.QmgrName != "QM1" {
		t.Fatalf("expected QM1, got %s", result.QmgrName)
	}
}

func TestHealthCheckQM2(t *testing.T) {
	ctx := context.Background()
	result := examples.CheckHealth(ctx, qm2Session(t))

	if !result.Reachable {
		t.Fatal("QM2 should be reachable")
	}
	if !result.Passed {
		t.Fatal("QM2 health check should pass")
	}
	if result.QmgrName != "QM2" {
		t.Fatalf("expected QM2, got %s", result.QmgrName)
	}
}

// ---------------------------------------------------------------------------
// Queue depth monitor
// ---------------------------------------------------------------------------

func TestQueueDepthMonitor(t *testing.T) {
	ctx := context.Background()
	results, err := examples.MonitorQueueDepths(ctx, qm1Session(t), 80.0)
	if err != nil {
		t.Fatalf("monitor queue depths: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("should find local queues")
	}

	found := false
	for _, q := range results {
		if q.Name == "DEV.QLOCAL" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("should include DEV.QLOCAL")
	}
}

// ---------------------------------------------------------------------------
// Channel status
// ---------------------------------------------------------------------------

func TestChannelStatusReport(t *testing.T) {
	ctx := context.Background()
	results, err := examples.ReportChannelStatus(ctx, qm1Session(t))
	if err != nil {
		t.Fatalf("report channel status: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("should find channels")
	}

	found := false
	for _, c := range results {
		if c.Name == "DEV.SVRCONN" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("should include DEV.SVRCONN")
	}
}

// ---------------------------------------------------------------------------
// DLQ inspector
// ---------------------------------------------------------------------------

func TestDLQInspector(t *testing.T) {
	ctx := context.Background()
	report, err := examples.InspectDLQ(ctx, qm1Session(t))
	if err != nil {
		t.Fatalf("inspect DLQ: %v", err)
	}

	if !report.Configured {
		t.Fatal("DLQ should be configured")
	}
	if report.DLQName != "DEV.DEAD.LETTER" {
		t.Fatalf("expected DEV.DEAD.LETTER, got %s", report.DLQName)
	}
	if report.CurrentDepth != 0 {
		t.Fatalf("expected depth 0, got %d", report.CurrentDepth)
	}
}

// ---------------------------------------------------------------------------
// Queue status
// ---------------------------------------------------------------------------

func TestQueueStatusHandles(t *testing.T) {
	ctx := context.Background()
	handles := examples.ReportQueueHandles(ctx, qm1Session(t))

	if handles == nil {
		t.Fatal("should return a slice (possibly empty), not nil")
	}
}

func TestConnectionHandles(t *testing.T) {
	ctx := context.Background()
	handles := examples.ReportConnectionHandles(ctx, qm1Session(t))

	if handles == nil {
		t.Fatal("should return a slice (possibly empty), not nil")
	}
}

// ---------------------------------------------------------------------------
// Provision and teardown
// ---------------------------------------------------------------------------

func TestProvisionAndTeardown(t *testing.T) {
	ctx := context.Background()
	qm1 := qm1Session(t)
	qm2 := qm2Session(t)

	result := examples.Provision(ctx, qm1, qm2)

	if len(result.ObjectsCreated) == 0 {
		t.Fatal("should create objects")
	}
	if !result.Verified {
		t.Fatal("verification should pass")
	}

	failures := examples.Teardown(ctx, qm1, qm2)
	if len(failures) > 0 {
		t.Fatalf("teardown failures: %v", failures)
	}
}
