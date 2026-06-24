package dto

import (
	"time"
	"timetable_api/internal/types"

	"github.com/google/uuid"
)

type TeacherAssignmentCreationRequest struct {
	ProjectID    uuid.UUID                           `json:"-"`
	ClassID      uuid.UUID                           `json:"class_id"`
	TeacherID    uuid.UUID                           `json:"teacher_id"`
	SubjectID    uuid.UUID                           `json:"subject_id"`
	TargetRoomID *uuid.UUID                          `json:"target_room_id" binding:"omitempty"`
	Constraints  *types.TeacherAssignmentConstraints `json:"constraints" binding:"required"`
}

type TeacherAssignmentUpdationRequest struct {
	ClassID      *uuid.UUID                          `json:"class_id" binding:"omitempty"`
	TeacherID    *uuid.UUID                          `json:"teacher_id" binding:"omitempty"`
	SubjectID    *uuid.UUID                          `json:"subject_id" binding:"omitempty"`
	TargetRoomID *types.NullableUUID                 `json:"target_room_id" binding:"omitempty"`
	Constraints  *types.TeacherAssignmentConstraints `json:"constraints" binding:"omitempty"`
}

type TeacherAssignmentResponse struct {
	ID           uuid.UUID                          `json:"id"`
	ProjectID    uuid.UUID                          `json:"project_id"`
	ClassID      uuid.UUID                          `json:"class_id"`
	TeacherID    uuid.UUID                          `json:"teacher_id"`
	SubjectID    uuid.UUID                          `json:"subject_id"`
	TargetRoomID *uuid.UUID                         `json:"target_room_id"`
	Constraints  types.TeacherAssignmentConstraints `json:"constraints"`
	CreatedAt    time.Time                          `json:"created_at"`
	UpdatedAt    time.Time                          `json:"updated_at"`
}
