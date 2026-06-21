package repository

import (
	"context"
	"fmt"
	"strings"
	"time"
	"timetable_api/internal/dto"
	"timetable_api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, userRequest *dto.UserRegisterRequest) (*models.User, error) {

	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var query string = `
		INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3) RETURNING *
	`
	var user models.User

	var err error = pool.QueryRow(ctx, query, userRequest.Email, userRequest.UserName, userRequest.Password).Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.PasswordHash,
		&user.Disabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "SQLSTATE 23505") || strings.Contains(errMsg, "unique constraint") {
			return nil, fmt.Errorf("Email already exists!")
		}
		return nil, err
	}

	return &user, err
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {

	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var query string = `
		SELECT * FROM users WHERE email = $1
	`
	var user models.User

	var err error = pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.PasswordHash,
		&user.Disabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, err
}

func GetUserById(pool *pgxpool.Pool, id uuid.UUID) (*models.User, error) {

	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var query string = `
		SELECT * FROM users WHERE id = $1
	`
	var user models.User

	var err error = pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.PasswordHash,
		&user.Disabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, err
}
