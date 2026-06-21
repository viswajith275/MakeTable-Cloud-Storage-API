package dto

import (
	"time"
	"timetable_api/internal/types"

	"github.com/google/uuid"
)

type ClassCreationRequest struct {
	Name        string                  `json:"name" binding:"required"`
	RoomID      uuid.UUID               `json:"room_id" binding:"required"`
	Constraints *types.ClassConstraints `json:"constraints" binding:"required"`
}

type ClassUpdationRequest struct {
	Name        *string                 `json:"name" binding:"omitempty"`
	RoomID      *uuid.UUID              `json:"room_id" binding:"omitempty"`
	Constraints *types.ClassConstraints `json:"constraints" binding:"omitempty"`
}

type ClassResponse struct {
	ID          uuid.UUID              `json:"id"`
	ProjectID   uuid.UUID              `json:"project_id"`
	Name        string                 `json:"name"`
	RoomID      uuid.UUID              `json:"room_id"`
	Constraints types.ClassConstraints `json:"constraints"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}
