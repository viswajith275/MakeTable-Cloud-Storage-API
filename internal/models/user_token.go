package models

import (
	"time"

	"github.com/google/uuid"
)

type UserToken struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	IsRevoked bool      `db:"is_required"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
