package examples

import (
	"context"
	"fmt"

	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

const prefix = "PROV"

// ProvisionResult holds the result of the provisioning operation.
type ProvisionResult struct {
	ObjectsCreated []string
	ObjectsFailed  []string
	Verified       bool
}

// Provision creates cross-QM objects on both queue managers.
func Provision(ctx context.Context, qm1, qm2 *mqrestadmin.Session) ProvisionResult {
	result := ProvisionResult{}

	defineObject(ctx, &result, qm1, "DefineQlocal", prefix+".QM1.LOCAL", map[string]any{
		"replace": "yes", "default_persistence": "yes",
		"description": "provisioned local queue on QM1",
	})
	defineObject(ctx, &result, qm2, "DefineQlocal", prefix+".QM2.LOCAL", map[string]any{
		"replace": "yes", "default_persistence": "yes",
		"description": "provisioned local queue on QM2",
	})

	defineObject(ctx, &result, qm1, "DefineQlocal", prefix+".QM1.TO.QM2.XMITQ", map[string]any{
		"replace": "yes", "usage": "XMITQ",
		"description": "xmit queue QM1 to QM2",
	})
	defineObject(ctx, &result, qm2, "DefineQlocal", prefix+".QM2.TO.QM1.XMITQ", map[string]any{
		"replace": "yes", "usage": "XMITQ",
		"description": "xmit queue QM2 to QM1",
	})

	defineObject(ctx, &result, qm1, "DefineQremote", prefix+".REMOTE.TO.QM2", map[string]any{
		"replace": "yes", "remote_queue_name": prefix + ".QM2.LOCAL",
		"remote_queue_manager_name": "QM2",
		"transmission_queue_name":   prefix + ".QM1.TO.QM2.XMITQ",
		"description":               "remote queue QM1 to QM2",
	})
	defineObject(ctx, &result, qm2, "DefineQremote", prefix+".REMOTE.TO.QM1", map[string]any{
		"replace": "yes", "remote_queue_name": prefix + ".QM1.LOCAL",
		"remote_queue_manager_name": "QM1",
		"transmission_queue_name":   prefix + ".QM2.TO.QM1.XMITQ",
		"description":               "remote queue QM2 to QM1",
	})

	defineObject(ctx, &result, qm1, "DefineChannel", prefix+".QM1.TO.QM2", map[string]any{
		"replace": "yes", "channel_type": "SDR", "transport_type": "TCP",
		"connection_name":           "qm2(1414)",
		"transmission_queue_name":   prefix + ".QM1.TO.QM2.XMITQ",
		"description":               "sender QM1 to QM2",
	})
	defineObject(ctx, &result, qm2, "DefineChannel", prefix+".QM1.TO.QM2", map[string]any{
		"replace": "yes", "channel_type": "RCVR", "transport_type": "TCP",
		"description": "receiver QM1 to QM2",
	})
	defineObject(ctx, &result, qm2, "DefineChannel", prefix+".QM2.TO.QM1", map[string]any{
		"replace": "yes", "channel_type": "SDR", "transport_type": "TCP",
		"connection_name":           "qm1(1414)",
		"transmission_queue_name":   prefix + ".QM2.TO.QM1.XMITQ",
		"description":               "sender QM2 to QM1",
	})
	defineObject(ctx, &result, qm1, "DefineChannel", prefix+".QM2.TO.QM1", map[string]any{
		"replace": "yes", "channel_type": "RCVR", "transport_type": "TCP",
		"description": "receiver QM2 to QM1",
	})

	qm1Queues, err1 := qm1.DisplayQueue(ctx, prefix+".*")
	qm2Queues, err2 := qm2.DisplayQueue(ctx, prefix+".*")
	result.Verified = err1 == nil && err2 == nil && len(qm1Queues) >= 3 && len(qm2Queues) >= 3

	return result
}

// Teardown removes all provisioned objects from both queue managers.
func Teardown(ctx context.Context, qm1, qm2 *mqrestadmin.Session) []string {
	var failures []string

	// QM1 objects
	deleteObject(ctx, &failures, qm1, qm1.DeleteChannel, prefix+".QM1.TO.QM2")
	deleteObject(ctx, &failures, qm1, qm1.DeleteChannel, prefix+".QM2.TO.QM1")
	deleteObject(ctx, &failures, qm1, qm1.DeleteQremote, prefix+".REMOTE.TO.QM2")
	deleteObject(ctx, &failures, qm1, qm1.DeleteQlocal, prefix+".QM1.TO.QM2.XMITQ")
	deleteObject(ctx, &failures, qm1, qm1.DeleteQlocal, prefix+".QM1.LOCAL")

	// QM2 objects
	deleteObject(ctx, &failures, qm2, qm2.DeleteChannel, prefix+".QM1.TO.QM2")
	deleteObject(ctx, &failures, qm2, qm2.DeleteChannel, prefix+".QM2.TO.QM1")
	deleteObject(ctx, &failures, qm2, qm2.DeleteQremote, prefix+".REMOTE.TO.QM1")
	deleteObject(ctx, &failures, qm2, qm2.DeleteQlocal, prefix+".QM2.TO.QM1.XMITQ")
	deleteObject(ctx, &failures, qm2, qm2.DeleteQlocal, prefix+".QM2.LOCAL")

	return failures
}

// PrintProvision provisions, reports, and tears down the environment.
func PrintProvision(ctx context.Context, qm1, qm2 *mqrestadmin.Session) ProvisionResult {
	fmt.Println("\n=== Provisioning environment ===")
	result := Provision(ctx, qm1, qm2)

	fmt.Printf("\nCreated: %d\n", len(result.ObjectsCreated))
	for _, obj := range result.ObjectsCreated {
		fmt.Printf("  + %s\n", obj)
	}
	if len(result.ObjectsFailed) > 0 {
		fmt.Printf("\nFailed: %d\n", len(result.ObjectsFailed))
		for _, obj := range result.ObjectsFailed {
			fmt.Printf("  ! %s\n", obj)
		}
	}
	fmt.Printf("\nVerified: %t\n", result.Verified)

	fmt.Println("\n=== Tearing down ===")
	failures := Teardown(ctx, qm1, qm2)
	if len(failures) > 0 {
		fmt.Printf("Teardown failures: %v\n", failures)
	} else {
		fmt.Println("Teardown complete.")
	}

	return result
}

type defineFunc func(ctx context.Context, name string, opts ...mqrestadmin.CommandOption) error
type deleteFunc func(ctx context.Context, name string, opts ...mqrestadmin.CommandOption) error

func defineObject(ctx context.Context, result *ProvisionResult, session *mqrestadmin.Session, method, name string, params map[string]any) {
	label := session.QmgrName() + "/" + name

	var fn defineFunc
	switch method {
	case "DefineQlocal":
		fn = session.DefineQlocal
	case "DefineQremote":
		fn = session.DefineQremote
	case "DefineChannel":
		fn = session.DefineChannel
	default:
		result.ObjectsFailed = append(result.ObjectsFailed, label)
		return
	}

	if err := fn(ctx, name, mqrestadmin.WithRequestParameters(params)); err != nil {
		result.ObjectsFailed = append(result.ObjectsFailed, label)
	} else {
		result.ObjectsCreated = append(result.ObjectsCreated, label)
	}
}

func deleteObject(ctx context.Context, failures *[]string, session *mqrestadmin.Session, fn deleteFunc, name string) {
	if err := fn(ctx, name); err != nil {
		*failures = append(*failures, session.QmgrName()+"/"+name)
	}
}
