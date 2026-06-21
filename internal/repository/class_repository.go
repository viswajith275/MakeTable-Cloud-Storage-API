package repository

import (
	"context"
	"fmt"
	"strings"
	"time"
	"timetable_api/internal/dto"
	"timetable_api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateClass(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID, classRequest *dto.ClassCreationRequest) (*models.Class, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO classes (user_id, project_id, name, room_id, constraints) VALUES ($1, $2, $3, $4, $5::JSONB) RETURNING id, user_id, project_id, name, room_id, constraints, created_at, updated_at
	`

	var newClass models.Class

	err := pool.QueryRow(ctx, query, UserID, ProjectID, classRequest.Name, classRequest.RoomID, classRequest.Constraints).Scan(
		&newClass.ID,
		&newClass.UserID,
		&newClass.ProjectID,
		&newClass.Name,
		&newClass.RoomID,
		&newClass.Constraints,
		&newClass.CreatedAt,
		&newClass.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &newClass, err
}

func GetAllClassess(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID) ([]models.Class, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, name, room_id, constraints, created_at, updated_at FROM classes WHERE user_id = $1 AND project_id = $2 ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx, query, UserID, ProjectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class = []models.Class{}

	for rows.Next() {
		var class models.Class

		err = rows.Scan(
			&class.ID,
			&class.UserID,
			&class.ProjectID,
			&class.Name,
			&class.RoomID,
			&class.Constraints,
			&class.CreatedAt,
			&class.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		classes = append(classes, class)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return classes, err
}

func GetClassByID(pool *pgxpool.Pool, UserID uuid.UUID, ClassID uuid.UUID) (*models.Class, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, name, room_id, constraints, created_at, updated_at FROM classes WHERE user_id = $1 AND id = $2
	`
	var class models.Class
	err := pool.QueryRow(ctx, query, UserID, ClassID).Scan(
		&class.ID,
		&class.UserID,
		&class.ProjectID,
		&class.Name,
		&class.RoomID,
		&class.Constraints,
		&class.CreatedAt,
		&class.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("No class found")
		}
		return nil, err
	}

	return &class, err
}

func UpdateClass(pool *pgxpool.Pool, UserID uuid.UUID, ClassID uuid.UUID, classRequest *dto.ClassUpdationRequest) (*models.Class, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var setClauses []string
	args := pgx.NamedArgs{
		"user_id": UserID,
		"id":      ClassID,
	}

	if classRequest.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *classRequest.Name
	}

	if classRequest.RoomID != nil {
		setClauses = append(setClauses, "room_id = @room_id")
		args["room_id"] = *classRequest.RoomID
	}

	if classRequest.Constraints != nil {

		constraintsUpdate := make(map[string]any)

		// Add more constraints if added in the future

		if len(constraintsUpdate) > 0 {
			setClauses = append(setClauses, "constraints = COALESCE(constraints, '{}'::jsonb) || @constraints::jsonb")
			args["constraints"] = constraintsUpdate
		}

	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("Nothing to update!")
	}

	query := fmt.Sprintf(
		"UPDATE classes SET %s WHERE user_id = @user_id AND id = @id RETURNING id, user_id, project_id, name, room_id, constraints, created_at, updated_at",
		strings.Join(setClauses, ", "),
	)

	var class models.Class

	err := pool.QueryRow(ctx, query, args).Scan(
		&class.ID,
		&class.UserID,
		&class.ProjectID,
		&class.Name,
		&class.RoomID,
		&class.Constraints,
		&class.CreatedAt,
		&class.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &class, err
}

func DeleteClass(pool *pgxpool.Pool, UserID uuid.UUID, ClassID uuid.UUID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM classes WHERE user_id = $1 AND id = $2
	`
	commandTag, err := pool.Exec(ctx, query, UserID, ClassID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Class with id %v was not found", ClassID)
	}

	return err
}
