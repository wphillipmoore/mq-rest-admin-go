package mqrestadmin

import (
	"context"
	"errors"
	"testing"
)

func TestEnsureQlocal_Created(t *testing.T) {
	transport := newMockTransport()
	// DISPLAY returns command error (object not found)
	transport.addCommandErrorResponse(2, 2085)
	// DEFINE succeeds
	transport.addSuccessResponse()

	session := newTestSession(transport)

	result, err := session.EnsureQlocal(context.Background(), "NEW.QUEUE",
		map[string]any{"MAXDEPTH": "5000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Action != EnsureCreated {
		t.Errorf("Action = %v, want EnsureCreated", result.Action)
	}

	// Should have made 2 calls: DISPLAY + DEFINE
	if transport.callCount() != 2 {
		t.Errorf("expected 2 transport calls, got %d", transport.callCount())
	}

	// Verify DEFINE payload
	defineCall := transport.calls[1]
	if defineCall.Payload["command"] != "DEFINE" {
		t.Errorf("second call command = %v, want DEFINE", defineCall.Payload["command"])
	}
}

func TestEnsureQlocal_Unchanged(t *testing.T) {
	transport := newMockTransport()
	// DISPLAY returns existing object with matching attributes
	transport.addSuccessResponse(map[string]any{
		"QNAME":    "EXISTING.QUEUE",
		"MAXDEPTH": "5000",
	})

	session := newTestSession(transport)

	result, err := session.EnsureQlocal(context.Background(), "EXISTING.QUEUE",
		map[string]any{"MAXDEPTH": "5000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Action != EnsureUnchanged {
		t.Errorf("Action = %v, want EnsureUnchanged", result.Action)
	}

	// Should have made only 1 call: DISPLAY
	if transport.callCount() != 1 {
		t.Errorf("expected 1 transport call, got %d", transport.callCount())
	}
}

func TestEnsureQlocal_Updated(t *testing.T) {
	transport := newMockTransport()
	// DISPLAY returns existing object with different attributes
	transport.addSuccessResponse(map[string]any{
		"QNAME":    "EXISTING.QUEUE",
		"MAXDEPTH": "5000",
	})
	// ALTER succeeds
	transport.addSuccessResponse()

	session := newTestSession(transport)

	result, err := session.EnsureQlocal(context.Background(), "EXISTING.QUEUE",
		map[string]any{"MAXDEPTH": "10000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Action != EnsureUpdated {
		t.Errorf("Action = %v, want EnsureUpdated", result.Action)
	}

	if len(result.Changed) != 1 || result.Changed[0] != "MAXDEPTH" {
		t.Errorf("Changed = %v, want [MAXDEPTH]", result.Changed)
	}

	// Verify ALTER payload contains only changed attributes
	alterCall := transport.calls[1]
	params := alterCall.Payload["parameters"].(map[string]any)
	if params["MAXDEPTH"] != "10000" {
		t.Errorf("ALTER MAXDEPTH = %v, want 10000", params["MAXDEPTH"])
	}
}

func TestEnsureQlocal_UnchangedNoParams(t *testing.T) {
	transport := newMockTransport()
	// DISPLAY returns existing object
	transport.addSuccessResponse(map[string]any{
		"QNAME": "EXISTING.QUEUE",
	})

	session := newTestSession(transport)

	result, err := session.EnsureQlocal(context.Background(), "EXISTING.QUEUE", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Action != EnsureUnchanged {
		t.Errorf("Action = %v, want EnsureUnchanged", result.Action)
	}
}

func TestEnsureQmgr_Updated(t *testing.T) {
	transport := newMockTransport()
	// DISPLAY returns current qmgr state
	transport.addSuccessResponse(map[string]any{
		"QMNAME": "QM1",
		"DESCR":  "old description",
	})
	// ALTER succeeds
	transport.addSuccessResponse()

	session := newTestSession(transport)

	result, err := session.EnsureQmgr(context.Background(),
		map[string]any{"DESCR": "new description"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Action != EnsureUpdated {
		t.Errorf("Action = %v, want EnsureUpdated", result.Action)
	}
}

func TestEnsureQmgr_UnchangedNoParams(t *testing.T) {
	session := newTestSession(newMockTransport())

	result, err := session.EnsureQmgr(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Action != EnsureUnchanged {
		t.Errorf("Action = %v, want EnsureUnchanged", result.Action)
	}
}

func TestValuesMatch_CaseInsensitive(t *testing.T) {
	tests := []struct {
		desired  any
		current  any
		expected bool
	}{
		{"5000", "5000", true},
		{"RUNNING", "running", true},
		{"  trimmed  ", "trimmed", true},
		{"5000", "10000", false},
		{5000, "5000", true},
		{"yes", "YES", true},
	}

	for _, test := range tests {
		result := valuesMatch(test.desired, test.current)
		if result != test.expected {
			t.Errorf("valuesMatch(%v, %v) = %v, want %v",
				test.desired, test.current, result, test.expected)
		}
	}
}

func TestEnsureQlocal_DisplayNonCommandError(t *testing.T) {
	transport := newMockTransport()
	// DISPLAY returns a transport error (not a CommandError)
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection refused"),
	})

	session := newTestSession(transport)

	_, err := session.EnsureQlocal(context.Background(), "NEW.QUEUE",
		map[string]any{"MAXDEPTH": "5000"})
	if err == nil {
		t.Fatal("expected error for transport failure")
	}

	// Should propagate the transport error, not treat as "not found"
	var transportErr *TransportError
	if !errors.As(err, &transportErr) {
		t.Errorf("expected TransportError to be wrapped, got: %v", err)
	}
}

func TestEnsureQlocal_DefineError(t *testing.T) {
	transport := newMockTransport()
	// DISPLAY returns command error (object not found)
	transport.addCommandErrorResponse(2, 2085)
	// DEFINE fails
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection lost"),
	})

	session := newTestSession(transport)

	_, err := session.EnsureQlocal(context.Background(), "NEW.QUEUE",
		map[string]any{"MAXDEPTH": "5000"})
	if err == nil {
		t.Fatal("expected error when DEFINE fails")
	}
}

func TestEnsureQlocal_AlterError(t *testing.T) {
	transport := newMockTransport()
	// DISPLAY returns existing object with different attributes
	transport.addSuccessResponse(map[string]any{
		"QNAME":    "EXISTING.QUEUE",
		"MAXDEPTH": "5000",
	})
	// ALTER fails
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection lost"),
	})

	session := newTestSession(transport)

	_, err := session.EnsureQlocal(context.Background(), "EXISTING.QUEUE",
		map[string]any{"MAXDEPTH": "10000"})
	if err == nil {
		t.Fatal("expected error when ALTER fails")
	}
}

func TestEnsureQmgr_DisplayError(t *testing.T) {
	transport := newMockTransport()
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection refused"),
	})

	session := newTestSession(transport)

	_, err := session.EnsureQmgr(context.Background(),
		map[string]any{"DESCR": "new"})
	if err == nil {
		t.Fatal("expected error when DISPLAY QMGR fails")
	}
}

