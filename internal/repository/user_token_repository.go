package repository

import (
	"context"
	"fmt"
	"strings"
	"time"
	"timetable_api/internal/dto"
	"timetable_api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUserToken(pool *pgxpool.Pool, tokenRequest *dto.UserTokenCreateRequest) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO user_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING *
	`

	var newToken models.UserToken

	err := pool.QueryRow(ctx, query, tokenRequest.UserID, tokenRequest.TokenHash, tokenRequest.ExpiresAT).Scan(
		&newToken.ID,
		&newToken.UserID,
		&newToken.TokenHash,
		&newToken.IsRevoked,
		&newToken.ExpiresAt,
		&newToken.CreatedAt,
		&newToken.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return err
}

func GetUserTokenByTokenHash(pool *pgxpool.Pool, tokenHash string) (*models.UserToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT * FROM user_tokens WHERE token_hash = $1
	`

	var newToken models.UserToken

	err := pool.QueryRow(ctx, query, tokenHash).Scan(
		&newToken.ID,
		&newToken.UserID,
		&newToken.TokenHash,
		&newToken.IsRevoked,
		&newToken.ExpiresAt,
		&newToken.CreatedAt,
		&newToken.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &newToken, err

}

func UpdateUserToken(pool *pgxpool.Pool, tokenRequest *dto.UserTokenUpdateRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var setClauses []string
	args := pgx.NamedArgs{
		"token_hash": tokenRequest.TokenHash,
	}

	if tokenRequest.IsRevoked != nil {
		setClauses = append(setClauses, "is_revoked = @is_revoked")
		args["is_revoked"] = *tokenRequest.IsRevoked
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("Nothing to update!")
	}

	query := fmt.Sprintf(
		"UPDATE user_tokens SET %s WHERE token_hash = @token_hash RETURNING *",
		strings.Join(setClauses, ", "),
	)

	var token models.UserToken

	err := pool.QueryRow(ctx, query, args).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.IsRevoked,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return err

}
