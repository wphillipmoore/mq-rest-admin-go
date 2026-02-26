package mqrestadmin

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// EnsureQmgr ensures the queue manager has the specified attributes. Since
// the queue manager always exists, the result is never EnsureCreated.
func (session *Session) EnsureQmgr(ctx context.Context, requestParameters map[string]any) (EnsureResult, error) {
	if len(requestParameters) == 0 {
		return EnsureResult{Action: EnsureUnchanged}, nil
	}

	// DISPLAY current state
	current, err := session.DisplayQmgr(ctx, WithResponseParameters([]string{"all"}))
	if err != nil {
		return EnsureResult{}, fmt.Errorf("ensure qmgr display: %w", err)
	}

	// Compare and alter if needed
	changed, changedParams := diffAttributes(requestParameters, current)
	if len(changed) == 0 {
		return EnsureResult{Action: EnsureUnchanged}, nil
	}

	err = session.AlterQmgr(ctx, WithRequestParameters(changedParams))
	if err != nil {
		return EnsureResult{}, fmt.Errorf("ensure qmgr alter: %w", err)
	}

	return EnsureResult{Action: EnsureUpdated, Changed: changed}, nil
}

// EnsureQlocal ensures a local queue exists with the specified attributes.
func (session *Session) EnsureQlocal(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "QUEUE", "QLOCAL", "QLOCAL")
}

// EnsureQremote ensures a remote queue exists with the specified attributes.
func (session *Session) EnsureQremote(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "QUEUE", "QREMOTE", "QREMOTE")
}

// EnsureQalias ensures an alias queue exists with the specified attributes.
func (session *Session) EnsureQalias(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "QUEUE", "QALIAS", "QALIAS")
}

// EnsureQmodel ensures a model queue exists with the specified attributes.
func (session *Session) EnsureQmodel(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "QUEUE", "QMODEL", "QMODEL")
}

// EnsureChannel ensures a channel exists with the specified attributes.
func (session *Session) EnsureChannel(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "CHANNEL", "CHANNEL", "CHANNEL")
}

// EnsureAuthinfo ensures an authentication information object exists with the
// specified attributes.
func (session *Session) EnsureAuthinfo(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "AUTHINFO", "AUTHINFO", "AUTHINFO")
}

// EnsureListener ensures a listener exists with the specified attributes.
func (session *Session) EnsureListener(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "LISTENER", "LISTENER", "LISTENER")
}

// EnsureNamelist ensures a namelist exists with the specified attributes.
func (session *Session) EnsureNamelist(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "NAMELIST", "NAMELIST", "NAMELIST")
}

// EnsureProcess ensures a process exists with the specified attributes.
func (session *Session) EnsureProcess(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "PROCESS", "PROCESS", "PROCESS")
}

// EnsureService ensures a service exists with the specified attributes.
func (session *Session) EnsureService(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "SERVICE", "SERVICE", "SERVICE")
}

// EnsureTopic ensures a topic exists with the specified attributes.
func (session *Session) EnsureTopic(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "TOPIC", "TOPIC", "TOPIC")
}

// EnsureSub ensures a subscription exists with the specified attributes.
func (session *Session) EnsureSub(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "SUB", "SUB", "SUB")
}

// EnsureStgclass ensures a storage class exists with the specified attributes.
func (session *Session) EnsureStgclass(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "STGCLASS", "STGCLASS", "STGCLASS")
}

// EnsureComminfo ensures a communication information object exists with the
// specified attributes.
func (session *Session) EnsureComminfo(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "COMMINFO", "COMMINFO", "COMMINFO")
}

// EnsureCfstruct ensures a CF structure exists with the specified attributes.
func (session *Session) EnsureCfstruct(ctx context.Context, name string, requestParameters map[string]any) (EnsureResult, error) {
	return session.ensureObject(ctx, name, requestParameters, "CFSTRUCT", "CFSTRUCT", "CFSTRUCT")
}

// ensureObject implements the idempotent upsert pattern:
// 1. DISPLAY to check existence (command error treated as "not found")
// 2. DEFINE if missing
// 3. Compare attributes if found
// 4. ALTER if changed
func (session *Session) ensureObject(ctx context.Context, name string,
	requestParameters map[string]any, displayQualifier, defineQualifier, alterQualifier string,
) (EnsureResult, error) {
	// Step 1: DISPLAY to check existence
	currentObjects, err := session.mqscCommand(ctx, "DISPLAY", displayQualifier, &name,
		nil, []string{"all"}, nil, true)
	if err != nil {
		// Command errors mean the object doesn't exist
		var cmdErr *CommandError
		if !errors.As(err, &cmdErr) {
			return EnsureResult{}, fmt.Errorf("ensure %s display: %w", strings.ToLower(defineQualifier), err)
		}
		currentObjects = nil
	}

	// Step 2: Not found -> DEFINE
	if len(currentObjects) == 0 {
		_, err := session.mqscCommand(ctx, "DEFINE", defineQualifier, &name,
			requestParameters, nil, nil, false)
		if err != nil {
			return EnsureResult{}, fmt.Errorf("ensure %s define: %w", strings.ToLower(defineQualifier), err)
		}
		return EnsureResult{Action: EnsureCreated}, nil
	}

	// Step 3: No params to check -> UNCHANGED
	if len(requestParameters) == 0 {
		return EnsureResult{Action: EnsureUnchanged}, nil
	}

	// Step 4: Compare attributes
	current := currentObjects[0]
	changed, changedParams := diffAttributes(requestParameters, current)
	if len(changed) == 0 {
		return EnsureResult{Action: EnsureUnchanged}, nil
	}

	// Step 5: ALTER with only the changed attributes
	_, err = session.mqscCommand(ctx, "ALTER", alterQualifier, &name,
		changedParams, nil, nil, false)
	if err != nil {
		return EnsureResult{}, fmt.Errorf("ensure %s alter: %w", strings.ToLower(alterQualifier), err)
	}

	return EnsureResult{Action: EnsureUpdated, Changed: changed}, nil
}

// diffAttributes compares desired attributes against current values and
// returns the list of changed attribute names and a map of only the changed
// key-value pairs.
func diffAttributes(desired, current map[string]any) (changed []string, changedParams map[string]any) {
	changedParams = make(map[string]any)

	for key, desiredValue := range desired {
		currentValue, exists := current[key]
		if !exists || !valuesMatch(desiredValue, currentValue) {
			changed = append(changed, key)
			changedParams[key] = desiredValue
		}
	}

	return changed, changedParams
}

// valuesMatch compares two attribute values using case-insensitive string
// comparison after trimming whitespace, matching the Java port's behavior.
func valuesMatch(desired, current any) bool {
	desiredStr := strings.TrimSpace(fmt.Sprintf("%v", desired))
	currentStr := strings.TrimSpace(fmt.Sprintf("%v", current))
	return strings.EqualFold(desiredStr, currentStr)
}
