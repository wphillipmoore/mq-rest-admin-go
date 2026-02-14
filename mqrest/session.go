package mqrest

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	ltpaLoginPath  = "/login"
	ltpaCookieName = "LtpaToken2"
	mqscEndpoint   = "/admin/action/qmgr/%s/mqsc"

	defaultTimeout       = 30 * time.Second
	defaultCSRFToken     = "local"
	defaultSyncTimeout   = 30 * time.Second
	defaultPollInterval  = 1 * time.Second
)

// Session manages communication with the IBM MQ administrative REST API.
type Session struct {
	restBaseURL   string
	qmgrName      string
	credentials   Credentials
	transport     Transport
	gatewayQmgr   string
	verifyTLS     bool
	timeout       time.Duration
	mapAttributes bool
	mappingStrict bool
	csrfToken     *string
	mapper        *attributeMapper
	ltpaToken     string
	clock         clock

	// LastHTTPStatus is the HTTP status code from the most recent command.
	LastHTTPStatus int
	// LastResponseText is the raw response body from the most recent command.
	LastResponseText string
	// LastResponsePayload is the parsed JSON response from the most recent command.
	LastResponsePayload map[string]any
	// LastCommandPayload is the request payload sent for the most recent command.
	LastCommandPayload map[string]any
}

// clock abstracts time operations for testability.
type clock interface {
	now() time.Time
	sleep(duration time.Duration)
}

type systemClock struct{}

func (systemClock) now() time.Time              { return time.Now() }
func (systemClock) sleep(duration time.Duration) { time.Sleep(duration) }

// Option configures a Session during construction.
type Option func(*sessionConfig)

type sessionConfig struct {
	transport            Transport
	gatewayQmgr          string
	verifyTLS            bool
	timeout              time.Duration
	mapAttributes        bool
	mappingStrict        bool
	csrfToken            *string
	mappingOverrides     map[string]any
	mappingOverridesMode MappingOverrideMode
}

func defaultConfig() sessionConfig {
	csrfToken := defaultCSRFToken
	return sessionConfig{
		verifyTLS:     true,
		timeout:       defaultTimeout,
		mapAttributes: true,
		mappingStrict: true,
		csrfToken:     &csrfToken,
	}
}

// WithTransport sets a custom Transport implementation. If not provided,
// a default HTTPTransport is used.
func WithTransport(transport Transport) Option {
	return func(config *sessionConfig) {
		config.transport = transport
	}
}

// WithGatewayQmgr sets the gateway queue manager name for routing commands
// through a gateway in multi-queue-manager deployments.
func WithGatewayQmgr(name string) Option {
	return func(config *sessionConfig) {
		config.gatewayQmgr = name
	}
}

// WithVerifyTLS controls TLS certificate verification. Set to false to allow
// self-signed certificates (useful for development/testing).
func WithVerifyTLS(verify bool) Option {
	return func(config *sessionConfig) {
		config.verifyTLS = verify
	}
}

// WithTimeout sets the HTTP request timeout. Use zero for no timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(config *sessionConfig) {
		config.timeout = timeout
	}
}

// WithMapAttributes controls whether attribute names are translated between
// snake_case and MQSC parameter names. Defaults to true.
func WithMapAttributes(enabled bool) Option {
	return func(config *sessionConfig) {
		config.mapAttributes = enabled
	}
}

// WithMappingStrict controls whether unknown attributes cause an error (true)
// or are passed through silently (false). Defaults to true.
func WithMappingStrict(strict bool) Option {
	return func(config *sessionConfig) {
		config.mappingStrict = strict
	}
}

// WithCSRFToken sets the CSRF token value. The IBM MQ REST API requires the
// ibm-mq-rest-csrf-token header. Set to nil to omit the header entirely.
func WithCSRFToken(token *string) Option {
	return func(config *sessionConfig) {
		config.csrfToken = token
	}
}

