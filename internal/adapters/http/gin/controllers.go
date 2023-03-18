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
	"net/http"
	"os"
	"time"

	"github.com/AntonyIS/portfolio-be/internal/adapters/middleware"
	"github.com/AntonyIS/portfolio-be/internal/core/domain"
	"github.com/AntonyIS/portfolio-be/internal/core/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	Home(ctx *gin.Context)
	Authorize(ctx *gin.Context)
	Login(ctx *gin.Context)
	Signup(ctx *gin.Context)
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

func (h handler) Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Antony Injila Portfolio",
	})
}

func (h handler) Login(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	dbUser, err := h.svc.ReadUser(user.Email)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid email or password",
		})
	}

	if dbUser.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	if dbUser.CheckPasswordHarsh(dbUser.Password) {

		tokenString, err := middleware.GenerateToken(user.Email)

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		}
		ctx.SetSameSite(http.SameSiteLaxMode)
		ctx.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
		ctx.JSON(http.StatusOK, gin.H{
			"user": dbUser,
		})

	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid email or password",
		})
	}
}

func (h handler) Signup(ctx *gin.Context) {

	var user domain.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	dbUser, err := h.svc.ReadUser(user.Email)

	if dbUser.Email != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	password, err := user.GenerateHashPassord()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to harsh password",
		})
		return
	}

	user.Password = password
	newUser, err := h.svc.CreateUser(&user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}
	ctx.JSON(http.StatusOK, newUser)
}

func (h handler) Authorize(ctx *gin.Context) {
	tokenString := ctx.Query("token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userEmail := fmt.Sprintf("%s", claims["email"])
		user, err := h.svc.ReadUser(userEmail)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if user.Email == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
