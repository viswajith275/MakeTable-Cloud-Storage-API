package handlers

import (
	"net/http"
	"timetable_api/internal/dto"
	"timetable_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTeacherAssignmentHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		var req dto.TeacherAssignmentCreationRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err = req.Constraints.ValidateForPost(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		teacher, err := repository.GetTeacherByID(pool, UserID, req.TeacherID)

		if err != nil {
			if err.Error() == "No teacher found" {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		class, err := repository.GetClassByID(pool, UserID, req.ClassID)

		if err != nil {
			if err.Error() == "No class found" || class.ProjectID != teacher.ProjectID {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		subject, err := repository.GetSubjectByID(pool, UserID, req.SubjectID)

		if err != nil {
			if err.Error() == "No subject found" || subject.ProjectID != teacher.ProjectID {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if req.TargetRoomID != nil {

			if class.RoomID == *req.TargetRoomID {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Room is already the home room of class",
				})
				return
			}

			target_room, err := repository.GetRoomByID(pool, UserID, *req.TargetRoomID)

			if err != nil {
				if err.Error() == "No room found" || target_room.ProjectID != teacher.ProjectID {
					ctx.JSON(http.StatusNotFound, gin.H{
						"error": err.Error(),
					})
					return
				}
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		if req.Constraints != nil && *req.Constraints.IsClassTeacher {

			teacher_assignments, err := repository.GetAssignmentsByTeacher(pool, UserID, teacher.ID)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Class Fetch to database error: " + err.Error(),
				})
				return
			}

			for _, teacher_assignment := range teacher_assignments {

				if *teacher_assignment.Constraints.IsClassTeacher {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"error": "Teacher is already a class teacher!",
					})
					return
				}
			}
		}

		exists, err := repository.CheckAssignmentExists(pool, req.ClassID, req.TeacherID, req.SubjectID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if exists {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Teacher assignment already exists!",
			})
			return
		}

		req.ProjectID = teacher.ProjectID

		teacher_assignment, err := repository.CreateTeacherAssignment(pool, UserID, &req)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Teacher creation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dto.TeacherAssignmentResponse{
			ID:           teacher_assignment.ID,
			ProjectID:    teacher_assignment.ProjectID,
			ClassID:      teacher_assignment.ClassID,
			TeacherID:    teacher_assignment.TeacherID,
			SubjectID:    teacher_assignment.SubjectID,
			TargetRoomID: teacher_assignment.TargetRoomID,
			Constraints:  teacher_assignment.Constraints,
			CreatedAt:    teacher_assignment.CreatedAt,
			UpdatedAt:    teacher_assignment.UpdatedAt,
		})
	}
}

func GetTeacherAssignmentHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		ProjectID, err := uuid.Parse(ctx.Param("project_id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid project id",
			})
			return
		}

		project, err := repository.GetProjectByID(pool, UserID, ProjectID)

		if err != nil {
			if err.Error() == "Project not found" {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		teacher_assignments, err := repository.GetAssignmentsByProject(pool, UserID, project.ID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Teacher Assignment Fetch to database error: " + err.Error(),
			})
			return
		}

		var res []dto.TeacherAssignmentResponse = []dto.TeacherAssignmentResponse{}

		for _, teacher_assignment := range teacher_assignments {
			res = append(res, dto.TeacherAssignmentResponse{
				ID:           teacher_assignment.ID,
				ProjectID:    teacher_assignment.ProjectID,
				ClassID:      teacher_assignment.ClassID,
				TeacherID:    teacher_assignment.TeacherID,
				SubjectID:    teacher_assignment.SubjectID,
				TargetRoomID: teacher_assignment.TargetRoomID,
				Constraints:  teacher_assignment.Constraints,
				CreatedAt:    teacher_assignment.CreatedAt,
				UpdatedAt:    teacher_assignment.UpdatedAt,
			})
		}

		ctx.JSON(http.StatusOK, res)
	}
}

func UpdateTeacherAssignmentHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		AssignmentID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid room id",
			})
			return
		}

		teacher_assignment, err := repository.GetAssignmentsByID(pool, UserID, AssignmentID)

		if err != nil {
			if err.Error() == "No assignment found" {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var req dto.TeacherAssignmentUpdationRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err = req.Constraints.Validate(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var CurrentRoomID *uuid.UUID
		var CurrentClassID *uuid.UUID
		var CurrentTeacherID *uuid.UUID
		var CurrentSubjectID *uuid.UUID

		if req.TeacherID != nil {

			teacher, err := repository.GetTeacherByID(pool, UserID, *req.TeacherID)

			if err != nil {
				if err.Error() == "No teacher found" || teacher.ProjectID != teacher_assignment.ProjectID {
					ctx.JSON(http.StatusNotFound, gin.H{
						"error": err.Error(),
					})
					return
				}
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			CurrentTeacherID = &teacher.ID

		} else {
			CurrentTeacherID = &teacher_assignment.TeacherID
		}

		if req.ClassID != nil {
			class, err := repository.GetClassByID(pool, UserID, *req.ClassID)

			if err != nil {
				if err.Error() == "No class found" || class.ProjectID != teacher_assignment.ProjectID {
					ctx.JSON(http.StatusNotFound, gin.H{
						"error": err.Error(),
					})
					return
				}
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			CurrentRoomID = &class.RoomID
			CurrentClassID = &class.ID

		} else {

			existingClass, err := repository.GetClassByID(pool, UserID, teacher_assignment.ClassID)

			if err != nil {
				if err.Error() == "No class found" {
					ctx.JSON(http.StatusNotFound, gin.H{
						"error": err.Error(),
					})
					return
				}
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			CurrentRoomID = &existingClass.RoomID
			CurrentClassID = &teacher_assignment.ClassID
		}

		if req.SubjectID != nil {

			subject, err := repository.GetSubjectByID(pool, UserID, *req.SubjectID)

			if err != nil {
				if err.Error() == "No subject found" || subject.ProjectID != teacher_assignment.ProjectID {
					ctx.JSON(http.StatusNotFound, gin.H{
						"error": err.Error(),
					})
					return
				}
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			CurrentSubjectID = &subject.ID
		} else {

			CurrentSubjectID = &teacher_assignment.SubjectID
		}

		if req.TargetRoomID != nil {

			if *CurrentRoomID == req.TargetRoomID.Value {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Room is already the home room of class",
				})
				return
			}

			target_room, err := repository.GetRoomByID(pool, UserID, req.TargetRoomID.Value)

			if err != nil {
				if err.Error() == "No room found" || target_room.ProjectID != teacher_assignment.ProjectID {
					ctx.JSON(http.StatusNotFound, gin.H{
						"error": err.Error(),
					})
					return
				}
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		if req.Constraints != nil && *req.Constraints.IsClassTeacher {

			teacher_assignments, err := repository.GetAssignmentsByTeacher(pool, UserID, *CurrentTeacherID)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Teacher Assignment Fetch to database error: " + err.Error(),
				})
				return
			}

			for _, teacher_assignment := range teacher_assignments {

				if *teacher_assignment.Constraints.IsClassTeacher {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"error": "Teacher is already a class teacher!",
					})
					return
				}
			}
		}

		if req.ClassID != nil || req.TeacherID != nil || req.SubjectID != nil {

			exists, err := repository.CheckAssignmentExists(pool, *CurrentClassID, *CurrentTeacherID, *CurrentSubjectID)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			if exists {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Teacher assignment already exists!",
				})
				return
			}
		}

		teacher_assignment, err = repository.UpdateTeacherAssignment(pool, UserID, AssignmentID, &req)

		if err != nil {
			if err.Error() == "Nothing to update!" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Teacher Assignment updation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dto.TeacherAssignmentResponse{
			ID:           teacher_assignment.ID,
			ProjectID:    teacher_assignment.ProjectID,
			ClassID:      teacher_assignment.ClassID,
			TeacherID:    teacher_assignment.TeacherID,
			SubjectID:    teacher_assignment.SubjectID,
			TargetRoomID: teacher_assignment.TargetRoomID,
			Constraints:  teacher_assignment.Constraints,
			CreatedAt:    teacher_assignment.CreatedAt,
			UpdatedAt:    teacher_assignment.UpdatedAt,
		})

	}
}

func DeleteTeacherAssignmentHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		AssignmentID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid room id",
			})
			return
		}

		err = repository.DeleteTeacherAssignment(pool, UserID, AssignmentID)

		if err != nil {
			if err.Error() == "Assignment with id "+AssignmentID.String()+" was not found" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Teacher Assignment deletion to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Deleted successfully",
		})
	}
}
