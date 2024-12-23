package util

import (
	"encoding/json"
	"fmt"
)

func MapToStruct(m map[string]interface{}, s interface{}) error {
	jsonbody, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("invalid argument: %+v, %w", m, err)
	}

	if err := json.Unmarshal(jsonbody, &s); err != nil {
		return fmt.Errorf("invalid argument: %+v, %w", m, err)
	}

	return nil
}

func AnySliceToTyped[T any](input any) ([]T, error) {
	switch x := input.(type) {
	case []any:
		var result []T
		for i, v := range x {
			elem, ok := v.(T)
			if !ok {
				return nil, fmt.Errorf("element %v at index %d is not a string", v, i)
			}
			result = append(result, elem)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("expecting slice")
	}
}
