package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProjectCreationRequest struct {
	Name  string   `json:"name" binding:"required,min=3"`
	Slots int      `json:"slots" binding:"required,gt=0"`
	Days  []string `json:"days" binding:"required,min=1,unique,dive,oneof=Mon Tue Wed Thu Fri Sat Sun"`
}

type ProjectUpdationRequest struct {
	Name  *string   `json:"name" binding:"omitempty,min=1"`
	Slots *int      `json:"slots" binding:"omitempty,gt=0"`
	Days  *[]string `json:"days" binding:"omitempty,min=1,unique,dive,oneof=Mon Tue Wed Thu Fri Sat Sun"`
}

type ProjectResponse struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Slots              int       `json:"slots"`
	Days               []string  `json:"days"`
	RoomsVersion       int       `json:"rooms_version"`
	ClassesVersion     int       `json:"classes_version"`
	TeachersVersion    int       `json:"teachers_version"`
	SubjectsVersion    int       `json:"subjects_version"`
	AssignmentsVersion int       `json:"assignments_version"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
