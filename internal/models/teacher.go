package models

import (
	"time"
	"timetable_api/internal/types"

	"github.com/google/uuid"
)

type Teacher struct {
	ID          uuid.UUID                `db:"id"`
	UserID      uuid.UUID                `db:"user_id"`
	ProjectID   uuid.UUID                `db:"project_id"`
	Name        string                   `db:"name"`
	Constraints types.TeacherConstraints `db:"constraints"`
	CreatedAt   time.Time                `db:"created_at"`
	UpdatedAt   time.Time                `db:"updated_at"`
}
