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

	"github.com/AntonyIS/portfolio-be/internal/adapters/middleware"
	"github.com/AntonyIS/portfolio-be/internal/core/services"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(svc services.PortfolioService) {
	handler := NewGinHandler(svc)
	router := gin.Default()
	middleware := middleware.NewMiddleware(&svc)
	usersRoutes := router.Group("/v1/users")
	projectsRoutes := router.Group("/v1/projects")
	router.GET("/", handler.Home)
	router.POST("/login", handler.Login)
	router.POST("/signup", handler.Signup)
	usersRoutes.Use(middleware.Authorize)
	projectsRoutes.Use(middleware.Authorize)
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
	port := os.Getenv("SERVER_PORT")
	router.Run(fmt.Sprintf(":%s", port))
}