// WithMappingOverrides provides custom mapping data that is overlaid on or
// replaces the default mapping definitions.
func WithMappingOverrides(overrides map[string]any, mode MappingOverrideMode) Option {
	return func(config *sessionConfig) {
		config.mappingOverrides = overrides
		config.mappingOverridesMode = mode
	}
}

// WithBasicAuth configures HTTP Basic authentication.
func WithBasicAuth(_, _ string) Option {
	return func(_ *sessionConfig) {
		// Credentials are set via a special path; see NewSession.
	}
}

// NewSession creates a new Session connected to the specified MQ REST API
// endpoint. The credentials parameter is required. Options configure
// transport, timeouts, TLS, attribute mapping, and other behavior.
//
// If LTPAAuth credentials are provided, the session performs a login at
// construction time to obtain an LTPA token.
func NewSession(restBaseURL, qmgrName string, credentials Credentials, opts ...Option) (*Session, error) {
	config := defaultConfig()
	for _, opt := range opts {
		opt(&config)
	}

	// Default transport
	transport := config.transport
	if transport == nil {
		httpTransport := &HTTPTransport{}
		// Configure mTLS if using certificate auth
		if certAuth, isCert := credentials.(CertificateAuth); isCert {
			certificate, err := certAuth.loadTLSCertificate()
			if err != nil {
				return nil, fmt.Errorf("load client certificate: %w", err)
			}
			httpTransport.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{*certificate},
			}
		}
		transport = httpTransport
	}

	// Attribute mapper
	var mapper *attributeMapper
	var err error
	if config.mapAttributes {
		if config.mappingOverrides != nil {
			mapper, err = newAttributeMapperWithOverrides(config.mappingOverrides, config.mappingOverridesMode)
		} else {
			mapper, err = newAttributeMapper()
		}
		if err != nil {
			return nil, fmt.Errorf("initialize attribute mapper: %w", err)
		}
	}

	session := &Session{
		restBaseURL:   strings.TrimRight(restBaseURL, "/"),
		qmgrName:      qmgrName,
		credentials:   credentials,
		transport:     transport,
		gatewayQmgr:   config.gatewayQmgr,
		verifyTLS:     config.verifyTLS,
		timeout:       config.timeout,
		mapAttributes: config.mapAttributes,
		mappingStrict: config.mappingStrict,
		csrfToken:     config.csrfToken,
		mapper:        mapper,
		clock:         systemClock{},
	}

	// LTPA login
	if ltpaAuth, isLTPA := credentials.(LTPAAuth); isLTPA {
		if err := session.performLTPALogin(ltpaAuth); err != nil {
			return nil, err
		}
	}

	return session, nil
}

// QmgrName returns the queue manager name for this session.
func (session *Session) QmgrName() string {
	return session.qmgrName
}

// GatewayQmgr returns the gateway queue manager name, or empty string if none.
func (session *Session) GatewayQmgr() string {
	return session.gatewayQmgr
}

