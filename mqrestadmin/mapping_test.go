package mqrestadmin

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

func TestCopyMap_PermissiveUnknownQualifier(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]any{"key1": "val1", "key2": "val2"}
	result, issues := mapper.mapRequestAttributes("nonexistent_qualifier", input, false)
	if len(issues) == 0 {
		t.Error("expected issues for unknown qualifier")
	}
	// In permissive mode, copyMap is used — result should be a copy
	if result["key1"] != "val1" || result["key2"] != "val2" {
		t.Errorf("expected values to pass through, got %v", result)
	}
}

func TestMergeNestedStringMap_NewKey(t *testing.T) {
	overrides := map[string]any{
		"qualifiers": map[string]any{
			"queue": map[string]any{
				"response_value_map": map[string]any{
					"NEW_KEY": map[string]any{
						"VAL1": "val_one",
					},
				},
			},
		},
	}

	mapper, err := newAttributeMapperWithOverrides(overrides, MappingOverrideMerge)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	qualifierData := mapper.data.Qualifiers["queue"]
	if qualifierData.ResponseValueMap["NEW_KEY"] == nil {
		t.Error("expected NEW_KEY in response value map after merge")
	}
	if qualifierData.ResponseValueMap["NEW_KEY"]["VAL1"] != "val_one" {
		t.Error("expected VAL1 = val_one")
	}
}

func TestMapValue_UnknownStringValue(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find a key that has a value map in response
	qualifierData := mapper.data.Qualifiers["queue"]
	if len(qualifierData.ResponseValueMap) == 0 {
		t.Skip("no response value mappings for queue qualifier")
	}

	for key := range qualifierData.ResponseValueMap {
		input := map[string]any{key: "TOTALLY_UNKNOWN_VALUE_XYZ"}
		_, issues := mapper.mapResponseAttributes("queue", input, false)
		foundUnknownValue := false
		for _, issue := range issues {
			if issue.Reason == MappingUnknownValue {
				foundUnknownValue = true
				break
			}
		}
		if !foundUnknownValue {
			t.Errorf("expected MappingUnknownValue issue for key %s", key)
		}
		break
	}
}

func TestMapValue_ListWithNonStringItems(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	qualifierData := mapper.data.Qualifiers["queue"]
	if len(qualifierData.ResponseValueMap) == 0 {
		t.Skip("no response value mappings")
	}

	for key := range qualifierData.ResponseValueMap {
		input := map[string]any{key: []any{42, true}}
		result, _ := mapper.mapResponseAttributes("queue", input, false)
		mappedKey := key
		if snakeName, found := qualifierData.ResponseKeyMap[key]; found {
			mappedKey = snakeName
		}
		if resultList, isList := result[mappedKey].([]any); isList {
			if resultList[0] != 42 || resultList[1] != true {
				t.Errorf("non-string items should pass through, got %v", resultList)
			}
		}
		break
	}
}

func TestMapValue_ListWithUnknownString(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	qualifierData := mapper.data.Qualifiers["queue"]
	if len(qualifierData.ResponseValueMap) == 0 {
		t.Skip("no response value mappings")
	}

	for key := range qualifierData.ResponseValueMap {
		input := map[string]any{key: []any{"TOTALLY_UNKNOWN_XYZ"}}
		_, issues := mapper.mapResponseAttributes("queue", input, false)
		foundUnknownValue := false
		for _, issue := range issues {
			if issue.Reason == MappingUnknownValue {
				foundUnknownValue = true
				break
			}
		}
		if !foundUnknownValue {
			t.Errorf("expected MappingUnknownValue issue for list item in key %s", key)
		}
		break
	}
}

func TestMapValue_NonStringNonList(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	qualifierData := mapper.data.Qualifiers["queue"]
	if len(qualifierData.ResponseValueMap) == 0 {
		t.Skip("no response value mappings")
	}

	for key := range qualifierData.ResponseValueMap {
		input := map[string]any{key: 42}
		result, _ := mapper.mapResponseAttributes("queue", input, false)
		mappedKey := key
		if snakeName, found := qualifierData.ResponseKeyMap[key]; found {
			mappedKey = snakeName
		}
		if result[mappedKey] != 42 {
			t.Errorf("non-string/non-list value should pass through, got %v", result[mappedKey])
		}
		break
	}
}

func TestNewAttributeMapperWithOverrides_MergeNewQualifier(t *testing.T) {
	overrides := map[string]any{
		"qualifiers": map[string]any{
			"brand_new_qualifier": map[string]any{
				"request_key_map": map[string]any{
					"my_attr": "MYATTR",
				},
			},
		},
	}

	mapper, err := newAttributeMapperWithOverrides(overrides, MappingOverrideMerge)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]any{"my_attr": "value"}
	result, _ := mapper.mapRequestAttributes("brand_new_qualifier", input, false)
	if result["MYATTR"] != "value" {
		t.Errorf("expected MYATTR in result, got %v", result)
	}
}

