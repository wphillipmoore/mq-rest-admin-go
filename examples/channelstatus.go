// Package examples provides runnable example functions that demonstrate
// common MQ administration tasks using mqrestadmin.
package examples

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/wphillipmoore/mq-rest-admin-go/mqrestadmin"
)

// ChannelInfo holds combined channel definition and status information.
type ChannelInfo struct {
	Name           string
	ChannelType    string
	ConnectionName string
	Defined        bool
	Status         string
}

// ReportChannelStatus reports channel definitions and live status.
func ReportChannelStatus(ctx context.Context, session *mqrestadmin.Session) ([]ChannelInfo, error) {
	channels, err := session.DisplayChannel(ctx, "*")
	if err != nil {
		return nil, err
	}

	definitions := make(map[string]map[string]any)
	for _, ch := range channels {
		name := strings.TrimSpace(fmt.Sprint(ch["channel_name"]))
		if name != "" {
			definitions[name] = ch
		}
	}

	liveStatus := make(map[string]string)
	statuses, err := session.DisplayChstatus(ctx, "*")
	if err == nil {
		for _, entry := range statuses {
			name := strings.TrimSpace(fmt.Sprint(entry["channel_name"]))
			status := strings.TrimSpace(fmt.Sprint(entry["status"]))
			if name != "" {
				liveStatus[name] = status
			}
		}
	}

	var results []ChannelInfo

	names := make([]string, 0, len(definitions))
	for name := range definitions {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		defn := definitions[name]
		ctype := strings.TrimSpace(fmt.Sprint(defn["channel_type"]))
		conname := strings.TrimSpace(fmt.Sprint(defn["connection_name"]))
		status := "INACTIVE"
		if s, ok := liveStatus[name]; ok {
			status = s
		}
		results = append(results, ChannelInfo{
			Name:           name,
			ChannelType:    ctype,
			ConnectionName: conname,
			Defined:        true,
			Status:         status,
		})
	}

	statusNames := make([]string, 0, len(liveStatus))
	for name := range liveStatus {
		statusNames = append(statusNames, name)
	}
	sort.Strings(statusNames)

	for _, name := range statusNames {
		if _, ok := definitions[name]; !ok {
			results = append(results, ChannelInfo{
				Name:    name,
				Defined: false,
				Status:  liveStatus[name],
			})
		}
	}

	return results, nil
}

// PrintChannelStatus runs the channel status report and prints formatted output.
func PrintChannelStatus(ctx context.Context, session *mqrestadmin.Session) ([]ChannelInfo, error) {
	results, err := ReportChannelStatus(ctx, session)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n%-30s %-12s %-25s %-8s %s\n",
		"Channel", "Type", "Connection", "Defined", "Status")
	fmt.Println(strings.Repeat("-", 90))

	for _, info := range results {
		defined := "Yes"
		if !info.Defined {
			defined = "No"
		}
		fmt.Printf("%-30s %-12s %-25s %-8s %s\n",
			info.Name, info.ChannelType, info.ConnectionName, defined, info.Status)
	}

	var inactive []string
	for _, c := range results {
		if c.Defined && c.Status == "INACTIVE" {
			inactive = append(inactive, c.Name)
		}
	}
	if len(inactive) > 0 {
		fmt.Printf("\nDefined but inactive: %s\n", strings.Join(inactive, ", "))
	}

	return results, nil
}
