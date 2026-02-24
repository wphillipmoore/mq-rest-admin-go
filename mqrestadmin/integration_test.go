//go:build integration

package mqrestadmin_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

// ---------------------------------------------------------------------------
// Seeded object names (created by mq_seed.sh)
// ---------------------------------------------------------------------------

var seededQueues = []string{
	"DEV.DEAD.LETTER",
	"DEV.QLOCAL",
	"DEV.QREMOTE",
	"DEV.QALIAS",
	"DEV.QMODEL",
	"DEV.XMITQ",
}

var seededChannels = []string{
	"DEV.SVRCONN",
	"DEV.SDR",
	"DEV.RCVR",
}

const (
	seededListener = "DEV.LSTR"
	seededTopic    = "DEV.TOPIC"
	seededNamelist = "DEV.NAMELIST"
	seededProcess  = "DEV.PROC"

	testQlocal   = "DEV.TEST.QLOCAL"
	testQremote  = "DEV.TEST.QREMOTE"
	testQalias   = "DEV.TEST.QALIAS"
	testQmodel   = "DEV.TEST.QMODEL"
	testChannel  = "DEV.TEST.SVRCONN"
	testListener = "DEV.TEST.LSTR"
	testProcess  = "DEV.TEST.PROC"
	testTopic    = "DEV.TEST.TOPIC"
	testNamelist = "DEV.TEST.NAMELIST"

	testEnsureQlocal  = "DEV.ENSURE.QLOCAL"
	testEnsureChannel = "DEV.ENSURE.CHL"
)

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

type integrationConfig struct {
	restBaseURL    string
	restBaseURLQM2 string
	adminUser      string
	adminPassword  string
	qmgrName       string
	qmgrNameQM2    string
	verifyTLS      bool
}

