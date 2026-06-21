package handlers

import (
	"net/http"
	"timetable_api/internal/dto"
	"timetable_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateRoomHandler(pool *pgxpool.Pool) gin.HandlerFunc {
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

		var req dto.RoomCreationRequest

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

		room, err := repository.CreateRoom(pool, UserID, project.ID, &req)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Room creation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dto.RoomResponse{
			ID:          room.ID,
			ProjectID:   room.ProjectID,
			Name:        room.Name,
			IsLab:       room.IsLab,
			Constraints: room.Constraints,
			CreatedAt:   room.CreatedAt,
			UpdatedAt:   room.UpdatedAt,
		})

	}
}

func GetRoomsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
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

		rooms, err := repository.GetAllRooms(pool, UserID, project.ID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Room Fetch to database error: " + err.Error(),
			})
			return
		}

		var res []dto.RoomResponse = []dto.RoomResponse{}

		for _, room := range rooms {
			res = append(res, dto.RoomResponse{
				ID:          room.ID,
				ProjectID:   room.ProjectID,
				Name:        room.Name,
				IsLab:       room.IsLab,
				Constraints: room.Constraints,
				CreatedAt:   room.CreatedAt,
				UpdatedAt:   room.UpdatedAt,
			})
		}

		ctx.JSON(http.StatusCreated, res)
	}
}

func UpdateRoomHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		RoomID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid room id",
			})
			return
		}

		room, err := repository.GetRoomByID(pool, UserID, RoomID)

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

		var req dto.RoomUpdationRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		room, err = repository.UpdateRoom(pool, UserID, RoomID, &req)

		if err != nil {
			if err.Error() == "Nothing to update!" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Room updation to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, dto.RoomResponse{
			ID:          room.ID,
			ProjectID:   room.ProjectID,
			Name:        room.Name,
			IsLab:       room.IsLab,
			Constraints: room.Constraints,
			CreatedAt:   room.CreatedAt,
			UpdatedAt:   room.UpdatedAt,
		})

	}
}

func DeleteRoomHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		RoomID, err := uuid.Parse(ctx.Param("id"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid room id",
			})
			return
		}

		err = repository.DeleteRoom(pool, UserID, RoomID)

		if err != nil {
			if err.Error() == "Room with id "+RoomID.String()+" was not found" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Room deletion to database error: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Deleted successfully",
		})

	}
}
