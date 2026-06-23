package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type TeacherConstraints struct {
	MaxPerDay      *int `json:"max_per_day"  binding:"omitempty,gte=1,lte=100"`
	MaxPerWeek     *int `json:"max_per_week"  binding:"omitempty,gte=1,lte=100"`
	MaxConsecutive *int `json:"max_consecutive"  binding:"omitempty,gte=1,lte=100"`
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

func (c *TeacherConstraints) ValidateForPost() error {
	if c == nil || c.MaxPerDay == nil || c.MaxPerWeek == nil || c.MaxConsecutive == nil {
		return fmt.Errorf("constraints cannot be none!")
	}
	return nil
}
