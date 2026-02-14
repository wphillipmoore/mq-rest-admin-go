package mqrest

import (
	"context"
	"testing"
)

func TestNewAttributeMapper_LoadsSuccessfully(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mapper == nil {
		t.Fatal("expected non-nil mapper")
	}
	if len(mapper.data.Commands) == 0 {
		t.Error("expected commands to be loaded")
	}
	if len(mapper.data.Qualifiers) == 0 {
		t.Error("expected qualifiers to be loaded")
	}
}

func TestResolveMappingQualifier(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		command   string
		qualifier string
		expected  string
	}{
		{"DISPLAY", "QUEUE", "queue"},
		{"DISPLAY", "CHANNEL", "channel"},
		{"ALTER", "QMGR", "qmgr"},
		{"DEFINE", "CHANNEL", "channel"},
		{"CLEAR", "QLOCAL", "queue"},
		{"NONEXISTENT", "THING", ""},
	}

	for _, test := range tests {
		result := mapper.resolveMappingQualifier(test.command, test.qualifier)
		if result != test.expected {
			t.Errorf("resolveMappingQualifier(%q, %q) = %q, want %q",
				test.command, test.qualifier, result, test.expected)
		}
	}
}

func TestMapRequestAttributes_KeyMap(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]any{
		"max_queue_depth": "5000",
	}

	result, issues := mapper.mapRequestAttributes("queue", input, false)
	if len(issues) > 0 {
		t.Logf("mapping issues (permissive): %v", issues)
	}

	// The key should be mapped to the MQSC equivalent
	if _, hasOriginal := result["max_queue_depth"]; hasOriginal {
		t.Error("original key should be replaced")
	}
	if result["MAXDEPTH"] == nil {
		t.Error("expected MAXDEPTH key in mapped result")
	}
}

func TestMapResponseAttributes_KeyMap(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]any{
		"MAXDEPTH": "5000",
		"QUEUE":    "TEST.Q",
	}

	result, issues := mapper.mapResponseAttributes("queue", input, false)
	if len(issues) > 0 {
		t.Logf("mapping issues (permissive): %v", issues)
	}

	if result["max_queue_depth"] == nil {
		t.Error("expected max_queue_depth key in mapped result")
	}
	if result["queue_name"] == nil {
		t.Error("expected queue_name key in mapped result")
	}
}

func TestMapRequestAttributes_StrictMode_UnknownKey(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]any{
		"nonexistent_attribute": "value",
	}

	_, issues := mapper.mapRequestAttributes("queue", input, true)
	if len(issues) == 0 {
		t.Error("expected mapping issues for unknown key in strict mode")
	}

	foundUnknownKey := false
	for _, issue := range issues {
		if issue.Reason == MappingUnknownKey {
			foundUnknownKey = true
			break
		}
	}
	if !foundUnknownKey {
		t.Error("expected MappingUnknownKey issue")
	}
}

func TestMapRequestAttributes_PermissiveMode_UnknownKey(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]any{
		"nonexistent_attribute": "value",
	}

	result, issues := mapper.mapRequestAttributes("queue", input, false)
	if len(issues) == 0 {
		t.Error("expected issues even in permissive mode")
	}
	// Value should pass through
	if result["nonexistent_attribute"] != "value" {
		t.Error("unknown key should pass through in permissive mode")
	}
}

func TestMapAttributes_UnknownQualifier(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]any{"key": "value"}
	_, issues := mapper.mapRequestAttributes("nonexistent_qualifier", input, true)

	if len(issues) == 0 {
		t.Error("expected issues for unknown qualifier")
	}

	foundQualifierIssue := false
	for _, issue := range issues {
		if issue.Reason == MappingUnknownQualifier {
			foundQualifierIssue = true
			break
		}
	}
	if !foundQualifierIssue {
		t.Error("expected MappingUnknownQualifier issue")
	}
}

func TestMapResponseList(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	objects := []map[string]any{
		{"QUEUE": "Q1", "MAXDEPTH": "5000"},
		{"QUEUE": "Q2", "MAXDEPTH": "10000"},
	}

	result, issues := mapper.mapResponseList("queue", objects, false)
	if len(issues) > 0 {
		t.Logf("mapping issues: %v", issues)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0]["queue_name"] == nil {
		t.Error("expected queue_name in first object")
	}
	if result[1]["queue_name"] == nil {
		t.Error("expected queue_name in second object")
	}
}

