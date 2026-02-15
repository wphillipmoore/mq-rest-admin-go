# The runCommandJSON Endpoint

--8<-- "design/runcommand-endpoint.md"

## Nested object flattening

Some commands (`DISPLAY CONN TYPE(HANDLE)`, `DISPLAY QSTATUS
TYPE(HANDLE)`) return responses where each `commandResponse` item
contains an `objects` array of per-handle attributes alongside
parent-scoped attributes. `mqrestadmin` automatically detects and flattens
these structures so that every command returns uniform flat maps:

```json
{"conn": "A1B2C3D4E5F6", "objname": "MY.QUEUE", "hstate": "ACTIVE"}
```

See [nested object flattening](nested-object-flattening.md) for the full algorithm, edge cases,
and before/after examples.

## Go implementation notes

In `mqrestadmin`, the error handling described above is implemented using
typed error structs:

- Non-zero overall codes produce a `*CommandError` with the full response
  payload attached.
- DISPLAY commands with no matches (reason code 2085) return an empty
  slice and nil error.
- The CSRF token defaults to `"local"` and can be overridden via
  `WithCSRFToken()`, or omitted with `WithCSRFToken("")`.
- Authentication is configured via functional options: `WithBasicAuth()`,
  `WithCertificateAuth()`, or `WithLTPAToken()`.
