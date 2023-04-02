/*
Package name : services
File name : services.go
Author : Antony Injila
Description :
	- Host code for portfilio logic
*/

package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AntonyIS/portfolio-be/internal/core/domain"
	"github.com/AntonyIS/portfolio-be/internal/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type PortfolioService struct {
	repo ports.PortfolioRepository
}

func NewPortfolioService(repo *ports.PortfolioRepository) *PortfolioService {
	return &PortfolioService{
		repo: *repo,
	}
}

func (svc *PortfolioService) CreateUser(user *domain.User) (*domain.User, error) {
	user.Id = uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	return svc.repo.CreateUser(user)
}

func (svc *PortfolioService) ReadUser(id string) (*domain.User, error) {
	return svc.repo.ReadUser(id)
}

func (svc *PortfolioService) ReadUserWithEmail(email string) (*domain.User, error) {
	return svc.repo.ReadUserWithEmail(email)
}

func (svc *PortfolioService) ReadUsers() ([]*domain.User, error) {
	return svc.repo.ReadUsers()
}

func (svc *PortfolioService) UpdateUser(user *domain.User) (*domain.User, error) {
	return svc.repo.UpdateUser(user)
}

func (svc *PortfolioService) DeleteUser(email string) error {
	user, err := svc.repo.ReadUserWithEmail(email)
	if err != nil {
		return err
	}
	userProjects := user.Projects
	for _, project := range userProjects {
		err = svc.repo.DeleteProject(project.Id)
		if err != nil {
			return err
		}
	}
	return svc.repo.DeleteUser(email)
}

func (svc *PortfolioService) CreateProject(project *domain.Project) (*domain.Project, error) {
	project.Id = uuid.New().String()
	project.CreateAt = time.Now().UTC().Unix()
	userID := project.UserID
	user, err := svc.ReadUser(userID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("User with id %s not found", userID))
	}
	if user.Projects == nil {
		user.Projects = map[string]*domain.Project{
			project.Id: project,
		}
	} else {
		user.Projects[project.Id] = project
	}
	svc.repo.UpdateUser(user)
	return svc.repo.CreateProject(project)
}

func (svc *PortfolioService) ReadProject(id string) (*domain.Project, error) {
	return svc.repo.ReadProject(id)
}

func (svc *PortfolioService) ReadProjects() ([]*domain.Project, error) {
	return svc.repo.ReadProjects()
}

func (svc *PortfolioService) UpdateProject(project *domain.Project) (*domain.Project, error) {
	return svc.repo.UpdateProject(project)
}

func (svc *PortfolioService) DeleteProject(id string) error {
	return svc.repo.DeleteProject(id)
}
