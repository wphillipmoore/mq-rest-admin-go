package mqrestadmin

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// objectTypeConfig defines the MQSC qualifiers and status keys for a
// specific MQ object type used in sync operations.
type objectTypeConfig struct {
	startQualifier    string
	stopQualifier     string
	statusQualifier   string
	statusKeys        []string
	emptyMeansStopped bool
}

var (
	channelConfig = objectTypeConfig{
		startQualifier:    "CHANNEL",
		stopQualifier:     "CHANNEL",
		statusQualifier:   "CHSTATUS",
		statusKeys:        []string{"channel_status", "STATUS"},
		emptyMeansStopped: true,
	}
	listenerConfig = objectTypeConfig{
		startQualifier:    "LISTENER",
		stopQualifier:     "LISTENER",
		statusQualifier:   "LSSTATUS",
		statusKeys:        []string{"status", "STATUS"},
		emptyMeansStopped: false,
	}
	serviceConfig = objectTypeConfig{
		startQualifier:    "SERVICE",
		stopQualifier:     "SERVICE",
		statusQualifier:   "SVSTATUS",
		statusKeys:        []string{"status", "STATUS"},
		emptyMeansStopped: false,
	}
)

var runningValues = map[string]bool{
	"RUNNING": true,
	"running": true,
}

var stoppedValues = map[string]bool{
	"STOPPED":  true,
	"stopped":  true,
	"INACTIVE": true,
	"inactive": true,
}

// StartChannelSync starts a channel and polls until it is running.
func (session *Session) StartChannelSync(ctx context.Context, name string, config SyncConfig) (SyncResult, error) {
	return session.startAndPoll(ctx, name, &channelConfig, config)
}

// StopChannelSync stops a channel and polls until it is stopped.
func (session *Session) StopChannelSync(ctx context.Context, name string, config SyncConfig) (SyncResult, error) {
	return session.stopAndPoll(ctx, name, &channelConfig, config)
}

// RestartChannel stops and then starts a channel, polling at each step.
func (session *Session) RestartChannel(ctx context.Context, name string, config SyncConfig) (SyncResult, error) {
	return session.restartObject(ctx, name, &channelConfig, config)
}

// StartListenerSync starts a listener and polls until it is running.
func (session *Session) StartListenerSync(ctx context.Context, name string, config SyncConfig) (SyncResult, error) {
	return session.startAndPoll(ctx, name, &listenerConfig, config)
}

// StopListenerSync stops a listener and polls until it is stopped.
func (session *Session) StopListenerSync(ctx context.Context, name string, config SyncConfig) (SyncResult, error) {
	return session.stopAndPoll(ctx, name, &listenerConfig, config)
}

// RestartListener stops and then starts a listener, polling at each step.
func (session *Session) RestartListener(ctx context.Context, name string, config SyncConfig) (SyncResult, error) {
	return session.restartObject(ctx, name, &listenerConfig, config)
}

// StartServiceSync starts a service and polls until it is running.
func (session *Session) StartServiceSync(ctx context.Context, name string, config SyncConfig) (SyncResult, error) {
	return session.startAndPoll(ctx, name, &serviceConfig, config)
}

// StopServiceSync stops a service and polls until it is stopped.
func (session *Session) StopServiceSync(ctx context.Context, name string, config SyncConfig) (SyncResult, error) {
	return session.stopAndPoll(ctx, name, &serviceConfig, config)
}

// RestartService stops and then starts a service, polling at each step.
func (session *Session) RestartService(ctx context.Context, name string, config SyncConfig) (SyncResult, error) {
	return session.restartObject(ctx, name, &serviceConfig, config)
}

func (session *Session) startAndPoll(ctx context.Context, name string,
	objectConfig *objectTypeConfig, syncConfig SyncConfig,
) (SyncResult, error) {
	syncConfig, err := normalizeSyncConfig(syncConfig)
	if err != nil {
		return SyncResult{}, err
	}

	// Issue START command
	_, err = session.mqscCommand(ctx, "START", objectConfig.startQualifier, &name,
		nil, nil, nil, false)
	if err != nil {
		return SyncResult{}, err
	}

	// Poll for RUNNING status
	startTime := session.clock.now()
	polls := 0

	for {
		session.clock.sleep(syncConfig.PollInterval)

		var statusRows []map[string]any
		statusRows, err = session.queryStatus(ctx, name, objectConfig)
		if err != nil {
			return SyncResult{}, err
		}
		polls++

		if hasStatus(statusRows, objectConfig.statusKeys, runningValues) {
			elapsed := session.clock.now().Sub(startTime).Seconds()
			return SyncResult{Operation: SyncStarted, Polls: polls, ElapsedSeconds: elapsed}, nil
		}

		elapsed := session.clock.now().Sub(startTime).Seconds()
		if elapsed >= syncConfig.Timeout.Seconds() {
			return SyncResult{}, &TimeoutError{
				Name:           name,
				Operation:      SyncStarted,
				ElapsedSeconds: elapsed,
			}
		}
	}
}

