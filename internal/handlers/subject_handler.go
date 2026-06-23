package handlers

import (
	"net/http"
	"timetable_api/internal/dto"
	"timetable_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateSubjectHandler(pool *pgxpool.Pool) gin.HandlerFunc {
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

		var req dto.SubjectCreationRequest

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

		subject, err := repository.CreateSubject(pool, UserID, project.ID, &req)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Subject creation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dto.SubjectResponse{
			ID:          subject.ID,
			ProjectID:   subject.ProjectID,
			Name:        subject.Name,
			Constraints: subject.Constraints,
			CreatedAt:   subject.CreatedAt,
			UpdatedAt:   subject.UpdatedAt,
		})

	}
}

func GetSubjectHandler(pool *pgxpool.Pool) gin.HandlerFunc {
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

		subjects, err := repository.GetAllSubjects(pool, UserID, project.ID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Subject Fetch to database error: " + err.Error(),
			})
			return
		}

		var res []dto.SubjectResponse = []dto.SubjectResponse{}

		for _, subject := range subjects {
			res = append(res, dto.SubjectResponse{
				ID:          subject.ID,
				ProjectID:   subject.ProjectID,
				Name:        subject.Name,
				Constraints: subject.Constraints,
				CreatedAt:   subject.CreatedAt,
				UpdatedAt:   subject.UpdatedAt,
			})
		}

		ctx.JSON(http.StatusOK, res)
	}
}

func UpdateSubjectHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		SubjectID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid room id",
			})
			return
		}

		subject, err := repository.GetSubjectByID(pool, UserID, SubjectID)

		if err != nil {
			if err.Error() == "No subjects found" {
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

		var req dto.SubjectUpdationRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		subject, err = repository.UpdateSubject(pool, UserID, SubjectID, &req)

		if err != nil {
			if err.Error() == "Nothing to update!" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Subject updation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, dto.SubjectResponse{
			ID:          subject.ID,
			ProjectID:   subject.ProjectID,
			Name:        subject.Name,
			Constraints: subject.Constraints,
			CreatedAt:   subject.CreatedAt,
			UpdatedAt:   subject.UpdatedAt,
		})

	}
}

func DeleteSubjectHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		SubjectID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid room id",
			})
			return
		}

		err = repository.DeleteSubject(pool, UserID, SubjectID)

		if err != nil {
			if err.Error() == "Subject with id "+SubjectID.String()+" was not found" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Subject deletion to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Deleted successfully",
		})

	}
}
