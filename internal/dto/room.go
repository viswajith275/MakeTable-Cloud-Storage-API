package dto

import (
	"time"
	"timetable_api/internal/types"

	"github.com/google/uuid"
)

type RoomCreationRequest struct {
	Name        string                 `json:"name" binding:"required"`
	IsLab       *bool                  `json:"is_lab" binding:"required"`
	Constraints *types.RoomConstraints `json:"constraints" binding:"required"`
}

type RoomUpdationRequest struct {
	Name        *string                `json:"name" binding:"omitempty"`
	IsLab       *bool                  `json:"is_lab" binding:"omitempty"`
	Constraints *types.RoomConstraints `json:"constraints" binding:"omitempty"`
}

type RoomResponse struct {
	ID          uuid.UUID             `json:"id"`
	ProjectID   uuid.UUID             `json:"project_id"`
	Name        string                `json:"name"`
	IsLab       bool                  `json:"is_lab"`
	Constraints types.RoomConstraints `json:"constraints"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}
