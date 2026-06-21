package dto

import (
	"time"
	"timetable_api/internal/types"

	"github.com/google/uuid"
)

type SubjectCreationRequest struct {
	Name        string                    `json:"name" binding:"required"`
	Constraints *types.SubjectConstraints `json:"constraints" binding:"required"`
}

type SubjectUpdationRequest struct {
	Name        *string                   `json:"name" binding:"omitempty"`
	Constraints *types.SubjectConstraints `json:"constraints" binding:"omitempty"`
}

type SubjectResponse struct {
	ID          uuid.UUID                `json:"id"`
	ProjectID   uuid.UUID                `json:"project_id"`
	Name        string                   `json:"name"`
	Constraints types.SubjectConstraints `json:"constraints"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}
