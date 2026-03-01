package mqrestadmin

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestStartChannelSync_Success(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// START command succeeds
	transport.addSuccessResponse()
	// First poll: not yet running
	transport.addSuccessResponse(map[string]any{
		"CHANNEL": "TO.REMOTE",
		"STATUS":  "STARTING",
	})
	// Second poll: running
	transport.addSuccessResponse(map[string]any{
		"CHANNEL": "TO.REMOTE",
		"STATUS":  "RUNNING",
	})

	session := newTestSessionWithClock(transport, clock)

	result, err := session.StartChannelSync(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Operation != SyncStarted {
		t.Errorf("Operation = %v, want SyncStarted", result.Operation)
	}
	if result.Polls != 2 {
		t.Errorf("Polls = %d, want 2", result.Polls)
	}
}

func TestStopChannelSync_EmptyMeansStopped(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// STOP command succeeds
	transport.addSuccessResponse()
	// Poll: command error (no status = stopped for channels)
	transport.addCommandErrorResponse(2, 2085)

	session := newTestSessionWithClock(transport, clock)

	result, err := session.StopChannelSync(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Operation != SyncStopped {
		t.Errorf("Operation = %v, want SyncStopped", result.Operation)
	}
}

func TestStopListenerSync_StoppedStatus(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// STOP command succeeds
	transport.addSuccessResponse()
	// Poll: stopped status
	transport.addSuccessResponse(map[string]any{
		"LISTENER": "LIS1",
		"STATUS":   "STOPPED",
	})

	session := newTestSessionWithClock(transport, clock)

	result, err := session.StopListenerSync(context.Background(), "LIS1",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Operation != SyncStopped {
		t.Errorf("Operation = %v, want SyncStopped", result.Operation)
	}
}

func TestStartChannelSync_Timeout(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// START command succeeds
	transport.addSuccessResponse()

	// All polls return non-running status — will timeout
	for range 50 {
		transport.addSuccessResponse(map[string]any{
			"CHANNEL": "TO.REMOTE",
			"STATUS":  "STARTING",
		})
	}

	session := newTestSessionWithClock(transport, clock)

	_, err := session.StartChannelSync(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 3 * time.Second, PollInterval: 1 * time.Second})
	if err == nil {
		t.Fatal("expected timeout error")
	}

	var timeoutErr *TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("expected TimeoutError, got %T: %v", err, err)
	}
	if timeoutErr.Operation != SyncStarted {
		t.Errorf("Operation = %v, want SyncStarted", timeoutErr.Operation)
	}
}

func TestRestartChannel_Success(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// STOP command succeeds
	transport.addSuccessResponse()
	// Poll: channel stopped (empty = stopped for channels)
	transport.addCommandErrorResponse(2, 2085)
	// START command succeeds
	transport.addSuccessResponse()
	// Poll: running
	transport.addSuccessResponse(map[string]any{
		"CHANNEL": "TO.REMOTE",
		"STATUS":  "RUNNING",
	})

	session := newTestSessionWithClock(transport, clock)

	result, err := session.RestartChannel(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Operation != SyncRestarted {
		t.Errorf("Operation = %v, want SyncRestarted", result.Operation)
	}
	if result.Polls != 2 {
		t.Errorf("Polls = %d, want 2 (1 stop + 1 start)", result.Polls)
	}
}

func TestSyncConfig_Defaults(t *testing.T) {
	config, err := normalizeSyncConfig(SyncConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.Timeout != defaultSyncTimeout {
		t.Errorf("Timeout = %v, want %v", config.Timeout, defaultSyncTimeout)
	}
	if config.PollInterval != defaultPollInterval {
		t.Errorf("PollInterval = %v, want %v", config.PollInterval, defaultPollInterval)
	}
}

func TestSyncConfig_PreservesExplicitValues(t *testing.T) {
	config, err := normalizeSyncConfig(SyncConfig{
		Timeout:      60 * time.Second,
		PollInterval: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want 60s", config.Timeout)
	}
	if config.PollInterval != 5*time.Second {
		t.Errorf("PollInterval = %v, want 5s", config.PollInterval)
	}
}

func TestSyncConfig_NegativeTimeout(t *testing.T) {
	_, err := normalizeSyncConfig(SyncConfig{Timeout: -1 * time.Second})
	if err == nil {
		t.Fatal("expected error for negative Timeout")
	}
	if !strings.Contains(err.Error(), "timeout must not be negative") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestSyncConfig_NegativePollInterval(t *testing.T) {
	_, err := normalizeSyncConfig(SyncConfig{PollInterval: -1 * time.Second})
	if err == nil {
		t.Fatal("expected error for negative PollInterval")
	}
	if !strings.Contains(err.Error(), "poll interval must not be negative") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestStartChannelSync_NegativeConfigError(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()
	session := newTestSessionWithClock(transport, clock)

	_, err := session.StartChannelSync(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: -1 * time.Second})
	if err == nil {
		t.Fatal("expected error for negative Timeout")
	}
}

func TestStopChannelSync_NegativeConfigError(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()
	session := newTestSessionWithClock(transport, clock)

	_, err := session.StopChannelSync(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: -1 * time.Second})
	if err == nil {
		t.Fatal("expected error for negative Timeout")
	}
}

func TestHasStatus(t *testing.T) {
	tests := []struct {
		name     string
		rows     []map[string]any
		keys     []string
		values   map[string]bool
		expected bool
	}{
		{
			name:     "running status found",
			rows:     []map[string]any{{"STATUS": "RUNNING"}},
			keys:     []string{"STATUS"},
			values:   runningValues,
			expected: true,
		},
		{
			name:     "running status lowercase",
			rows:     []map[string]any{{"STATUS": "running"}},
			keys:     []string{"STATUS"},
			values:   runningValues,
			expected: true,
		},
		{
			name:     "stopped status",
			rows:     []map[string]any{{"STATUS": "STOPPED"}},
			keys:     []string{"STATUS"},
			values:   stoppedValues,
			expected: true,
		},
		{
			name:     "inactive means stopped",
			rows:     []map[string]any{{"STATUS": "INACTIVE"}},
			keys:     []string{"STATUS"},
			values:   stoppedValues,
			expected: true,
		},
		{
			name:     "no match",
			rows:     []map[string]any{{"STATUS": "STARTING"}},
			keys:     []string{"STATUS"},
			values:   runningValues,
			expected: false,
		},
		{
			name:     "empty rows",
			rows:     nil,
			keys:     []string{"STATUS"},
			values:   runningValues,
			expected: false,
		},
		{
			name:     "alternate key name",
			rows:     []map[string]any{{"channel_status": "RUNNING"}},
			keys:     []string{"channel_status", "STATUS"},
			values:   runningValues,
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := hasStatus(test.rows, test.keys, test.values)
			if result != test.expected {
				t.Errorf("hasStatus() = %v, want %v", result, test.expected)
			}
		})
	}
}

func TestStopChannelSync_Timeout(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// STOP command succeeds
	transport.addSuccessResponse()
	// All polls return running status
	for range 50 {
		transport.addSuccessResponse(map[string]any{
			"CHANNEL": "TO.REMOTE",
			"STATUS":  "RUNNING",
		})
	}

	session := newTestSessionWithClock(transport, clock)

	_, err := session.StopChannelSync(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 3 * time.Second, PollInterval: 1 * time.Second})
	if err == nil {
		t.Fatal("expected timeout error")
	}

	var timeoutErr *TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("expected TimeoutError, got %T: %v", err, err)
	}
	if timeoutErr.Operation != SyncStopped {
		t.Errorf("Operation = %v, want SyncStopped", timeoutErr.Operation)
	}
}

func TestRestartChannel_StopTimeout(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// STOP command succeeds
	transport.addSuccessResponse()
	// All polls return running - stop will timeout
	for range 50 {
		transport.addSuccessResponse(map[string]any{
			"CHANNEL": "TO.REMOTE",
			"STATUS":  "RUNNING",
		})
	}

	session := newTestSessionWithClock(transport, clock)

	_, err := session.RestartChannel(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 3 * time.Second, PollInterval: 1 * time.Second})
	if err == nil {
		t.Fatal("expected error when stop phase times out")
	}

	var timeoutErr *TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("expected TimeoutError, got %T: %v", err, err)
	}
}

