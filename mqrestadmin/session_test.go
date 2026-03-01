package mqrestadmin

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewSession_BasicAuth(t *testing.T) {
	transport := newMockTransport()

	session, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		BasicAuth{Username: "admin", Password: "pass"},
		WithTransport(transport),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.QmgrName() != "QM1" {
		t.Errorf("QmgrName() = %q, want %q", session.QmgrName(), "QM1")
	}
	if session.GatewayQmgr() != "" {
		t.Errorf("GatewayQmgr() = %q, want empty", session.GatewayQmgr())
	}
}

func TestNewSession_WithOptions(t *testing.T) {
	transport := newMockTransport()

	session, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2/",
		"QM1",
		BasicAuth{Username: "admin", Password: "pass"},
		WithTransport(transport),
		WithGatewayQmgr("GATEWAY"),
		WithVerifyTLS(false),
		WithTimeout(60*time.Second),
		WithMapAttributes(false),
		WithMappingStrict(false),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.GatewayQmgr() != "GATEWAY" {
		t.Errorf("GatewayQmgr() = %q, want %q", session.GatewayQmgr(), "GATEWAY")
	}
	if session.verifyTLS {
		t.Error("verifyTLS should be false")
	}
	if session.timeout != 60*time.Second {
		t.Errorf("timeout = %v, want 60s", session.timeout)
	}
	// URL should have trailing slash stripped
	if session.restBaseURL != "https://localhost:9443/ibmmq/rest/v2" {
		t.Errorf("restBaseURL = %q, trailing slash not stripped", session.restBaseURL)
	}
}

func TestNewSession_LTPAAuth(t *testing.T) {
	transport := newMockTransport()
	transport.addResponse(200, map[string]any{}, map[string]string{
		"Set-Cookie": "LtpaToken2=abc123token; Path=/; Secure",
	})

	session, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		LTPAAuth{Username: "admin", Password: "pass"},
		WithTransport(transport),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.ltpaCookieName != "LtpaToken2" {
		t.Errorf("ltpaCookieName = %q, want %q", session.ltpaCookieName, "LtpaToken2")
	}
	if session.ltpaToken != "abc123token" {
		t.Errorf("ltpaToken = %q, want %q", session.ltpaToken, "abc123token")
	}

	// Verify login call
	call := transport.calls[0]
	if call.URL != "https://localhost:9443/ibmmq/rest/v2/login" {
		t.Errorf("login URL = %q, want /login endpoint", call.URL)
	}
	if call.Payload["username"] != "admin" {
		t.Errorf("login username = %v, want admin", call.Payload["username"])
	}
}

func TestNewSession_LTPAAuth_LoginFailure(t *testing.T) {
	transport := newMockTransport()
	transport.addResponse(401, map[string]any{}, nil)

	_, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		LTPAAuth{Username: "admin", Password: "wrong"},
		WithTransport(transport),
	)
	if err == nil {
		t.Fatal("expected error for failed LTPA login")
	}

	var authErr *AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthError, got %T: %v", err, err)
	}
	if authErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", authErr.StatusCode)
	}
}

func TestNewSession_LTPAAuth_MissingCookie(t *testing.T) {
	transport := newMockTransport()
	transport.addResponse(200, map[string]any{}, map[string]string{})

	_, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		LTPAAuth{Username: "admin", Password: "pass"},
		WithTransport(transport),
	)
	if err == nil {
		t.Fatal("expected error when LtpaToken2 cookie is missing")
	}
}

func TestMqscCommand_BuildsCorrectPayload(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{"QNAME": "TEST.Q"})
	session := newTestSession(transport)

	name := "TEST.Q"
	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QLOCAL", &name,
		map[string]any{"MAXDEPTH": "5000"}, []string{"all"}, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := transport.lastCall()
	payload := call.Payload

	if payload["type"] != "runCommandJSON" {
		t.Errorf("type = %v, want runCommandJSON", payload["type"])
	}
	if payload["command"] != "DISPLAY" {
		t.Errorf("command = %v, want DISPLAY", payload["command"])
	}
	if payload["qualifier"] != "QLOCAL" {
		t.Errorf("qualifier = %v, want QLOCAL", payload["qualifier"])
	}
	if payload["name"] != "TEST.Q" {
		t.Errorf("name = %v, want TEST.Q", payload["name"])
	}

	params, hasParams := payload["parameters"].(map[string]any)
	if !hasParams {
		t.Fatal("expected parameters in payload")
	}
	if params["MAXDEPTH"] != "5000" {
		t.Errorf("MAXDEPTH = %v, want 5000", params["MAXDEPTH"])
	}
}

