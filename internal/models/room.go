package models

import (
	"time"
	"timetable_api/internal/types"

	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID             `db:"id"`
	UserID      uuid.UUID             `db:"user_id"`
	ProjectID   uuid.UUID             `db:"project_id"`
	Name        string                `db:"name"`
	IsLab       bool                  `db:"is_lab"`
	Constraints types.RoomConstraints `db:"constraints"`
	CreatedAt   time.Time             `db:"created_at"`
	UpdatedAt   time.Time             `db:"updated_at"`
}
