/*
Package name : http
File name : gin.go
Author : Antony Injila
Description :
	- Host Go Gin webframework for handling HTTP requests
*/
package http

import (
	"net/http"

	"github.com/AntonyIS/portfolio-be/internal/adapters/repostitory"
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
	res := h.svc.DeleteUser(id)
	ctx.JSON(http.StatusCreated, res)
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
	res := h.svc.DeleteProject(id)
	ctx.JSON(http.StatusCreated, res)
}

func home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Antony Injila Portfolio",
	})
}

func InitGinRoutes() {
	router := gin.Default()
	// DynamoDB repository
	repo := repostitory.NewDynaDBRepository()
	// Portifolio service
	svc := services.NewPortfolioService(&repo)
	// Gin Route Handler
	handler := NewGinHandler(*svc)

	// Home route
	router.GET("/", home)
	// Users routes
	usersRoutes := router.Group("/v1/users")
	// Projects Routes
	projectsRoutes := router.Group("/v1/pojects")

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

	router.Run()
}
