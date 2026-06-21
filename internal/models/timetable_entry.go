package models

import (
	"time"

	"github.com/google/uuid"
)

type TimeTableEntry struct {
	ID           uuid.UUID  `db:"id"`
	VersionID    uuid.UUID  `db:"version_id"`
	Slot         int        `db:"slot"`
	Day          string     `db:"day"`
	AssignmentID *uuid.UUID `db:"assignment_id"`
	// Add snapshots

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
