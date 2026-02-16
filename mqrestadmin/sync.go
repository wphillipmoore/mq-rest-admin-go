package mqrestadmin

import "time"

// SyncOperation describes the type of state transition performed by a sync
// method.
type SyncOperation int

const (
	// SyncStarted indicates the object transitioned to a running state.
	SyncStarted SyncOperation = iota
	// SyncStopped indicates the object transitioned to a stopped state.
	SyncStopped
	// SyncRestarted indicates the object was stopped then started.
	SyncRestarted
)

func (operation SyncOperation) String() string {
	switch operation {
	case SyncStarted:
		return "started"
	case SyncStopped:
		return "stopped"
	case SyncRestarted:
		return "restarted"
	default:
		return "unknown"
	}
}

// SyncConfig configures polling behavior for synchronous operations.
type SyncConfig struct {
	// Timeout is the maximum duration to wait for the operation to complete.
	// Defaults to 30 seconds if zero.
	Timeout time.Duration
	// PollInterval is the duration between status checks.
	// Defaults to 1 second if zero.
	PollInterval time.Duration
}

// SyncResult describes the outcome of a synchronous polling operation.
type SyncResult struct {
	// Operation indicates whether the object was started, stopped, or restarted.
	Operation SyncOperation
	// Polls is the number of status checks performed.
	Polls int
	// ElapsedSeconds is the wall-clock time taken.
	ElapsedSeconds float64
}