// mqscCommand is the core dispatch method. It builds the MQSC command payload,
// sends it to the REST API, parses the response, and optionally maps attribute
// names.
func (session *Session) mqscCommand(ctx context.Context, command, mqscQualifier string,
	name *string, requestParameters map[string]any, responseParameters []string,
	_ *string, isDisplay bool,
) ([]map[string]any, error) {
	upperCommand := strings.ToUpper(command)
	upperQualifier := strings.ToUpper(mqscQualifier)

	// Copy request parameters to avoid mutating caller's map
	params := make(map[string]any)
	for key, value := range requestParameters {
		params[key] = value
	}

	// Default responseParameters for DISPLAY commands
	if isDisplay && responseParameters == nil {
		responseParameters = []string{"all"}
	}

	// Resolve mapping qualifier
	var mappingQualifier string
	if session.mapAttributes && session.mapper != nil {
		mappingQualifier = session.mapper.resolveMappingQualifier(upperCommand, upperQualifier)

		// Map request attributes
		if len(params) > 0 && mappingQualifier != "" {
			mapped, issues := session.mapper.mapRequestAttributes(mappingQualifier, params, session.mappingStrict)
			if session.mappingStrict && len(issues) > 0 {
				return nil, &MappingError{Issues: issues}
			}
			params = mapped
		}

		// Map response parameter names
		if len(responseParameters) > 0 && mappingQualifier != "" {
			responseParameters = session.mapResponseParameterNames(mappingQualifier, responseParameters)
		}

		// Expand response parameter macros
		if len(responseParameters) > 0 {
			responseParameters = session.mapper.resolveResponseParameterMacros(
				upperCommand, upperQualifier, responseParameters)
		}
	}

	// Build payload
	payload := session.buildCommandPayload(upperCommand, upperQualifier, name, params, responseParameters)
	session.LastCommandPayload = payload

	// Build URL and headers
	url := session.buildMQSCURL()
	headers := session.buildHeaders()

	// Send request
	response, err := session.transport.PostJSON(ctx, url, payload, headers, session.timeout, session.verifyTLS)
	if err != nil {
		return nil, err
	}

	session.LastHTTPStatus = response.StatusCode
	session.LastResponseText = response.Body

	// Check for auth errors
	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return nil, &AuthError{URL: url, StatusCode: response.StatusCode}
	}

	// Parse response JSON
	responsePayload, err := parseResponsePayload(response.Body)
	if err != nil {
		return nil, &ResponseError{ResponseText: response.Body, StatusCode: response.StatusCode}
	}
	session.LastResponsePayload = responsePayload

	// Check for command errors
	if err := checkCommandErrors(responsePayload, response.StatusCode); err != nil {
		return nil, err
	}

	// Extract command response objects
	objects := extractCommandResponseObjects(responsePayload)

	// Map response attributes
	if session.mapAttributes && session.mapper != nil && mappingQualifier != "" && len(objects) > 0 {
		mapped, issues := session.mapper.mapResponseList(mappingQualifier, objects, session.mappingStrict)
		if session.mappingStrict && len(issues) > 0 {
			return nil, &MappingError{Issues: issues}
		}
		objects = mapped
	}

	return objects, nil
}

func (session *Session) buildMQSCURL() string {
	return session.restBaseURL + fmt.Sprintf(mqscEndpoint, session.qmgrName)
}

func (session *Session) buildHeaders() map[string]string {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// Apply auth
	fakeRequest := &http.Request{Header: make(http.Header)}
	session.credentials.applyAuth(fakeRequest, session)
	for key, values := range fakeRequest.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// CSRF token
	if session.csrfToken != nil {
		headers["ibm-mq-rest-csrf-token"] = *session.csrfToken
	}

	// Gateway queue manager
	if session.gatewayQmgr != "" {
		headers["ibm-mq-rest-gateway-qmgr"] = session.gatewayQmgr
	}

	return headers
}

func (session *Session) buildCommandPayload(command, qualifier string, name *string,
	requestParameters map[string]any, responseParameters []string,
) map[string]any {
	payload := map[string]any{
		"type":    "runCommandJSON",
		"command": command,
	}

	payload["qualifier"] = qualifier

	if name != nil {
		payload["name"] = *name
	}

	if len(requestParameters) > 0 {
		payload["parameters"] = requestParameters
	}

	if len(responseParameters) > 0 {
		payload["responseParameters"] = responseParameters
	}

	return payload
}

func (session *Session) mapResponseParameterNames(qualifier string, params []string) []string {
	qualifierData, exists := session.mapper.data.Qualifiers[qualifier]
	if !exists {
		return params
	}

	mapped := make([]string, len(params))
	for idx, param := range params {
		if mqscName, found := qualifierData.RequestKeyMap[param]; found {
			mapped[idx] = mqscName
		} else {
			mapped[idx] = param
		}
	}
	return mapped
}

