# Mapping Pipeline

## The three-namespace problem

IBM MQ uses multiple naming conventions depending on the interface:

**MQSC names** (e.g. `CURDEPTH`, `DEFPSIST`)
: Short, uppercase tokens used in MQSC commands and the REST API's
  `runCommandJSON` endpoint.

**PCF names** (e.g. `CurrentQDepth`, `DefPersistence`)
: CamelCase names from the Programmable Command Formats. Not used
  directly by `mqrestadmin`, but they form the intermediate namespace in
  the mapping pipeline.

**Developer names** (e.g. `current_queue_depth`, `default_persistence`)
: Human-readable `snake_case` names for use in application code.

The mapping pipeline translates between MQSC and developer names. PCF names
were used as an intermediate reference during the original extraction
process that bootstrapped the mapping tables but do not appear at
runtime.

In Go, developer-friendly names use `snake_case` (e.g. `current_queue_depth`,
`default_persistence`), matching the convention across the mq-rest-admin library
family. The mapping tables are loaded from an embedded JSON resource via Go's
`embed` package and are identical across language implementations.

```go
ctx := context.Background()

// With mapping enabled (default) -- developer-friendly names
queues, err := session.DisplayQueue(ctx, "MY.QUEUE",
    mqrestadmin.WithResponseParameters([]string{"current_queue_depth", "max_queue_depth"}))
// Returns: [{"queue_name": "MY.QUEUE", "current_queue_depth": 0, "max_queue_depth": 5000}]

// With mapping disabled -- native MQSC names
// (create session with WithMapAttributes(false))
queues, err = session.DisplayQueue(ctx, "MY.QUEUE",
    mqrestadmin.WithResponseParameters([]string{"CURDEPTH", "MAXDEPTH"}))
// Returns: [{"queue": "MY.QUEUE", "curdepth": 0, "maxdepth": 5000}]
```

## Qualifier-based mapping

Mappings are organized by **qualifier** (e.g. `queue`, `channel`, `qmgr`),
not by command. A single qualifier's mapping tables serve all commands
that operate on that object type. For example, the `queue` qualifier
covers `DISPLAY QUEUE`, `DEFINE QLOCAL`, `DELETE QALIAS`, and all other
queue-related commands.

This design avoids duplicating mapping data across commands and reflects
how MQSC attributes are shared across command verbs.

See the [Qualifier Mapping Reference](mappings/index.md) for the complete
per-qualifier documentation, including every key map, value map, and key-value
map entry.

## Request mapping flow

When mapping is enabled, request attributes are translated before sending
to the MQ REST API:

1. **Key mapping**: Each `snake_case` attribute name is looked up in the
   qualifier's `request_key_map`. If found, the key is replaced with the
   MQSC parameter name.

2. **Value mapping**: For attributes with enumerated values, the
   qualifier's `request_value_map` translates developer values to MQSC
   values (e.g. `"yes"` to `"YES"`).

3. **Key-value mapping**: Some attributes require both key and value to
   change simultaneously. The `request_key_value_map` handles cases
   where a single developer attribute expands to a different MQSC key+value
   pair (e.g. `channel_type="server_connection"` to
   `CHLTYPE("SVRCONN")`).

### Go example

```go
ctx := context.Background()

// Developer-friendly request parameters
_, err := session.DefineQlocal(ctx, "MY.QUEUE",
    mqrestadmin.WithRequestParameters(map[string]any{
        "max_queue_depth":      50000,
        "default_persistence":  "yes",
        "description":          "Application queue",
    }))

// After request mapping, the JSON payload sent to MQ contains:
// { "MAXDEPTH": 50000, "DEFPSIST": "YES", "DESCR": "Application queue" }
```

## Response mapping flow

Response attributes are translated after receiving the MQ REST response:

1. **Key mapping**: Each MQSC parameter name from the response is looked
   up in the qualifier's `response_key_map`. If found, the key is
   replaced with the `snake_case` name.

2. **Value mapping**: Enumerated MQSC values are translated to
   developer-friendly values via the `response_value_map` (e.g. `"YES"` to
   `"yes"`).

### Go example

```go
ctx := context.Background()

// Request specific response attributes using developer-friendly names
queues, err := session.DisplayQueue(ctx, "MY.QUEUE",
    mqrestadmin.WithResponseParameters([]string{"current_queue_depth", "max_queue_depth"}))

// MQ returns MQSC names in the response:
// { "queue": "MY.QUEUE", "curdepth": 0, "maxdepth": 5000 }

// After response mapping, the caller receives developer-friendly names:
// { "queue_name": "MY.QUEUE", "current_queue_depth": 0, "max_queue_depth": 5000 }
```

## Response parameter mapping

When the caller specifies response parameters (the list of attributes
to return), those names are also mapped from `snake_case` to MQSC before
being sent in the request. This allows callers to request specific
attributes using developer-friendly names.

Response parameter macros (like `CFCONLOS` for channel status) are
recognized and passed through without mapping.

## WHERE keyword mapping

The `where` parameter on DISPLAY methods accepts a filter expression
like `"current_queue_depth GT 100"`. The first token (the keyword) is mapped
from `snake_case` to the MQSC name. The rest of the expression is
passed through unchanged.

