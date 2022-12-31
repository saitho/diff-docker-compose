package lib

import (
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
	Diffs     []YamlDiffEntry
	Structure map[string]*YamlDiffStructure
}

// HasChanged returns true when there differences for the node
func (r YamlDiffResult) HasChanged(path []string) bool {
	results := r.GetAll(path)
	return len(results) > 0
}

// Get is deprecated. Use GetAll instead.
// @deprecated
func (r YamlDiffResult) Get(path []string) []YamlDiffEntry {
	return r.GetAll(path)
}

// Get returns the differences for the node identified by the path
func (r YamlDiffResult) GetAll(path []string) []YamlDiffEntry {
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

// GetStructure returns the structure for the node identified by the path
// Access the diff for that node or its children for diff information.
func (r YamlDiffResult) GetStructure(path []string) *YamlDiffStructure {
	currentMap := r.Structure
	pathLevel := 0
	running := true
	for running {
		running = false
		for _, structure := range currentMap {
			if len(path) <= pathLevel || len(structure.diff.Path) <= pathLevel {
				break
			}
			if path[pathLevel] != structure.diff.Path[pathLevel] {
				continue
			}
			pathLevel++
			if pathLevel == len(path) {
				return structure
			}
			currentMap = structure.children
			running = true // keep running with next map level
		}
	}
	return nil
}

// DiffYaml compares two map[string]interface{} objects (e.g. read via YAML parser).
func DiffYaml(oldYaml map[string]interface{}, newYaml map[string]interface{}) YamlDiffResult {
	result := YamlDiffResult{}
	result.Diffs = diffYaml(oldYaml, newYaml, []string{})
	result.Structure = diffStructure(oldYaml, newYaml, []string{})
	return result
}

func diffYaml(oldYaml map[string]interface{}, newYaml map[string]interface{}, currentPath []string) []YamlDiffEntry {
	var diffs []YamlDiffEntry

	for key, newVal := range newYaml {
		diff := YamlDiffEntry{
			Path:     append(currentPath, key),
			ValueNew: newVal,
		}
		if oldVal, ok := oldYaml[key]; ok {
			// key found in oldYaml -> node maybe changed
			if reflect.DeepEqual(newVal, oldVal) {
				continue // value did not change
			}
			if reflect.TypeOf(oldVal).Kind() == reflect.Map && reflect.TypeOf(oldVal).Kind() == reflect.Map {
				diffs = append(diffs, diffYaml(EnsureStringMap(oldVal), EnsureStringMap(newVal), append(currentPath, key))...)
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
			diff := YamlDiffEntry{
				Path:     append(currentPath, key),
				ValueOld: oldVal,
				ValueNew: nil,
			}
			diffs = append(diffs, diff)
		}
	}

	return diffs
}
