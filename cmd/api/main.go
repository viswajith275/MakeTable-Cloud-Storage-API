package main

import (
	"log"
	"net/http"
	"timetable_api/internal/config"
	"timetable_api/internal/db"
	"timetable_api/internal/handlers"
	"timetable_api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatal("Enviornment variables not loaded successfully!")
		return
	}

	pool, err := db.ConnectDatabase(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Database Connection Error!")
		return
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"health":   "Ok",
			"Database": "connected",
		})
	})
	auth_router := router.Group("/auth")
	{
		auth_router.POST("/register", handlers.RegisterHandler(pool))
		auth_router.POST("/login", handlers.LoginHandler(pool, cfg))
		auth_router.POST("/refresh", handlers.RefreshHandler(pool, cfg))
	}

	user_routes := router.Group("/users")
	user_routes.Use(middleware.AuthMiddleWare(cfg))

	{
		user_routes.POST("/logout", handlers.LogoutHandler(pool))
		user_routes.GET("/me", handlers.GetUserHandler(pool))
	}

	project_router := router.Group("/projects")
	project_router.Use(middleware.AuthMiddleWare(cfg))

	{
		project_router.GET("/", handlers.GetAllProjectsHandler(pool))
		project_router.POST("/", handlers.CreateProjectHandler(pool))
		project_router.GET("/:id", handlers.GetProjectByIDHandler(pool))
		project_router.PATCH("/:id", handlers.UpdateProjectHandler(pool))
		project_router.DELETE("/:id", handlers.DeleteProjectHandler(pool))
	}

	room_router := router.Group("/rooms")
	room_router.Use(middleware.AuthMiddleWare(cfg))

	{
		room_router.GET("/:project_id", handlers.GetRoomsHandler(pool))
		room_router.POST("/:project_id", handlers.CreateRoomHandler(pool))
		room_router.PATCH("/:id", handlers.UpdateRoomHandler(pool))
		room_router.DELETE("/:id", handlers.DeleteRoomHandler(pool))
	}

	class_router := router.Group("/classes")
	class_router.Use(middleware.AuthMiddleWare(cfg))

	{
		class_router.GET("/:project_id", handlers.GetClassesHandler(pool))
		class_router.POST("/:project_id", handlers.CreateClassHandler(pool))
		class_router.PATCH("/:id", handlers.UpdateClassHandler(pool))
		class_router.DELETE("/:id", handlers.DeleteClassHandler(pool))
	}

	teacher_router := router.Group("/teachers")
	teacher_router.Use(middleware.AuthMiddleWare(cfg))

	{
		teacher_router.GET("/:project_id", handlers.GetTeachersHandler(pool))
		teacher_router.POST("/:project_id", handlers.CreateTeacherHandler(pool))
		teacher_router.PATCH("/:id", handlers.UpdateTeacherHandler(pool))
		teacher_router.DELETE("/:id", handlers.DeleteTeacherHandler(pool))
	}

	subject_router := router.Group("/subjects")
	subject_router.Use(middleware.AuthMiddleWare(cfg))

	{
		subject_router.GET("/:project_id", handlers.GetSubjectHandler(pool))
		subject_router.POST("/:project_id", handlers.CreateSubjectHandler(pool))
		subject_router.PATCH("/:id", handlers.UpdateSubjectHandler(pool))
		subject_router.DELETE("/:id", handlers.DeleteSubjectHandler(pool))
	}

	teacher_assignment_router := router.Group("/assignments")
	teacher_assignment_router.Use(middleware.AuthMiddleWare(cfg))

	{
		teacher_assignment_router.GET("/:project_id", handlers.GetTeacherAssignmentHandler(pool))
		teacher_assignment_router.POST("/", handlers.CreateTeacherAssignmentHandler(pool))
		teacher_assignment_router.PATCH("/:id", handlers.UpdateTeacherAssignmentHandler(pool))
		teacher_assignment_router.DELETE("/:id", handlers.DeleteTeacherAssignmentHandler(pool))
	}

	router.Run(":" + cfg.Port)

}
