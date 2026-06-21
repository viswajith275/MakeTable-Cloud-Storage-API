package dto

import (
	"time"
	"timetable_api/internal/types"

	"github.com/google/uuid"
)

type TeacherCreationRequest struct {
	Name        string                    `json:"name" binding:"required"`
	Constraints *types.TeacherConstraints `json:"constraints" binding:"required"`
}

type TeacherUpdationRequest struct {
	Name        *string                   `json:"name" binding:"omitempty"`
	Constraints *types.TeacherConstraints `json:"constraints" binding:"omitempty"`
}

type TeacherResponse struct {
	ID          uuid.UUID                `json:"id"`
	ProjectID   uuid.UUID                `json:"project_id"`
	Name        string                   `json:"name"`
	Constraints types.TeacherConstraints `json:"constraints"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}