func (session *Session) stopAndPoll(ctx context.Context, name string,
	objectConfig *objectTypeConfig, syncConfig SyncConfig,
) (SyncResult, error) {
	syncConfig, err := normalizeSyncConfig(syncConfig)
	if err != nil {
		return SyncResult{}, err
	}

	// Issue STOP command
	_, err = session.mqscCommand(ctx, "STOP", objectConfig.stopQualifier, &name,
		nil, nil, nil, false)
	if err != nil {
		return SyncResult{}, err
	}

	// Poll for STOPPED status
	startTime := session.clock.now()
	polls := 0

	for {
		session.clock.sleep(syncConfig.PollInterval)

		var statusRows []map[string]any
		statusRows, err = session.queryStatus(ctx, name, objectConfig)
		if err != nil {
			return SyncResult{}, err
		}
		polls++

		// Empty status means stopped for channels
		if len(statusRows) == 0 && objectConfig.emptyMeansStopped {
			elapsed := session.clock.now().Sub(startTime).Seconds()
			return SyncResult{Operation: SyncStopped, Polls: polls, ElapsedSeconds: elapsed}, nil
		}

		if hasStatus(statusRows, objectConfig.statusKeys, stoppedValues) {
			elapsed := session.clock.now().Sub(startTime).Seconds()
			return SyncResult{Operation: SyncStopped, Polls: polls, ElapsedSeconds: elapsed}, nil
		}

		elapsed := session.clock.now().Sub(startTime).Seconds()
		if elapsed >= syncConfig.Timeout.Seconds() {
			return SyncResult{}, &TimeoutError{
				Name:           name,
				Operation:      SyncStopped,
				ElapsedSeconds: elapsed,
			}
		}
	}
}

func (session *Session) restartObject(ctx context.Context, name string,
	objectConfig *objectTypeConfig, syncConfig SyncConfig,
) (SyncResult, error) {
	stopResult, err := session.stopAndPoll(ctx, name, objectConfig, syncConfig)
	if err != nil {
		return SyncResult{}, err
	}

	startResult, err := session.startAndPoll(ctx, name, objectConfig, syncConfig)
	if err != nil {
		return SyncResult{}, err
	}

	return SyncResult{
		Operation:      SyncRestarted,
		Polls:          stopResult.Polls + startResult.Polls,
		ElapsedSeconds: stopResult.ElapsedSeconds + startResult.ElapsedSeconds,
	}, nil
}

func (session *Session) queryStatus(ctx context.Context, name string, objectConfig *objectTypeConfig) ([]map[string]any, error) {
	rows, err := session.mqscCommand(ctx, "DISPLAY", objectConfig.statusQualifier, &name,
		nil, []string{"all"}, nil, true)
	if err != nil {
		// Swallow command errors — object not found during polling is expected
		var cmdErr *CommandError
		if errors.As(err, &cmdErr) {
			return nil, nil
		}
		return nil, err
	}
	return rows, nil
}

func hasStatus(rows []map[string]any, statusKeys []string, targetValues map[string]bool) bool {
	for _, row := range rows {
		for _, key := range statusKeys {
			if value, exists := row[key]; exists {
				if strValue, isStr := value.(string); isStr {
					if targetValues[strings.TrimSpace(strValue)] {
						return true
					}
				}
			}
		}
	}
	return false
}

func normalizeSyncConfig(config SyncConfig) (SyncConfig, error) {
	if config.Timeout < 0 {
		return SyncConfig{}, fmt.Errorf("Timeout must not be negative, got %v", config.Timeout)
	}
	if config.PollInterval < 0 {
		return SyncConfig{}, fmt.Errorf("PollInterval must not be negative, got %v", config.PollInterval)
	}
	if config.Timeout == 0 {
		config.Timeout = defaultSyncTimeout
	}
	if config.PollInterval == 0 {
		config.PollInterval = defaultPollInterval
	}
	return config, nil
}
