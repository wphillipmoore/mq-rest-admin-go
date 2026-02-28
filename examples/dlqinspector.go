package examples

import (
	"context"
	"fmt"
	"strings"

	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

const criticalDepthPct = 90.0

// DLQReport holds the dead letter queue inspection result.
type DLQReport struct {
	QmgrName     string
	DLQName      string
	Configured   bool
	CurrentDepth int
	MaxDepth     int
	DepthPct     float64
	OpenInput    int
	OpenOutput   int
	Suggestion   string
}

// InspectDLQ inspects the dead letter queue for a queue manager.
func InspectDLQ(ctx context.Context, session *mqrestadmin.Session) (DLQReport, error) {
	qmgr, err := session.DisplayQmgr(ctx)
	if err != nil {
		return DLQReport{QmgrName: session.QmgrName()}, err
	}

	dlqName := strings.TrimSpace(fmt.Sprint(qmgr["dead_letter_queue_name"]))
	if dlqName == "" || dlqName == "<nil>" {
		return DLQReport{
			QmgrName:   session.QmgrName(),
			Suggestion: "No dead letter queue configured. Define one with ALTER QMGR DEADQ.",
		}, nil
	}

	queues, err := session.DisplayQueue(ctx, dlqName)
	if err != nil || len(queues) == 0 {
		return DLQReport{
			QmgrName:   session.QmgrName(),
			DLQName:    dlqName,
			Configured: true,
			Suggestion: fmt.Sprintf("DLQ '%s' is configured but the queue does not exist.", dlqName),
		}, nil
	}

	dlq := queues[0]
	currentDepth := toInt(dlq["current_queue_depth"])
	maxDepth := toInt(dlq["max_queue_depth"])
	depthPct := 0.0
	if maxDepth > 0 {
		depthPct = float64(currentDepth) / float64(maxDepth) * 100.0
	}

	suggestion := "DLQ is healthy."
	switch {
	case currentDepth == 0:
		suggestion = "DLQ is empty. No action needed."
	case depthPct >= criticalDepthPct:
		suggestion = "DLQ is near capacity. Investigate and clear undeliverable messages urgently."
	case currentDepth > 0:
		suggestion = "DLQ has messages. Investigate undeliverable messages."
	}

	return DLQReport{
		QmgrName:     session.QmgrName(),
		DLQName:      dlqName,
		Configured:   true,
		CurrentDepth: currentDepth,
		MaxDepth:     maxDepth,
		DepthPct:     depthPct,
		OpenInput:    toInt(dlq["open_input_count"]),
		OpenOutput:   toInt(dlq["open_output_count"]),
		Suggestion:   suggestion,
	}, nil
}

// PrintDLQInspection runs the DLQ inspection and prints formatted output.
func PrintDLQInspection(ctx context.Context, session *mqrestadmin.Session) (DLQReport, error) {
	report, err := InspectDLQ(ctx, session)
	if err != nil {
		return report, err
	}

	fmt.Printf("\n=== Dead Letter Queue: %s ===\n", report.QmgrName)
	fmt.Printf("  Configured: %t\n", report.Configured)
	dlqDisplay := report.DLQName
	if dlqDisplay == "" {
		dlqDisplay = "(none)"
	}
	fmt.Printf("  DLQ name:   %s\n", dlqDisplay)

	if report.Configured && report.DLQName != "" {
		fmt.Printf("  Depth:      %d / %d (%.1f%%)\n",
			report.CurrentDepth, report.MaxDepth, report.DepthPct)
		fmt.Printf("  Input:      %d\n", report.OpenInput)
		fmt.Printf("  Output:     %d\n", report.OpenOutput)
	}

	fmt.Printf("  Suggestion: %s\n", report.Suggestion)
	return report, nil
}
