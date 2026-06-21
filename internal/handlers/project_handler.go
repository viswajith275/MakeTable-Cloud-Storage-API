package handlers

import (
	"net/http"
	"timetable_api/internal/dto"
	"timetable_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateProjectHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		var req dto.ProjectCreationRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		project, err := repository.CreateProject(pool, UserID, &req)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Project creation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dto.ProjectResponse{
			ID:                 project.ID,
			Name:               project.Name,
			Slots:              project.Slots,
			Days:               project.Days,
			RoomsVersion:       project.RoomsVersion,
			ClassesVersion:     project.ClassesVersion,
			TeachersVersion:    project.TeachersVersion,
			SubjectsVersion:    project.SubjectsVersion,
			AssignmentsVersion: project.AssignmentsVersion,
			CreatedAt:          project.CreatedAt,
			UpdatedAt:          project.UpdatedAt,
		})
	}
}

func GetAllProjectsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}
		projects, err := repository.GetAllProject(pool, UserID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Project fetch database error: " + err.Error(),
			})
			return
		}

		var res []dto.ProjectResponse = []dto.ProjectResponse{}

		for _, project := range projects {
			res = append(res, dto.ProjectResponse{
				ID:                 project.ID,
				Name:               project.Name,
				Slots:              project.Slots,
				Days:               project.Days,
				RoomsVersion:       project.RoomsVersion,
				ClassesVersion:     project.ClassesVersion,
				TeachersVersion:    project.TeachersVersion,
				SubjectsVersion:    project.SubjectsVersion,
				AssignmentsVersion: project.AssignmentsVersion,
				CreatedAt:          project.CreatedAt,
				UpdatedAt:          project.UpdatedAt,
			})
		}

		ctx.JSON(http.StatusOK, res)
	}
}

func GetProjectByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}
		ProjectID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid project id",
			})
			return
		}

		project, err := repository.GetProjectByID(pool, UserID, ProjectID)

		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "Project not found",
				})
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Project fetch by id database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, dto.ProjectResponse{
			ID:                 project.ID,
			Name:               project.Name,
			Slots:              project.Slots,
			Days:               project.Days,
			RoomsVersion:       project.RoomsVersion,
			ClassesVersion:     project.ClassesVersion,
			TeachersVersion:    project.TeachersVersion,
			SubjectsVersion:    project.SubjectsVersion,
			AssignmentsVersion: project.AssignmentsVersion,
			CreatedAt:          project.CreatedAt,
			UpdatedAt:          project.UpdatedAt,
		})

	}
}

func UpdateProjectHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}
		ProjectID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid project id",
			})
			return
		}

		var req dto.ProjectUpdationRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		project, err := repository.UpdateProject(pool, UserID, ProjectID, &req)

		if err != nil {
			if err.Error() == "Nothing to update!" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Project updation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dto.ProjectResponse{
			ID:                 project.ID,
			Name:               project.Name,
			Slots:              project.Slots,
			Days:               project.Days,
			RoomsVersion:       project.RoomsVersion,
			ClassesVersion:     project.ClassesVersion,
			TeachersVersion:    project.TeachersVersion,
			SubjectsVersion:    project.SubjectsVersion,
			AssignmentsVersion: project.AssignmentsVersion,
			CreatedAt:          project.CreatedAt,
			UpdatedAt:          project.UpdatedAt,
		})

	}
}

func DeleteProjectHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}
		ProjectID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid project id",
			})
			return
		}

		err = repository.DeleteProject(pool, UserID, ProjectID)

		if err != nil {
			if err.Error() == "Project with id "+ProjectID.String()+" was not found" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Project deletion to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"error": "Deleted successfully!",
		})
	}
}
