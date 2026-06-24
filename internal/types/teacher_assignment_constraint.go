package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type TeacherAssignmentConstraints struct {
	IsClassTeacher *bool     `json:"is_class_teacher" binding:"omitempty"`
	FirstSlotDays  *[]string `json:"first_slot_days" binding:"omitempty,min=1,unique,dive,oneof=Mon Tue Wed Thu Fri Sat Sun"`
}

func (c TeacherAssignmentConstraints) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *TeacherAssignmentConstraints) Scan(value any) error {
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

func (c *TeacherAssignmentConstraints) ValidateForPost() error {
	if c == nil {
		return errors.New("constraints is required")
	}

	if c.IsClassTeacher == nil {
		return errors.New("field 'is_class_teacher' is required for creation")
	}

	if *c.IsClassTeacher && (c.FirstSlotDays == nil || len(*c.FirstSlotDays) == 0) {
		return errors.New("field 'first_slot_days' cannot be empty when 'is_class_teacher' is true")
	}

	return c.Validate()
}

func (c *TeacherAssignmentConstraints) Validate() error {
	if c == nil {
		return nil
	}

	if c.IsClassTeacher != nil && !*c.IsClassTeacher {
		if c.FirstSlotDays != nil && len(*c.FirstSlotDays) > 0 {
			return errors.New("'first_slot_days' can only be specified when 'is_class_teacher' is true")
		}
	}

	return nil
}
