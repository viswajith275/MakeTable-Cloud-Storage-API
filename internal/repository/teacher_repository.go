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

func CreateTeacher(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID, teacherRequest *dto.TeacherCreationRequest) (*models.Teacher, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO teachers (user_id, project_id, name, constraints) VALUES ($1, $2, $3, $4::JSONB) RETURNING id, user_id, project_id, name, constraints, created_at, updated_at
	`

	var newTeacher models.Teacher

	err := pool.QueryRow(ctx, query, UserID, ProjectID, teacherRequest.Name, teacherRequest.Constraints).Scan(
		&newTeacher.ID,
		&newTeacher.UserID,
		&newTeacher.ProjectID,
		&newTeacher.Name,
		&newTeacher.Constraints,
		&newTeacher.CreatedAt,
		&newTeacher.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &newTeacher, err
}

func GetAllTeachers(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID) ([]models.Teacher, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, name, constraints, created_at, updated_at FROM teachers WHERE user_id = $1 AND project_id = $2 ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx, query, UserID, ProjectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher = []models.Teacher{}

	for rows.Next() {
		var teacher models.Teacher

		err = rows.Scan(
			&teacher.ID,
			&teacher.UserID,
			&teacher.ProjectID,
			&teacher.Name,
			&teacher.Constraints,
			&teacher.CreatedAt,
			&teacher.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		teachers = append(teachers, teacher)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teachers, err
}

func GetTeacherByID(pool *pgxpool.Pool, UserID uuid.UUID, TeacherID uuid.UUID) (*models.Teacher, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, name, constraints, created_at, updated_at FROM teachers WHERE user_id = $1 AND id = $2
	`
	var teacher models.Teacher
	err := pool.QueryRow(ctx, query, UserID, TeacherID).Scan(
		&teacher.ID,
		&teacher.UserID,
		&teacher.ProjectID,
		&teacher.Name,
		&teacher.Constraints,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("No teacher found")
		}
		return nil, err
	}

	return &teacher, err
}

func UpdateTeacher(pool *pgxpool.Pool, UserID uuid.UUID, TeacherID uuid.UUID, teacherRequest *dto.TeacherUpdationRequest) (*models.Teacher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var setClauses []string
	args := pgx.NamedArgs{
		"user_id": UserID,
		"id":      TeacherID,
	}

	if teacherRequest.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *teacherRequest.Name
	}

	if teacherRequest.Constraints != nil {

		constraintsUpdate := make(map[string]any)

		if teacherRequest.Constraints.MaxPerDay != nil {
			constraintsUpdate["max_per_day"] = *teacherRequest.Constraints.MaxPerDay
		}

		if teacherRequest.Constraints.MaxPerWeek != nil {
			constraintsUpdate["max_per_week"] = *teacherRequest.Constraints.MaxPerWeek
		}

		if teacherRequest.Constraints.MaxConsecutive != nil {
			constraintsUpdate["max_consecutive"] = *teacherRequest.Constraints.MaxConsecutive
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
		"UPDATE teachers SET %s WHERE user_id = @user_id AND id = @id RETURNING id, user_id, project_id, name, constraints, created_at, updated_at",
		strings.Join(setClauses, ", "),
	)

	var teacher models.Teacher

	err := pool.QueryRow(ctx, query, args).Scan(
		&teacher.ID,
		&teacher.UserID,
		&teacher.ProjectID,
		&teacher.Name,
		&teacher.Constraints,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &teacher, err
}

func DeleteTeacher(pool *pgxpool.Pool, UserID uuid.UUID, TeacherID uuid.UUID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM teachers WHERE user_id = $1 AND id = $2
	`
	commandTag, err := pool.Exec(ctx, query, UserID, TeacherID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Teacher with id %v was not found", TeacherID)
	}

	return err
}
