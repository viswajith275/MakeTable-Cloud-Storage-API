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

func CreateSubject(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID, subjectRequest *dto.SubjectCreationRequest) (*models.Subject, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO subjects (user_id, project_id, name, constraints) VALUES ($1, $2, $3, $4::JSONB) RETURNING id, user_id, project_id, name, constraints, created_at, updated_at
	`

	var newSubject models.Subject

	err := pool.QueryRow(ctx, query, UserID, ProjectID, subjectRequest.Name, subjectRequest.Constraints).Scan(
		&newSubject.ID,
		&newSubject.UserID,
		&newSubject.ProjectID,
		&newSubject.Name,
		&newSubject.Constraints,
		&newSubject.CreatedAt,
		&newSubject.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &newSubject, err
}

func GetAllSubjects(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID) ([]models.Subject, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, name, constraints, created_at, updated_at FROM subjects WHERE user_id = $1 AND project_id = $2 ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx, query, UserID, ProjectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject = []models.Subject{}

	for rows.Next() {
		var subject models.Subject

		err = rows.Scan(
			&subject.ID,
			&subject.UserID,
			&subject.ProjectID,
			&subject.Name,
			&subject.Constraints,
			&subject.CreatedAt,
			&subject.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		subjects = append(subjects, subject)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subjects, err
}

func GetSubjectByID(pool *pgxpool.Pool, UserID uuid.UUID, SubjectID uuid.UUID) (*models.Subject, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, name, constraints, created_at, updated_at FROM subjects WHERE user_id = $1 AND id = $2
	`
	var subject models.Subject
	err := pool.QueryRow(ctx, query, UserID, SubjectID).Scan(
		&subject.ID,
		&subject.UserID,
		&subject.ProjectID,
		&subject.Name,
		&subject.Constraints,
		&subject.CreatedAt,
		&subject.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("No subject found")
		}
		return nil, err
	}

	return &subject, err
}

func UpdateSubject(pool *pgxpool.Pool, UserID uuid.UUID, SubjectID uuid.UUID, subjectRequest *dto.SubjectUpdationRequest) (*models.Subject, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var setClauses []string
	args := pgx.NamedArgs{
		"user_id": UserID,
		"id":      SubjectID,
	}

	if subjectRequest.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *subjectRequest.Name
	}

	if subjectRequest.Constraints != nil {

		constraintsUpdate := make(map[string]any)

		if subjectRequest.Constraints.MorningTendency != nil {
			constraintsUpdate["morning_tendency"] = *subjectRequest.Constraints.MorningTendency
		}

		if subjectRequest.Constraints.MinPerDay != nil {
			constraintsUpdate["min_per_day"] = *subjectRequest.Constraints.MinPerDay
		}

		if subjectRequest.Constraints.MinPerWeek != nil {
			constraintsUpdate["min_per_week"] = *subjectRequest.Constraints.MinPerWeek
		}

		if subjectRequest.Constraints.MinConsecutive != nil {
			constraintsUpdate["min_consecutive"] = *subjectRequest.Constraints.MinConsecutive
		}

		if subjectRequest.Constraints.MaxPerDay != nil {
			constraintsUpdate["max_per_day"] = *subjectRequest.Constraints.MaxPerDay
		}

		if subjectRequest.Constraints.MaxPerWeek != nil {
			constraintsUpdate["max_per_week"] = *subjectRequest.Constraints.MaxPerWeek
		}

		if subjectRequest.Constraints.MaxConsecutive != nil {
			constraintsUpdate["max_consecutive"] = *subjectRequest.Constraints.MaxConsecutive
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
		"UPDATE subjects SET %s WHERE user_id = @user_id AND id = @id RETURNING id, user_id, project_id, name, constraints, created_at, updated_at",
		strings.Join(setClauses, ", "),
	)

	var subject models.Subject

	err := pool.QueryRow(ctx, query, args).Scan(
		&subject.ID,
		&subject.UserID,
		&subject.ProjectID,
		&subject.Name,
		&subject.Constraints,
		&subject.CreatedAt,
		&subject.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &subject, err
}

func DeleteSubject(pool *pgxpool.Pool, UserID uuid.UUID, SubjectID uuid.UUID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM subjects WHERE user_id = $1 AND id = $2
	`
	commandTag, err := pool.Exec(ctx, query, UserID, SubjectID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Subject with id %v was not found", SubjectID)
	}

	return err
}
