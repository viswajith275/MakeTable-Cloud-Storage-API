package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type NullableUUID struct {
	Value   uuid.UUID
	Present bool
	Valid   bool
}

func (n *NullableUUID) UnmarshalJSON(data []byte) error {
	n.Present = true
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(data, &n.Value)
}
