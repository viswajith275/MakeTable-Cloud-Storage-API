package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"timetable_api/internal/dto"
	"timetable_api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateRoom(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID, roomRequest *dto.RoomCreationRequest) (*models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("DEBUG 1 [Input]: %+v\n", roomRequest.Constraints)

	// THE FIX: Use an interface{} so we can pass either nil or a string
	var constraintsParam interface{}
	if roomRequest.Constraints != nil {
		b, err := json.Marshal(roomRequest.Constraints)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal constraints: %w", err)
		}

		// DEBUG 2: Did json.Marshal work properly?
		fmt.Printf("DEBUG 2 [Marshaled]: %s\n", string(b))

		// Force it to a string so pgx sends it as pure text, not binary bytea
		constraintsParam = string(b)
	}

	var query string = `
		INSERT INTO rooms (user_id, project_id, name, is_lab, constraints) VALUES ($1, $2, $3, $4, $5::JSONB) RETURNING id, user_id, project_id, name, is_lab, constraints, created_at, updated_at
	`

	var newRoom models.Room

	err := pool.QueryRow(ctx, query, UserID, ProjectID, roomRequest.Name, roomRequest.IsLab, constraintsParam).Scan(
		&newRoom.ID,
		&newRoom.UserID,
		&newRoom.ProjectID,
		&newRoom.Name,
		&newRoom.IsLab,
		&newRoom.Constraints,
		&newRoom.CreatedAt,
		&newRoom.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	fmt.Printf("DEBUG 3 [Scanned Output]: %+v\n", newRoom.Constraints)
	return &newRoom, err
}

func GetAllRooms(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID) ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, name, is_lab, constraints, created_at, updated_at FROM rooms WHERE user_id = $1 AND project_id = $2 ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx, query, UserID, ProjectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room = []models.Room{}

	for rows.Next() {
		var room models.Room

		err = rows.Scan(
			&room.ID,
			&room.UserID,
			&room.ProjectID,
			&room.Name,
			&room.IsLab,
			&room.Constraints,
			&room.CreatedAt,
			&room.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rooms, err
}

func GetRoomByID(pool *pgxpool.Pool, UserID uuid.UUID, RoomID uuid.UUID) (*models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, name, is_lab, constraints, created_at, updated_at FROM rooms WHERE user_id = $1 AND id = $2
	`
	var room models.Room
	err := pool.QueryRow(ctx, query, UserID, RoomID).Scan(
		&room.ID,
		&room.UserID,
		&room.ProjectID,
		&room.Name,
		&room.IsLab,
		&room.Constraints,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("No room found")
		}
		return nil, err
	}

	return &room, err
}

func UpdateRoom(pool *pgxpool.Pool, UserID uuid.UUID, RoomID uuid.UUID, roomRequest *dto.RoomUpdationRequest) (*models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var setClauses []string
	args := pgx.NamedArgs{
		"user_id": UserID,
		"id":      RoomID,
	}

	if roomRequest.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *roomRequest.Name
	}

	if roomRequest.IsLab != nil {
		setClauses = append(setClauses, "is_lab = @is_lab")
		args["is_lab"] = *roomRequest.IsLab
	}

	if roomRequest.Constraints != nil {
		constraintsUpdate := make(map[string]any)

		if roomRequest.Constraints.Capacity != nil {
			constraintsUpdate["capacity"] = *roomRequest.Constraints.Capacity
		}

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
		"UPDATE rooms SET %s WHERE user_id = @user_id AND id = @id RETURNING id, user_id, project_id, name, is_lab, constraints, created_at, updated_at",
		strings.Join(setClauses, ", "),
	)

	var room models.Room

	err := pool.QueryRow(ctx, query, args).Scan(
		&room.ID,
		&room.UserID,
		&room.ProjectID,
		&room.Name,
		&room.IsLab,
		&room.Constraints,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &room, err
}

func DeleteRoom(pool *pgxpool.Pool, UserID uuid.UUID, RoomID uuid.UUID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM rooms WHERE user_id = $1 AND id = $2
	`
	commandTag, err := pool.Exec(ctx, query, UserID, RoomID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Room with id %v was not found", RoomID)
	}

	return err
}