func loadIntegrationConfig() integrationConfig {
	return integrationConfig{
		restBaseURL:    envOrDefault("MQ_REST_BASE_URL", "https://localhost:9443/ibmmq/rest/v2"),
		restBaseURLQM2: envOrDefault("MQ_REST_BASE_URL_QM2", "https://localhost:9444/ibmmq/rest/v2"),
		adminUser:      envOrDefault("MQ_ADMIN_USER", "mqadmin"),
		adminPassword:  envOrDefault("MQ_ADMIN_PASSWORD", "mqadmin"),
		qmgrName:       envOrDefault("MQ_QMGR_NAME", "QM1"),
		qmgrNameQM2:    envOrDefault("MQ_QMGR_NAME_QM2", "QM2"),
		verifyTLS:      parseBool(envOrDefault("MQ_REST_VERIFY_TLS", "false")),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseBool(value string) bool {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// Session builders
// ---------------------------------------------------------------------------

func buildSession(t *testing.T, cfg integrationConfig) *mqrestadmin.Session {
	t.Helper()
	session, err := mqrestadmin.NewSession(
		cfg.restBaseURL,
		cfg.qmgrName,
		mqrestadmin.BasicAuth{Username: cfg.adminUser, Password: cfg.adminPassword},
		mqrestadmin.WithVerifyTLS(cfg.verifyTLS),
		mqrestadmin.WithMappingStrict(false),
	)
	if err != nil {
		t.Fatalf("NewSession: %v", err)
	}
	return session
}

func buildGatewaySession(t *testing.T, cfg integrationConfig, targetQmgr, gatewayQmgr, restBaseURL string) *mqrestadmin.Session {
	t.Helper()
	session, err := mqrestadmin.NewSession(
		restBaseURL,
		targetQmgr,
		mqrestadmin.BasicAuth{Username: cfg.adminUser, Password: cfg.adminPassword},
		mqrestadmin.WithGatewayQmgr(gatewayQmgr),
		mqrestadmin.WithVerifyTLS(cfg.verifyTLS),
		mqrestadmin.WithMappingStrict(false),
	)
	if err != nil {
		t.Fatalf("NewSession (gateway): %v", err)
	}
	return session
}

// ---------------------------------------------------------------------------
// Assertion helpers
// ---------------------------------------------------------------------------

func containsStringValue(obj map[string]any, expected string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(expected))
	for _, v := range obj {
		if s, ok := v.(string); ok && strings.ToUpper(strings.TrimSpace(s)) == normalized {
			return true
		}
	}
	return false
}

func anyContainsValue(results []map[string]any, expected string) bool {
	for _, obj := range results {
		if containsStringValue(obj, expected) {
			return true
		}
	}
	return false
}

func findMatchingObject(results []map[string]any, expected string) map[string]any {
	for _, obj := range results {
		if containsStringValue(obj, expected) {
			return obj
		}
	}
	return nil
}

func getAttributeCaseInsensitive(obj map[string]any, name string) (any, bool) {
	upper := strings.ToUpper(name)
	for k, v := range obj {
		if strings.ToUpper(k) == upper {
			return v, true
		}
	}
	return nil, false
}

// silentDelete suppresses errors from delete operations used for cleanup.
func silentDelete(fn func() error) {
	_ = fn()
}

// skipIfLifecycleDisabled skips tests that create/modify MQ objects when the
// MQ_SKIP_LIFECYCLE env var is set.
func skipIfLifecycleDisabled(t *testing.T) {
	t.Helper()
	if os.Getenv("MQ_SKIP_LIFECYCLE") == "1" {
		t.Skip("lifecycle tests disabled (MQ_SKIP_LIFECYCLE=1)")
	}
}

// ---------------------------------------------------------------------------
// Display tests — singletons
// ---------------------------------------------------------------------------

func TestDisplayQmgr(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	result, err := session.DisplayQmgr(ctx)
	if err != nil {
		t.Fatalf("DisplayQmgr: %v", err)
	}
	if result == nil {
		t.Fatal("DisplayQmgr returned nil")
	}
	if !containsStringValue(result, cfg.qmgrName) {
		t.Errorf("DisplayQmgr result does not contain qmgr name %q", cfg.qmgrName)
	}
}

func TestDisplayQmstatus(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	result, err := session.DisplayQmstatus(ctx)
	if err != nil {
		t.Fatalf("DisplayQmstatus: %v", err)
	}
	// Result may be nil or a map — both are acceptable.
	if result != nil {
		if _, ok := result["dummy"]; ok {
			_ = ok // just exercising the map type
		}
	}
}

func TestDisplayCmdserv(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	result, err := session.DisplayCmdserv(ctx)
	if err != nil {
		t.Fatalf("DisplayCmdserv: %v", err)
	}
	// Result may be nil or a map — both are acceptable.
	_ = result
}

// ---------------------------------------------------------------------------
// Display tests — seeded queues
// ---------------------------------------------------------------------------

func TestDisplaySeededQueues(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	for _, queueName := range seededQueues {
		t.Run(queueName, func(t *testing.T) {
			results, err := session.DisplayQueue(ctx, queueName)
			if err != nil {
				t.Fatalf("DisplayQueue(%s): %v", queueName, err)
			}
			if len(results) == 0 {
				t.Fatalf("DisplayQueue(%s) returned empty results", queueName)
			}
			if !anyContainsValue(results, queueName) {
				t.Errorf("DisplayQueue results do not contain %q", queueName)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Display tests — qstatus
// ---------------------------------------------------------------------------

func TestDisplayQstatus(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	results, err := session.DisplayQstatus(ctx, "DEV.QLOCAL")
	if err != nil {
		t.Fatalf("DisplayQstatus: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("DisplayQstatus returned empty results")
	}
	if !anyContainsValue(results, "DEV.QLOCAL") {
		t.Error("DisplayQstatus results do not contain DEV.QLOCAL")
	}
}

// ---------------------------------------------------------------------------
// Display tests — seeded channels
// ---------------------------------------------------------------------------

func TestDisplaySeededChannels(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	for _, channelName := range seededChannels {
		t.Run(channelName, func(t *testing.T) {
			results, err := session.DisplayChannel(ctx, channelName)
			if err != nil {
				t.Fatalf("DisplayChannel(%s): %v", channelName, err)
			}
			if len(results) == 0 {
				t.Fatalf("DisplayChannel(%s) returned empty results", channelName)
			}
			if !anyContainsValue(results, channelName) {
				t.Errorf("DisplayChannel results do not contain %q", channelName)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Display tests — seeded objects (listener, topic, namelist, process)
// ---------------------------------------------------------------------------

func TestDisplaySeededListener(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	results, err := session.DisplayListener(ctx, seededListener)
	if err != nil {
		t.Fatalf("DisplayListener: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("DisplayListener returned empty results")
	}
	if !anyContainsValue(results, seededListener) {
		t.Errorf("DisplayListener results do not contain %q", seededListener)
	}
}

func TestDisplaySeededTopic(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	results, err := session.DisplayTopic(ctx, seededTopic)
	if err != nil {
		t.Fatalf("DisplayTopic: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("DisplayTopic returned empty results")
	}
	if !anyContainsValue(results, seededTopic) {
		t.Errorf("DisplayTopic results do not contain %q", seededTopic)
	}
}

func TestDisplaySeededNamelist(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	results, err := session.DisplayNamelist(ctx, seededNamelist)
	if err != nil {
		t.Fatalf("DisplayNamelist: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("DisplayNamelist returned empty results")
	}
	if !anyContainsValue(results, seededNamelist) {
		t.Errorf("DisplayNamelist results do not contain %q", seededNamelist)
	}
}

func TestDisplaySeededProcess(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	results, err := session.DisplayProcess(ctx, seededProcess)
	if err != nil {
		t.Fatalf("DisplayProcess: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("DisplayProcess returned empty results")
	}
	if !anyContainsValue(results, seededProcess) {
		t.Errorf("DisplayProcess results do not contain %q", seededProcess)
	}
}

// ---------------------------------------------------------------------------
// Lifecycle CRUD tests
// ---------------------------------------------------------------------------

type lifecycleCase struct {
	name             string
	objectName       string
	define           func(ctx context.Context, name string, opts ...mqrestadmin.CommandOption) error
	display          func(ctx context.Context, name string, opts ...mqrestadmin.CommandOption) ([]map[string]any, error)
	delete           func(ctx context.Context, name string, opts ...mqrestadmin.CommandOption) error
	defineParams     map[string]any
	alter            func(ctx context.Context, name string, opts ...mqrestadmin.CommandOption) error
	alterParams      map[string]any
	alterDescription string
}

func TestMutatingObjectLifecycle(t *testing.T) {
	skipIfLifecycleDisabled(t)
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	cases := []lifecycleCase{
		{
			name:       "qlocal",
			objectName: testQlocal,
			define:     session.DefineQlocal,
			display:    session.DisplayQueue,
			delete:     session.DeleteQueue,
			defineParams: map[string]any{
				"default_persistence": "yes",
				"description":         "dev test qlocal",
			},
		},
		{
			name:       "qremote",
			objectName: testQremote,
			define:     session.DefineQremote,
			display:    session.DisplayQueue,
			delete:     session.DeleteQueue,
			defineParams: map[string]any{
				"remote_queue_name":          "DEV.TARGET",
				"remote_queue_manager_name":  cfg.qmgrName,
				"transmission_queue_name":    "DEV.XMITQ",
				"description":               "dev test qremote",
			},
		},
		{
			name:       "qalias",
			objectName: testQalias,
			define:     session.DefineQalias,
			display:    session.DisplayQueue,
			delete:     session.DeleteQueue,
			defineParams: map[string]any{
				"target_queue_name": "DEV.QLOCAL",
				"description":       "dev test qalias",
			},
		},
		{
			name:       "qmodel",
			objectName: testQmodel,
			define:     session.DefineQmodel,
			display:    session.DisplayQueue,
			delete:     session.DeleteQueue,
			defineParams: map[string]any{
				"definition_type":            "TEMPDYN",
				"default_input_open_option":  "SHARED",
				"description":               "dev test qmodel",
			},
		},
		{
			name:       "channel",
			objectName: testChannel,
			define:     session.DefineChannel,
			display:    session.DisplayChannel,
			delete:     session.DeleteChannel,
			defineParams: map[string]any{
				"channel_type":   "SVRCONN",
				"transport_type": "TCP",
				"description":    "dev test channel",
			},
			alter: session.AlterChannel,
			alterParams: map[string]any{
				"channel_type": "SVRCONN",
				"description":  "dev test channel updated",
			},
			alterDescription: "dev test channel updated",
		},
		{
			name:       "listener",
			objectName: testListener,
			define:     session.DefineListener,
			display:    session.DisplayListener,
			delete:     session.DeleteListener,
			defineParams: map[string]any{
				"transport_type": "TCP",
				"port":           1416,
				"start_mode":     "QMGR",
				"description":    "dev test listener",
			},
			alter: session.AlterListener,
			alterParams: map[string]any{
				"transport_type": "TCP",
				"description":    "dev test listener updated",
			},
			alterDescription: "dev test listener updated",
		},
		{
			name:       "process",
			objectName: testProcess,
			define:     session.DefineProcess,
			display:    session.DisplayProcess,
			delete:     session.DeleteProcess,
			defineParams: map[string]any{
				"application_id": "/bin/true",
				"description":    "dev test process",
			},
			alter: session.AlterProcess,
			alterParams: map[string]any{
				"description": "dev test process updated",
			},
			alterDescription: "dev test process updated",
		},
		{
			name:       "topic",
			objectName: testTopic,
			define:     session.DefineTopic,
			display:    session.DisplayTopic,
			delete:     session.DeleteTopic,
			defineParams: map[string]any{
				"topic_string": "dev/test",
				"description":  "dev test topic",
			},
			alter: session.AlterTopic,
			alterParams: map[string]any{
				"description": "dev test topic updated",
			},
			alterDescription: "dev test topic updated",
		},
		{
			name:       "namelist",
			objectName: testNamelist,
			define:     session.DefineNamelist,
			display:    session.DisplayNamelist,
			delete:     session.DeleteNamelist,
			defineParams: map[string]any{
				"names":       []string{"DEV.QLOCAL"},
				"description": "dev test namelist",
			},
			alter: session.AlterNamelist,
			alterParams: map[string]any{
				"description": "dev test namelist updated",
			},
			alterDescription: "dev test namelist updated",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runLifecycleCase(t, ctx, tc)
		})
	}
}

func runLifecycleCase(t *testing.T, ctx context.Context, tc lifecycleCase) {
	t.Helper()

	// Clean up from any prior failed run.
	silentDelete(func() error { return tc.delete(ctx, tc.objectName) })

	// Define.
	err := tc.define(ctx, tc.objectName, mqrestadmin.WithRequestParameters(tc.defineParams))
	if err != nil {
		t.Fatalf("define %s: %v", tc.objectName, err)
	}

	// Display and verify.
	results, err := tc.display(ctx, tc.objectName)
	if err != nil {
		t.Fatalf("display after define %s: %v", tc.objectName, err)
	}
	if !anyContainsValue(results, tc.objectName) {
		t.Errorf("display after define: results do not contain %q", tc.objectName)
	}

	// Alter (if applicable).
	if tc.alter != nil {
		verifyAlter(t, ctx, tc)
	}

	// Delete.
	err = tc.delete(ctx, tc.objectName)
	if err != nil {
		t.Fatalf("delete %s: %v", tc.objectName, err)
	}

	// Verify deletion.
	deleted, err := tc.display(ctx, tc.objectName)
	if err != nil {
		// Error on display after delete is acceptable (object not found).
		return
	}
	if anyContainsValue(deleted, tc.objectName) {
		t.Errorf("object %q still visible after delete", tc.objectName)
	}
}

func verifyAlter(t *testing.T, ctx context.Context, tc lifecycleCase) {
	t.Helper()

	err := tc.alter(ctx, tc.objectName, mqrestadmin.WithRequestParameters(tc.alterParams))
	if err != nil {
		t.Fatalf("alter %s: %v", tc.objectName, err)
	}

	updated, err := tc.display(ctx, tc.objectName)
	if err != nil {
		t.Fatalf("display after alter %s: %v", tc.objectName, err)
	}

	if tc.alterDescription == "" {
		return
	}

	matched := findMatchingObject(updated, tc.objectName)
	if matched == nil {
		t.Fatalf("display after alter: could not find %q", tc.objectName)
	}

	desc, found := getAttributeCaseInsensitive(matched, "description")
	if !found {
		desc, found = getAttributeCaseInsensitive(matched, "DESCR")
	}

	if descStr, ok := desc.(string); ok {
		if descStr != tc.alterDescription {
			t.Errorf("alter description: got %q, want %q", descStr, tc.alterDescription)
		}
	} else if found {
		t.Errorf("alter description: unexpected type %T", desc)
	}
}

// ---------------------------------------------------------------------------
// Ensure idempotent tests
// ---------------------------------------------------------------------------

func TestEnsureQmgrLifecycle(t *testing.T) {
	skipIfLifecycleDisabled(t)
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	// Read current description so we can restore it.
	qmgr, err := session.DisplayQmgr(ctx)
	if err != nil {
		t.Fatalf("DisplayQmgr: %v", err)
	}
	originalDescr, _ := qmgr["description"].(string)

	testDescr := "dev ensure_qmgr test"

	// Alter to test value.
	result, err := session.EnsureQmgr(ctx, map[string]any{"description": testDescr})
	if err != nil {
		t.Fatalf("EnsureQmgr (set): %v", err)
	}
	if result.Action != mqrestadmin.EnsureUpdated && result.Action != mqrestadmin.EnsureUnchanged {
		t.Errorf("EnsureQmgr (set): got %v, want Updated or Unchanged", result.Action)
	}

	// Idempotent — same attributes should be unchanged.
	result, err = session.EnsureQmgr(ctx, map[string]any{"description": testDescr})
	if err != nil {
		t.Fatalf("EnsureQmgr (unchanged): %v", err)
	}
	if result.Action != mqrestadmin.EnsureUnchanged {
		t.Errorf("EnsureQmgr (unchanged): got %v, want Unchanged", result.Action)
	}

	// Restore original description.
	_, err = session.EnsureQmgr(ctx, map[string]any{"description": originalDescr})
	if err != nil {
		t.Fatalf("EnsureQmgr (restore): %v", err)
	}
}

func TestEnsureQlocalLifecycle(t *testing.T) {
	skipIfLifecycleDisabled(t)
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	// Clean up from any prior failed run.
	silentDelete(func() error { return session.DeleteQueue(ctx, testEnsureQlocal) })

	// Create.
	result, err := session.EnsureQlocal(ctx, testEnsureQlocal, map[string]any{"description": "ensure test"})
	if err != nil {
		t.Fatalf("EnsureQlocal (create): %v", err)
	}
	if result.Action != mqrestadmin.EnsureCreated {
		t.Errorf("EnsureQlocal (create): got %v, want Created", result.Action)
	}

	// Unchanged (same attributes).
	result, err = session.EnsureQlocal(ctx, testEnsureQlocal, map[string]any{"description": "ensure test"})
	if err != nil {
		t.Fatalf("EnsureQlocal (unchanged): %v", err)
	}
	if result.Action != mqrestadmin.EnsureUnchanged {
		t.Errorf("EnsureQlocal (unchanged): got %v, want Unchanged", result.Action)
	}

	// Updated (different attribute).
	result, err = session.EnsureQlocal(ctx, testEnsureQlocal, map[string]any{"description": "ensure updated"})
	if err != nil {
		t.Fatalf("EnsureQlocal (update): %v", err)
	}
	if result.Action != mqrestadmin.EnsureUpdated {
		t.Errorf("EnsureQlocal (update): got %v, want Updated", result.Action)
	}

	// Cleanup.
	if err := session.DeleteQueue(ctx, testEnsureQlocal); err != nil {
		t.Fatalf("cleanup delete %s: %v", testEnsureQlocal, err)
	}
}

func TestEnsureChannelLifecycle(t *testing.T) {
	skipIfLifecycleDisabled(t)
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	// Clean up from any prior failed run.
	silentDelete(func() error { return session.DeleteChannel(ctx, testEnsureChannel) })

	// Create.
	result, err := session.EnsureChannel(ctx, testEnsureChannel, map[string]any{
		"channel_type": "SVRCONN",
		"description":  "ensure test",
	})
	if err != nil {
		t.Fatalf("EnsureChannel (create): %v", err)
	}
	if result.Action != mqrestadmin.EnsureCreated {
		t.Errorf("EnsureChannel (create): got %v, want Created", result.Action)
	}

	// Unchanged.
	result, err = session.EnsureChannel(ctx, testEnsureChannel, map[string]any{
		"channel_type": "SVRCONN",
		"description":  "ensure test",
	})
	if err != nil {
		t.Fatalf("EnsureChannel (unchanged): %v", err)
	}
	if result.Action != mqrestadmin.EnsureUnchanged {
		t.Errorf("EnsureChannel (unchanged): got %v, want Unchanged", result.Action)
	}

	// Updated.
	result, err = session.EnsureChannel(ctx, testEnsureChannel, map[string]any{
		"channel_type": "SVRCONN",
		"description":  "ensure updated",
	})
	if err != nil {
		t.Fatalf("EnsureChannel (update): %v", err)
	}
	if result.Action != mqrestadmin.EnsureUpdated {
		t.Errorf("EnsureChannel (update): got %v, want Updated", result.Action)
	}

	// Cleanup.
	if err := session.DeleteChannel(ctx, testEnsureChannel); err != nil {
		t.Fatalf("cleanup delete %s: %v", testEnsureChannel, err)
	}
}

// ---------------------------------------------------------------------------
// Gateway multi-QM tests
// ---------------------------------------------------------------------------

func TestGatewayDisplayQmgrQM2ViaQM1(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildGatewaySession(t, cfg, cfg.qmgrNameQM2, cfg.qmgrName, cfg.restBaseURL)
	ctx := context.Background()

	result, err := session.DisplayQmgr(ctx)
	if err != nil {
		t.Fatalf("DisplayQmgr (QM2 via QM1): %v", err)
	}
	if result == nil {
		t.Fatal("DisplayQmgr (QM2 via QM1) returned nil")
	}
	if !containsStringValue(result, cfg.qmgrNameQM2) {
		t.Errorf("DisplayQmgr (QM2 via QM1) does not contain %q", cfg.qmgrNameQM2)
	}
}

func TestGatewayDisplayQmgrQM1ViaQM2(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildGatewaySession(t, cfg, cfg.qmgrName, cfg.qmgrNameQM2, cfg.restBaseURLQM2)
	ctx := context.Background()

	result, err := session.DisplayQmgr(ctx)
	if err != nil {
		t.Fatalf("DisplayQmgr (QM1 via QM2): %v", err)
	}
	if result == nil {
		t.Fatal("DisplayQmgr (QM1 via QM2) returned nil")
	}
	if !containsStringValue(result, cfg.qmgrName) {
		t.Errorf("DisplayQmgr (QM1 via QM2) does not contain %q", cfg.qmgrName)
	}
}

func TestGatewayDisplayQueueQM2ViaQM1(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildGatewaySession(t, cfg, cfg.qmgrNameQM2, cfg.qmgrName, cfg.restBaseURL)
	ctx := context.Background()

	results, err := session.DisplayQueue(ctx, "DEV.QLOCAL")
	if err != nil {
		t.Fatalf("DisplayQueue (QM2 via QM1): %v", err)
	}
	if len(results) == 0 {
		t.Fatal("DisplayQueue (QM2 via QM1) returned empty results")
	}
	if !anyContainsValue(results, "DEV.QLOCAL") {
		t.Error("DisplayQueue (QM2 via QM1) results do not contain DEV.QLOCAL")
	}
}

func TestGatewaySessionProperties(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildGatewaySession(t, cfg, cfg.qmgrNameQM2, cfg.qmgrName, cfg.restBaseURL)

	if got := session.QmgrName(); got != cfg.qmgrNameQM2 {
		t.Errorf("QmgrName(): got %q, want %q", got, cfg.qmgrNameQM2)
	}
	if got := session.GatewayQmgr(); got != cfg.qmgrName {
		t.Errorf("GatewayQmgr(): got %q, want %q", got, cfg.qmgrName)
	}
}

// ---------------------------------------------------------------------------
// Session state test
// ---------------------------------------------------------------------------

func TestSessionStatePopulatedAfterCommand(t *testing.T) {
	cfg := loadIntegrationConfig()
	session := buildSession(t, cfg)
	ctx := context.Background()

	_, err := session.DisplayQmgr(ctx)
	if err != nil {
		t.Fatalf("DisplayQmgr: %v", err)
	}

	if session.LastHTTPStatus == 0 {
		t.Error("LastHTTPStatus is 0 after command")
	}
	if session.LastResponseText == "" {
		t.Error("LastResponseText is empty after command")
	}
}

// ---------------------------------------------------------------------------
// LTPA auth test (expected to fail on dev containers)
// ---------------------------------------------------------------------------

func TestLTPAAuthDisplayQmgr(t *testing.T) {
	cfg := loadIntegrationConfig()

	session, err := mqrestadmin.NewSession(
		cfg.restBaseURL,
		cfg.qmgrName,
		mqrestadmin.LTPAAuth{Username: cfg.adminUser, Password: cfg.adminPassword},
		mqrestadmin.WithVerifyTLS(cfg.verifyTLS),
	)
	if err != nil {
		t.Skipf("LTPA session creation failed (expected on dev containers): %v", err)
		return
	}

	result, err := session.DisplayQmgr(context.Background())
	if err != nil {
		t.Skipf("LTPA DisplayQmgr failed (expected on dev containers): %v", err)
		return
	}

	if result == nil {
		t.Error("LTPA DisplayQmgr returned nil")
	}
	if !containsStringValue(result, cfg.qmgrName) {
		t.Errorf("LTPA DisplayQmgr does not contain %q", cfg.qmgrName)
	}
}
