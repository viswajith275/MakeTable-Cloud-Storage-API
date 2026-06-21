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

func CreateProject(pool *pgxpool.Pool, UserID uuid.UUID, projectRequest *dto.ProjectCreationRequest) (*models.Project, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO projects (user_id, name, slots, days) VALUES ($1, $2, $3, $4::week_days[]) RETURNING id, user_id, name, slots, days, rooms_version, classes_version, teachers_version, subjects_version, assignments_version, created_at, updated_at
	`

	var newProject models.Project

	err := pool.QueryRow(ctx, query, UserID, projectRequest.Name, projectRequest.Slots, projectRequest.Days).Scan(
		&newProject.ID,
		&newProject.UserID,
		&newProject.Name,
		&newProject.Slots,
		&newProject.Days,
		&newProject.RoomsVersion,
		&newProject.ClassesVersion,
		&newProject.TeachersVersion,
		&newProject.SubjectsVersion,
		&newProject.AssignmentsVersion,
		&newProject.CreatedAt,
		&newProject.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &newProject, err
}

func GetAllProject(pool *pgxpool.Pool, UserID uuid.UUID) ([]models.Project, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, name, slots, days, rooms_version, classes_version, teachers_version, subjects_version, assignments_version, created_at, updated_at FROM projects WHERE user_id = $1 ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx, query, UserID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project = []models.Project{}

	for rows.Next() {
		var project models.Project

		err = rows.Scan(
			&project.ID,
			&project.UserID,
			&project.Name,
			&project.Slots,
			&project.Days,
			&project.RoomsVersion,
			&project.ClassesVersion,
			&project.TeachersVersion,
			&project.SubjectsVersion,
			&project.AssignmentsVersion,
			&project.CreatedAt,
			&project.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		projects = append(projects, project)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, err
}

func GetProjectByID(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID) (*models.Project, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, name, slots, days, rooms_version, classes_version, teachers_version, subjects_version, assignments_version, created_at, updated_at FROM projects WHERE user_id = $1 AND id = $2
	`
	var project models.Project
	err := pool.QueryRow(ctx, query, UserID, ProjectID).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.Slots,
		&project.Days,
		&project.RoomsVersion,
		&project.ClassesVersion,
		&project.TeachersVersion,
		&project.SubjectsVersion,
		&project.AssignmentsVersion,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &project, err
}

func UpdateProject(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID, projectRequest *dto.ProjectUpdationRequest) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var setClauses []string
	args := pgx.NamedArgs{
		"user_id": UserID,
		"id":      ProjectID,
	}

	if projectRequest.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *projectRequest.Name
	}

	if projectRequest.Slots != nil {
		setClauses = append(setClauses, "slots = @slots")
		args["slots"] = *projectRequest.Slots
	}

	if projectRequest.Days != nil {
		setClauses = append(setClauses, "days = @days::week_days[]")
		args["days"] = *projectRequest.Days
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("Nothing to update!")
	}

	query := fmt.Sprintf(
		"UPDATE projects SET %s WHERE user_id = @user_id AND id = @id RETURNING id, user_id, name, slots, days, rooms_version, classes_version, teachers_version, subjects_version, assignments_version, created_at, updated_at",
		strings.Join(setClauses, ", "),
	)

	var project models.Project

	err := pool.QueryRow(ctx, query, args).Scan(
		&project.ID,
		&project.UserID,
		&project.Name,
		&project.Slots,
		&project.Days,
		&project.RoomsVersion,
		&project.ClassesVersion,
		&project.TeachersVersion,
		&project.SubjectsVersion,
		&project.AssignmentsVersion,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &project, err
}

func DeleteProject(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM projects WHERE user_id = $1 AND id = $2
	`
	commandTag, err := pool.Exec(ctx, query, UserID, ProjectID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Project with id %v was not found", ProjectID)
	}

	return err
}