func TestMqscCommand_BuildsCorrectURL(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QMGR", nil,
		nil, nil, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := transport.lastCall()
	expected := "https://localhost:9443/ibmmq/rest/v2/admin/action/qmgr/QM1/mqsc"
	if call.URL != expected {
		t.Errorf("URL = %q, want %q", call.URL, expected)
	}
}

func TestMqscCommand_IncludesHeaders(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)
	session.gatewayQmgr = "GATEWAY_QM"

	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QMGR", nil,
		nil, nil, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	headers := transport.lastCall().Headers
	if headers["Accept"] != "application/json" {
		t.Errorf("Accept = %q", headers["Accept"])
	}
	if headers["ibm-mq-rest-csrf-token"] != "local" {
		t.Errorf("csrf-token = %q", headers["ibm-mq-rest-csrf-token"])
	}
	if headers["ibm-mq-rest-gateway-qmgr"] != "GATEWAY_QM" {
		t.Errorf("gateway-qmgr = %q", headers["ibm-mq-rest-gateway-qmgr"])
	}
}

func TestMqscCommand_AuthError401(t *testing.T) {
	transport := newMockTransport()
	transport.addResponse(401, map[string]any{}, nil)
	session := newTestSession(transport)

	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QMGR", nil,
		nil, nil, nil, true)
	if err == nil {
		t.Fatal("expected error for 401 response")
	}

	var authErr *AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthError, got %T", err)
	}
	if authErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", authErr.StatusCode)
	}
}

func TestMqscCommand_AuthError403(t *testing.T) {
	transport := newMockTransport()
	transport.addResponse(403, map[string]any{}, nil)
	session := newTestSession(transport)

	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QMGR", nil,
		nil, nil, nil, true)

	var authErr *AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthError, got %T: %v", err, err)
	}
}

func TestMqscCommand_CommandError(t *testing.T) {
	transport := newMockTransport()
	transport.addCommandErrorResponse(2, 2085)
	session := newTestSession(transport)

	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QLOCAL", nil,
		nil, nil, nil, true)
	if err == nil {
		t.Fatal("expected error for command error response")
	}

	var cmdErr *CommandError
	if !errors.As(err, &cmdErr) {
		t.Fatalf("expected CommandError, got %T", err)
	}
}

func TestMqscCommand_TransportError(t *testing.T) {
	transport := newMockTransport()
	transport.addErrorResponse(&TransportError{URL: "https://localhost", Err: errors.New("connection refused")})
	session := newTestSession(transport)

	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QMGR", nil,
		nil, nil, nil, true)
	if err == nil {
		t.Fatal("expected error for transport failure")
	}

	var transportErr *TransportError
	if !errors.As(err, &transportErr) {
		t.Fatalf("expected TransportError, got %T", err)
	}
}

func TestMqscCommand_InvalidJSON(t *testing.T) {
	transport := newMockTransport()
	transport.responses = append(transport.responses, mockResponse{
		Response: &TransportResponse{
			StatusCode: 200,
			Body:       "not json",
			Headers:    map[string]string{},
		},
	})
	session := newTestSession(transport)

	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QMGR", nil,
		nil, nil, nil, true)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}

	var respErr *ResponseError
	if !errors.As(err, &respErr) {
		t.Fatalf("expected ResponseError, got %T", err)
	}
}

