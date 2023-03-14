/*
Package name : http
File name : gin.go
Author : Antony Injila
Description :
	- Host Go Gin webframework for handling HTTP requests
	- Routes HTTP requests to thier respective handlers
*/
package http

import (
	"fmt"
	"net/http"

	"github.com/AntonyIS/portfolio-be/config"
	"github.com/AntonyIS/portfolio-be/internal/adapters/repository"
	"github.com/AntonyIS/portfolio-be/internal/core/domain"
	"github.com/AntonyIS/portfolio-be/internal/core/services"
	"github.com/gin-gonic/gin"
)

type GinHandler interface {
	PostUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	GetUsers(ctx *gin.Context)
	PutUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	PostProject(ctx *gin.Context)
	GetProject(ctx *gin.Context)
	GetProjects(ctx *gin.Context)
	PutProject(ctx *gin.Context)
	DeleteProject(ctx *gin.Context)
}

type handler struct {
	svc services.PortfolioService
}

func NewGinHandler(svc services.PortfolioService) GinHandler {
	return handler{
		svc: svc,
	}
}

func (h handler) PostUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, err := h.svc.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h handler) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := h.svc.ReadUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
	return
}

func (h handler) GetUsers(ctx *gin.Context) {
	users, err := h.svc.ReadUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (h handler) PutUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, err := h.svc.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
	return
}

func (h handler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.svc.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User deleted successfully",
	})
}

func (h handler) PostProject(ctx *gin.Context) {
	var project domain.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, err := h.svc.CreateProject(&project)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h handler) GetProject(ctx *gin.Context) {
	id := ctx.Param("id")
	project, err := h.svc.ReadProject(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	if project == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "project not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, project)
	return
}

func (h handler) GetProjects(ctx *gin.Context) {
	projects, err := h.svc.ReadProjects()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, projects)
}

func (h handler) PutProject(ctx *gin.Context) {
	var project domain.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, err := h.svc.UpdateProject(&project)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
	return
}

func (h handler) DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.svc.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Project deleted successfully",
	})
}

func home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Antony Injila Portfolio",
	})
}

func InitGinRoutes() {
	// Load application configuration
	config := config.NewConfiguration()
	// DynamoDB repository
	repo := repository.NewDynamoDBRepository(config)
	// Portifolio service
	svc := services.NewPortfolioService(&repo)
	// Gin Route Handler
	handler := NewGinHandler(*svc)
	// Initilize Gin
	router := gin.Default()
	// Home route
	router.GET("/", home)
	// Users routes
	usersRoutes := router.Group("/v1/users")
	// Projects Routes
	projectsRoutes := router.Group("/v1/projects")

	{
		// Group users routes
		usersRoutes.GET("/", handler.GetUsers)
		usersRoutes.GET("/:id", handler.GetUser)
		usersRoutes.POST("/", handler.PostUser)
		usersRoutes.PUT("/:id", handler.PutUser)
		usersRoutes.DELETE("/:id", handler.DeleteUser)
	}
	{
		// Group projects routes
		projectsRoutes.GET("/", handler.GetProjects)
		projectsRoutes.GET("/:id", handler.GetProject)
		projectsRoutes.POST("/", handler.PostProject)
		projectsRoutes.PUT("/:id", handler.PutProject)
		projectsRoutes.DELETE("/:id", handler.DeleteProject)
	}
	// Run Gin web server
	router.Run(fmt.Sprintf(":%s", config.Port))
}
