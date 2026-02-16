package mqrestadmin

import (
	"crypto/tls"
	"encoding/base64"
	"net/http"
)

// Credentials represents authentication credentials for the MQ REST API.
// The interface is sealed via an unexported method; only BasicAuth,
// LTPAAuth, and CertificateAuth implement it.
type Credentials interface {
	// applyAuth configures authentication on an HTTP request. This unexported
	// method restricts implementations to this package.
	applyAuth(request *http.Request, session *Session)
	// sealed prevents external implementations.
	sealed()
}

// BasicAuth provides HTTP Basic authentication credentials.
type BasicAuth struct {
	Username string
	Password string
}

func (auth BasicAuth) applyAuth(request *http.Request, _ *Session) {
	encoded := base64.StdEncoding.EncodeToString(
		[]byte(auth.Username + ":" + auth.Password),
	)
	request.Header.Set("Authorization", "Basic "+encoded)
}

//coverage:ignore
func (BasicAuth) sealed() {}

// LTPAAuth provides LTPA token-based authentication. The session performs a
// login at construction time to obtain an LtpaToken2 cookie, which is
// included in all subsequent requests.
type LTPAAuth struct {
	Username string
	Password string
}

func (auth LTPAAuth) applyAuth(request *http.Request, session *Session) {
	if session.ltpaToken != "" {
		request.Header.Set("Cookie", "LtpaToken2="+session.ltpaToken)
	}
}

//coverage:ignore
func (LTPAAuth) sealed() {}

// CertificateAuth provides mutual TLS (mTLS) authentication using a client
// certificate. The certificate is configured on the transport's TLS settings.
type CertificateAuth struct {
	// CertPath is the path to the client certificate PEM file.
	CertPath string
	// KeyPath is the path to the private key PEM file. If empty, the
	// certificate file is expected to contain both the certificate and key.
	KeyPath string
}

//coverage:ignore
func (CertificateAuth) applyAuth(_ *http.Request, _ *Session) {
	// Certificate auth is handled at the TLS transport level, not per-request.
}

//coverage:ignore
func (CertificateAuth) sealed() {}

// loadTLSCertificate loads the client certificate for mTLS authentication.
func (auth CertificateAuth) loadTLSCertificate() (*tls.Certificate, error) {
	keyPath := auth.KeyPath
	if keyPath == "" {
		keyPath = auth.CertPath
	}
	certificate, err := tls.LoadX509KeyPair(auth.CertPath, keyPath)
	if err != nil {
		return nil, err
	}
	return &certificate, nil
}
