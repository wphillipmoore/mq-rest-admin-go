package examples

import (
	"context"
	"testing"

	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

// ---------------------------------------------------------------------------
// healthcheck.go
// ---------------------------------------------------------------------------

func TestCheckHealth_HappyPath(t *testing.T) {
	transport := &mockTransport{}
	// DisplayQmgr
	transport.addSuccessResponse(map[string]any{"queue_manager_name": "QM1"})
	// DisplayQmstatus
	transport.addSuccessResponse(map[string]any{"ha_status": "ACTIVE"})
	// DisplayCmdserv
	transport.addSuccessResponse(map[string]any{"status": "RUNNING"})
	// DisplayListener
	transport.addSuccessResponse(
		map[string]any{"listener_name": "LISTENER.TCP", "start_mode": "MANUAL"},
		map[string]any{"listener_name": "LISTENER.LU62", "start_mode": "AUTO"},
	)

	session := newTestSession(t, transport)
	result := CheckHealth(context.Background(), session)

	if !result.Reachable {
		t.Error("expected Reachable to be true")
	}
	if result.Status != "ACTIVE" {
		t.Errorf("Status = %q, want %q", result.Status, "ACTIVE")
	}
	if result.CommandServer != "RUNNING" {
		t.Errorf("CommandServer = %q, want %q", result.CommandServer, "RUNNING")
	}
	if len(result.Listeners) != 2 {
		t.Fatalf("len(Listeners) = %d, want 2", len(result.Listeners))
	}
	if result.Listeners[0].Name != "LISTENER.TCP" {
		t.Errorf("Listener[0].Name = %q, want %q", result.Listeners[0].Name, "LISTENER.TCP")
	}
	if !result.Passed {
		t.Error("expected Passed to be true")
	}
}

func TestCheckHealth_DisplayQmgrError(t *testing.T) {
	transport := &mockTransport{}
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	result := CheckHealth(context.Background(), session)

	if result.Reachable {
		t.Error("expected Reachable to be false")
	}
	if result.Status != "UNKNOWN" {
		t.Errorf("Status = %q, want %q", result.Status, "UNKNOWN")
	}
	if result.Passed {
		t.Error("expected Passed to be false")
	}
}

func TestCheckHealth_QmstatusAndCmdservError(t *testing.T) {
	transport := &mockTransport{}
	// DisplayQmgr succeeds
	transport.addSuccessResponse(map[string]any{"queue_manager_name": "QM1"})
	// DisplayQmstatus fails
	transport.addCommandErrorResponse(2, 2085)
	// DisplayCmdserv fails
	transport.addCommandErrorResponse(2, 2085)
	// DisplayListener returns empty
	transport.addSuccessResponse()

	session := newTestSession(t, transport)
	result := CheckHealth(context.Background(), session)

	if !result.Reachable {
		t.Error("expected Reachable to be true")
	}
	if result.Status != "UNKNOWN" {
		t.Errorf("Status = %q, want %q", result.Status, "UNKNOWN")
	}
	if result.CommandServer != "UNKNOWN" {
		t.Errorf("CommandServer = %q, want %q", result.CommandServer, "UNKNOWN")
	}
	if len(result.Listeners) != 0 {
		t.Errorf("len(Listeners) = %d, want 0", len(result.Listeners))
	}
	if result.Passed {
		t.Error("expected Passed to be false")
	}
}

func TestCheckHealth_NoListeners(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(map[string]any{"queue_manager_name": "QM1"})
	transport.addSuccessResponse(map[string]any{"ha_status": "ACTIVE"})
	transport.addSuccessResponse(map[string]any{"status": "RUNNING"})
	// DisplayListener returns no entries
	transport.addSuccessResponse()

	session := newTestSession(t, transport)
	result := CheckHealth(context.Background(), session)

	if len(result.Listeners) != 0 {
		t.Errorf("len(Listeners) = %d, want 0", len(result.Listeners))
	}
	if !result.Passed {
		t.Error("expected Passed to be true")
	}
}

func TestCheckHealth_ListenerError(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(map[string]any{"queue_manager_name": "QM1"})
	transport.addSuccessResponse(map[string]any{"ha_status": "ACTIVE"})
	transport.addSuccessResponse(map[string]any{"status": "RUNNING"})
	// DisplayListener error
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	result := CheckHealth(context.Background(), session)

	if len(result.Listeners) != 0 {
		t.Errorf("len(Listeners) = %d, want 0", len(result.Listeners))
	}
	if !result.Passed {
		t.Error("expected Passed to be true")
	}
}

func TestPrintHealthCheck_PassAndFail(t *testing.T) {
	transport := &mockTransport{}
	// Session 1: healthy
	transport.addSuccessResponse(map[string]any{"queue_manager_name": "QM1"})
	transport.addSuccessResponse(map[string]any{"ha_status": "ACTIVE"})
	transport.addSuccessResponse(map[string]any{"status": "RUNNING"})
	transport.addSuccessResponse(map[string]any{"listener_name": "L1", "start_mode": "AUTO"})

	session1 := newTestSession(t, transport)

	transport2 := &mockTransport{}
	// Session 2: unreachable
	transport2.addCommandErrorResponse(2, 2085)
	session2 := newTestSession(t, transport2)

	results := PrintHealthCheck(context.Background(), []*mqrestadmin.Session{session1, session2})

	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	if !results[0].Passed {
		t.Error("expected first result to pass")
	}
	if results[1].Passed {
		t.Error("expected second result to fail")
	}
}

// ---------------------------------------------------------------------------
// channelstatus.go
// ---------------------------------------------------------------------------

func TestReportChannelStatus_HappyPath(t *testing.T) {
	transport := &mockTransport{}
	// DisplayChannel
	transport.addSuccessResponse(
		map[string]any{"channel_name": "CHAN.A", "channel_type": "SDR", "connection_name": "host(1414)"},
		map[string]any{"channel_name": "CHAN.B", "channel_type": "RCVR", "connection_name": ""},
	)
	// DisplayChstatus
	transport.addSuccessResponse(
		map[string]any{"channel_name": "CHAN.A", "status": "RUNNING"},
	)

	session := newTestSession(t, transport)
	results, err := ReportChannelStatus(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	// CHAN.A is first alphabetically
	if results[0].Name != "CHAN.A" || results[0].Status != "RUNNING" || !results[0].Defined {
		t.Errorf("CHAN.A result unexpected: %+v", results[0])
	}
	if results[1].Name != "CHAN.B" || results[1].Status != "INACTIVE" || !results[1].Defined {
		t.Errorf("CHAN.B result unexpected: %+v", results[1])
	}
}

func TestReportChannelStatus_DisplayChannelError(t *testing.T) {
	transport := &mockTransport{}
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	_, err := ReportChannelStatus(context.Background(), session)
	if err == nil {
		t.Error("expected error from DisplayChannel")
	}
}

func TestReportChannelStatus_LiveStatusNoDefinition(t *testing.T) {
	transport := &mockTransport{}
	// DisplayChannel returns one channel
	transport.addSuccessResponse(
		map[string]any{"channel_name": "CHAN.A", "channel_type": "SDR", "connection_name": "host(1414)"},
	)
	// DisplayChstatus returns status for CHAN.A and an extra CHAN.X
	transport.addSuccessResponse(
		map[string]any{"channel_name": "CHAN.A", "status": "RUNNING"},
		map[string]any{"channel_name": "CHAN.X", "status": "RETRYING"},
	)

	session := newTestSession(t, transport)
	results, err := ReportChannelStatus(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	// CHAN.X should appear with Defined=false
	found := false
	for _, r := range results {
		if r.Name == "CHAN.X" {
			found = true
			if r.Defined {
				t.Error("CHAN.X should have Defined=false")
			}
			if r.Status != "RETRYING" {
				t.Errorf("CHAN.X Status = %q, want %q", r.Status, "RETRYING")
			}
		}
	}
	if !found {
		t.Error("CHAN.X not in results")
	}
}

func TestReportChannelStatus_ChstatusError(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(
		map[string]any{"channel_name": "CHAN.A", "channel_type": "SDR", "connection_name": "host(1414)"},
	)
	// DisplayChstatus error — channels should still appear with INACTIVE status
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	results, err := ReportChannelStatus(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].Status != "INACTIVE" {
		t.Errorf("Status = %q, want %q", results[0].Status, "INACTIVE")
	}
}

func TestPrintChannelStatus_InactiveChannels(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(
		map[string]any{"channel_name": "CHAN.A", "channel_type": "SDR", "connection_name": "host(1414)"},
		map[string]any{"channel_name": "CHAN.B", "channel_type": "RCVR", "connection_name": ""},
	)
	// No live statuses — both inactive
	transport.addSuccessResponse()

	session := newTestSession(t, transport)
	results, err := PrintChannelStatus(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	for _, r := range results {
		if r.Status != "INACTIVE" {
			t.Errorf("%s Status = %q, want %q", r.Name, r.Status, "INACTIVE")
		}
	}
}

func TestPrintChannelStatus_Error(t *testing.T) {
	transport := &mockTransport{}
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	_, err := PrintChannelStatus(context.Background(), session)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------------------------------------------------------------------------
// depthmonitor.go
// ---------------------------------------------------------------------------

func TestMonitorQueueDepths_HappyPath(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(
		map[string]any{"queue_name": "Q.LOCAL.1", "type": "QLOCAL", "current_queue_depth": float64(50), "max_queue_depth": float64(100), "open_input_count": float64(1), "open_output_count": float64(2)},
		map[string]any{"queue_name": "Q.ALIAS", "type": "QALIAS", "current_queue_depth": float64(0), "max_queue_depth": float64(100)},
		map[string]any{"queue_name": "Q.LOCAL.2", "type": "LOCAL", "current_queue_depth": float64(10), "max_queue_depth": float64(200), "open_input_count": float64(0), "open_output_count": float64(0)},
	)

	session := newTestSession(t, transport)
	results, err := MonitorQueueDepths(context.Background(), session, 40.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only QLOCAL and LOCAL, not QALIAS
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	// Sorted by depth% descending: Q.LOCAL.1 (50%) before Q.LOCAL.2 (5%)
	if results[0].Name != "Q.LOCAL.1" {
		t.Errorf("first result Name = %q, want %q", results[0].Name, "Q.LOCAL.1")
	}
	if results[0].DepthPct != 50.0 {
		t.Errorf("DepthPct = %f, want 50.0", results[0].DepthPct)
	}
	if !results[0].Warning {
		t.Error("Q.LOCAL.1 should have Warning=true at 40% threshold")
	}
	if results[1].Warning {
		t.Error("Q.LOCAL.2 should have Warning=false at 40% threshold")
	}
}

func TestMonitorQueueDepths_Error(t *testing.T) {
	transport := &mockTransport{}
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	_, err := MonitorQueueDepths(context.Background(), session, 80.0)
	if err == nil {
		t.Error("expected error from DisplayQueue")
	}
}

func TestMonitorQueueDepths_ZeroMaxDepth(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(
		map[string]any{"queue_name": "Q.ZERO", "type": "QLOCAL", "current_queue_depth": float64(0), "max_queue_depth": float64(0)},
	)

	session := newTestSession(t, transport)
	results, err := MonitorQueueDepths(context.Background(), session, 80.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].DepthPct != 0.0 {
		t.Errorf("DepthPct = %f, want 0.0", results[0].DepthPct)
	}
}

func TestPrintQueueDepths_HappyPath(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(
		map[string]any{"queue_name": "Q.LOCAL", "type": "QLOCAL", "current_queue_depth": float64(90), "max_queue_depth": float64(100), "open_input_count": float64(1), "open_output_count": float64(0)},
	)

	session := newTestSession(t, transport)
	results, err := PrintQueueDepths(context.Background(), session, 80.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if !results[0].Warning {
		t.Error("expected Warning at 80% threshold with 90% depth")
	}
}

func TestPrintQueueDepths_Error(t *testing.T) {
	transport := &mockTransport{}
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	_, err := PrintQueueDepths(context.Background(), session, 80.0)
	if err == nil {
		t.Error("expected error")
	}
}

func TestToInt_Int(t *testing.T) {
	if got := toInt(42); got != 42 {
		t.Errorf("toInt(42) = %d, want 42", got)
	}
}

func TestToInt_Float64(t *testing.T) {
	if got := toInt(float64(3.7)); got != 3 {
		t.Errorf("toInt(3.7) = %d, want 3", got)
	}
}

func TestToInt_String(t *testing.T) {
	if got := toInt("123"); got != 123 {
		t.Errorf("toInt(\"123\") = %d, want 123", got)
	}
}

func TestToInt_StringInvalid(t *testing.T) {
	if got := toInt("abc"); got != 0 {
		t.Errorf("toInt(\"abc\") = %d, want 0", got)
	}
}

func TestToInt_Default(t *testing.T) {
	if got := toInt(nil); got != 0 {
		t.Errorf("toInt(nil) = %d, want 0", got)
	}
}

func TestToInt_DefaultValidSprint(t *testing.T) {
	// int64 falls through to default case (not int, float64, or string)
	if got := toInt(int64(99)); got != 99 {
		t.Errorf("toInt(int64(99)) = %d, want 99", got)
	}
}

// ---------------------------------------------------------------------------
// dlqinspector.go
// ---------------------------------------------------------------------------

func TestInspectDLQ_HappyPath_Empty(t *testing.T) {
	transport := &mockTransport{}
	// DisplayQmgr
	transport.addSuccessResponse(map[string]any{"dead_letter_queue_name": "SYSTEM.DEAD.LETTER.QUEUE"})
	// DisplayQueue
	transport.addSuccessResponse(map[string]any{
		"queue_name": "SYSTEM.DEAD.LETTER.QUEUE", "current_queue_depth": float64(0),
		"max_queue_depth": float64(5000), "open_input_count": float64(0), "open_output_count": float64(0),
	})

	session := newTestSession(t, transport)
	report, err := InspectDLQ(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Configured {
		t.Error("expected Configured=true")
	}
	if report.CurrentDepth != 0 {
		t.Errorf("CurrentDepth = %d, want 0", report.CurrentDepth)
	}
	if report.Suggestion != "DLQ is empty. No action needed." {
		t.Errorf("Suggestion = %q", report.Suggestion)
	}
}

func TestInspectDLQ_DisplayQmgrError(t *testing.T) {
	transport := &mockTransport{}
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	_, err := InspectDLQ(context.Background(), session)
	if err == nil {
		t.Error("expected error from DisplayQmgr")
	}
}

func TestInspectDLQ_NoDLQConfigured(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(map[string]any{"dead_letter_queue_name": ""})

	session := newTestSession(t, transport)
	report, err := InspectDLQ(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Configured {
		t.Error("expected Configured=false")
	}
	if report.Suggestion == "" {
		t.Error("expected non-empty Suggestion")
	}
}

func TestInspectDLQ_NoDLQConfigured_Nil(t *testing.T) {
	transport := &mockTransport{}
	// When dead_letter_queue_name is nil, Sprint produces "<nil>"
	transport.addSuccessResponse(map[string]any{})

	session := newTestSession(t, transport)
	report, err := InspectDLQ(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Configured {
		t.Error("expected Configured=false for nil DLQ name")
	}
}

func TestInspectDLQ_QueueNotFound(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(map[string]any{"dead_letter_queue_name": "MY.DLQ"})
	// DisplayQueue returns empty
	transport.addSuccessResponse()

	session := newTestSession(t, transport)
	report, err := InspectDLQ(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Configured {
		t.Error("expected Configured=true")
	}
	if report.Suggestion == "" {
		t.Error("expected non-empty Suggestion about queue not existing")
	}
}

func TestInspectDLQ_WithMessages(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(map[string]any{"dead_letter_queue_name": "MY.DLQ"})
	transport.addSuccessResponse(map[string]any{
		"queue_name": "MY.DLQ", "current_queue_depth": float64(15),
		"max_queue_depth": float64(5000), "open_input_count": float64(0), "open_output_count": float64(0),
	})

	session := newTestSession(t, transport)
	report, err := InspectDLQ(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.CurrentDepth != 15 {
		t.Errorf("CurrentDepth = %d, want 15", report.CurrentDepth)
	}
	if report.Suggestion != "DLQ has messages. Investigate undeliverable messages." {
		t.Errorf("Suggestion = %q", report.Suggestion)
	}
}

func TestInspectDLQ_NearCapacity(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(map[string]any{"dead_letter_queue_name": "MY.DLQ"})
	transport.addSuccessResponse(map[string]any{
		"queue_name": "MY.DLQ", "current_queue_depth": float64(95),
		"max_queue_depth": float64(100), "open_input_count": float64(0), "open_output_count": float64(0),
	})

	session := newTestSession(t, transport)
	report, err := InspectDLQ(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.DepthPct != 95.0 {
		t.Errorf("DepthPct = %f, want 95.0", report.DepthPct)
	}
	if report.Suggestion != "DLQ is near capacity. Investigate and clear undeliverable messages urgently." {
		t.Errorf("Suggestion = %q", report.Suggestion)
	}
}

func TestPrintDLQInspection_Configured(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(map[string]any{"dead_letter_queue_name": "MY.DLQ"})
	transport.addSuccessResponse(map[string]any{
		"queue_name": "MY.DLQ", "current_queue_depth": float64(0),
		"max_queue_depth": float64(5000), "open_input_count": float64(0), "open_output_count": float64(0),
	})

	session := newTestSession(t, transport)
	report, err := PrintDLQInspection(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Configured {
		t.Error("expected Configured=true")
	}
}

func TestPrintDLQInspection_NotConfigured(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(map[string]any{"dead_letter_queue_name": ""})

	session := newTestSession(t, transport)
	report, err := PrintDLQInspection(context.Background(), session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Configured {
		t.Error("expected Configured=false")
	}
}

func TestPrintDLQInspection_Error(t *testing.T) {
	transport := &mockTransport{}
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	_, err := PrintDLQInspection(context.Background(), session)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------------------------------------------------------------------------
// queuestatus.go
// ---------------------------------------------------------------------------

func TestReportQueueHandles_HappyPath(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(
		map[string]any{"queue_name": "Q1", "handle_state": "ACTIVE", "connection_id": "CONN1", "open_options": "INPUT"},
	)

	session := newTestSession(t, transport)
	results := ReportQueueHandles(context.Background(), session)
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].QueueName != "Q1" {
		t.Errorf("QueueName = %q, want %q", results[0].QueueName, "Q1")
	}
}

func TestReportQueueHandles_Error(t *testing.T) {
	transport := &mockTransport{}
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	results := ReportQueueHandles(context.Background(), session)
	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}

func TestReportQueueHandles_Empty(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse()

	session := newTestSession(t, transport)
	results := ReportQueueHandles(context.Background(), session)
	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}

func TestReportConnectionHandles_HappyPath(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(
		map[string]any{"connection_id": "CONN1", "object_name": "Q1", "handle_state": "ACTIVE", "object_type": "QUEUE"},
	)

	session := newTestSession(t, transport)
	results := ReportConnectionHandles(context.Background(), session)
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].ConnectionID != "CONN1" {
		t.Errorf("ConnectionID = %q, want %q", results[0].ConnectionID, "CONN1")
	}
}

func TestReportConnectionHandles_Error(t *testing.T) {
	transport := &mockTransport{}
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSession(t, transport)
	results := ReportConnectionHandles(context.Background(), session)
	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}

func TestReportConnectionHandles_Empty(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse()

	session := newTestSession(t, transport)
	results := ReportConnectionHandles(context.Background(), session)
	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}

func TestPrintQueueStatus_EmptyHandles(t *testing.T) {
	transport := &mockTransport{}
	// DisplayQstatus returns empty
	transport.addSuccessResponse()
	// DisplayConn returns empty
	transport.addSuccessResponse()

	session := newTestSession(t, transport)
	// Should not panic
	PrintQueueStatus(context.Background(), session)
}

func TestPrintQueueStatus_WithHandles(t *testing.T) {
	transport := &mockTransport{}
	transport.addSuccessResponse(
		map[string]any{"queue_name": "Q1", "handle_state": "ACTIVE", "connection_id": "CONN1", "open_options": "INPUT"},
	)
	transport.addSuccessResponse(
		map[string]any{"connection_id": "CONN1", "object_name": "Q1", "handle_state": "ACTIVE", "object_type": "QUEUE"},
	)

	session := newTestSession(t, transport)
	PrintQueueStatus(context.Background(), session)
}

// ---------------------------------------------------------------------------
// provisionenv.go
// ---------------------------------------------------------------------------

func TestProvision_HappyPath(t *testing.T) {
	t1 := &mockTransport{}
	t2 := &mockTransport{}

	// QM1: 5 defines + 1 DisplayQueue for verification
	for range 5 {
		t1.addSuccessResponse()
	}
	t1.addSuccessResponse(
		map[string]any{"queue_name": "PROV.QM1.LOCAL"},
		map[string]any{"queue_name": "PROV.QM1.TO.QM2.XMITQ"},
		map[string]any{"queue_name": "PROV.REMOTE.TO.QM2"},
	)

	// QM2: 5 defines + 1 DisplayQueue for verification
	for range 5 {
		t2.addSuccessResponse()
	}
	t2.addSuccessResponse(
		map[string]any{"queue_name": "PROV.QM2.LOCAL"},
		map[string]any{"queue_name": "PROV.QM2.TO.QM1.XMITQ"},
		map[string]any{"queue_name": "PROV.REMOTE.TO.QM1"},
	)

	qm1 := newTestSession(t, t1)
	qm2 := newTestSession(t, t2)
	result := Provision(context.Background(), qm1, qm2)

	if len(result.ObjectsCreated) != 10 {
		t.Errorf("ObjectsCreated = %d, want 10", len(result.ObjectsCreated))
	}
	if len(result.ObjectsFailed) != 0 {
		t.Errorf("ObjectsFailed = %d, want 0", len(result.ObjectsFailed))
	}
	if !result.Verified {
		t.Error("expected Verified=true")
	}
}

func TestProvision_DefineFailures(t *testing.T) {
	t1 := &mockTransport{}
	t2 := &mockTransport{}

	// QM1: first define fails, rest succeed
	t1.addCommandErrorResponse(2, 2085)
	for range 4 {
		t1.addSuccessResponse()
	}
	t1.addSuccessResponse(
		map[string]any{"queue_name": "PROV.QM1.LOCAL"},
		map[string]any{"queue_name": "PROV.QM1.TO.QM2.XMITQ"},
		map[string]any{"queue_name": "PROV.REMOTE.TO.QM2"},
	)

	for range 5 {
		t2.addSuccessResponse()
	}
	t2.addSuccessResponse(
		map[string]any{"queue_name": "PROV.QM2.LOCAL"},
		map[string]any{"queue_name": "PROV.QM2.TO.QM1.XMITQ"},
		map[string]any{"queue_name": "PROV.REMOTE.TO.QM1"},
	)

	qm1 := newTestSession(t, t1)
	qm2 := newTestSession(t, t2)
	result := Provision(context.Background(), qm1, qm2)

	if len(result.ObjectsFailed) != 1 {
		t.Errorf("ObjectsFailed = %d, want 1", len(result.ObjectsFailed))
	}
	if len(result.ObjectsCreated) != 9 {
		t.Errorf("ObjectsCreated = %d, want 9", len(result.ObjectsCreated))
	}
}

func TestProvision_VerificationFailure(t *testing.T) {
	t1 := &mockTransport{}
	t2 := &mockTransport{}

	for range 5 {
		t1.addSuccessResponse()
	}
	// Verification returns fewer than 3 queues
	t1.addSuccessResponse(map[string]any{"queue_name": "PROV.QM1.LOCAL"})

	for range 5 {
		t2.addSuccessResponse()
	}
	t2.addSuccessResponse(
		map[string]any{"queue_name": "PROV.QM2.LOCAL"},
		map[string]any{"queue_name": "PROV.QM2.TO.QM1.XMITQ"},
		map[string]any{"queue_name": "PROV.REMOTE.TO.QM1"},
	)

	qm1 := newTestSession(t, t1)
	qm2 := newTestSession(t, t2)
	result := Provision(context.Background(), qm1, qm2)

	if result.Verified {
		t.Error("expected Verified=false when fewer than 3 queues found")
	}
}

func TestDefineObject_UnknownMethod(t *testing.T) {
	transport := &mockTransport{}
	session := newTestSession(t, transport)
	result := ProvisionResult{}

	defineObject(context.Background(), &result, session, "DefineQmodel", "TEST.Q", map[string]any{})

	if len(result.ObjectsFailed) != 1 {
		t.Errorf("ObjectsFailed = %d, want 1", len(result.ObjectsFailed))
	}
	if len(result.ObjectsCreated) != 0 {
		t.Errorf("ObjectsCreated = %d, want 0", len(result.ObjectsCreated))
	}
}

func TestTeardown_HappyPath(t *testing.T) {
	t1 := &mockTransport{}
	t2 := &mockTransport{}

	// QM1: 5 deletes
	for range 5 {
		t1.addSuccessResponse()
	}
	// QM2: 5 deletes
	for range 5 {
		t2.addSuccessResponse()
	}

	qm1 := newTestSession(t, t1)
	qm2 := newTestSession(t, t2)
	failures := Teardown(context.Background(), qm1, qm2)

	if len(failures) != 0 {
		t.Errorf("failures = %v, want empty", failures)
	}
}

func TestTeardown_WithFailures(t *testing.T) {
	t1 := &mockTransport{}
	t2 := &mockTransport{}

	// QM1: first delete fails, rest succeed
	t1.addCommandErrorResponse(2, 2085)
	for range 4 {
		t1.addSuccessResponse()
	}
	for range 5 {
		t2.addSuccessResponse()
	}

	qm1 := newTestSession(t, t1)
	qm2 := newTestSession(t, t2)
	failures := Teardown(context.Background(), qm1, qm2)

	if len(failures) != 1 {
		t.Errorf("len(failures) = %d, want 1", len(failures))
	}
}

func TestPrintProvision_HappyPath(t *testing.T) {
	t1 := &mockTransport{}
	t2 := &mockTransport{}

	// Provision: 5 defines each + verification
	for range 5 {
		t1.addSuccessResponse()
	}
	t1.addSuccessResponse(
		map[string]any{"queue_name": "PROV.QM1.LOCAL"},
		map[string]any{"queue_name": "PROV.QM1.TO.QM2.XMITQ"},
		map[string]any{"queue_name": "PROV.REMOTE.TO.QM2"},
	)
	for range 5 {
		t2.addSuccessResponse()
	}
	t2.addSuccessResponse(
		map[string]any{"queue_name": "PROV.QM2.LOCAL"},
		map[string]any{"queue_name": "PROV.QM2.TO.QM1.XMITQ"},
		map[string]any{"queue_name": "PROV.REMOTE.TO.QM1"},
	)

	// Teardown: 5 deletes each
	for range 5 {
		t1.addSuccessResponse()
	}
	for range 5 {
		t2.addSuccessResponse()
	}

	qm1 := newTestSession(t, t1)
	qm2 := newTestSession(t, t2)
	result := PrintProvision(context.Background(), qm1, qm2)

	if len(result.ObjectsCreated) != 10 {
		t.Errorf("ObjectsCreated = %d, want 10", len(result.ObjectsCreated))
	}
	if !result.Verified {
		t.Error("expected Verified=true")
	}
}

func TestPrintProvision_WithFailures(t *testing.T) {
	t1 := &mockTransport{}
	t2 := &mockTransport{}

	// Provision: first define fails, rest succeed
	t1.addCommandErrorResponse(2, 2085)
	for range 4 {
		t1.addSuccessResponse()
	}
	t1.addSuccessResponse(
		map[string]any{"queue_name": "PROV.QM1.LOCAL"},
		map[string]any{"queue_name": "PROV.QM1.TO.QM2.XMITQ"},
		map[string]any{"queue_name": "PROV.REMOTE.TO.QM2"},
	)

	for range 5 {
		t2.addSuccessResponse()
	}
	t2.addSuccessResponse(
		map[string]any{"queue_name": "PROV.QM2.LOCAL"},
		map[string]any{"queue_name": "PROV.QM2.TO.QM1.XMITQ"},
		map[string]any{"queue_name": "PROV.REMOTE.TO.QM1"},
	)

	// Teardown: all succeed
	for range 5 {
		t1.addSuccessResponse()
	}
	for range 5 {
		t2.addSuccessResponse()
	}

	qm1 := newTestSession(t, t1)
	qm2 := newTestSession(t, t2)
	result := PrintProvision(context.Background(), qm1, qm2)

	if len(result.ObjectsFailed) != 1 {
		t.Errorf("ObjectsFailed = %d, want 1", len(result.ObjectsFailed))
	}
}

