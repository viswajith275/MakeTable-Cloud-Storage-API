package handlers

import (
	"net/http"
	"timetable_api/internal/dto"
	"timetable_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTeacherHandler(pool *pgxpool.Pool) gin.HandlerFunc {
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

		var req dto.TeacherCreationRequest

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

		teacher, err := repository.CreateTeacher(pool, UserID, project.ID, &req)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Teacher creation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dto.TeacherResponse{
			ID:          teacher.ID,
			ProjectID:   teacher.ProjectID,
			Name:        teacher.Name,
			Constraints: teacher.Constraints,
			CreatedAt:   teacher.CreatedAt,
			UpdatedAt:   teacher.UpdatedAt,
		})

	}
}

func GetTeachersHandler(pool *pgxpool.Pool) gin.HandlerFunc {
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

		teachers, err := repository.GetAllTeachers(pool, UserID, project.ID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Class Fetch to database error: " + err.Error(),
			})
			return
		}

		var res []dto.TeacherResponse = []dto.TeacherResponse{}

		for _, teacher := range teachers {
			res = append(res, dto.TeacherResponse{
				ID:          teacher.ID,
				ProjectID:   teacher.ProjectID,
				Name:        teacher.Name,
				Constraints: teacher.Constraints,
				CreatedAt:   teacher.CreatedAt,
				UpdatedAt:   teacher.UpdatedAt,
			})
		}

		ctx.JSON(http.StatusOK, res)
	}
}

func UpdateTeacherHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		TeacherID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid room id",
			})
			return
		}

		teacher, err := repository.GetTeacherByID(pool, UserID, TeacherID)

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

		var req dto.TeacherUpdationRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		teacher, err = repository.UpdateTeacher(pool, UserID, TeacherID, &req)

		if err != nil {
			if err.Error() == "Nothing to update!" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Teacher updation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, dto.TeacherResponse{
			ID:          teacher.ID,
			ProjectID:   teacher.ProjectID,
			Name:        teacher.Name,
			Constraints: teacher.Constraints,
			CreatedAt:   teacher.CreatedAt,
			UpdatedAt:   teacher.UpdatedAt,
		})

	}
}

func DeleteTeacherHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		TeacherID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid room id",
			})
			return
		}

		err = repository.DeleteTeacher(pool, UserID, TeacherID)

		if err != nil {
			if err.Error() == "Teacher with id "+TeacherID.String()+" was not found" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Teacher deletion to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Deleted successfully",
		})

	}
}
