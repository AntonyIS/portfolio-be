/*
Package name : http
File name : gin.go
Author : Antony Injila
Description :
	- Host Go Gin webframework for handling HTTP requests
	- Routes HTTP requests to thier respective handlers
*/
package gin

import (
	"fmt"
	"os"

	"github.com/AntonyIS/portfolio-be/config"
	// "github.com/AntonyIS/portfolio-be/internal/adapters/middleware"
	"github.com/AntonyIS/portfolio-be/internal/core/services"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func InitGinRoutes(svc services.PortfolioService, config config.AppConfig) {
	// Enable detailed error responses
	gin.SetMode(gin.DebugMode)

	// Setup Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Setup application route handlers
	handler := NewGinHandler(svc)
	router.GET("/", handler.Home)
	router.POST("/api/v1/login", handler.Login)
	router.POST("/api/v1/signup", handler.Signup)

	// Group users API
	usersRoutes := router.Group("/api/v1/users")

	// Group projects API
	projectsRoutes := router.Group("/api/v1/projects")

	// Add middleware in production
	// if config.Env == "pro" {
	// 	middleware := middleware.NewMiddleware(&svc)
	// 	usersRoutes.Use(middleware.Authorize)
	// 	projectsRoutes.Use(middleware.Authorize)
	// }

	{
		usersRoutes.GET("/", handler.GetUsers)
		usersRoutes.GET("/:id", handler.GetUser)
		usersRoutes.POST("/", handler.PostUser)
		usersRoutes.PUT("/:id", handler.PutUser)
		usersRoutes.DELETE("/:id", handler.DeleteUser)
	}
	{
		projectsRoutes.GET("/", handler.GetProjects)
		projectsRoutes.GET("/:id", handler.GetProject)
		projectsRoutes.POST("/", handler.PostProject)
		projectsRoutes.PUT("/:id", handler.PutProject)
		projectsRoutes.DELETE("/:id", handler.DeleteProject)
	}

	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))

	router.Run(port)
}
