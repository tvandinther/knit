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