// performLTPALogin authenticates with the MQ REST API using LTPA credentials
// and stores the resulting LtpaToken2 cookie for subsequent requests.
func (session *Session) performLTPALogin(auth LTPAAuth) error {
	loginURL := session.restBaseURL + ltpaLoginPath

	loginPayload := map[string]any{
		"username": auth.Username,
		"password": auth.Password,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}
	if session.csrfToken != nil {
		headers["ibm-mq-rest-csrf-token"] = *session.csrfToken
	}

	response, err := session.transport.PostJSON(
		context.Background(), loginURL, loginPayload, headers, session.timeout, session.verifyTLS)
	if err != nil {
		return fmt.Errorf("LTPA login request failed: %w", err)
	}

	if response.StatusCode >= 400 {
		return &AuthError{URL: loginURL, StatusCode: response.StatusCode}
	}

	token := extractLTPAToken(response.Headers)
	if token == "" {
		return &AuthError{URL: loginURL, StatusCode: response.StatusCode}
	}

	session.ltpaToken = token
	return nil
}

func extractLTPAToken(headers map[string]string) string {
	for key, value := range headers {
		if !strings.EqualFold(key, "Set-Cookie") {
			continue
		}
		for _, part := range strings.Split(value, ";") {
			trimmed := strings.TrimSpace(part)
			if strings.HasPrefix(trimmed, ltpaCookieName+"=") {
				return strings.TrimPrefix(trimmed, ltpaCookieName+"=")
			}
		}
	}
	return ""
}

func parseResponsePayload(body string) (map[string]any, error) {
	var result map[string]any
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func checkCommandErrors(payload map[string]any, httpStatus int) error {
	// Check overall completion and reason codes
	if hasErrorCodes(payload["overallCompletionCode"], payload["overallReasonCode"]) {
		return &CommandError{Payload: payload, StatusCode: httpStatus}
	}

	// Check per-item codes in commandResponse
	if commandResponse, exists := payload["commandResponse"]; exists {
		if items, isList := commandResponse.([]any); isList {
			for _, item := range items {
				if itemMap, isMap := item.(map[string]any); isMap {
					if hasErrorCodes(itemMap["completionCode"], itemMap["reasonCode"]) {
						return &CommandError{Payload: payload, StatusCode: httpStatus}
					}
				}
			}
		}
	}

	return nil
}

func hasErrorCodes(completionCode, reasonCode any) bool {
	return isNonZeroNumber(completionCode) || isNonZeroNumber(reasonCode)
}

func isNonZeroNumber(value any) bool {
	if value == nil {
		return false
	}
	switch typed := value.(type) {
	case float64:
		return typed != 0
	case int:
		return typed != 0
	default:
		return false
	}
}

func extractCommandResponseObjects(payload map[string]any) []map[string]any {
	commandResponse, exists := payload["commandResponse"]
	if !exists {
		return nil
	}

	items, isList := commandResponse.([]any)
	if !isList {
		return nil
	}

	var result []map[string]any
	for _, item := range items {
		itemMap, isMap := item.(map[string]any)
		if !isMap {
			continue
		}

		params, hasParams := itemMap["parameters"]
		if !hasParams {
			continue
		}

		paramsMap, isParamsMap := params.(map[string]any)
		if !isParamsMap {
			continue
		}

		// Check for nested objects array (multi-row results like QSTATUS HANDLE)
		if objectsRaw, hasObjects := paramsMap["objects"]; hasObjects {
			if objects, isObjectList := objectsRaw.([]any); isObjectList {
				result = append(result, flattenNestedObjects(paramsMap, objects)...)
				continue
			}
		}

		result = append(result, paramsMap)
	}

	return result
}

// flattenNestedObjects merges parent-level fields into each nested object.
func flattenNestedObjects(parent map[string]any, objects []any) []map[string]any {
	var result []map[string]any

	for _, object := range objects {
		objectMap, isMap := object.(map[string]any)
		if !isMap {
			continue
		}

		flattened := make(map[string]any)
		// Copy parent fields (except "objects")
		for key, value := range parent {
			if key != "objects" {
				flattened[key] = value
			}
		}
		// Overlay nested object fields
		for key, value := range objectMap {
			flattened[key] = value
		}
		result = append(result, flattened)
	}

	return result
}
