package mqrestadmin

// EnsureAction describes the action taken by an ensure operation.
type EnsureAction int

const (
	// EnsureCreated indicates the object was newly defined.
	EnsureCreated EnsureAction = iota
	// EnsureUpdated indicates the object existed and attributes were altered.
	EnsureUpdated
	// EnsureUnchanged indicates the object existed and all attributes already matched.
	EnsureUnchanged
)

func (action EnsureAction) String() string {
	switch action {
	case EnsureCreated:
		return "created"
	case EnsureUpdated:
		return "updated"
	case EnsureUnchanged:
		return "unchanged"
	default:
		return "unknown"
	}
}

// EnsureResult describes the outcome of an idempotent ensure operation.
type EnsureResult struct {
	// Action indicates whether the object was created, updated, or unchanged.
	Action EnsureAction
	// Changed lists the attribute names that triggered an ALTER, in the
	// caller's namespace (snake_case if mapping is enabled).
	Changed []string
}
