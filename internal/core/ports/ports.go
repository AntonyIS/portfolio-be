/*
Package name : ports
File name : ports.go
Author : Antony Injila
Description :
	- Host code the describe the purpose of the application
	- Has the Portifolio service and repository interfaces
*/
package ports

import "github.com/AntonyIS/portfolio-be/internal/core/domain"

type PortfolioService interface {
	CreateUser(user *domain.User) (*domain.User, error)
	ReadUser(id string) (*domain.User, error)
	ReadUsers() ([]*domain.User, error)
	UpdateUser(user *domain.User) (*domain.User, error)
	DeleteUser(id string) error
	CreateProject(Project *domain.Project) (*domain.Project, error)
	ReadProject(id string) (*domain.Project, error)
	ReadProjects() ([]*domain.Project, error)
	UpdateProject(Project *domain.Project) (*domain.Project, error)
	DeleteProject(id string) error
}

type PortfolioRepository interface {
	CreateUser(user *domain.User) (*domain.User, error)
	ReadUser(id string) (*domain.User, error)
	ReadUsers() ([]*domain.User, error)
	UpdateUser(user *domain.User) (*domain.User, error)
	DeleteUser(id string) error
	CreateProject(Project *domain.Project) (*domain.Project, error)
	ReadProject(id string) (*domain.Project, error)
	ReadProjects() ([]*domain.Project, error)
	UpdateProject(Project *domain.Project) (*domain.Project, error)
	DeleteProject(id string) error
}