## Qualifier resolution

When a command is executed, the mapping qualifier is resolved by:

1. Looking up the command key (e.g. `"DISPLAY QUEUE"`) in
   the commands section of the mapping data for an explicit qualifier.
2. Falling back to a hardcoded default map (e.g. `QLOCAL` to `queue`,
   `CHANNEL` to `channel`).
3. As a last resort, lowercasing the MQSC qualifier.

This means `DEFINE QLOCAL`, `DEFINE QREMOTE`, and `DISPLAY QUEUE` all
resolve to the `queue` qualifier and share the same mapping tables.

## Strict vs lenient mode

**Strict mode** (default): Any attribute name or value that cannot be
mapped causes a `*MappingError` to be returned. This catches typos and
unsupported attributes early.

**Lenient mode** (`WithMappingStrict(false)`): Unknown attribute names and
values pass through unchanged. This is useful when working with
attributes not yet covered by the mapping tables.

The mode is set at session construction and applies to all mapping
operations. It cannot be overridden per-call.

In Go, the mode is set via a session option:

```go
// Strict mode (default) -- returns *MappingError on unknown attributes
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
    mqrestadmin.WithMappingStrict(true), // default
)

// Lenient mode -- unknown attributes pass through unchanged
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
    mqrestadmin.WithMappingStrict(false),
)
```

## Custom mapping overrides

The built-in mapping tables cover all standard MQSC attributes, but sites may
use different `snake_case` conventions. The `WithMappingOverrides` option on
`NewSession` lets you layer sparse changes on top of the built-in data
without replacing it.

### How merging works

Overrides are merged at the **key level** within each sub-map. You only specify
the entries you want to change -- all other mappings remain intact:

```go
overrides := map[string]any{
    "qualifiers": map[string]any{
        "queue": map[string]any{
            "response_key_map": map[string]any{
                "CURDEPTH": "queue_depth",     // override built-in mapping
                "MAXDEPTH": "queue_max_depth",  // override built-in mapping
            },
        },
    },
}

session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
    mqrestadmin.WithMappingOverrides(overrides, mqrestadmin.MappingOverrideMerge),
)

ctx := context.Background()
queues, err := session.DisplayQueue(ctx, "MY.QUEUE")
// Returns: [{"queue_depth": 0, "queue_max_depth": 5000, ...}]
// All other queue attributes keep their default names
```

When this override is applied:

1. The built-in mapping data is deep-copied (the original is never mutated).
2. The `queue` qualifier's `response_key_map` is updated: the entry for
   `CURDEPTH` changes from `"current_queue_depth"` to `"queue_depth"`.
3. All other entries in `response_key_map` (and all other sub-maps) remain
   unchanged.

### Override modes

The `WithMappingOverrides` option accepts a `MappingOverrideMode` that
controls how overrides are applied:

- **`MappingOverrideMerge`**: Overlays overrides onto the default mapping
  data, adding or replacing individual entries within each sub-map. This is
  the most common mode.
- **`MappingOverrideReplace`**: Replaces entire qualifier sections with the
  override data. Use this when you want to completely redefine a qualifier's
  mappings.

### Supported override keys

The top level of the overrides map accepts two keys:

- **`commands`**: Override command-level metadata (e.g. which qualifier a
  command resolves to). Each command entry is shallow-merged.
- **`qualifiers`**: Override qualifier mapping tables. Each qualifier supports
  five sub-maps:
  - `request_key_map` -- `snake_case` to MQSC key mapping for requests
  - `request_value_map` -- value translations for request attributes
  - `request_key_value_map` -- combined key+value translations for requests
  - `response_key_map` -- MQSC to `snake_case` key mapping for responses
  - `response_value_map` -- value translations for response attributes

### Adding new qualifiers

You can add mappings for qualifiers not yet covered by the built-in data:

```go
overrides := map[string]any{
    "qualifiers": map[string]any{
        "custom_object": map[string]any{
            "request_key_map":   map[string]any{"my_attr": "MYATTR"},
            "response_key_map":  map[string]any{"MYATTR": "my_attr"},
            "request_value_map":  map[string]any{},
            "response_value_map": map[string]any{},
        },
    },
}
```

### Validation

The override structure is validated at session construction time. Invalid
shapes cause `NewSession` to return an error immediately, so problems are
caught before any commands are sent.

## Per-call opt-out

Mapping can be disabled entirely at the session level by passing
`WithMapAttributes(false)` when creating the session:

```go
// All commands use native MQSC names
session, err := mqrestadmin.NewSession(
    "https://localhost:9443/ibmmq/rest/v2",
    "QM1",
    mqrestadmin.BasicAuth{Username: "admin", Password: "passw0rd"},
    mqrestadmin.WithMapAttributes(false),
)

ctx := context.Background()
queues, err := session.DisplayQueue(ctx, "MY.QUEUE",
    mqrestadmin.WithResponseParameters([]string{"CURDEPTH", "MAXDEPTH"}))
// Returns native MQSC names: [{"queue": "MY.QUEUE", "curdepth": 0, ...}]
```

When mapping is disabled, attributes pass through in their native MQSC
form.
