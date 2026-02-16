# Mapping

## Overview

The mapping module provides bidirectional attribute translation between
developer-friendly `snake_case` names and native MQSC parameter names. The
mapper is created internally by `Session` and is not typically used directly.

See [Mapping Pipeline](../mapping-pipeline.md) for a conceptual overview of
how mapping works.

## Attribute mapping

The internal `attributeMapper` translates attribute names and values between the
developer-friendly namespace and the MQSC namespace. The mapper performs three
types of translation in each direction:

- **Key mapping**: Attribute name translation (e.g. `current_queue_depth` to
  `CURDEPTH`)
- **Value mapping**: Enumerated value translation (e.g. `"yes"` to `"YES"`,
  `"server_connection"` to `"SVRCONN"`)
- **Key-value mapping**: Combined name+value translation for cases where both
  key and value change together (e.g. `channel_type="server_connection"` to
  `CHLTYPE("SVRCONN")`)

The mapper is qualifier-aware: it selects the correct mapping tables based on
the MQSC command's qualifier (e.g. `queue`, `channel`, `qmgr`).

## Controlling mapping

### WithMapAttributes

Controls whether attribute names are translated. Defaults to `true`.

```go
// Disable mapping -- use raw MQSC parameter names
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2", "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
    mqrestadmin.WithMapAttributes(false),
)

// With mapping disabled, use MQSC names directly
queues, err := session.DisplayQueue(ctx, "*",
    mqrestadmin.WithRequestParameters(map[string]any{"CURDEPTH": 0}),
)
```

### WithMappingStrict

Controls whether unknown attributes cause an error or pass through silently.
Defaults to `true` (strict mode).

```go
// Permissive mode -- unknown attributes pass through without error
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2", "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
    mqrestadmin.WithMappingStrict(false),
)
```

**Strict mode** (default): Any attribute name or value that cannot be mapped
causes a `*MappingError` to be returned. This catches typos and ensures all
attributes are correctly translated.

**Permissive mode**: Unknown attributes pass through unchanged. A
`*MappingError` is not returned, but `MappingIssue` values are still tracked
internally. This is useful when working with custom or version-specific
attributes not covered by the built-in mapping data.

## MappingOverrideMode

Controls how custom overrides are merged with built-in mapping data:

```go
const (
    MappingOverrideMerge   MappingOverrideMode = iota  // overlay at key level
    MappingOverrideReplace                             // replace entire qualifier section
)
```

- **MappingOverrideMerge** (default): Override entries are merged at the key
  level within each sub-map. Existing entries not mentioned in the override are
  preserved. This is the common case for changing a few attribute names without
  losing the rest.
- **MappingOverrideReplace**: The override completely replaces the specified
  qualifier section. Use when you need full control over a qualifier's mapping.

```go
overrides := map[string]any{
    "qualifiers": map[string]any{
        "queue": map[string]any{
            "request_key_map": map[string]string{
                "my_custom_attr": "MYCUSTOM",
            },
        },
    },
}

session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2", "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
    mqrestadmin.WithMappingOverrides(overrides, mqrestadmin.MappingOverrideMerge),
)
```

## Mapping data

The mapping definitions are embedded in the compiled binary from
`mapping-data.json` using Go's `//go:embed` directive. The data is organized by
qualifier (e.g. `queue`, `channel`, `qmgr`) with separate maps for request and
response directions. Each qualifier contains:

- `request_key_map` -- developer-friendly to MQSC key mapping for requests
- `request_value_map` -- value translations for request attributes
- `request_key_value_map` -- combined key+value translations for requests
- `response_key_map` -- MQSC to developer-friendly key mapping for responses
- `response_value_map` -- value translations for response attributes

The mapping data was originally bootstrapped from IBM MQ 9.4 documentation and
covers all standard MQSC attributes across 42 qualifiers.

## MappingIssue

Tracks mapping problems encountered during translation:

```go
type MappingIssue struct {
    Direction      MappingDirection  // MappingRequest or MappingResponse
    Reason         MappingReason     // MappingUnknownKey, MappingUnknownValue, or MappingUnknownQualifier
    AttributeName  string            // The attribute that failed translation
    AttributeValue any               // The value that failed translation (if applicable)
    ObjectIndex    *int              // Index within a response list (if applicable)
    Qualifier      string            // The mapping qualifier (if applicable)
}
```

### MappingDirection

```go
const (
    MappingRequest  MappingDirection = iota  // Failure during request translation
    MappingResponse                          // Failure during response translation
)
```

### MappingReason

```go
const (
    MappingUnknownKey       MappingReason = iota  // Unrecognized attribute name
    MappingUnknownValue                           // Unrecognized attribute value
    MappingUnknownQualifier                       // Unrecognized qualifier
)
```

## MappingError

Returned when attribute mapping fails in strict mode. Separate from the main
error types (it does not wrap any of the transport/command error types). Contains
the list of `MappingIssue` instances that caused the failure.

```go
type MappingError struct {
    Issues []MappingIssue
}
```

```go
queues, err := session.DisplayQueue(ctx, "*",
    mqrestadmin.WithRequestParameters(map[string]any{
        "invalid_attribute_name": "value",
    }),
)
if err != nil {
    var mappingErr *mqrestadmin.MappingError
    if errors.As(err, &mappingErr) {
        for _, issue := range mappingErr.Issues {
            fmt.Printf("Mapping issue: %s %s: %s\n",
                issue.Direction, issue.Reason, issue.AttributeName)
        }
    }
}
```

See [Errors](errors.md) for the complete error type reference.