func TestNewAttributeMapperWithOverrides_RequestKeyValueMap(t *testing.T) {
	overrides := map[string]any{
		"qualifiers": map[string]any{
			"queue": map[string]any{
				"request_key_value_map": map[string]any{
					"custom_flag": map[string]any{
						"key":   "FLAGATTR",
						"value": "YES",
					},
				},
			},
		},
	}

	mapper, err := newAttributeMapperWithOverrides(overrides, MappingOverrideMerge)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]any{"custom_flag": "ignored"}
	result, _ := mapper.mapRequestAttributes("queue", input, false)
	if result["FLAGATTR"] != "YES" {
		t.Errorf("expected FLAGATTR=YES from key-value map, got %v", result)
	}
}

func TestResolveResponseParameterMacros_NoMacros(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// A non-existent command should return params unchanged
	result := mapper.resolveResponseParameterMacros("NONEXISTENT", "THING", []string{"all"})
	if len(result) != 1 || result[0] != "all" {
		t.Errorf("expected unchanged params, got %v", result)
	}
}

func TestResolveResponseParameterMacros_NotAll(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Specific params (not "all") should not trigger macro expansion
	result := mapper.resolveResponseParameterMacros("DISPLAY", "QUEUE", []string{"MAXDEPTH"})
	if len(result) != 1 || result[0] != "MAXDEPTH" {
		t.Errorf("expected unchanged params, got %v", result)
	}
}

func TestMapResponseList_StrictWithIssues(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	objects := []map[string]any{
		{"UNKNOWN_ATTR_XYZ": "val"},
	}

	_, issues := mapper.mapResponseList("queue", objects, true)
	if len(issues) == 0 {
		t.Error("expected issues for unknown attribute in strict mode")
	}
}

func TestMergeNestedStringMap_ExistingKey(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find an existing key in some qualifier's response_value_map
	var existingKey string
	for key := range mapper.data.Qualifiers["queue"].ResponseValueMap {
		existingKey = key
		break
	}
	if existingKey == "" {
		t.Skip("no existing response value map keys for queue")
	}

	overrides := map[string]any{
		"qualifiers": map[string]any{
			"queue": map[string]any{
				"response_value_map": map[string]any{
					existingKey: map[string]any{
						"CUSTOM_OVERRIDE_VAL": "custom_override",
					},
				},
			},
		},
	}

	mergedMapper, err := newAttributeMapperWithOverrides(overrides, MappingOverrideMerge)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	qualifierData := mergedMapper.data.Qualifiers["queue"]
	if qualifierData.ResponseValueMap[existingKey]["CUSTOM_OVERRIDE_VAL"] != "custom_override" {
		t.Error("expected merged value in existing key")
	}
}

func TestMapResponseList_ObjectIndexInIssues(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find a key with value mappings and use an unknown value through mapResponseList
	var valueMapKey string
	for key := range mapper.data.Qualifiers["queue"].ResponseValueMap {
		valueMapKey = key
		break
	}
	if valueMapKey == "" {
		t.Skip("no response value map keys")
	}

	// Use unknown value to trigger MappingUnknownValue with objectIndex set
	objects := []map[string]any{
		{valueMapKey: "UNKNOWN_TEST_VALUE_XYZ"},
	}

	_, issues := mapper.mapResponseList("queue", objects, true)
	foundWithIndex := false
	for _, issue := range issues {
		if issue.ObjectIndex != nil && *issue.ObjectIndex == 0 {
			foundWithIndex = true
			break
		}
	}
	if !foundWithIndex {
		t.Error("expected issue with ObjectIndex = 0")
	}
}

func TestMapResponseList_ObjectIndexInListIssues(t *testing.T) {
	mapper, err := newAttributeMapper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var valueMapKey string
	for key := range mapper.data.Qualifiers["queue"].ResponseValueMap {
		valueMapKey = key
		break
	}
	if valueMapKey == "" {
		t.Skip("no response value map keys")
	}

	objects := []map[string]any{
		{valueMapKey: []any{"UNKNOWN_LIST_VAL_XYZ"}},
	}

	_, issues := mapper.mapResponseList("queue", objects, true)
	foundWithIndex := false
	for _, issue := range issues {
		if issue.ObjectIndex != nil {
			foundWithIndex = true
			break
		}
	}
	if !foundWithIndex {
		t.Error("expected issue with ObjectIndex set for list value")
	}
}

func TestMapResponseParameterNames_UnknownQualifier_ViaSession(t *testing.T) {
	transport := newMockTransport()
	transport.addSuccessResponse()
	session := newTestSessionWithMapping(transport)
	session.mappingStrict = false

	// Call mqscCommand with a known command but unknown qualifier to trigger
	// mapResponseParameterNames returning early for unknown qualifier
	name := "TEST"
	_, _ = session.mqscCommand(context.Background(), "DISPLAY", "QUEUE", &name,
		nil, []string{"unknown_param"}, nil, true)
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
