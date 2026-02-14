package mqrestadmin

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Transport defines the interface for sending HTTP requests to the MQ REST
// API. The default implementation uses net/http. Custom implementations can
// be provided for testing or specialized HTTP handling.
type Transport interface {
	// PostJSON sends a JSON POST request and returns the response.
	PostJSON(ctx context.Context, url string, payload map[string]any,
		headers map[string]string, timeout time.Duration, verifyTLS bool,
	) (*TransportResponse, error)
}

// TransportResponse holds the HTTP response data from a transport call.
type TransportResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// HTTPTransport is the default Transport implementation using net/http.
type HTTPTransport struct {
	// TLSConfig is an optional TLS configuration for mTLS or custom CA
	// certificates. If nil, the default TLS configuration is used.
	TLSConfig *tls.Config
}

// PostJSON sends a JSON POST request using net/http.
func (transport *HTTPTransport) PostJSON(ctx context.Context, url string,
	payload map[string]any, headers map[string]string, timeout time.Duration,
	verifyTLS bool,
) (*TransportResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, &TransportError{URL: url, Err: fmt.Errorf("marshal payload: %w", err)}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, &TransportError{URL: url, Err: fmt.Errorf("create request: %w", err)}
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := transport.buildClient(timeout, verifyTLS)

	response, err := client.Do(request)
	if err != nil {
		return nil, &TransportError{URL: url, Err: err}
	}
	defer func() { _ = response.Body.Close() }()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, &TransportError{URL: url, Err: fmt.Errorf("read response: %w", err)}
	}

	responseHeaders := make(map[string]string, len(response.Header))
	for key, values := range response.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	return &TransportResponse{
		StatusCode: response.StatusCode,
		Body:       string(responseBody),
		Headers:    responseHeaders,
	}, nil
}

func (transport *HTTPTransport) buildClient(timeout time.Duration, verifyTLS bool) *http.Client {
	tlsConfiguration := transport.TLSConfig
	if tlsConfiguration == nil {
		tlsConfiguration = &tls.Config{}
	} else {
		tlsConfiguration = tlsConfiguration.Clone()
	}

	if !verifyTLS {
		tlsConfiguration.InsecureSkipVerify = true
	}

	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfiguration,
		},
	}
}