func TestMapAttributes_EmptyInput(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, issues := mapper.mapRequestAttributes("queue", map[string]any{}, true)
	if len(issues) != 0 {
		t.Errorf("expected no issues for empty input, got %d", len(issues))
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestNewAttributeMapperWithOverrides_Merge(t *testing.T) {
	overrides := map[string]any{
		"qualifiers": map[string]any{
			"queue": map[string]any{
				"request_key_map": map[string]any{
					"custom_attr": "CUSTOMATTR",
				},
			},
		},
	}

	mapper, err := newAttributeMapperWithOverrides(overrides, MappingOverrideMerge)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Custom attr should be available
	input := map[string]any{"custom_attr": "value"}
	result, _ := mapper.mapRequestAttributes("queue", input, false)
	if result["CUSTOMATTR"] != "value" {
		t.Error("expected custom override to be applied")
	}

	// Original attrs should still work
	input2 := map[string]any{"max_queue_depth": "5000"}
	result2, _ := mapper.mapRequestAttributes("queue", input2, false)
	if result2["MAXDEPTH"] == nil {
		t.Error("expected original mapping to still work after merge override")
	}
}

func TestNewAttributeMapperWithOverrides_Replace(t *testing.T) {
	overrides := map[string]any{
		"qualifiers": map[string]any{
			"queue": map[string]any{
				"request_key_map": map[string]any{
					"custom_attr": "CUSTOMATTR",
				},
			},
		},
	}

	mapper, err := newAttributeMapperWithOverrides(overrides, MappingOverrideReplace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Custom attr should be available
	input := map[string]any{"custom_attr": "value"}
	result, _ := mapper.mapRequestAttributes("queue", input, false)
	if result["CUSTOMATTR"] != "value" {
		t.Error("expected custom override to be applied")
	}
}

func TestMqscCommand_WithMapping(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse(map[string]any{
		"QUEUE":    "TEST.Q",
		"MAXDEPTH": float64(5000),
	})
	session := newTestSessionWithMapping(transport)

	queues, err := session.DisplayQueue(context.Background(), "TEST.Q")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(queues) != 1 {
		t.Fatalf("expected 1 queue, got %d", len(queues))
	}

	// Response attributes should be mapped to snake_case
	if queues[0]["queue_name"] == nil {
		t.Error("expected queue_name (mapped from QUEUE)")
	}
	if queues[0]["max_queue_depth"] == nil {
		t.Error("expected max_queue_depth (mapped from MAXDEPTH)")
	}
}

func TestMapValue_StringMapping(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find a qualifier that has value mappings
	qualifierData, exists := mapper.data.Qualifiers["queue"]
	if !exists {
		t.Skip("queue qualifier not found in mapping data")
	}

	// Check if there's a value map for response
	if len(qualifierData.ResponseValueMap) == 0 {
		t.Skip("no response value mappings for queue qualifier")
	}

	// Value mapping should translate known values
	for key, valueMap := range qualifierData.ResponseValueMap {
		for mqscValue, snakeValue := range valueMap {
			input := map[string]any{key: mqscValue}
			result, _ := mapper.mapResponseAttributes("queue", input, false)
			// After key mapping, the result key is the snake_case name
			mappedKey := key
			if snakeName, found := qualifierData.ResponseKeyMap[key]; found {
				mappedKey = snakeName
			}
			if result[mappedKey] != snakeValue {
				t.Errorf("value mapping for %s=%s: got %v, want %s",
					key, mqscValue, result[mappedKey], snakeValue)
			}
			break // just test one
		}
		break // just test one key
	}
}

func TestMapValue_ListMapping(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that list values are mapped element by element
	qualifierData := mapper.data.Qualifiers["queue"]
	if len(qualifierData.ResponseValueMap) == 0 {
		t.Skip("no value mappings to test")
	}

	for key, valueMap := range qualifierData.ResponseValueMap {
		for mqscValue, expectedValue := range valueMap {
			input := map[string]any{key: []any{mqscValue}}
			result, _ := mapper.mapResponseAttributes("queue", input, false)
			// After key mapping, the result key is the snake_case name
			mappedKey := key
			if snakeName, found := qualifierData.ResponseKeyMap[key]; found {
				mappedKey = snakeName
			}
			if resultList, isList := result[mappedKey].([]any); isList {
				if len(resultList) != 1 || resultList[0] != expectedValue {
					t.Errorf("list value mapping: got %v, want [%s]", resultList, expectedValue)
				}
			}
			break
		}
		break
	}
}
