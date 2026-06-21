package handlers

import (
	"net/http"
	"timetable_api/internal/dto"
	"timetable_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		UserID, err := uuid.Parse(ctx.GetString("user_id"))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User id has not been setted in the context",
			})
			return
		}

		user, err := repository.GetUserById(pool, UserID)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		ctx.JSON(http.StatusOK, &dto.UserResponse{
			Email:     user.Email,
			Username:  user.UserName,
			Disabled:  user.Disabled,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}
