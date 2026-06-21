package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	"timetable_api/internal/config"
	"timetable_api/internal/dto"
	"timetable_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func HashRefreshToken(token string) string {
	hash := sha256.New()

	hash.Write([]byte(token))

	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateToken(UserID uuid.UUID, cfg *config.Config) (string, string, error) {

	accessClaims := jwt.MapClaims{
		"user_id": UserID,
		"exp":     time.Now().Add(time.Duration(cfg.AccessTokenTTL) * time.Minute).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(cfg.SecretKEY))

	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"user_id": UserID,
		"exp":     time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Hour * 24).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte(cfg.SecretKEY))

	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, err
}

func RegisterHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req dto.UserRegisterRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password" + err.Error(),
			})
			return
		}

		req.Password = string(hashed_password)

		user, err := repository.CreateUser(pool, &req)

		if err != nil {
			if err.Error() == "Email already exists!" {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "user creation repository error : " + err.Error(),
			})
			return
		}

		res := dto.UserResponse{
			Email:     user.Email,
			Username:  user.UserName,
			Disabled:  user.Disabled,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		ctx.JSON(http.StatusCreated, res)
	}
}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req dto.UserLoginRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := repository.GetUserByEmail(pool, req.Email)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		accessToken, refreshToken, err := GenerateToken(user.ID, cfg)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "filed to generate tokens: " + err.Error(),
			})
			return
		}

		refreshTokenRequest := dto.UserTokenCreateRequest{
			UserID:    user.ID,
			ExpiresAT: time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Hour * 24),
			TokenHash: HashRefreshToken(refreshToken),
		}

		err = repository.CreateUserToken(pool, &refreshTokenRequest)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while token creation in db: " + err.Error(),
			})
			return
		}

		ctx.SetCookie("access_token", accessToken, cfg.AccessTokenTTL*60, "/", "", false, true)
		ctx.SetCookie("refresh_token", refreshToken, cfg.RefreshTokenTTL*60*60*24, "/", "", false, true)

		ctx.JSON(http.StatusOK, dto.UserResponse{
			Email:     user.Email,
			Username:  user.UserName,
			Disabled:  user.Disabled,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})

	}
}

func RefreshHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		refreshTokenString, err := ctx.Cookie("refresh_token")

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "No refresh token found",
			})
			return
		}

		token, err := jwt.Parse(refreshTokenString, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("Unexpected signing method %v", t.Header["alg"])
			}
			return []byte(cfg.SecretKEY), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		refreshToken, err := repository.GetUserTokenByTokenHash(pool, HashRefreshToken(refreshTokenString))

		if err != nil || refreshToken.IsRevoked || time.Now().After(refreshToken.ExpiresAt) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Refresh token invalid, revoked and expired",
			})
			return
		}

		newRefreshToken, newAccessToken, err := GenerateToken(refreshToken.UserID, cfg)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "filed to generate tokens: " + err.Error(),
			})
			return
		}

		err = repository.CreateUserToken(pool, &dto.UserTokenCreateRequest{
			UserID:    refreshToken.UserID,
			TokenHash: HashRefreshToken(newRefreshToken),
			ExpiresAT: time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Hour * 24),
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while token creation in db: " + err.Error(),
			})
			return
		}

		var isRevoked bool = true

		err = repository.UpdateUserToken(pool, &dto.UserTokenUpdateRequest{
			TokenHash: refreshToken.TokenHash,
			IsRevoked: &isRevoked,
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.SetCookie("access_token", newAccessToken, cfg.AccessTokenTTL*60, "/", "", false, true)
		ctx.SetCookie("refresh_token", newRefreshToken, cfg.RefreshTokenTTL*60*60*24, "/", "", false, true)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Tokens refreshed successfully!",
		})

	}
}

func LogoutHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		refreshTokenString, err := ctx.Cookie("refresh_token")

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "No refresh token found",
			})
			return
		}

		var isRevoked bool = true

		err = repository.UpdateUserToken(pool, &dto.UserTokenUpdateRequest{
			TokenHash: HashRefreshToken(refreshTokenString),
			IsRevoked: &isRevoked,
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.SetCookie("access_token", "", -1, "/", "", false, true)
		ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Logged out successfully",
		})
	}
}
