package lib

import (
	"fmt"
	"reflect"
)

type YamlDiffEntry struct {
	Path     []string
	ValueOld interface{}
	ValueNew interface{}
}

type DiffType string

const (
	Unknown DiffType = "unknown"
	Added            = "added"
	Removed          = "removed"
	Changed          = "changed"
)

func (e YamlDiffEntry) GetType() DiffType {
	if e.ValueOld == nil && e.ValueNew == nil {
		return Unknown
	}
	if e.ValueOld == nil && e.ValueNew != nil {
		return Added
	}
	if e.ValueOld != nil && e.ValueNew == nil {
		return Removed
	}
	return Changed
}

type YamlDiffResult struct {
	Diffs []YamlDiffEntry
}

func (r YamlDiffResult) HasChanged(path []string) bool {
	results := r.Get(path)
	return len(results) > 0
}

func (r YamlDiffResult) Get(path []string) []YamlDiffEntry {
	var entries []YamlDiffEntry
	for _, diff := range r.Diffs {
		exited := false
		for i, p := range path {
			if diff.Path[i] != p {
				exited = true
				break
			}
		}
		if !exited {
			entries = append(entries, diff)
		}
	}
	return entries
}

func DiffYaml(oldYaml map[string]interface{}, newYaml map[string]interface{}) YamlDiffResult {
	result := YamlDiffResult{}
	result.Diffs = diffYaml(oldYaml, newYaml, []string{})
	return result
}

func diffYaml(oldYaml map[string]interface{}, newYaml map[string]interface{}, currentPath []string) []YamlDiffEntry {
	var diffs []YamlDiffEntry

	for key, newVal := range newYaml {
		diff := YamlDiffEntry{}
		diff.Path = append(currentPath, key)
		diff.ValueNew = newVal
		if oldVal, ok := oldYaml[key]; ok {
			// key found in oldYaml -> node maybe changed
			if reflect.DeepEqual(newVal, oldVal) {
				continue // value did not change
			}
			if reflect.TypeOf(oldVal).Kind() == reflect.Map && reflect.TypeOf(oldVal).Kind() == reflect.Map {
				oldValMap := oldVal
				if reflect.TypeOf(oldVal).String() == "map[interface {}]interface {}" {
					oldValMap = cleanUpInterfaceMap(oldVal.(map[interface{}]interface{}))
				}
				newValMap := newVal
				if reflect.TypeOf(newVal).String() == "map[interface {}]interface {}" {
					newValMap = cleanUpInterfaceMap(newVal.(map[interface{}]interface{}))
				}

				diffs = append(diffs, diffYaml(oldValMap.(map[string]interface{}), newValMap.(map[string]interface{}), append(currentPath, key))...)
			} else {
				diff.ValueOld = oldVal
				diffs = append(diffs, diff)
			}
		} else {
			// key not in oldYaml -> new node
			diff.ValueOld = nil
			diffs = append(diffs, diff)
		}
	}

	// look for removed nodes
	for key, oldVal := range oldYaml {
		if _, ok := newYaml[key]; !ok {
			diff := YamlDiffEntry{}
			diff.Path = append(currentPath, key)
			diff.ValueOld = oldVal
			diff.ValueNew = nil
			diffs = append(diffs, diff)
		}
	}

	return diffs
}

// copied from https://github.com/elastic/beats/blob/6435194af9f42cbf778ca0a1a92276caf41a0da8/libbeat/common/mapstr.go
func cleanUpInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = cleanUpMapValue(v)
	}
	return result
}

func cleanUpMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanUpInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanUpInterfaceMap(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
func cleanUpInterfaceArray(in []interface{}) []interface{} {
	result := make([]interface{}, len(in))
	for i, v := range in {
		result[i] = cleanUpMapValue(v)
	}
	return result
}
