package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserTokenCreateRequest struct {
	UserID    uuid.UUID
	TokenHash string
	ExpiresAT time.Time
}

type UserTokenUpdateRequest struct {
	TokenHash string
	IsRevoked *bool
}
