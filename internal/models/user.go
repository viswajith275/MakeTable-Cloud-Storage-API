package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	UserName     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Disabled     bool      `db:"disabled"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
