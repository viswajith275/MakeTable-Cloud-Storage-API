package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID                 uuid.UUID `db:"id"`
	UserID             uuid.UUID `db:"user_id"`
	Name               string    `db:"name"`
	Slots              int       `db:"slots"`
	Days               []string  `db:"days"`
	RoomsVersion       int       `db:"rooms_version"`
	ClassesVersion     int       `db:"classes_version"`
	TeachersVersion    int       `db:"teachers_version"`
	SubjectsVersion    int       `db:"subjects_version"`
	AssignmentsVersion int       `db:"assignments_version"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}
