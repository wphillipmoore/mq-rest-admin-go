# API Reference

## Core

- [Session](session.md) -- `Session` struct and `NewSession()` constructor
- [Commands](commands.md) -- MQSC command methods
- [Transport](transport.md) -- `Transport` interface and `HTTPTransport`

## Authentication

- [Auth](auth.md) -- `Credentials` sealed interface and implementations

## Mapping

- [Mapping](mapping.md) -- Attribute mapping pipeline and override modes

## Errors

- [Errors](errors.md) -- Error types for `errors.As()` matching

## Patterns

- [Ensure](ensure.md) -- `EnsureResult`, `EnsureAction`
- [Ensure Methods](../ensure-methods.md) -- Per-object-type ensure convenience methods
- [Sync](sync.md) -- `SyncConfig`, `SyncResult`, `SyncOperation`
- [Sync Methods](../sync-methods.md) -- Per-object-type sync convenience methods
