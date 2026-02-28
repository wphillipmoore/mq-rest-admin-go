package examples

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

// QueueDepthInfo holds depth information for a single queue.
type QueueDepthInfo struct {
	Name         string
	CurrentDepth int
	MaxDepth     int
	DepthPct     float64
	OpenInput    int
	OpenOutput   int
	Warning      bool
}

// MonitorQueueDepths returns depth information for all local queues.
func MonitorQueueDepths(ctx context.Context, session *mqrestadmin.Session, thresholdPct float64) ([]QueueDepthInfo, error) {
	queues, err := session.DisplayQueue(ctx, "*")
	if err != nil {
		return nil, err
	}

	var results []QueueDepthInfo

	for _, queue := range queues {
		qtype := strings.ToUpper(strings.TrimSpace(fmt.Sprint(queue["type"])))
		if qtype != "QLOCAL" && qtype != "LOCAL" {
			continue
		}

		currentDepth := toInt(queue["current_queue_depth"])
		maxDepth := toInt(queue["max_queue_depth"])
		depthPct := 0.0
		if maxDepth > 0 {
			depthPct = float64(currentDepth) / float64(maxDepth) * 100.0
		}

		results = append(results, QueueDepthInfo{
			Name:         strings.TrimSpace(fmt.Sprint(queue["queue_name"])),
			CurrentDepth: currentDepth,
			MaxDepth:     maxDepth,
			DepthPct:     depthPct,
			OpenInput:    toInt(queue["open_input_count"]),
			OpenOutput:   toInt(queue["open_output_count"]),
			Warning:      depthPct >= thresholdPct,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].DepthPct > results[j].DepthPct
	})

	return results, nil
}

// PrintQueueDepths runs the queue depth monitor and prints formatted output.
func PrintQueueDepths(ctx context.Context, session *mqrestadmin.Session, thresholdPct float64) ([]QueueDepthInfo, error) {
	results, err := MonitorQueueDepths(ctx, session, thresholdPct)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n%-40s %8s %8s %6s %4s %4s %s\n",
		"Queue", "Depth", "Max", "%", "In", "Out", "Status")
	fmt.Println(strings.Repeat("-", 90))

	warningCount := 0
	for _, info := range results {
		status := "OK"
		if info.Warning {
			status = "WARNING"
			warningCount++
		}
		fmt.Printf("%-40s %8d %8d %5.1f%% %4d %4d %s\n",
			info.Name, info.CurrentDepth, info.MaxDepth,
			info.DepthPct, info.OpenInput, info.OpenOutput, status)
	}

	fmt.Printf("\nTotal queues: %d, warnings: %d\n", len(results), warningCount)
	return results, nil
}

func toInt(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0
		}
		return n
	default:
		n, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(v)))
		if err != nil {
			return 0
		}
		return n
	}
}