func TestStartChannelSync_StartCommandError(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// START command fails with transport error
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection refused"),
	})

	session := newTestSessionWithClock(transport, clock)

	_, err := session.StartChannelSync(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err == nil {
		t.Fatal("expected error when START command fails")
	}
}

func TestStopChannelSync_StopCommandError(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// STOP command fails
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection refused"),
	})

	session := newTestSessionWithClock(transport, clock)

	_, err := session.StopChannelSync(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err == nil {
		t.Fatal("expected error when STOP command fails")
	}
}

func TestStartServiceSync_Success(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// START succeeds
	transport.addSuccessResponse()
	// Poll: running
	transport.addSuccessResponse(map[string]any{
		"SERVICE": "SVC1",
		"STATUS":  "RUNNING",
	})

	session := newTestSessionWithClock(transport, clock)

	result, err := session.StartServiceSync(context.Background(), "SVC1",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Operation != SyncStarted {
		t.Errorf("Operation = %v, want SyncStarted", result.Operation)
	}
}

func TestStartListenerSync_Success(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	transport.addSuccessResponse()
	transport.addSuccessResponse(map[string]any{
		"LISTENER": "LIS1",
		"STATUS":   "RUNNING",
	})

	session := newTestSessionWithClock(transport, clock)

	result, err := session.StartListenerSync(context.Background(), "LIS1",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Operation != SyncStarted {
		t.Errorf("Operation = %v, want SyncStarted", result.Operation)
	}
}

func TestStopServiceSync_Success(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	transport.addSuccessResponse()
	transport.addSuccessResponse(map[string]any{
		"SERVICE": "SVC1",
		"STATUS":  "STOPPED",
	})

	session := newTestSessionWithClock(transport, clock)

	result, err := session.StopServiceSync(context.Background(), "SVC1",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Operation != SyncStopped {
		t.Errorf("Operation = %v, want SyncStopped", result.Operation)
	}
}

func TestRestartListener_Success(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// STOP succeeds
	transport.addSuccessResponse()
	// Poll: stopped
	transport.addSuccessResponse(map[string]any{
		"LISTENER": "LIS1",
		"STATUS":   "STOPPED",
	})
	// START succeeds
	transport.addSuccessResponse()
	// Poll: running
	transport.addSuccessResponse(map[string]any{
		"LISTENER": "LIS1",
		"STATUS":   "RUNNING",
	})

	session := newTestSessionWithClock(transport, clock)

	result, err := session.RestartListener(context.Background(), "LIS1",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Operation != SyncRestarted {
		t.Errorf("Operation = %v, want SyncRestarted", result.Operation)
	}
}

func TestRestartService_Success(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// STOP succeeds
	transport.addSuccessResponse()
	// Poll: stopped
	transport.addSuccessResponse(map[string]any{
		"SERVICE": "SVC1",
		"STATUS":  "STOPPED",
	})
	// START succeeds
	transport.addSuccessResponse()
	// Poll: running
	transport.addSuccessResponse(map[string]any{
		"SERVICE": "SVC1",
		"STATUS":  "RUNNING",
	})

	session := newTestSessionWithClock(transport, clock)

	result, err := session.RestartService(context.Background(), "SVC1",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Operation != SyncRestarted {
		t.Errorf("Operation = %v, want SyncRestarted", result.Operation)
	}
}

func TestStopListenerSync_InactiveStatus(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	transport.addSuccessResponse()
	transport.addSuccessResponse(map[string]any{
		"LISTENER": "LIS1",
		"status":   "inactive",
	})

	session := newTestSessionWithClock(transport, clock)

	result, err := session.StopListenerSync(context.Background(), "LIS1",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Operation != SyncStopped {
		t.Errorf("Operation = %v, want SyncStopped", result.Operation)
	}
}

func TestRestartChannel_StartPhaseError(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// STOP succeeds
	transport.addSuccessResponse()
	// Poll: channel stopped
	transport.addCommandErrorResponse(2, 2085)
	// START fails
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection refused"),
	})

	session := newTestSessionWithClock(transport, clock)

	_, err := session.RestartChannel(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err == nil {
		t.Fatal("expected error when start phase fails during restart")
	}
}

func TestQueryStatus_NonCommandError_StartPoll(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// START succeeds
	transport.addSuccessResponse()
	// Status poll returns transport error (non-command) — must propagate
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection refused"),
	})

	session := newTestSessionWithClock(transport, clock)

	_, err := session.StartChannelSync(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err == nil {
		t.Fatal("expected error when transport fails during polling")
	}

	var transportErr *TransportError
	if !errors.As(err, &transportErr) {
		t.Fatalf("expected TransportError, got %T: %v", err, err)
	}
}

func TestQueryStatus_NonCommandError_StopPoll(t *testing.T) {
	transport := newMockTransport()
	clock := newMockClock()

	// STOP succeeds
	transport.addSuccessResponse()
	// Status poll returns auth error (non-command) — must propagate
	transport.addErrorResponse(&AuthError{
		URL:        "https://localhost:9443",
		StatusCode: 401,
	})

	session := newTestSessionWithClock(transport, clock)

	_, err := session.StopChannelSync(context.Background(), "TO.REMOTE",
		SyncConfig{Timeout: 30 * time.Second, PollInterval: 1 * time.Second})
	if err == nil {
		t.Fatal("expected error when auth fails during polling")
	}

	var authErr *AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthError, got %T: %v", err, err)
	}
}

func TestHasStatus_NonStringValue(t *testing.T) {
	rows := []map[string]any{{"STATUS": 42}}
	result := hasStatus(rows, []string{"STATUS"}, runningValues)
	if result {
		t.Error("expected false for non-string status value")
	}
}
