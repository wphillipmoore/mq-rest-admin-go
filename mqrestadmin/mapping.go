package mqrestadmin

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MappingOverrideMode controls how mapping overrides are applied.
type MappingOverrideMode int

const (
	// MappingOverrideMerge overlays overrides onto the default mapping data,
	// adding or replacing individual entries.
	MappingOverrideMerge MappingOverrideMode = iota
	// MappingOverrideReplace replaces entire qualifier sections with the
	// override data.
	MappingOverrideReplace
)

// mappingData holds the parsed mapping definitions loaded from the embedded
// JSON resource.
type mappingData struct {
	Commands   map[string]commandMapping   `json:"commands"`
	Qualifiers map[string]qualifierMapping `json:"qualifiers"`
	Version    int                         `json:"version"`
}

type commandMapping struct {
	Qualifier               string   `json:"qualifier"`
	ResponseParameterMacros []string `json:"response_parameter_macros,omitempty"`
}

type qualifierMapping struct {
	RequestKeyMap      map[string]string            `json:"request_key_map"`
	RequestValueMap    map[string]map[string]string `json:"request_value_map"`
	RequestKeyValueMap map[string]map[string]keyValueEntry `json:"request_key_value_map"`
	ResponseKeyMap     map[string]string            `json:"response_key_map"`
	ResponseValueMap   map[string]map[string]string `json:"response_value_map"`
}

type keyValueEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// attributeMapper translates attribute names and values between snake_case
// (caller-facing) and MQSC parameter names (API-facing).
type attributeMapper struct {
	data *mappingData
}

// newAttributeMapper creates an attribute mapper from the default embedded
// mapping data.
func newAttributeMapper() (*attributeMapper, error) {
	var data mappingData
	if err := json.Unmarshal(mappingDataJSON, &data); err != nil { // coverage-ignore -- embedded JSON is valid by construction
		return nil, fmt.Errorf("parse mapping data: %w", err)
	}
	return &attributeMapper{data: &data}, nil
}

// newAttributeMapperWithOverrides creates an attribute mapper with custom
// overrides applied to the default mapping data.
func newAttributeMapperWithOverrides(overrides map[string]any, mode MappingOverrideMode) (*attributeMapper, error) {
	mapper, err := newAttributeMapper()
	if err != nil { // coverage-ignore -- newAttributeMapper only fails on invalid embedded data
		return nil, err
	}

	overrideBytes, err := json.Marshal(overrides)
	if err != nil { // coverage-ignore -- json.Marshal on map[string]any cannot fail
		return nil, fmt.Errorf("marshal mapping overrides: %w", err)
	}

	var overrideData mappingData
	if err := json.Unmarshal(overrideBytes, &overrideData); err != nil { // coverage-ignore -- re-parsed from valid Marshal output
		return nil, fmt.Errorf("parse mapping overrides: %w", err)
	}

	if mode == MappingOverrideReplace {
		for qualifier, override := range overrideData.Qualifiers {
			mapper.data.Qualifiers[qualifier] = override
		}
	} else {
		for qualifier, override := range overrideData.Qualifiers {
			existing, exists := mapper.data.Qualifiers[qualifier]
			if !exists {
				mapper.data.Qualifiers[qualifier] = override
				continue
			}
			mergeStringMap(existing.RequestKeyMap, override.RequestKeyMap)
			mergeNestedStringMap(existing.RequestValueMap, override.RequestValueMap)
			for key, valueEntries := range override.RequestKeyValueMap {
				if existingEntries, exists := existing.RequestKeyValueMap[key]; exists {
					for val, entry := range valueEntries {
						existingEntries[val] = entry
					}
				} else {
					existing.RequestKeyValueMap[key] = valueEntries
				}
			}
			mergeStringMap(existing.ResponseKeyMap, override.ResponseKeyMap)
			mergeNestedStringMap(existing.ResponseValueMap, override.ResponseValueMap)
			mapper.data.Qualifiers[qualifier] = existing
		}
	}

	return mapper, nil
}

