package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Violation struct {
}

type ViolationList []Violation

func (v ViolationList) Value() (driver.Value, error) {
	if v == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(v)
}

func (v *ViolationList) Scan(value any) error {
	if value == nil {
		return nil
	}

	var bytes []byte

	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSON value: %v", value)
	}

	if string(bytes) == "null" {
		return nil
	}

	return json.Unmarshal(bytes, v)
}

func (v *ViolationList) Validate() error {
	return nil
}
