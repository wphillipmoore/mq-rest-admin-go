package mqrestadmin

import (
	"fmt"
	"strings"
)

// TransportError indicates a network or connection failure during an HTTP
// request to the MQ REST API.
type TransportError struct {
	URL string
	Err error
}

func (e *TransportError) Error() string {
	return fmt.Sprintf("mqrestadmin transport error for %s: %v", e.URL, e.Err)
}

func (e *TransportError) Unwrap() error {
	return e.Err
}

// ResponseError indicates the MQ REST API returned a response that could not
// be parsed as valid JSON or had an unexpected structure.
type ResponseError struct {
	ResponseText string
	StatusCode   int
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("mqrestadmin response error (HTTP %d): %s", e.StatusCode, e.ResponseText)
}

// AuthError indicates an authentication or authorization failure (HTTP 401 or
// 403), or an LTPA login failure.
type AuthError struct {
	URL        string
	StatusCode int
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("mqrestadmin auth error (HTTP %d) for %s", e.StatusCode, e.URL)
}

// CommandError indicates the MQSC command returned a non-zero completion code
// or reason code.
type CommandError struct {
	Payload    map[string]any
	StatusCode int
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("mqrestadmin command error (HTTP %d): %v", e.StatusCode, e.Payload)
}

// TimeoutError indicates a synchronous polling operation exceeded its
// configured timeout.
type TimeoutError struct {
	Name           string
	Operation      SyncOperation
	ElapsedSeconds float64
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("mqrestadmin timeout: %s %s after %.1fs", e.Operation, e.Name, e.ElapsedSeconds)
}

// MappingError indicates one or more attribute translation failures in strict
// mode.
type MappingError struct {
	Issues []MappingIssue
}

func (e *MappingError) Error() string {
	descriptions := make([]string, len(e.Issues))
	for idx, issue := range e.Issues {
		descriptions[idx] = issue.String()
	}
	return fmt.Sprintf("mqrestadmin mapping error: %s", strings.Join(descriptions, "; "))
}

// MappingIssue describes a single attribute translation failure.
type MappingIssue struct {
	Direction     MappingDirection
	Reason        MappingReason
	AttributeName string
	// AttributeValue is the value that failed translation, if applicable.
	AttributeValue any
	// ObjectIndex is the index within a response list, if applicable.
	ObjectIndex *int
	// Qualifier is the mapping qualifier, if applicable.
	Qualifier string
}

func (issue MappingIssue) String() string {
	return fmt.Sprintf("%s %s: %s", issue.Direction, issue.Reason, issue.AttributeName)
}

// MappingDirection indicates whether a mapping failure occurred during request
// or response translation.
type MappingDirection int

const (
	// MappingRequest indicates the failure occurred translating request attributes.
	MappingRequest MappingDirection = iota
	// MappingResponse indicates the failure occurred translating response attributes.
	MappingResponse
)

func (direction MappingDirection) String() string {
	switch direction {
	case MappingRequest:
		return "request"
	case MappingResponse:
		return "response"
	default:
		return "unknown"
	}
}

// MappingReason indicates the cause of a mapping failure.
type MappingReason int

const (
	// MappingUnknownKey indicates an unrecognized attribute name.
	MappingUnknownKey MappingReason = iota
	// MappingUnknownValue indicates an unrecognized attribute value.
	MappingUnknownValue
	// MappingUnknownQualifier indicates an unrecognized qualifier.
	MappingUnknownQualifier
)

func (reason MappingReason) String() string {
	switch reason {
	case MappingUnknownKey:
		return "unknown_key"
	case MappingUnknownValue:
		return "unknown_value"
	case MappingUnknownQualifier:
		return "unknown_qualifier"
	default:
		return "unknown"
	}
}
