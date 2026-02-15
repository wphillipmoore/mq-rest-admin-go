# Authentication

## Overview

The authentication module provides credential types for the three
authentication modes supported by the IBM MQ REST API: mutual TLS (mTLS)
client certificates, LTPA token, and HTTP Basic.

Pass a credential value to `NewSession` as the third argument. Always use TLS
(`https://`) for production deployments to protect credentials and data in
transit.

```go
import "github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"

// mTLS client certificate auth -- strongest; no shared secrets
session, err := mqrestadmin.NewSession(
    "https://mq-host:9443/ibmmq/rest/v2", "QM1",
    mqrestadmin.CertificateAuth{CertPath: "/path/to/cert.pem", KeyPath: "/path/to/key.pem"},
)

// LTPA token auth -- credentials sent once at login, then cookie-based
session, err := mqrestadmin.NewSession(
    "https://mq-host:9443/ibmmq/rest/v2", "QM1",
    mqrestadmin.LTPAAuth{Username: "user", Password: "pass"},
)

// Basic auth -- credentials sent with every request
session, err := mqrestadmin.NewSession(
    "https://mq-host:9443/ibmmq/rest/v2", "QM1",
    mqrestadmin.BasicAuth{Username: "user", Password: "pass"},
)
```

## Credentials

`Credentials` is a sealed interface representing authentication credentials for
the MQ REST API. The interface is sealed via an unexported method; only
`BasicAuth`, `LTPAAuth`, and `CertificateAuth` implement it.

```go
type Credentials interface {
    applyAuth(request *http.Request, session *Session)
    sealed()
}
```

Because `Credentials` is sealed, external packages cannot create new
implementations. This ensures the session can exhaustively handle all
credential types.

## CertificateAuth

Client certificate authentication via TLS mutual authentication (mTLS). This is
the strongest authentication mode -- no shared secrets cross the wire.

```go
type CertificateAuth struct {
    CertPath string  // path to client certificate PEM file
    KeyPath  string  // path to private key PEM file (empty if combined)
}
```

```go
// Separate certificate and key files
creds := mqrestadmin.CertificateAuth{
    CertPath: "/path/to/cert.pem",
    KeyPath:  "/path/to/key.pem",
}

// Combined cert+key file (omit KeyPath)
creds := mqrestadmin.CertificateAuth{
    CertPath: "/path/to/combined.pem",
}
```

No `Authorization` header is sent; authentication is handled at the TLS layer.
When `CertificateAuth` is provided and no custom transport is set, `NewSession`
automatically configures the default `HTTPTransport` with the loaded
`tls.Certificate`.

## LTPAAuth

LTPA token-based authentication. Credentials are sent once during a `/login`
request at session construction; subsequent API calls carry only the LTPA
cookie.

```go
type LTPAAuth struct {
    Username string
    Password string
}
```

```go
creds := mqrestadmin.LTPAAuth{Username: "mqadmin", Password: "passw0rd"}
```

The session performs the login automatically in `NewSession` and extracts the
`LtpaToken2` cookie for subsequent requests. If the login fails, `NewSession`
returns an `*AuthError`.

## BasicAuth

HTTP Basic authentication. The `Authorization` header is constructed from the
username and password and sent with every request.

```go
type BasicAuth struct {
    Username string
    Password string
}
```

```go
creds := mqrestadmin.BasicAuth{Username: "mqadmin", Password: "passw0rd"}
```

!!! note
    All examples and documentation in this project use LTPA as the default
    authentication method. If you see `LTPAAuth{...}` in an example, you
    can substitute `BasicAuth{...}` or `CertificateAuth{...}` based on
    your environment.

## Choosing between LTPA and Basic authentication

Both LTPA and Basic authentication use a username and password. The key
difference is how often those credentials cross the wire.

**LTPA is the recommended choice for username/password authentication.**
Credentials are sent once during the `/login` request; subsequent API calls
carry only the LTPA cookie. This reduces credential exposure and is more
efficient for sessions that issue many commands.

**Use Basic authentication as a fallback when:**

- The mqweb configuration does not enable the `/login` endpoint (for example,
  minimal container images that only expose the REST API).
- A reverse proxy or API gateway handles authentication and forwards a Basic
  auth header; cookie-based flows may not survive the proxy.
- Single-command scripts where the login round-trip doubles the request count
  for no security benefit.
- Long-running sessions where LTPA token expiry (typically two hours) could
  cause mid-operation failures; the library does not currently re-authenticate
  automatically.
- Local development or CI against a `localhost` container, where transport
  security is not a concern.
