package utils

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// Returns the map that recursively matches the given keys
// Returns an empty map if the key does not exist
func GetSubMap(
	inputMap interface{},
	keys ...string,
) map[string]interface{} {
	subMap := GetSubMapOrNil(inputMap, keys...)
	if subMap == nil {
		return make(map[string]interface{})
	}

	return subMap
}

// Returns the map that recursively matches the given keys
// Returns nil if the key does not exist
func GetSubMapOrNil(
	inputMap interface{},
	keys ...string,
) map[string]interface{} {
	if inputMap == nil {
		return nil
	}

	if _, ok := inputMap.(map[string]interface{}); !ok {
		return nil
	}

	currentMap := inputMap.(map[string]interface{})
	for _, key := range keys {
		value, ok := currentMap[key]
		if ok {
			currentMap = value.(map[string]interface{})
		} else {
			return nil
		}
	}

	return currentMap
}

// Returns the array that recursively matches the given keys
// Returns an empty array if the key does not exist
func GetSubArray(
	inputMap interface{},
	keys ...string,
) []map[string]interface{} {
	if inputMap == nil {
		return []map[string]interface{}{}
	}

	var currentMap interface{} = inputMap
	for _, key := range keys {
		value, ok := currentMap.(map[string]interface{})[key]
		if ok {
			currentMap = value
		} else {
			return []map[string]interface{}{}
		}
	}

	switch value := currentMap.(type) {
	case []map[string]interface{}:
		return value
	case []interface{}:
		arrayAsMaps := make([]map[string]interface{}, len(value))
		for i, v := range value {
			arrayAsMaps[i] = v.(map[string]interface{})
		}

		return arrayAsMaps
	default:
		return []map[string]interface{}{}
	}
}

// Returns the integer that recursively matches the given keys
// Returns nil if the key does not exist or it not an integer
func GetSubInteger(
	inputMap interface{},
	keys ...string,
) *int {
	if inputMap == nil {
		return nil
	}

	if _, ok := inputMap.(map[string]interface{}); !ok {
		return nil
	}

	currentMap := inputMap.(map[string]interface{})
	for i := 0; i < len(keys); i++ {
		value, ok := currentMap[keys[i]]
		if !ok {
			return nil
		}

		if i == len(keys)-1 {
			switch value := value.(type) {
			case int:
				intValue := value
				return &intValue
			case float64:
				intValue := int(value)
				return &intValue
			case float32:
				intValue := int(value)
				return &intValue
			default:
				return nil
			}
		} else if currentMap, ok = value.(map[string]interface{}); !ok {
			return nil
		}
	}

	return nil
}

// Returns the same map, without the non-scalar values
func GetOnlyScalar(
	inputMap map[string]interface{},
) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range inputMap {
		if v == nil {
			newMap[k] = v
			continue
		}

		switch reflect.ValueOf(v).Type().Kind() {
		case
			reflect.Invalid,
			reflect.Array,
			reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Slice,
			reflect.Struct:
			continue
		default:
			newMap[k] = v
		}
	}

	return newMap
}

// Assign map values to a struct using mapstructure
// This requires the mapstructure tags on the target entity
// This one returns errors for unmapped properties
// and assigns the nil values
func MapToStruct(
	input map[string]interface{},
	output interface{},
) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		ZeroFields:  true,
		Result:      output,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func StructToMap(
	input interface{},
	output interface{},
) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		ZeroFields:  true,
		Result:      output,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

// Returns the array that recursively matches the given keys
// Returns nil if the key does not exist
func GetSubArrayOrNil(
	inputMap interface{},
	keys ...string,
) []map[string]interface{} {
	if inputMap == nil {
		return nil
	}

	var currentMap interface{} = inputMap
	for _, key := range keys {
		value, ok := currentMap.(map[string]interface{})[key]
		if ok {
			currentMap = value
		} else {
			return nil
		}
	}

	if currentMap == nil {
		return nil
	}

	switch value := currentMap.(type) {
	case []map[string]interface{}:
		return value
	case []interface{}:
		arrayAsMaps := make([]map[string]interface{}, len(value))
		for i, v := range value {
			arrayAsMaps[i] = v.(map[string]interface{})
		}

		return arrayAsMaps
	default:
		return []map[string]interface{}{}
	}
}

func GetOnlyScalarWithouNilValue(
	inputMap map[string]interface{},
) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range inputMap {
		if v == nil {
			continue
		}

		switch reflect.ValueOf(v).Type().Kind() {
		case
			reflect.Invalid,
			reflect.Array,
			reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Slice,
			reflect.Struct:
			continue
		default:
			newMap[k] = v
		}
	}

	return newMap
}
