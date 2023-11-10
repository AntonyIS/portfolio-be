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
	// Check if user already exist in the database
	users, err := svc.repo.ReadUsers()
	if err != nil {
		return nil, err
	}
	// Loop through current users and search for user with email
	for _, item := range users {
		if item.Email == user.Email {
			// User found, return error message
			return nil, errors.New("user with email exists!")
		}
	}
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
	users, err := svc.repo.ReadUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user with id not found!")
}

func (svc *PortfolioService) ReadUsers() ([]*domain.User, error) {
	return svc.repo.ReadUsers()
}

func (svc *PortfolioService) UpdateUser(user *domain.User) (*domain.User, error) {
	return svc.repo.UpdateUser(user)
}

func (svc *PortfolioService) DeleteUser(id string) error {
	// Check if user exists
	_, err := svc.ReadUser(id)
	if err != nil {
		return err
	}

	return svc.repo.DeleteUser(id)
}

func (svc *PortfolioService) CreateProject(project *domain.Project) (*domain.Project, error) {
	// Create Project ID
	project.Id = uuid.New().String()
	// Create project created at timestamp
	project.CreateAt = time.Now().UTC().Unix()
	// Get the id of the user
	userID := project.UserID

	// Get user with id
	user, err := svc.ReadUser(userID)
	if err != nil {
		return nil, err
	}

	project.UserName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	project.UserTitle = fmt.Sprintf("%s ", user.Title)
	// Add the new project to user
	user.Projects = append(user.Projects, project)
	// Save the new changes for user
	svc.repo.UpdateUser(user)
	// Add the new project into the database
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
	// Get the project to delete
	project, err := svc.repo.ReadProject(id)

	if err != nil {
		return err
	}

	// Delete project
	err = svc.repo.DeleteProject(id)
	if err != nil {
		return err
	}

	// Get the user ID from the project
	userID := project.UserID
	// Get user with the userID
	user, err := svc.repo.ReadUser(userID)

	if err != nil {
		return err
	}
	// Loop through current project of the user and get the project to delete
	for index, item := range user.Projects {
		if item.Id == id {
			// Delete project from user projects list
			user.Projects = append(user.Projects[:index], user.Projects[index+1:]...)
			// Update the user
			_, err := svc.repo.UpdateUser(user)

			if err != nil {
				return err
			}
			// Delete the project from the database
			return svc.repo.DeleteProject(id)
		}
	}

	return errors.New("Internal server error: Unable to delete project")
}