func TestMqscCommand_ExtractsParameters(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(
		map[string]any{"QNAME": "Q1", "CURDEPTH": float64(5)},
		map[string]any{"QNAME": "Q2", "CURDEPTH": float64(10)},
	)
	session := newTestSession(transport)

	name := "*"
	objects, err := session.mqscCommand(context.Background(), "DISPLAY", "QLOCAL", &name,
		nil, []string{"all"}, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(objects) != 2 {
		t.Fatalf("expected 2 objects, got %d", len(objects))
	}
	if objects[0]["QNAME"] != "Q1" {
		t.Errorf("objects[0][QNAME] = %v, want Q1", objects[0]["QNAME"])
	}
	if objects[1]["QNAME"] != "Q2" {
		t.Errorf("objects[1][QNAME] = %v, want Q2", objects[1]["QNAME"])
	}
}

func TestMqscCommand_FlattensNestedObjects(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{
		"CONN": "CONN1",
		"objects": []any{
			map[string]any{"OBJNAME": "Q1", "OBJTYPE": "QUEUE"},
			map[string]any{"OBJNAME": "Q2", "OBJTYPE": "QUEUE"},
		},
	})
	session := newTestSession(transport)

	name := "*"
	objects, err := session.mqscCommand(context.Background(), "DISPLAY", "CONN", &name,
		nil, []string{"all"}, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(objects) != 2 {
		t.Fatalf("expected 2 flattened objects, got %d", len(objects))
	}
	// Parent field merged into each child
	if objects[0]["CONN"] != "CONN1" {
		t.Errorf("objects[0][CONN] = %v, want CONN1", objects[0]["CONN"])
	}
	if objects[0]["OBJNAME"] != "Q1" {
		t.Errorf("objects[0][OBJNAME] = %v, want Q1", objects[0]["OBJNAME"])
	}
	if objects[1]["OBJNAME"] != "Q2" {
		t.Errorf("objects[1][OBJNAME] = %v, want Q2", objects[1]["OBJNAME"])
	}
}

func TestMqscCommand_SavesLastState(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{"QNAME": "Q1"})
	session := newTestSession(transport)

	name := "Q1"
	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QLOCAL", &name,
		nil, nil, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.LastHTTPStatus != 200 {
		t.Errorf("LastHTTPStatus = %d, want 200", session.LastHTTPStatus)
	}
	if session.LastResponsePayload == nil {
		t.Error("LastResponsePayload should not be nil")
	}
	if session.LastCommandPayload == nil {
		t.Error("LastCommandPayload should not be nil")
	}
	if session.LastResponseText == "" {
		t.Error("LastResponseText should not be empty")
	}
}

func TestDisplayQueue_DefaultsNameToWildcard(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{"QNAME": "Q1"})
	session := newTestSession(transport)

	_, err := session.DisplayQueue(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if payload["name"] != "*" {
		t.Errorf("name = %v, want *", payload["name"])
	}
}

func TestDisplayQmgr_ReturnsSingleObject(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{"QMNAME": "QM1"})
	session := newTestSession(transport)

	result, err := session.DisplayQmgr(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result["QMNAME"] != "QM1" {
		t.Errorf("QMNAME = %v, want QM1", result["QMNAME"])
	}
}

