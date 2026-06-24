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

func CreateTeacherAssignment(pool *pgxpool.Pool, UserID uuid.UUID, teacherAssignmentRequest *dto.TeacherAssignmentCreationRequest) (*models.TeacherAssignment, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO teacher_assignments (user_id, project_id, class_id, teacher_id, subject_id, target_room_id, constraints) VALUES ($1, $2, $3, $4, $5, $6, $7::JSONB) RETURNING id, user_id, project_id, teacher_id, class_id, subject_id, target_room_id, constraints, created_at, updated_at
	`

	var newTeacherAssignment models.TeacherAssignment

	err := pool.QueryRow(ctx, query, UserID, teacherAssignmentRequest.ProjectID, teacherAssignmentRequest.ClassID, teacherAssignmentRequest.TeacherID, teacherAssignmentRequest.SubjectID, teacherAssignmentRequest.TargetRoomID, teacherAssignmentRequest.Constraints).Scan(
		&newTeacherAssignment.ID,
		&newTeacherAssignment.UserID,
		&newTeacherAssignment.ProjectID,
		&newTeacherAssignment.TeacherID,
		&newTeacherAssignment.ClassID,
		&newTeacherAssignment.SubjectID,
		&newTeacherAssignment.TargetRoomID,
		&newTeacherAssignment.Constraints,
		&newTeacherAssignment.CreatedAt,
		&newTeacherAssignment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &newTeacherAssignment, err
}

func GetAssignmentsByProject(pool *pgxpool.Pool, UserID uuid.UUID, ProjectID uuid.UUID) ([]models.TeacherAssignment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, class_id, teacher_id, subject_id, target_room_id, constraints, created_at, updated_at FROM teacher_assignments WHERE user_id = $1 AND project_id = $2 ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx, query, UserID, ProjectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teacher_assignments []models.TeacherAssignment = []models.TeacherAssignment{}

	for rows.Next() {
		var teacher_assignment models.TeacherAssignment

		err = rows.Scan(
			&teacher_assignment.ID,
			&teacher_assignment.UserID,
			&teacher_assignment.ProjectID,
			&teacher_assignment.ClassID,
			&teacher_assignment.TeacherID,
			&teacher_assignment.SubjectID,
			&teacher_assignment.TargetRoomID,
			&teacher_assignment.Constraints,
			&teacher_assignment.CreatedAt,
			&teacher_assignment.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		teacher_assignments = append(teacher_assignments, teacher_assignment)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teacher_assignments, err
}

func CheckAssignmentExists(pool *pgxpool.Pool, ClassID uuid.UUID, TeacherID uuid.UUID, SubjectID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT EXISTS (SELECT 1 FROM teacher_assignments WHERE class_id = $1 AND teacher_id = $2 AND subject_id = $3)
	`
	var exists bool
	err := pool.QueryRow(ctx, query, ClassID, TeacherID, SubjectID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, err
}

func GetAssignmentsByTeacher(pool *pgxpool.Pool, UserID uuid.UUID, TeacherID uuid.UUID) ([]models.TeacherAssignment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, class_id, teacher_id, subject_id, target_room_id, constraints, created_at, updated_at FROM teacher_assignments WHERE user_id = $1 AND teacher_id = $2 ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx, query, UserID, TeacherID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teacher_assignments []models.TeacherAssignment = []models.TeacherAssignment{}

	for rows.Next() {
		var teacher_assignment models.TeacherAssignment

		err = rows.Scan(
			&teacher_assignment.ID,
			&teacher_assignment.UserID,
			&teacher_assignment.ProjectID,
			&teacher_assignment.ClassID,
			&teacher_assignment.TeacherID,
			&teacher_assignment.SubjectID,
			&teacher_assignment.TargetRoomID,
			&teacher_assignment.Constraints,
			&teacher_assignment.CreatedAt,
			&teacher_assignment.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		teacher_assignments = append(teacher_assignments, teacher_assignment)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teacher_assignments, err
}

func GetAssignmentsByID(pool *pgxpool.Pool, UserID uuid.UUID, AssignmentID uuid.UUID) (*models.TeacherAssignment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		SELECT id, user_id, project_id, class_id, teacher_id, subject_id, target_room_id, constraints, created_at, updated_at FROM teacher_assignments WHERE user_id = $1 AND id = $2
	`
	var teacher_assignment models.TeacherAssignment
	err := pool.QueryRow(ctx, query, UserID, AssignmentID).Scan(
		&teacher_assignment.ID,
		&teacher_assignment.UserID,
		&teacher_assignment.ProjectID,
		&teacher_assignment.ClassID,
		&teacher_assignment.TeacherID,
		&teacher_assignment.SubjectID,
		&teacher_assignment.TargetRoomID,
		&teacher_assignment.Constraints,
		&teacher_assignment.CreatedAt,
		&teacher_assignment.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("No assignment found")
		}
		return nil, err
	}

	return &teacher_assignment, err
}

func UpdateTeacherAssignment(pool *pgxpool.Pool, UserID uuid.UUID, AssignmentID uuid.UUID, teacherAssignmentRequest *dto.TeacherAssignmentUpdationRequest) (*models.TeacherAssignment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var setClauses []string
	args := pgx.NamedArgs{
		"user_id": UserID,
		"id":      AssignmentID,
	}

	if teacherAssignmentRequest.ClassID != nil {
		setClauses = append(setClauses, "class_id = @class_id")
		args["class_id"] = *teacherAssignmentRequest.ClassID
	}
	if teacherAssignmentRequest.TeacherID != nil {
		setClauses = append(setClauses, "teacher_id = @teaher_id")
		args["teacher_id"] = *teacherAssignmentRequest.TeacherID
	}
	if teacherAssignmentRequest.SubjectID != nil {
		setClauses = append(setClauses, "subject_id = @subject_id")
		args["subject_id"] = *teacherAssignmentRequest.SubjectID
	}

	if teacherAssignmentRequest.TargetRoomID != nil {
		setClauses = append(setClauses, "target_room_id = @target_room_id")

		if teacherAssignmentRequest.TargetRoomID.Valid {
			args["target_room_id"] = teacherAssignmentRequest.TargetRoomID.Value
		} else {
			args["target_room_id"] = nil
		}
	}

	if teacherAssignmentRequest.Constraints != nil {
		constraintsUpdate := make(map[string]any)

		// Add more constraints if added in the future

		constraintsUpdate["is_class_teacher"] = *teacherAssignmentRequest.Constraints.IsClassTeacher
		constraintsUpdate["first_slot_days"] = *teacherAssignmentRequest.Constraints.FirstSlotDays

		if len(constraintsUpdate) > 0 {
			setClauses = append(setClauses, "constraints = COALESCE(constraints, '{}'::jsonb) || @constraints::jsonb")
			args["constraints"] = constraintsUpdate
		}
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("Nothing to update!")
	}

	query := fmt.Sprintf(
		"UPDATE teacher_assignments SET %s WHERE user_id = @user_id AND id = @id RETURNING id, user_id, project_id, class_id, teacher_id, subject_id, target_room_id, constraints, created_at, updated_at",
		strings.Join(setClauses, ", "),
	)

	var teacher_assignment models.TeacherAssignment

	err := pool.QueryRow(ctx, query, args).Scan(
		&teacher_assignment.ID,
		&teacher_assignment.UserID,
		&teacher_assignment.ProjectID,
		&teacher_assignment.ClassID,
		&teacher_assignment.TeacherID,
		&teacher_assignment.SubjectID,
		&teacher_assignment.TargetRoomID,
		&teacher_assignment.Constraints,
		&teacher_assignment.CreatedAt,
		&teacher_assignment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &teacher_assignment, err
}

func DeleteTeacherAssignment(pool *pgxpool.Pool, UserID uuid.UUID, AssignmentID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM teacher_assignments WHERE user_id = $1 AND id = $2
	`
	commandTag, err := pool.Exec(ctx, query, UserID, AssignmentID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Assignment with id %v was not found", AssignmentID)
	}

	return err
}
