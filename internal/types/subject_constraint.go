package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type SubjectConstraints struct {
	MorningTendency *string `json:"morning_tendency" binding:"omitempty,oneof=Low Med High"`
	MinPerDay       *int    `json:"min_per_day"  binding:"omitempty,gte=1,lte=100"`
	MaxPerDay       *int    `json:"max_per_day"  binding:"omitempty,gte=1,lte=100"`
	MinPerWeek      *int    `json:"min_per_week"  binding:"omitempty,gte=1,lte=100"`
	MaxPerWeek      *int    `json:"max_per_week"  binding:"omitempty,gte=1,lte=100"`
	MinConsecutive  *int    `json:"min_consecutive"  binding:"omitempty,gte=1,lte=100"`
	MaxConsecutive  *int    `json:"max_consecutive"  binding:"omitempty,gte=1,lte=100"`
}

func (c SubjectConstraints) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *SubjectConstraints) Scan(value interface{}) error {
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

func (c *SubjectConstraints) ValidateForPost() error {
	if c == nil || c.MorningTendency == nil || c.MinPerDay == nil || c.MinPerWeek == nil || c.MinConsecutive == nil || c.MaxPerDay == nil || c.MaxPerWeek == nil || c.MaxConsecutive == nil {
		return fmt.Errorf("constraints cannot be none!")
	}
	return nil
}
