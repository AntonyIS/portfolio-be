package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AntonyIS/portfolio-be/config"
	"github.com/AntonyIS/portfolio-be/internal/adapters/repository"
	"github.com/AntonyIS/portfolio-be/internal/core/domain"
	"github.com/AntonyIS/portfolio-be/internal/core/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
func TestApplicationRoutes(t *testing.T) {
	config := &config.AppConfig{
		Env:              "Dev",
		Port:             "3000",
		UsersTable:       "Authors",
		ProjectTable:     "Projects",
		AWSDefaultRegion: "me-south-1",
	}
	repo := repository.NewDynamoDBRepository(config)
	svc := services.NewPortfolioService(&repo)
	handler := NewGinHandler(*svc)
	t.Run("Gin Post user", func(t *testing.T) {
		r := SetUpRouter()
		r.POST("/v1/users", handler.PostUser)

		newUser := domain.User{
			FirstName: "Marco",
			LastName:  "Injila",
			Email:     "marco@gmail.com",
			Title:     "Golang Software Engineer",
			Password:  "password",
			Projects:  nil,
		}
		jsonValue, _ := json.Marshal(newUser)
		req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonValue))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	t.Run("Gin Read all user", func(t *testing.T) {
		r := SetUpRouter()
		r.GET("/v1/users", handler.GetUsers)
		req, _ := http.NewRequest("GET", "/v1/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Gin Read user with email", func(t *testing.T) {
		r := SetUpRouter()
		r.GET("/v1/users/:email", handler.GetUserWithEmail)
		req, _ := http.NewRequest("GET", "/v1/users/marco@gmail.com", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var user domain.User

		json.Unmarshal(w.Body.Bytes(), &user)
		if user.Email != "marco@gmail.com" {
			t.Errorf("User with email %s does not match email marco@gmail.com", user.Email)
		}
		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Gin Read user with non existing id", func(t *testing.T) {
		r := SetUpRouter()
		r.GET("/v1/users/:id", handler.GetUserWithID)
		req, _ := http.NewRequest("GET", "/v1/users/1234ewe", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		fmt.Println(http.StatusNotFound, w.Code)
		assert.Equal(t, http.StatusNotFound, w.Code)

	})

}
