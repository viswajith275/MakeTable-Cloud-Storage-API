package models

import (
	"time"
	"timetable_api/internal/types"

	"github.com/google/uuid"
)

type Version struct {
	ID         uuid.UUID           `db:"id"`
	UserID     uuid.UUID           `db:"user_id"`
	ProjectID  uuid.UUID           `db:"project_id"`
	IsPublic   bool                `db:"is_public"`
	IsLive     bool                `db:"is_live"`
	Violations types.ViolationList `db:"violations"`
	CreatedAt  time.Time           `db:"created_at"`
	UpdatedAt  time.Time           `db:"updated_at"`
}
