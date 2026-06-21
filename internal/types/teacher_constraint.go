package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type TeacherConstraints struct {
}

func (c TeacherConstraints) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *TeacherConstraints) Scan(value interface{}) error {
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

	return json.Unmarshal(bytes, c)
}

func (c *TeacherConstraints) Validate() error {
	return nil
}
