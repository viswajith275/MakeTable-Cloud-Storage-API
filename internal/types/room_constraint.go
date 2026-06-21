package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type RoomConstraints struct {
	Capacity *int `json:"capacity" binding:"omitempty,gte=1,lte=100"`
}

func (c *RoomConstraints) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *RoomConstraints) Scan(value interface{}) error {

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

func (c *RoomConstraints) ValidateForPost() error {
	if c == nil || c.Capacity == nil {
		return fmt.Errorf("Capacity cannot be none!")
	}
	return nil
}
