/*
Package name : services
File name : services.go
Author : Antony Injila
Description :
	- Host code for portfilio logic
*/

package services

import (
	"github.com/AntonyIS/portfolio-be/internal/core/domain"
	"github.com/AntonyIS/portfolio-be/internal/core/ports"
)

type PortfolioService struct {
	repo *ports.PortfolioRepository
}

func NewPortfolioRepository(repo *ports.PortfolioRepository) *PortfolioService {
	return &PortfolioService{
		repo: repo,
	}
}

func (svc *PortfolioService) CreateUser(user *domain.User) (*domain.User, error) {
	return nil, nil
}

func (svc *PortfolioService) ReadUser(id string) (*domain.User, error) {
	return nil, nil
}

func (svc *PortfolioService) ReadUsers() ([]*domain.User, error) {
	return nil, nil
}

func (svc *PortfolioService) UpdateUser(user *domain.User) (*domain.User, error) {
	return nil, nil
}
func (svc *PortfolioService) DeleteUser(id string) error {
	return nil
}

func (svc *PortfolioService) CreateProject(Project *domain.Project) (*domain.Project, error) {
	return nil, nil
}

func (svc *PortfolioService) ReadProject(id string) (*domain.Project, error) {
	return nil, nil
}

func (svc *PortfolioService) ReadProjects() ([]*domain.Project, error) {
	return nil, nil
}

func (svc *PortfolioService) UpdateProject(Project *domain.Project) (*domain.Project, error) {
	return nil, nil
}
func (svc *PortfolioService) DeleteProject(id string) error {
	return nil
}
