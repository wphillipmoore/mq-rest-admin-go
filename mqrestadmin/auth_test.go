package mqrestadmin

import (
	"encoding/base64"
	"net/http"
	"strings"
	"testing"
)

func TestBasicAuth_ApplyAuth(t *testing.T) {
	auth := BasicAuth{Username: "admin", Password: "secret"}
	request := &http.Request{Header: make(http.Header)}

	auth.applyAuth(request, nil)

	authHeader := request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Basic ") {
		t.Fatalf("expected Basic auth header, got %q", authHeader)
	}

	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}
	if string(decoded) != "admin:secret" {
		t.Errorf("decoded credentials = %q, want %q", string(decoded), "admin:secret")
	}
}

func TestLTPAAuth_ApplyAuth_WithToken(t *testing.T) {
	auth := LTPAAuth{Username: "admin", Password: "secret"}
	session := &Session{ltpaToken: "test-token-123"}
	request := &http.Request{Header: make(http.Header)}

	auth.applyAuth(request, session)

	cookie := request.Header.Get("Cookie")
	if cookie != "LtpaToken2=test-token-123" {
		t.Errorf("Cookie = %q, want %q", cookie, "LtpaToken2=test-token-123")
	}
}

func TestLTPAAuth_ApplyAuth_NoToken(t *testing.T) {
	auth := LTPAAuth{Username: "admin", Password: "secret"}
	session := &Session{}
	request := &http.Request{Header: make(http.Header)}

	auth.applyAuth(request, session)

	cookie := request.Header.Get("Cookie")
	if cookie != "" {
		t.Errorf("Cookie should be empty when no token, got %q", cookie)
	}
}

func TestCertificateAuth_ApplyAuth_NoOp(t *testing.T) {
	auth := CertificateAuth{CertPath: "/path/to/cert.pem"}
	request := &http.Request{Header: make(http.Header)}

	// Should not panic or set any headers
	auth.applyAuth(request, nil)

	if request.Header.Get("Authorization") != "" {
		t.Error("CertificateAuth should not set Authorization header")
	}
}

func TestExtractLTPAToken(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected string
	}{
		{
			name:     "standard cookie",
			headers:  map[string]string{"Set-Cookie": "LtpaToken2=abc123; Path=/; Secure"},
			expected: "abc123",
		},
		{
			name:     "case insensitive header",
			headers:  map[string]string{"set-cookie": "LtpaToken2=abc123; Path=/"},
			expected: "abc123",
		},
		{
			name:     "no cookie header",
			headers:  map[string]string{},
			expected: "",
		},
		{
			name:     "wrong cookie name",
			headers:  map[string]string{"Set-Cookie": "SessionId=xyz; Path=/"},
			expected: "",
		},
		{
			name:     "token only",
			headers:  map[string]string{"Set-Cookie": "LtpaToken2=simple_token"},
			expected: "simple_token",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := extractLTPAToken(test.headers)
			if result != test.expected {
				t.Errorf("extractLTPAToken() = %q, want %q", result, test.expected)
			}
		})
	}
}

func TestBuildHeaders_BasicAuth(t *testing.T) {
	session := &Session{
		credentials: BasicAuth{Username: "admin", Password: "pass"},
		csrfToken:   ptrString("local"),
	}

	headers := session.buildHeaders()

	if headers["Accept"] != "application/json" {
		t.Errorf("Accept = %q", headers["Accept"])
	}
	if !strings.HasPrefix(headers["Authorization"], "Basic ") {
		t.Errorf("Authorization = %q, want Basic prefix", headers["Authorization"])
	}
	if headers["ibm-mq-rest-csrf-token"] != "local" {
		t.Errorf("csrf-token = %q", headers["ibm-mq-rest-csrf-token"])
	}
}

func TestBuildHeaders_NoCSRFToken(t *testing.T) {
	session := &Session{
		credentials: BasicAuth{Username: "admin", Password: "pass"},
		csrfToken:   nil,
	}

	headers := session.buildHeaders()

	if _, hasCSRF := headers["ibm-mq-rest-csrf-token"]; hasCSRF {
		t.Error("should not include CSRF token when nil")
	}
}

func TestBuildHeaders_GatewayQmgr(t *testing.T) {
	session := &Session{
		credentials: BasicAuth{Username: "admin", Password: "pass"},
		gatewayQmgr: "GATEWAY",
		csrfToken:   ptrString("local"),
	}

	headers := session.buildHeaders()

	if headers["ibm-mq-rest-gateway-qmgr"] != "GATEWAY" {
		t.Errorf("gateway-qmgr = %q, want GATEWAY", headers["ibm-mq-rest-gateway-qmgr"])
	}
}

func TestSealed_BasicAuth(_ *testing.T) {
	auth := BasicAuth{}
	auth.sealed()
}

func TestSealed_LTPAAuth(_ *testing.T) {
	auth := LTPAAuth{}
	auth.sealed()
}

func TestSealed_CertificateAuth(_ *testing.T) {
	auth := CertificateAuth{}
	auth.sealed()
}

func TestLoadTLSCertificate_Success(t *testing.T) {
	certPEM, keyPEM := generateSelfSignedCert(t)

	certFile := writeTempFile(t, "cert-*.pem", certPEM)
	keyFile := writeTempFile(t, "key-*.pem", keyPEM)

	auth := CertificateAuth{CertPath: certFile, KeyPath: keyFile}
	cert, err := auth.loadTLSCertificate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cert == nil {
		t.Fatal("expected non-nil certificate")
	}
}

func TestLoadTLSCertificate_CombinedFile(t *testing.T) {
	certPEM, keyPEM := generateSelfSignedCert(t)
	combined := certPEM
	combined = append(combined, keyPEM...)
	combinedFile := writeTempFile(t, "combined-*.pem", combined)

	auth := CertificateAuth{CertPath: combinedFile}
	cert, err := auth.loadTLSCertificate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cert == nil {
		t.Fatal("expected non-nil certificate")
	}
}

func TestLoadTLSCertificate_InvalidPath(t *testing.T) {
	auth := CertificateAuth{CertPath: "/nonexistent/cert.pem", KeyPath: "/nonexistent/key.pem"}
	_, err := auth.loadTLSCertificate()
	if err == nil {
		t.Fatal("expected error for nonexistent certificate files")
	}
}

func ptrString(value string) *string {
	return &value
}
