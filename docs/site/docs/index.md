# mqrestadmin

Go wrapper for the IBM MQ administrative REST API.

`mqrestadmin` provides typed Go methods for every MQSC command exposed
by the IBM MQ 9.4 `runCommandJSON` REST endpoint. Attribute names are
automatically translated between Go `snake_case` and native MQSC
parameter names, so you work with idiomatic Go identifiers throughout.

Zero external runtime dependencies — stdlib only.
