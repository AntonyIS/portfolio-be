package services

import (
	"testing"

	"github.com/AntonyIS/portfolio-be/config"
	"github.com/AntonyIS/portfolio-be/internal/adapters/repository"
	"github.com/AntonyIS/portfolio-be/internal/core/domain"
)

func TestApplicationService(t *testing.T) {
	config.LoadEnv()
	env := "Production"
	config := config.NewConfiguration(env)

	repo := repository.NewDynamoDBRepository(config)
	svc := NewPortfolioService(&repo)

	t.Run("Test create new user", func(t *testing.T) {
		newUser := domain.User{
			FirstName: "Antony",
			LastName:  "Injila",
			Email:     "antonyshikubu@gmail.com",
			Password:  "SuperPass@#@#$",
		}
		_, err := svc.repo.CreateUser(&newUser)
		if err != nil {
			t.Error(err)
		}

	})
}
