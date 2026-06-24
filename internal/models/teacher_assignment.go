package models

import (
	"time"
	"timetable_api/internal/types"

	"github.com/google/uuid"
)

type TeacherAssignment struct {
	ID           uuid.UUID                          `db:"id"`
	UserID       uuid.UUID                          `db:"user_id"`
	ProjectID    uuid.UUID                          `db:"project_id"`
	ClassID      uuid.UUID                          `db:"class_id"`
	TeacherID    uuid.UUID                          `db:"teacher_id"`
	SubjectID    uuid.UUID                          `db:"subject_id"`
	TargetRoomID *uuid.UUID                         `db:"target_room_id"`
	Constraints  types.TeacherAssignmentConstraints `db:"constraints"`
	CreatedAt    time.Time                          `db:"created_at"`
	UpdatedAt    time.Time                          `db:"updated_at"`
}
