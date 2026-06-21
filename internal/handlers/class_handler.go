package handlers

import (
	"net/http"
	"timetable_api/internal/dto"
	"timetable_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateClassHandler(pool *pgxpool.Pool) gin.HandlerFunc {
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

		var req dto.ClassCreationRequest

		room, err := repository.GetRoomByID(pool, UserID, req.RoomID)

		if err != nil {
			if err.Error() == "No room found" {
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

		if room.ProjectID != ProjectID {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "No room found",
			})
			return
		}

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

		class, err := repository.CreateClass(pool, UserID, ProjectID, &req)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Class creation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dto.ClassResponse{
			ID:          class.ID,
			ProjectID:   class.ProjectID,
			Name:        class.Name,
			RoomID:      room.ID,
			Constraints: class.Constraints,
			CreatedAt:   class.CreatedAt,
			UpdatedAt:   class.UpdatedAt,
		})

	}
}

func GetClassesHandler(pool *pgxpool.Pool) gin.HandlerFunc {
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

		classes, err := repository.GetAllClassess(pool, UserID, ProjectID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Class Fetch to database error: " + err.Error(),
			})
			return
		}

		var res []dto.ClassResponse = []dto.ClassResponse{}

		for _, class := range classes {
			res = append(res, dto.ClassResponse{
				ID:          class.ID,
				ProjectID:   class.ProjectID,
				Name:        class.Name,
				RoomID:      class.RoomID,
				Constraints: class.Constraints,
				CreatedAt:   class.CreatedAt,
				UpdatedAt:   class.UpdatedAt,
			})
		}

		ctx.JSON(http.StatusCreated, res)
	}
}

func UpdateClassHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		ClassID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid room id",
			})
			return
		}

		var req dto.ClassUpdationRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if req.RoomID != &uuid.Nil {
			class, err := repository.GetClassByID(pool, UserID, ClassID)

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

			room, err := repository.GetRoomByID(pool, UserID, *req.RoomID)

			if err != nil {
				if err.Error() == "No room found" {
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

			if room.ProjectID != class.ProjectID {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "No room found",
				})
				return
			}
		}

		class, err := repository.UpdateClass(pool, UserID, ClassID, &req)

		if err != nil {
			if err.Error() == "Nothing to update!" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Class updation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dto.ClassResponse{
			ID:          class.ID,
			ProjectID:   class.ProjectID,
			Name:        class.Name,
			RoomID:      class.RoomID,
			Constraints: class.Constraints,
			CreatedAt:   class.CreatedAt,
			UpdatedAt:   class.UpdatedAt,
		})

	}
}

func DeleteClassHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		ClassID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid room id",
			})
			return
		}

		err = repository.DeleteClass(pool, UserID, ClassID)

		if err != nil {
			if err.Error() == "Class with id "+ClassID.String()+" was not found" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Class deletion to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Deleted successfully",
		})

	}
}
