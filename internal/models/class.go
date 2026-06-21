package models

import (
	"time"
	"timetable_api/internal/types"

	"github.com/google/uuid"
)

type Class struct {
	ID          uuid.UUID              `db:"id"`
	UserID      uuid.UUID              `db:"user_id"`
	ProjectID   uuid.UUID              `db:"project_id"`
	Name        string                 `db:"name"`
	RoomID      uuid.UUID              `db:"room_id"`
	Constraints types.ClassConstraints `db:"constraints"`
	CreatedAt   time.Time              `db:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
}
