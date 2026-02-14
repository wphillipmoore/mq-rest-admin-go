package mqrest

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestTransportError_Error(t *testing.T) {
	err := &TransportError{
		URL: "https://localhost:9443",
		Err: fmt.Errorf("connection refused"),
	}

	msg := err.Error()
	if !strings.Contains(msg, "https://localhost:9443") {
		t.Errorf("error message should contain URL: %s", msg)
	}
	if !strings.Contains(msg, "connection refused") {
		t.Errorf("error message should contain cause: %s", msg)
	}
}

func TestTransportError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("connection refused")
	err := &TransportError{URL: "https://localhost", Err: cause}

	if !errors.Is(err, cause) {
		t.Error("Unwrap should expose underlying error")
	}
}

func TestResponseError_Error(t *testing.T) {
	err := &ResponseError{
		ResponseText: "not json",
		StatusCode:   200,
	}

	msg := err.Error()
	if !strings.Contains(msg, "200") {
		t.Errorf("error message should contain status code: %s", msg)
	}
}

func TestAuthError_Error(t *testing.T) {
	err := &AuthError{
		URL:        "https://localhost:9443/login",
		StatusCode: 401,
	}

	msg := err.Error()
	if !strings.Contains(msg, "401") {
		t.Errorf("error message should contain status code: %s", msg)
	}
	if !strings.Contains(msg, "https://localhost:9443/login") {
		t.Errorf("error message should contain URL: %s", msg)
	}
}

func TestCommandError_Error(t *testing.T) {
	err := &CommandError{
		Payload:    map[string]any{"overallReasonCode": float64(2085)},
		StatusCode: 200,
	}

	msg := err.Error()
	if !strings.Contains(msg, "200") {
		t.Errorf("error message should contain status code: %s", msg)
	}
}

func TestTimeoutError_Error(t *testing.T) {
	err := &TimeoutError{
		Name:           "TO.REMOTE",
		Operation:      SyncStarted,
		ElapsedSeconds: 30.5,
	}

	msg := err.Error()
	if !strings.Contains(msg, "TO.REMOTE") {
		t.Errorf("error message should contain name: %s", msg)
	}
	if !strings.Contains(msg, "started") {
		t.Errorf("error message should contain operation: %s", msg)
	}
}

func TestMappingError_Error(t *testing.T) {
	err := &MappingError{
		Issues: []MappingIssue{
			{
				Direction:     MappingRequest,
				Reason:        MappingUnknownKey,
				AttributeName: "bad_attr",
			},
			{
				Direction:     MappingResponse,
				Reason:        MappingUnknownValue,
				AttributeName: "status",
			},
		},
	}

	msg := err.Error()
	if !strings.Contains(msg, "bad_attr") {
		t.Errorf("error message should contain attribute name: %s", msg)
	}
	if !strings.Contains(msg, "status") {
		t.Errorf("error message should contain second attribute: %s", msg)
	}
}

func TestMappingIssue_String(t *testing.T) {
	issue := MappingIssue{
		Direction:     MappingRequest,
		Reason:        MappingUnknownKey,
		AttributeName: "bad_attr",
	}

	result := issue.String()
	if !strings.Contains(result, "request") {
		t.Errorf("String() should contain direction: %s", result)
	}
	if !strings.Contains(result, "unknown_key") {
		t.Errorf("String() should contain reason: %s", result)
	}
}

func TestErrorsAs_TransportError(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &TransportError{URL: "https://localhost", Err: fmt.Errorf("fail")})

	var transportErr *TransportError
	if !errors.As(err, &transportErr) {
		t.Error("errors.As should find TransportError")
	}
}

func TestErrorsAs_AuthError(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &AuthError{URL: "https://localhost", StatusCode: 401})

	var authErr *AuthError
	if !errors.As(err, &authErr) {
		t.Error("errors.As should find AuthError")
	}
}

func TestErrorsAs_CommandError(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &CommandError{
		Payload:    map[string]any{"reasonCode": float64(2085)},
		StatusCode: 200,
	})

	var cmdErr *CommandError
	if !errors.As(err, &cmdErr) {
		t.Error("errors.As should find CommandError")
	}
}

func TestMappingDirection_String(t *testing.T) {
	if MappingRequest.String() != "request" {
		t.Errorf("MappingRequest.String() = %q", MappingRequest.String())
	}
	if MappingResponse.String() != "response" {
		t.Errorf("MappingResponse.String() = %q", MappingResponse.String())
	}
}

func TestMappingReason_String(t *testing.T) {
	if MappingUnknownKey.String() != "unknown_key" {
		t.Errorf("MappingUnknownKey.String() = %q", MappingUnknownKey.String())
	}
	if MappingUnknownValue.String() != "unknown_value" {
		t.Errorf("MappingUnknownValue.String() = %q", MappingUnknownValue.String())
	}
	if MappingUnknownQualifier.String() != "unknown_qualifier" {
		t.Errorf("MappingUnknownQualifier.String() = %q", MappingUnknownQualifier.String())
	}
}

func TestEnsureAction_String(t *testing.T) {
	if EnsureCreated.String() != "created" {
		t.Errorf("EnsureCreated.String() = %q", EnsureCreated.String())
	}
	if EnsureUpdated.String() != "updated" {
		t.Errorf("EnsureUpdated.String() = %q", EnsureUpdated.String())
	}
	if EnsureUnchanged.String() != "unchanged" {
		t.Errorf("EnsureUnchanged.String() = %q", EnsureUnchanged.String())
	}
}

func TestSyncOperation_String(t *testing.T) {
	if SyncStarted.String() != "started" {
		t.Errorf("SyncStarted.String() = %q", SyncStarted.String())
	}
	if SyncStopped.String() != "stopped" {
		t.Errorf("SyncStopped.String() = %q", SyncStopped.String())
	}
	if SyncRestarted.String() != "restarted" {
		t.Errorf("SyncRestarted.String() = %q", SyncRestarted.String())
	}
}
