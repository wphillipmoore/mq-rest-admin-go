package examples

import (
	"context"
	"fmt"
	"strings"

	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

// QueueHandleInfo holds per-handle queue status information.
type QueueHandleInfo struct {
	QueueName   string
	HandleState string
	ConnectionID string
	OpenOptions string
}

// ConnectionHandleInfo holds per-handle connection information.
type ConnectionHandleInfo struct {
	ConnectionID string
	ObjectName   string
	HandleState  string
	ObjectType   string
}

// ReportQueueHandles returns per-handle queue status entries.
func ReportQueueHandles(ctx context.Context, session *mqrestadmin.Session) []QueueHandleInfo {
	entries, err := session.DisplayQstatus(ctx, "*",
		mqrestadmin.WithRequestParameters(map[string]any{"type": "HANDLE"}))
	if err != nil {
		return []QueueHandleInfo{}
	}

	results := make([]QueueHandleInfo, 0, len(entries))
	for _, entry := range entries {
		results = append(results, QueueHandleInfo{
			QueueName:    strings.TrimSpace(fmt.Sprint(entry["queue_name"])),
			HandleState:  strings.TrimSpace(fmt.Sprint(entry["handle_state"])),
			ConnectionID: strings.TrimSpace(fmt.Sprint(entry["connection_id"])),
			OpenOptions:  strings.TrimSpace(fmt.Sprint(entry["open_options"])),
		})
	}
	return results
}

// ReportConnectionHandles returns per-handle connection entries.
func ReportConnectionHandles(ctx context.Context, session *mqrestadmin.Session) []ConnectionHandleInfo {
	entries, err := session.DisplayConn(ctx, "*",
		mqrestadmin.WithRequestParameters(map[string]any{"connection_info_type": "HANDLE"}))
	if err != nil {
		return []ConnectionHandleInfo{}
	}

	results := make([]ConnectionHandleInfo, 0, len(entries))
	for _, entry := range entries {
		results = append(results, ConnectionHandleInfo{
			ConnectionID: strings.TrimSpace(fmt.Sprint(entry["connection_id"])),
			ObjectName:   strings.TrimSpace(fmt.Sprint(entry["object_name"])),
			HandleState:  strings.TrimSpace(fmt.Sprint(entry["handle_state"])),
			ObjectType:   strings.TrimSpace(fmt.Sprint(entry["object_type"])),
		})
	}
	return results
}

// PrintQueueStatus runs the queue status report and prints formatted output.
func PrintQueueStatus(ctx context.Context, session *mqrestadmin.Session) {
	queueHandles := ReportQueueHandles(ctx, session)

	fmt.Printf("\n%-30s %-15s %-30s %s\n",
		"Queue", "Handle State", "Connection ID", "Open Options")
	fmt.Println(strings.Repeat("-", 90))
	for _, info := range queueHandles {
		fmt.Printf("%-30s %-15s %-30s %s\n",
			info.QueueName, info.HandleState, info.ConnectionID, info.OpenOptions)
	}
	if len(queueHandles) == 0 {
		fmt.Println("  (no active queue handles)")
	}

	connHandles := ReportConnectionHandles(ctx, session)

	fmt.Printf("\n%-30s %-30s %-15s %s\n",
		"Connection ID", "Object Name", "Handle State", "Object Type")
	fmt.Println(strings.Repeat("-", 90))
	for _, info := range connHandles {
		fmt.Printf("%-30s %-30s %-15s %s\n",
			info.ConnectionID, info.ObjectName, info.HandleState, info.ObjectType)
	}
	if len(connHandles) == 0 {
		fmt.Println("  (no active connection handles)")
	}
}
