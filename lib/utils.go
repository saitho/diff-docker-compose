package lib

import (
	"fmt"
	"reflect"
	"strings"
)

func EnsureStringMap(value any) map[string]interface{} {
	var valMap map[string]interface{}
	if reflect.TypeOf(value).String() == "map[string]interface {}" {
		return value.(map[string]interface{})
	}
	if strings.HasPrefix(reflect.TypeOf(value).String(), "map[interface {}]") {
		valMap = cleanUpInterfaceMap(value.(map[interface{}]interface{}))
	} else {
		panic("invalid value passed")
	}
	return valMap
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
