package lib

import (
	"reflect"
)

type YamlDiffStructure struct {
	name     string
	diff     YamlDiffEntry
	children map[string]*YamlDiffStructure
}

func (y YamlDiffStructure) GetName() string {
	return y.name
}

func (y YamlDiffStructure) GetFullPath() []string {
	return y.diff.Path
}

func (y YamlDiffStructure) GetDiff() YamlDiffEntry {
	return y.diff
}

func (y YamlDiffStructure) GetChildren() map[string]*YamlDiffStructure {
	return y.children
}

func diffStructure(oldYaml map[string]interface{}, newYaml map[string]interface{}, currentPath []string) map[string]*YamlDiffStructure {
	structure := map[string]*YamlDiffStructure{}

	for key, newVal := range newYaml {
		diff := YamlDiffEntry{
			Path:     append(currentPath, key),
			ValueOld: nil,
			ValueNew: newVal,
		}
		if oldVal, ok := oldYaml[key]; ok {
			diff.ValueOld = oldVal
		}

		str := &YamlDiffStructure{
			name: key,
			diff: diff,
		}

		var diffMapOld map[string]interface{}
		if diff.ValueOld != nil && reflect.TypeOf(diff.ValueOld).Kind() == reflect.Map {
			diffMapOld = EnsureStringMap(diff.ValueOld)
		}
		var diffMapNew map[string]interface{}
		if diff.ValueNew != nil && reflect.TypeOf(diff.ValueNew).Kind() == reflect.Map {
			diffMapNew = EnsureStringMap(diff.ValueNew)
		}
		str.children = diffStructure(diffMapOld, diffMapNew, diff.Path)
		structure[key] = str
	}

	// look for removed nodes
	for key, oldVal := range oldYaml {
		if _, ok := newYaml[key]; !ok {
			str := &YamlDiffStructure{
				name: key,
				diff: YamlDiffEntry{
					Path:     append(currentPath, key),
					ValueOld: oldVal,
					ValueNew: nil,
				},
			}
			structure[key] = str
		}
	}

	return structure
}