func TestEnsureQmgr_AlterError(t *testing.T) {
	transport := newMockTransport()
	// DISPLAY succeeds
	transport.addSuccessResponse(map[string]any{
		"QMNAME": "QM1",
		"DESCR":  "old",
	})
	// ALTER fails
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection lost"),
	})

	session := newTestSession(transport)

	_, err := session.EnsureQmgr(context.Background(),
		map[string]any{"DESCR": "new"})
	if err == nil {
		t.Fatal("expected error when ALTER QMGR fails")
	}
}

func TestEnsureWrappers_Created(t *testing.T) {
	type ensureEntry struct {
		name string
		call func(*Session, context.Context, string, map[string]any) (EnsureResult, error)
	}

	entries := []ensureEntry{
		{"EnsureQremote", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureQremote(ctx, name, params) }},
		{"EnsureQalias", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureQalias(ctx, name, params) }},
		{"EnsureQmodel", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureQmodel(ctx, name, params) }},
		{"EnsureChannel", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureChannel(ctx, name, params) }},
		{"EnsureAuthinfo", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureAuthinfo(ctx, name, params) }},
		{"EnsureListener", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureListener(ctx, name, params) }},
		{"EnsureNamelist", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureNamelist(ctx, name, params) }},
		{"EnsureProcess", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureProcess(ctx, name, params) }},
		{"EnsureService", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureService(ctx, name, params) }},
		{"EnsureTopic", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureTopic(ctx, name, params) }},
		{"EnsureSub", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureSub(ctx, name, params) }},
		{"EnsureStgclass", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureStgclass(ctx, name, params) }},
		{"EnsureComminfo", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureComminfo(ctx, name, params) }},
		{"EnsureCfstruct", func(s *Session, ctx context.Context, name string, params map[string]any) (EnsureResult, error) { return s.EnsureCfstruct(ctx, name, params) }},
	}

	for _, entry := range entries {
		t.Run(entry.name, func(t *testing.T) {
			transport := newMockTransport()
			// DISPLAY returns command error (not found)
			transport.addCommandErrorResponse(2, 2085)
			// DEFINE succeeds
			transport.addSuccessResponse()

			session := newTestSession(transport)

			result, err := entry.call(session, context.Background(), "NEW.OBJ",
				map[string]any{"DESCR": "test"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Action != EnsureCreated {
				t.Errorf("Action = %v, want EnsureCreated", result.Action)
			}
			if transport.callCount() != 2 {
				t.Errorf("expected 2 transport calls, got %d", transport.callCount())
			}
		})
	}
}

func TestEnsureQmgr_Unchanged(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{
		"QMNAME": "QM1",
		"DESCR":  "same",
	})

	session := newTestSession(transport)

	result, err := session.EnsureQmgr(context.Background(),
		map[string]any{"DESCR": "same"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != EnsureUnchanged {
		t.Errorf("Action = %v, want EnsureUnchanged", result.Action)
	}
}

func TestDiffAttributes(t *testing.T) {
	desired := map[string]any{
		"MAXDEPTH": "10000",
		"DESCR":    "same",
		"NEWATTR":  "value",
	}
	current := map[string]any{
		"MAXDEPTH": "5000",
		"DESCR":    "same",
	}

	changed, changedParams := diffAttributes(desired, current)

	// MAXDEPTH changed, NEWATTR is new (missing from current)
	if len(changed) != 2 {
		t.Errorf("expected 2 changed, got %d: %v", len(changed), changed)
	}

	if changedParams["MAXDEPTH"] != "10000" {
		t.Errorf("MAXDEPTH = %v, want 10000", changedParams["MAXDEPTH"])
	}
	if changedParams["NEWATTR"] != "value" {
		t.Errorf("NEWATTR = %v, want value", changedParams["NEWATTR"])
	}
	if _, hasSame := changedParams["DESCR"]; hasSame {
		t.Error("DESCR should not be in changed params")
	}
}
