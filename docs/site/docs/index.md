# mqrestadmin

## Overview

**mqrestadmin** provides a Go-friendly interface to IBM MQ queue manager
administration via the `runCommandJSON` REST endpoint. It translates between
Go `snake_case` attribute names and native MQSC parameter names, wraps
every MQSC command as a typed method, and handles authentication, CSRF tokens,
and error propagation.

## Key features

- **~144 command methods** covering all MQSC verbs and qualifiers
- **Bidirectional attribute mapping** between developer-friendly names and MQSC parameters
- **Idempotent ensure methods** for declarative object management
- **Bulk sync operations** for configuration-as-code workflows
- **Zero runtime dependencies** — stdlib only
- **Transport abstraction** for easy testing with mock transports

## Installation

```bash
go get github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin
```

## Status

This project is in **pre-alpha** (initial setup). The API surface, mapping
tables, and return shapes are under active development.

## License

GNU General Public License v3.0