// defaultMappingQualifiers maps MQSC qualifier names to their mapping
// qualifier when the exact command+qualifier combination is not in the
// commands map.  This allows commands like DEFINE QLOCAL (not explicitly
// listed) to resolve to the "queue" mapping qualifier, matching the
// Python and Ruby SDKs.
var defaultMappingQualifiers = map[string]string{
	"QUEUE":    "queue",
	"QLOCAL":  "queue",
	"QREMOTE": "queue",
	"QALIAS":  "queue",
	"QMODEL":  "queue",
	"QMSTATUS": "qmstatus",
	"QSTATUS": "qstatus",
	"CHANNEL": "channel",
	"QMGR":   "qmgr",
}

// resolveMappingQualifier looks up the mapping qualifier for a given MQSC
// command and qualifier combination. For example, "DISPLAY" + "QLOCAL"
// resolves to "queue". If the exact command is not in the commands map,
// falls back to the default qualifier for the MQSC qualifier name.
func (mapper *attributeMapper) resolveMappingQualifier(command, mqscQualifier string) string {
	key := command + " " + mqscQualifier
	if cmdMapping, exists := mapper.data.Commands[key]; exists {
		return cmdMapping.Qualifier
	}
	if fallback, exists := defaultMappingQualifiers[mqscQualifier]; exists {
		return fallback
	}
	return strings.ToLower(mqscQualifier)
}

// mapRequestAttributes translates request attributes from snake_case to MQSC
// parameter names using the 3-layer pipeline.
func (mapper *attributeMapper) mapRequestAttributes(qualifier string,
	attributes map[string]any, strict bool,
) (map[string]any, []MappingIssue) {
	return mapper.mapAttributes(qualifier, attributes, strict, MappingRequest, -1)
}

// mapResponseAttributes translates response attributes from MQSC parameter
// names to snake_case using the key map and value map layers.
func (mapper *attributeMapper) mapResponseAttributes(qualifier string,
	attributes map[string]any, strict bool,
) (map[string]any, []MappingIssue) {
	return mapper.mapAttributes(qualifier, attributes, strict, MappingResponse, -1)
}

// mapResponseList translates a list of response attribute maps, tracking
// object indices in any mapping issues.
func (mapper *attributeMapper) mapResponseList(qualifier string,
	objects []map[string]any, strict bool,
) ([]map[string]any, []MappingIssue) {
	var allIssues []MappingIssue
	result := make([]map[string]any, len(objects))

	for idx, object := range objects {
		mapped, issues := mapper.mapAttributes(qualifier, object, false, MappingResponse, idx)
		result[idx] = mapped
		allIssues = append(allIssues, issues...)
	}

	if strict && len(allIssues) > 0 {
		return result, allIssues
	}
	return result, allIssues
}

func (mapper *attributeMapper) mapAttributes(qualifier string,
	attributes map[string]any, strict bool, direction MappingDirection, objectIndex int,
) (map[string]any, []MappingIssue) {
	if len(attributes) == 0 {
		return make(map[string]any), nil
	}

	qualifierData, exists := mapper.data.Qualifiers[qualifier]
	if !exists {
		issues := []MappingIssue{{
			Direction:     direction,
			Reason:        MappingUnknownQualifier,
			AttributeName: qualifier,
			Qualifier:     qualifier,
		}}
		if strict {
			return attributes, issues
		}
		return copyMap(attributes), issues
	}

	var keyMap map[string]string
	var valueMap map[string]map[string]string
	var keyValueMap map[string]map[string]keyValueEntry

	if direction == MappingRequest {
		keyMap = qualifierData.RequestKeyMap
		valueMap = qualifierData.RequestValueMap
		keyValueMap = qualifierData.RequestKeyValueMap
	} else {
		keyMap = qualifierData.ResponseKeyMap
		valueMap = qualifierData.ResponseValueMap
		keyValueMap = nil // key-value map is request-only
	}

	result := make(map[string]any, len(attributes))
	var issues []MappingIssue

	for key, value := range attributes {
		// The MQ REST API returns MQSC parameter names in lowercase, but the
		// mapping data uses uppercase. Normalize response keys to uppercase
		// for all lookups.
		lookupKey := key
		if direction == MappingResponse {
			lookupKey = strings.ToUpper(key)
		}

		// Layer 1: Key-value map (request only)
		if keyValueMap != nil {
			if valueEntries, found := keyValueMap[lookupKey]; found {
				if strValue, ok := value.(string); ok {
					if entry, found := valueEntries[strValue]; found {
						result[entry.Key] = entry.Value
						continue
					}
				}
			}
		}

		// Layer 2: Key map
		mappedKey := key
		if keyMap != nil {
			if mapped, found := keyMap[lookupKey]; found {
				mappedKey = mapped
			} else {
				issue := MappingIssue{
					Direction:     direction,
					Reason:        MappingUnknownKey,
					AttributeName: key,
					Qualifier:     qualifier,
				}
				if objectIndex >= 0 {
					issue.ObjectIndex = &objectIndex
				}
				issues = append(issues, issue)
			}
		}

		// Layer 3: Value map (lookup uses the normalized key, matching the
		// namespace of the mapping data)
		mappedValue := mapper.mapValue(lookupKey, value, valueMap, direction, qualifier, objectIndex, &issues)

		result[mappedKey] = mappedValue
	}

	return result, issues
}

