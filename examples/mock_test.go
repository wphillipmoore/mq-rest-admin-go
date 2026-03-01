package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

// mockTransport is a test Transport that records calls and returns
// pre-configured responses.
type mockTransport struct {
	calls     []mockCall
	responses []mockResponse
	callIndex int
}

type mockCall struct {
	URL       string
	Payload   map[string]any
	Headers   map[string]string
	Timeout   time.Duration
	VerifyTLS bool
}

type mockResponse struct {
	Response *mqrestadmin.TransportResponse
	Err      error
}

// addSuccessResponse queues a successful MQSC response with command responses.
func (transport *mockTransport) addSuccessResponse(commandResponses ...map[string]any) {
	body := map[string]any{
		"overallCompletionCode": float64(0),
		"overallReasonCode":     float64(0),
	}
	if len(commandResponses) > 0 {
		items := make([]any, len(commandResponses))
		for idx, params := range commandResponses {
			items[idx] = map[string]any{
				"completionCode": float64(0),
				"reasonCode":     float64(0),
				"parameters":     params,
			}
		}
		body["commandResponse"] = items
	}
	bodyBytes, _ := json.Marshal(body)
	transport.responses = append(transport.responses, mockResponse{
		Response: &mqrestadmin.TransportResponse{
			StatusCode: 200,
			Body:       string(bodyBytes),
			Headers:    map[string]string{},
		},
	})
}

// addCommandErrorResponse queues an MQSC command error response.
func (transport *mockTransport) addCommandErrorResponse(completionCode, reasonCode int) {
	body := map[string]any{
		"overallCompletionCode": float64(completionCode),
		"overallReasonCode":     float64(reasonCode),
	}
	bodyBytes, _ := json.Marshal(body)
	transport.responses = append(transport.responses, mockResponse{
		Response: &mqrestadmin.TransportResponse{
			StatusCode: 200,
			Body:       string(bodyBytes),
			Headers:    map[string]string{},
		},
	})
}

func (transport *mockTransport) PostJSON(_ context.Context, url string, payload map[string]any,
	headers map[string]string, timeout time.Duration, verifyTLS bool,
) (*mqrestadmin.TransportResponse, error) {
	transport.calls = append(transport.calls, mockCall{
		URL:       url,
		Payload:   payload,
		Headers:   headers,
		Timeout:   timeout,
		VerifyTLS: verifyTLS,
	})

	if transport.callIndex >= len(transport.responses) {
		return nil, fmt.Errorf("mock transport: no response configured for call %d", transport.callIndex)
	}

	response := transport.responses[transport.callIndex]
	transport.callIndex++

	return response.Response, response.Err
}

// newTestSession creates a Session with a mock transport for testing.
func newTestSession(t *testing.T, transport *mockTransport) *mqrestadmin.Session {
	t.Helper()
	token := "test"
	session, err := mqrestadmin.NewSession(
		"https://localhost:9443/ibmmq/rest/v2",
		"QM1",
		mqrestadmin.BasicAuth{Username: "admin", Password: "admin"},
		mqrestadmin.WithTransport(transport),
		mqrestadmin.WithCSRFToken(&token),
	)
	if err != nil {
		t.Fatalf("newTestSession: %v", err)
	}
	return session
}
