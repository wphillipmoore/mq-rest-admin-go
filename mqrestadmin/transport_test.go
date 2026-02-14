package mqrestadmin

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPTransport_PostJSON_Success(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		if r.Header.Get("X-Custom") != "header-value" {
			t.Errorf("X-Custom header = %q, want header-value", r.Header.Get("X-Custom"))
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}

		w.Header().Set("X-Response", "resp-value")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	}))
	defer server.Close()

	transport := &HTTPTransport{
		TLSConfig: server.TLS.Clone(),
	}

	response, err := transport.PostJSON(
		context.Background(),
		server.URL+"/test",
		map[string]any{"key": "value"},
		map[string]string{"X-Custom": "header-value"},
		30*time.Second,
		false,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", response.StatusCode)
	}
	if response.Headers["X-Response"] != "resp-value" {
		t.Errorf("X-Response header = %q, want resp-value", response.Headers["X-Response"])
	}
	if response.Body == "" {
		t.Error("expected non-empty response body")
	}
}

func TestHTTPTransport_PostJSON_NetworkError(t *testing.T) {
	transport := &HTTPTransport{}

	_, err := transport.PostJSON(
		context.Background(),
		"https://127.0.0.1:1/nonexistent",
		map[string]any{"key": "value"},
		nil,
		1*time.Second,
		false,
	)
	if err == nil {
		t.Fatal("expected error for unreachable host")
	}

	var transportErr *TransportError
	if !errors.As(err, &transportErr) {
		t.Fatalf("expected TransportError, got %T: %v", err, err)
	}
}

func TestBuildClient_NilTLSConfig(t *testing.T) {
	transport := &HTTPTransport{TLSConfig: nil}
	client := transport.buildClient(10*time.Second, true)
	if client.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s", client.Timeout)
	}
}

func TestBuildClient_WithTLSConfig(t *testing.T) {
	tlsConfig := &tls.Config{
		ServerName: "test-server",
	}
	transport := &HTTPTransport{TLSConfig: tlsConfig}
	client := transport.buildClient(10*time.Second, true)

	httpTransport := client.Transport.(*http.Transport)
	if httpTransport.TLSClientConfig.ServerName != "test-server" {
		t.Error("expected TLSConfig to be cloned with ServerName")
	}
	if httpTransport.TLSClientConfig.InsecureSkipVerify {
		t.Error("expected InsecureSkipVerify = false when verifyTLS = true")
	}
}

func TestHTTPTransport_PostJSON_InvalidURL(t *testing.T) {
	transport := &HTTPTransport{}

	_, err := transport.PostJSON(
		context.Background(),
		"://invalid-url",
		map[string]any{"key": "value"},
		nil,
		1*time.Second,
		false,
	)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}

	var transportErr *TransportError
	if !errors.As(err, &transportErr) {
		t.Fatalf("expected TransportError, got %T: %v", err, err)
	}
}

func TestBuildClient_VerifyTLSFalse(t *testing.T) {
	transport := &HTTPTransport{}
	client := transport.buildClient(10*time.Second, false)

	httpTransport := client.Transport.(*http.Transport)
	if !httpTransport.TLSClientConfig.InsecureSkipVerify {
		t.Error("expected InsecureSkipVerify = true when verifyTLS = false")
	}
}