func (mapper *attributeMapper) mapValue(key string, value any,
	valueMap map[string]map[string]string, direction MappingDirection,
	qualifier string, objectIndex int, issues *[]MappingIssue,
) any {
	if valueMap == nil {
		return value
	}

	keyValues, exists := valueMap[key]
	if !exists {
		return value
	}

	switch typedValue := value.(type) {
	case string:
		if mapped, found := keyValues[typedValue]; found {
			return mapped
		}
		issue := MappingIssue{
			Direction:      direction,
			Reason:         MappingUnknownValue,
			AttributeName:  key,
			AttributeValue: typedValue,
			Qualifier:      qualifier,
		}
		if objectIndex >= 0 {
			issue.ObjectIndex = &objectIndex
		}
		*issues = append(*issues, issue)
		return value

	case []any:
		mapped := make([]any, len(typedValue))
		for idx, item := range typedValue {
			if strItem, isStr := item.(string); isStr {
				if mappedItem, found := keyValues[strItem]; found {
					mapped[idx] = mappedItem
				} else {
					issue := MappingIssue{
						Direction:      direction,
						Reason:         MappingUnknownValue,
						AttributeName:  key,
						AttributeValue: strItem,
						Qualifier:      qualifier,
					}
					if objectIndex >= 0 {
						issue.ObjectIndex = &objectIndex
					}
					*issues = append(*issues, issue)
					mapped[idx] = item
				}
			} else {
				mapped[idx] = item
			}
		}
		return mapped

	default:
		return value
	}
}

// resolveResponseParameterMacros returns the command's defined response
// parameter macros if the caller requested "all", otherwise returns the
// caller's list as-is. The macros are additional parameter names to include
// in the response.
func (mapper *attributeMapper) resolveResponseParameterMacros(command, mqscQualifier string,
	responseParameters []string,
) []string {
	key := command + " " + mqscQualifier
	cmdMapping, exists := mapper.data.Commands[key]
	if !exists || len(cmdMapping.ResponseParameterMacros) == 0 {
		return responseParameters
	}

	// If "all" is requested, append the macro parameters
	for _, param := range responseParameters {
		if strings.EqualFold(param, "all") {
			return append(responseParameters, cmdMapping.ResponseParameterMacros...)
		}
	}

	return responseParameters
}

func mergeStringMap(target, source map[string]string) {
	for key, value := range source {
		target[key] = value
	}
}

func mergeNestedStringMap(target, source map[string]map[string]string) {
	for key, sourceInner := range source {
		if targetInner, exists := target[key]; exists {
			mergeStringMap(targetInner, sourceInner)
		} else {
			target[key] = sourceInner
		}
	}
}

func copyMap(source map[string]any) map[string]any {
	result := make(map[string]any, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}
