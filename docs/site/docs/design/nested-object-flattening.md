# Nested Object Flattening

--8<-- "design/nested-object-flattening.md"

When attribute mapping is enabled (the default), the MQSC names are then
translated to `snake_case`:

```go
// After flattening and mapping, each result is a flat map:
// map[string]any{
//     "object_name":  "MY.QUEUE",
//     "handle_state": "ACTIVE",
//     ...
// }
```

## Go implementation

In `mqrestadmin`, the flattening logic lives in the internal
`flattenNestedObjects()` function. It processes `[]map[string]any` and
returns a new `[]map[string]any` with all nesting resolved.

The merge uses `maps.Clone()` from the standard library to copy the shared
parent keys, then sets each nested-item key on the clone. Nested-item keys
override any same-named parent keys.

## See also

- [runCommandJSON endpoint](runcommand-endpoint.md) -- general `runCommandJSON` request/response
  structure
- [rationale](rationale.md) -- overall design rationale
- [DISPLAY CONN](https://www.ibm.com/docs/en/ibm-mq/9.4?topic=reference-display-conn) --
  IBM MQ documentation
- [DISPLAY QSTATUS](https://www.ibm.com/docs/en/ibm-mq/9.4?topic=reference-display-qstatus) --
  IBM MQ documentation
