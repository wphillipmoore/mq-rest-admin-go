// Package mqrestadmin provides a Go client for the IBM MQ administrative REST API.
//
// It wraps the IBM MQ 9.4 runCommandJSON endpoint, providing typed Go methods
// for every MQSC command (DISPLAY, DEFINE, ALTER, DELETE, START, STOP, etc.)
// with automatic attribute name translation between Go snake_case and native
// MQSC parameter names.
//
// This package has zero external runtime dependencies, using only the Go
// standard library for HTTP, JSON, TLS, and embedded resources.
//
// # Session Creation
//
// Create a session with functional options:
//
//	session, err := mqrestadmin.NewSession(
//	    "https://host:9443/ibmmq/rest/v2",
//	    "QM1",
//	    mqrestadmin.WithBasicAuth("user", "pass"),
//	    mqrestadmin.WithTimeout(30 * time.Second),
//	)
//
// # Command Methods
//
// Each MQSC command has a corresponding method:
//
//	queues, err := session.DisplayQlocal(ctx, "*")
//	err = session.DefineQlocal(ctx, "APP.REQUESTS", params)
//	err = session.AlterQlocal(ctx, "APP.REQUESTS", params)
//	err = session.DeleteQlocal(ctx, "APP.REQUESTS")
//
// # Idempotent Operations
//
// Ensure methods provide idempotent create-or-update semantics:
//
//	result, err := session.EnsureQlocal(ctx, "APP.REQUESTS", params)
//	// result.Action is EnsureCreated, EnsureUpdated, or EnsureUnchanged
//
// # Synchronous Polling
//
// Sync methods poll until an operation completes:
//
//	result, err := session.StartChannelSync(ctx, "TO.REMOTE", syncConfig)
//	// result.Operation is SyncStarted
package mqrestadmin
