package mqrestadmin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
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
	Response *TransportResponse
	Err      error
}

func newMockTransport() *mockTransport {
	return &mockTransport{}
}

// addResponse queues a response to be returned on the next PostJSON call.
func (transport *mockTransport) addResponse(statusCode int, body map[string]any, headers map[string]string) {
	bodyBytes, _ := json.Marshal(body)
	if headers == nil {
		headers = map[string]string{}
	}
	transport.responses = append(transport.responses, mockResponse{
		Response: &TransportResponse{
			StatusCode: statusCode,
			Body:       string(bodyBytes),
			Headers:    headers,
		},
	})
}

// addErrorResponse queues an error to be returned on the next PostJSON call.
func (transport *mockTransport) addErrorResponse(err error) {
	transport.responses = append(transport.responses, mockResponse{Err: err})
}

// addSuccessResponse queues a successful MQSC response with no command errors.
func (transport *mockTransport) addSuccessResponse(commandResponses ...map[string]any) {
	body := map[string]any{
		"overallCompletionCode": float64(0),
		"overallReasonCode":    float64(0),
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
	transport.addResponse(200, body, nil)
}

// addCommandErrorResponse queues an MQSC command error response.
func (transport *mockTransport) addCommandErrorResponse(completionCode, reasonCode int) {
	body := map[string]any{
		"overallCompletionCode": float64(completionCode),
		"overallReasonCode":    float64(reasonCode),
	}
	transport.addResponse(200, body, nil)
}

func (transport *mockTransport) PostJSON(_ context.Context, url string, payload map[string]any,
	headers map[string]string, timeout time.Duration, verifyTLS bool,
) (*TransportResponse, error) {
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

// lastCall returns the most recent call, or panics if none.
func (transport *mockTransport) lastCall() mockCall {
	if len(transport.calls) == 0 {
		panic("mock transport: no calls recorded")
	}
	return transport.calls[len(transport.calls)-1]
}

// callCount returns the number of calls made.
func (transport *mockTransport) callCount() int {
	return len(transport.calls)
}

// mockClock is a test clock that allows manual time control.
type mockClock struct {
	currentTime time.Time
	sleepCalls  []time.Duration
}

func newMockClock() *mockClock {
	return &mockClock{
		currentTime: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func (clock *mockClock) now() time.Time {
	return clock.currentTime
}

func (clock *mockClock) sleep(duration time.Duration) {
	clock.sleepCalls = append(clock.sleepCalls, duration)
	clock.currentTime = clock.currentTime.Add(duration)
}

// newTestSession creates a Session with a mock transport for testing.
// The session has mapping disabled by default to simplify test assertions.
func newTestSession(transport *mockTransport) *Session {
	csrfToken := "local"
	return &Session{
		restBaseURL:   "https://localhost:9443/ibmmq/rest/v2",
		qmgrName:      "QM1",
		credentials:   BasicAuth{Username: "admin", Password: "admin"},
		transport:     transport,
		verifyTLS:     true,
		timeout:       30 * time.Second,
		mapAttributes: false,
		mappingStrict: true,
		csrfToken:     &csrfToken,
		clock:         systemClock{},
	}
}

// newTestSessionWithMapping creates a Session with a mock transport and
// attribute mapping enabled.
func newTestSessionWithMapping(transport *mockTransport) *Session {
	session := newTestSession(transport)
	session.mapAttributes = true
	mapper, err := newAttributeMapper()
	if err != nil {
		panic(fmt.Sprintf("failed to create attribute mapper: %v", err))
	}
	session.mapper = mapper
	return session
}

// newTestSessionWithClock creates a Session with a mock transport and mock
// clock for sync testing.
func newTestSessionWithClock(transport *mockTransport, clock *mockClock) *Session {
	session := newTestSession(transport)
	session.clock = clock
	return session
}