func TestDisplayQmgr_ReturnsNilWhenEmpty(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	result, err := session.DisplayQmgr(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestDefineQlocal_SendsCorrectCommand(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	err := session.DefineQlocal(context.Background(), "MY.QUEUE",
		WithRequestParameters(map[string]any{"MAXDEPTH": "5000"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if payload["command"] != "DEFINE" {
		t.Errorf("command = %v, want DEFINE", payload["command"])
	}
	if payload["qualifier"] != "QLOCAL" {
		t.Errorf("qualifier = %v, want QLOCAL", payload["qualifier"])
	}
	if payload["name"] != "MY.QUEUE" {
		t.Errorf("name = %v, want MY.QUEUE", payload["name"])
	}
}

func TestAlterQmgr_NoNameInPayload(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	err := session.AlterQmgr(context.Background(),
		WithRequestParameters(map[string]any{"DESCR": "test"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if payload["command"] != "ALTER" {
		t.Errorf("command = %v, want ALTER", payload["command"])
	}
	if _, hasName := payload["name"]; hasName {
		t.Error("ALTER QMGR should not include name in payload")
	}
}

func TestDeleteQueue_SendsCorrectCommand(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	err := session.DeleteQueue(context.Background(), "MY.QUEUE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if payload["command"] != "DELETE" {
		t.Errorf("command = %v, want DELETE", payload["command"])
	}
	if payload["qualifier"] != "QUEUE" {
		t.Errorf("qualifier = %v, want QUEUE", payload["qualifier"])
	}
}

func TestDeleteQlocal_SendsCorrectCommand(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	err := session.DeleteQlocal(context.Background(), "MY.QUEUE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if payload["command"] != "DELETE" {
		t.Errorf("command = %v, want DELETE", payload["command"])
	}
	if payload["qualifier"] != "QLOCAL" {
		t.Errorf("qualifier = %v, want QLOCAL", payload["qualifier"])
	}
}

func TestDeleteQremote_SendsCorrectCommand(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	err := session.DeleteQremote(context.Background(), "MY.QUEUE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if payload["command"] != "DELETE" {
		t.Errorf("command = %v, want DELETE", payload["command"])
	}
	if payload["qualifier"] != "QREMOTE" {
		t.Errorf("qualifier = %v, want QREMOTE", payload["qualifier"])
	}
}

func TestDeleteQalias_SendsCorrectCommand(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	err := session.DeleteQalias(context.Background(), "MY.QUEUE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if payload["command"] != "DELETE" {
		t.Errorf("command = %v, want DELETE", payload["command"])
	}
	if payload["qualifier"] != "QALIAS" {
		t.Errorf("qualifier = %v, want QALIAS", payload["qualifier"])
	}
}

func TestDeleteQmodel_SendsCorrectCommand(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	err := session.DeleteQmodel(context.Background(), "MY.QUEUE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if payload["command"] != "DELETE" {
		t.Errorf("command = %v, want DELETE", payload["command"])
	}
	if payload["qualifier"] != "QMODEL" {
		t.Errorf("qualifier = %v, want QMODEL", payload["qualifier"])
	}
}

func TestStartChannel_SendsCorrectCommand(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	err := session.StartChannel(context.Background(), "TO.REMOTE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if payload["command"] != "START" {
		t.Errorf("command = %v, want START", payload["command"])
	}
	if payload["qualifier"] != "CHANNEL" {
		t.Errorf("qualifier = %v, want CHANNEL", payload["qualifier"])
	}
}

func TestVoidCommand_OptionalNameEmpty(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	// Empty name should result in no "name" field in payload
	err := session.StopTrace(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if _, hasName := payload["name"]; hasName {
		t.Error("empty name should not include name in payload")
	}
}

func TestMqscCommand_DisplayDefaultsResponseParamsToAll(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QMGR", nil,
		nil, nil, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	respParams, ok := payload["responseParameters"].([]string)
	if !ok {
		t.Fatal("expected responseParameters to be []string")
	}
	if len(respParams) != 1 || respParams[0] != "all" {
		t.Errorf("responseParameters = %v, want [all]", respParams)
	}
}

func TestMqscCommand_NonDisplayNoDefaultResponseParams(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	name := "Q1"
	_, err := session.mqscCommand(context.Background(), "DEFINE", "QLOCAL", &name,
		nil, nil, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if _, hasRP := payload["responseParameters"]; hasRP {
		t.Error("non-DISPLAY should not include default responseParameters")
	}
}

func TestDisplayChannel_DefaultsNameToWildcard(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{"CHANNEL": "TO.REMOTE"})
	session := newTestSession(transport)

	result, err := session.DisplayChannel(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if payload["name"] != "*" {
		t.Errorf("name = %v, want *", payload["name"])
	}
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
}

func TestDisplayList_OptionalName(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(
		map[string]any{"LISTENER": "LIS1"},
		map[string]any{"LISTENER": "LIS2"},
	)
	session := newTestSession(transport)

	result, err := session.DisplayListener(context.Background(), "LIS*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}

	payload := transport.lastCall().Payload
	if payload["command"] != "DISPLAY" {
		t.Errorf("command = %v, want DISPLAY", payload["command"])
	}
	if payload["qualifier"] != "LISTENER" {
		t.Errorf("qualifier = %v, want LISTENER", payload["qualifier"])
	}
	if payload["name"] != "LIS*" {
		t.Errorf("name = %v, want LIS*", payload["name"])
	}
}

func TestDisplayList_EmptyNameOmitsName(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSession(transport)

	_, err := session.DisplayService(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	if _, hasName := payload["name"]; hasName {
		t.Error("empty name should not include name in payload")
	}
}

func TestWithWhere_IncludedInPayload(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{"QUEUE": "TEST.Q"})
	session := newTestSession(transport)

	_, err := session.DisplayQueue(context.Background(), "*",
		WithWhere("CURDEPTH GT 0"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The where clause is passed through to mqscCommand which includes it
	// in the call arguments (not directly in the JSON payload structure)
	call := transport.lastCall()
	if call.Payload["command"] != "DISPLAY" {
		t.Errorf("command = %v, want DISPLAY", call.Payload["command"])
	}
}

func TestWithResponseParameters(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{"QUEUE": "Q1"})
	session := newTestSession(transport)

	_, err := session.DisplayQueue(context.Background(), "Q1",
		WithResponseParameters([]string{"MAXDEPTH", "CURDEPTH"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	respParams, ok := payload["responseParameters"].([]string)
	if !ok {
		t.Fatal("expected responseParameters to be []string")
	}
	if len(respParams) != 2 {
		t.Errorf("expected 2 responseParameters, got %d", len(respParams))
	}
}

func TestMqscCommand_MappingStrictResponseError(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{
		"QUEUE":       "TEST.Q",
		"UNKNOWNATTR": "value",
	})
	session := newTestSessionWithMapping(transport)
	session.mappingStrict = true

	_, err := session.DisplayQueue(context.Background(), "TEST.Q")

	var mappingErr *MappingError
	if !errors.As(err, &mappingErr) {
		t.Fatalf("expected MappingError for unknown attr in strict mode, got %T: %v", err, err)
	}
}

func TestMqscCommand_MappingStrictRequestError(t *testing.T) {
	transport := newMockTransport()
	session := newTestSessionWithMapping(transport)
	session.mappingStrict = true

	// Use DisplayQueue which resolves to mapping qualifier "queue"
	_, err := session.DisplayQueue(context.Background(), "Q1",
		WithRequestParameters(map[string]any{"unknown_attr": "value"}))

	var mappingErr *MappingError
	if !errors.As(err, &mappingErr) {
		t.Fatalf("expected MappingError for unknown attr in strict request, got %T: %v", err, err)
	}
}

func TestMqscCommand_MapResponseParameterNames(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{"QUEUE": "Q1", "MAXDEPTH": float64(5000)})
	session := newTestSessionWithMapping(transport)
	session.mappingStrict = false

	_, err := session.DisplayQueue(context.Background(), "Q1",
		WithResponseParameters([]string{"max_queue_depth"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := transport.lastCall().Payload
	respParams, ok := payload["responseParameters"].([]string)
	if !ok {
		t.Fatal("expected responseParameters to be []string")
	}
	// "max_queue_depth" should be mapped to MQSC name
	foundMQSC := false
	for _, param := range respParams {
		if param == "MAXDEPTH" {
			foundMQSC = true
		}
	}
	if !foundMQSC {
		t.Errorf("expected response parameter to be mapped to MAXDEPTH, got %v", respParams)
	}
}

func TestMqscCommand_PerItemCommandError(t *testing.T) {
	transport := newMockTransport()
	// Overall codes OK, but per-item has error
	body := map[string]any{
		"overallCompletionCode": float64(0),
		"overallReasonCode":     float64(0),
		"commandResponse": []any{
			map[string]any{
				"completionCode": float64(2),
				"reasonCode":     float64(2085),
				"parameters":     map[string]any{},
			},
		},
	}
	transport.addResponse(200, body, nil)
	session := newTestSession(transport)

	_, err := session.mqscCommand(context.Background(), "DISPLAY", "QLOCAL", nil,
		nil, nil, nil, true)
	if err == nil {
		t.Fatal("expected error for per-item command error")
	}

	var cmdErr *CommandError
	if !errors.As(err, &cmdErr) {
		t.Fatalf("expected CommandError, got %T", err)
	}
}

func TestNewSession_WithCSRFToken(t *testing.T) {
	transport := newMockTransport()

	token := "custom-token"
	session, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		BasicAuth{Username: "admin", Password: "pass"},
		WithTransport(transport),
		WithCSRFToken(&token),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.csrfToken == nil || *session.csrfToken != "custom-token" {
		t.Errorf("csrfToken = %v, want custom-token", session.csrfToken)
	}
}

func TestNewSession_WithMappingOverrides(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{"QUEUE": "Q1"})

	overrides := map[string]any{
		"qualifiers": map[string]any{
			"queue": map[string]any{
				"request_key_map": map[string]any{
					"custom_field": "CUSTOMFLD",
				},
			},
		},
	}

	session, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		BasicAuth{Username: "admin", Password: "pass"},
		WithTransport(transport),
		WithMapAttributes(true),
		WithMappingStrict(false),
		WithMappingOverrides(overrides, MappingOverrideMerge),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.mapper == nil {
		t.Fatal("expected mapper to be initialized")
	}
}

func TestNewSession_CertificateAuth_WithTransport(t *testing.T) {
	transport := newMockTransport()

	session, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		CertificateAuth{CertPath: "/path/to/cert.pem"},
		WithTransport(transport),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session == nil {
		t.Fatal("expected non-nil session")
	}
}

func TestNewSession_CertificateAuth_NoTransport_InvalidCert(t *testing.T) {
	_, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		CertificateAuth{CertPath: "/nonexistent/cert.pem", KeyPath: "/nonexistent/key.pem"},
	)
	if err == nil {
		t.Fatal("expected error for invalid certificate path")
	}
}

func TestNewSession_CertificateAuth_NoTransport_ValidCert(t *testing.T) {
	certPEM, keyPEM := generateSelfSignedCert(&testing.T{})
	// Use a real *testing.T for TempDir
	certFile := writeTempFile(t, "cert.pem", certPEM)
	keyFile := writeTempFile(t, "key.pem", keyPEM)

	session, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		CertificateAuth{CertPath: certFile, KeyPath: keyFile},
		WithMapAttributes(false),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session == nil {
		t.Fatal("expected non-nil session")
	}
}

func TestWithBasicAuth_Option(t *testing.T) {
	transport := newMockTransport()

	session, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		BasicAuth{Username: "admin", Password: "pass"},
		WithTransport(transport),
		WithBasicAuth("user", "pass"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session == nil {
		t.Fatal("expected non-nil session")
	}
}

func TestSystemClock_Now(t *testing.T) {
	clock := systemClock{}
	before := time.Now()
	result := clock.now()
	after := time.Now()

	if result.Before(before) || result.After(after) {
		t.Errorf("systemClock.now() = %v, expected between %v and %v", result, before, after)
	}
}

func TestSystemClock_Sleep(t *testing.T) {
	clock := systemClock{}
	start := time.Now()
	clock.sleep(1 * time.Millisecond)
	elapsed := time.Since(start)

	if elapsed < 1*time.Millisecond {
		t.Errorf("systemClock.sleep() returned too quickly: %v", elapsed)
	}
}

func TestPerformLTPALogin_TransportError(t *testing.T) {
	transport := newMockTransport()
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection refused"),
	})

	_, err := NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		LTPAAuth{Username: "admin", Password: "pass"},
		WithTransport(transport),
	)
	if err == nil {
		t.Fatal("expected error for transport failure during LTPA login")
	}
}

func TestMapResponseParameterNames_UnknownQualifier(t *testing.T) {
	session := newTestSessionWithMapping(newMockTransport())

	// Call mapResponseParameterNames directly with unknown qualifier
	result := session.mapResponseParameterNames("nonexistent_qualifier", []string{"param1", "param2"})
	if len(result) != 2 || result[0] != "param1" || result[1] != "param2" {
		t.Errorf("expected params to pass through unchanged, got %v", result)
	}
}

func TestDisplayQmstatus_Error(t *testing.T) {
	transport := newMockTransport()
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection refused"),
	})
	session := newTestSession(transport)

	_, err := session.DisplayQmstatus(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDisplayCmdserv_Error(t *testing.T) {
	transport := newMockTransport()
	transport.addErrorResponse(&TransportError{
		URL: "https://localhost:9443",
		Err: errors.New("connection refused"),
	})
	session := newTestSession(transport)

	_, err := session.DisplayCmdserv(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExtractCommandResponseObjects_NonListCommandResponse(t *testing.T) {
	payload := map[string]any{
		"commandResponse": "not a list",
	}
	result := extractCommandResponseObjects(payload)
	if result != nil {
		t.Errorf("expected nil for non-list commandResponse, got %v", result)
	}
}

func TestExtractCommandResponseObjects_NonMapItem(t *testing.T) {
	payload := map[string]any{
		"commandResponse": []any{"not a map"},
	}
	result := extractCommandResponseObjects(payload)
	if len(result) != 0 {
		t.Errorf("expected 0 results for non-map item, got %d", len(result))
	}
}

func TestExtractCommandResponseObjects_MissingParameters(t *testing.T) {
	payload := map[string]any{
		"commandResponse": []any{
			map[string]any{"completionCode": float64(0)},
		},
	}
	result := extractCommandResponseObjects(payload)
	if len(result) != 0 {
		t.Errorf("expected 0 results when parameters missing, got %d", len(result))
	}
}

func TestExtractCommandResponseObjects_NonMapParameters(t *testing.T) {
	payload := map[string]any{
		"commandResponse": []any{
			map[string]any{
				"completionCode": float64(0),
				"parameters":     "not a map",
			},
		},
	}
	result := extractCommandResponseObjects(payload)
	if len(result) != 0 {
		t.Errorf("expected 0 results for non-map parameters, got %d", len(result))
	}
}

func TestFlattenNestedObjects_NonMapChild(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{
		"CONN":    "CONN1",
		"objects": []any{"not a map", map[string]any{"OBJNAME": "Q1"}},
	})
	session := newTestSession(transport)

	name := "*"
	objects, err := session.mqscCommand(context.Background(), "DISPLAY", "CONN", &name,
		nil, []string{"all"}, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only the valid map should be flattened
	if len(objects) != 1 {
		t.Errorf("expected 1 flattened object, got %d", len(objects))
	}
}

func TestIsNonZeroNumber_Nil(t *testing.T) {
	if isNonZeroNumber(nil) {
		t.Error("expected nil to not be a non-zero number")
	}
}

func TestIsNonZeroNumber_Float64NonZero(t *testing.T) {
	if !isNonZeroNumber(float64(42)) {
		t.Error("expected float64(42) to be non-zero")
	}
}

func TestExtractLTPAToken_MultipleHeaders(t *testing.T) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"Set-Cookie":   "LtpaToken2=abc123; Path=/; Secure",
	}
	name, val := extractLTPAToken(headers)
	if name != "LtpaToken2" {
		t.Errorf("extractLTPAToken() name = %q, want LtpaToken2", name)
	}
	if val != "abc123" {
		t.Errorf("extractLTPAToken() value = %q, want abc123", val)
	}
}

func TestMqscCommand_MappingPermissiveRequest(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{"QUEUE": "Q1"})
	session := newTestSessionWithMapping(transport)
	session.mappingStrict = false

	_, err := session.DisplayQueue(context.Background(), "Q1",
		WithRequestParameters(map[string]any{"unknown_attr": "value"}))
	if err != nil {
		t.Fatalf("expected no error in permissive mode, got %v", err)
	}
}

func TestIsNonZeroNumber_Float64Zero(t *testing.T) {
	if isNonZeroNumber(float64(0)) {
		t.Error("expected float64(0) to be zero")
	}
}

func TestIsNonZeroNumber_IntType(t *testing.T) {
	if !isNonZeroNumber(42) {
		t.Error("expected 42 (int) to be non-zero")
	}
	if isNonZeroNumber(0) {
		t.Error("expected 0 (int) to be zero")
	}
	if isNonZeroNumber("not a number") {
		t.Error("expected string to not be a number")
	}
}
