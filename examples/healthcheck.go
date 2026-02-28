package examples

import (
	"context"
	"fmt"
	"strings"

	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

// ListenerResult holds the health status for a single listener.
type ListenerResult struct {
	Name   string
	Status string
}

// QMHealthResult holds the health check result for a single queue manager.
type QMHealthResult struct {
	QmgrName      string
	Reachable     bool
	Status        string
	CommandServer string
	Listeners     []ListenerResult
	Passed        bool
}

// CheckHealth runs a health check against a single queue manager.
func CheckHealth(ctx context.Context, session *mqrestadmin.Session) QMHealthResult {
	result := QMHealthResult{
		QmgrName:      session.QmgrName(),
		Status:        "UNKNOWN",
		CommandServer: "UNKNOWN",
	}

	qmgr, err := session.DisplayQmgr(ctx)
	if err != nil {
		return result
	}

	result.Reachable = true

	if name, ok := qmgr["queue_manager_name"]; ok {
		if s := strings.TrimSpace(fmt.Sprint(name)); s != "" {
			result.QmgrName = s
		}
	}

	qmstatus, err := session.DisplayQmstatus(ctx)
	if err == nil && qmstatus != nil {
		result.Status = strings.TrimSpace(fmt.Sprint(qmstatus["ha_status"]))
	}

	cmdserv, err := session.DisplayCmdserv(ctx)
	if err == nil && cmdserv != nil {
		result.CommandServer = strings.TrimSpace(fmt.Sprint(cmdserv["status"]))
	}

	listeners, err := session.DisplayListener(ctx, "*")
	if err == nil {
		for _, listener := range listeners {
			result.Listeners = append(result.Listeners, ListenerResult{
				Name:   strings.TrimSpace(fmt.Sprint(listener["listener_name"])),
				Status: strings.TrimSpace(fmt.Sprint(listener["start_mode"])),
			})
		}
	}

	result.Passed = result.Reachable && result.Status != "UNKNOWN"
	return result
}

// PrintHealthCheck runs health checks and prints formatted output.
func PrintHealthCheck(ctx context.Context, sessions []*mqrestadmin.Session) []QMHealthResult {
	var results []QMHealthResult

	for _, session := range sessions {
		r := CheckHealth(ctx, session)
		results = append(results, r)

		verdict := "FAIL"
		if r.Passed {
			verdict = "PASS"
		}
		fmt.Printf("\n=== %s: %s ===\n", r.QmgrName, verdict)
		fmt.Printf("  Reachable:      %t\n", r.Reachable)
		fmt.Printf("  Status:         %s\n", r.Status)
		fmt.Printf("  Command server: %s\n", r.CommandServer)
		fmt.Printf("  Listeners:      %d\n", len(r.Listeners))
		for _, l := range r.Listeners {
			fmt.Printf("    %s: %s\n", l.Name, l.Status)
		}
	}

	return results
}
