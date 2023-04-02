package services

import (
	"reflect"
	"testing"

	"github.com/AntonyIS/portfolio-be/config"
	"github.com/AntonyIS/portfolio-be/internal/adapters/repository"
	"github.com/AntonyIS/portfolio-be/internal/core/domain"
)

func TestApplicationService(t *testing.T) {
	config := &config.AppConfig{
		Env:          "Dev",
		Port:         "3000",
		UsersTable:   "Authors",
		ProjectTable: "Projects",
		Region:       "me-south-1",
	}
	repo := repository.NewDynamoDBRepository(config)
	svc := NewPortfolioService(&repo)

	t.Run("Test create new user", func(t *testing.T) {
		newUser := domain.User{
			FirstName: "Antony",
			LastName:  "Injila",
			Email:     "antony@gmail.com",
			Password:  "password",
		}
		_, err := svc.ReadUserWithEmail(newUser.Email)

		if err != nil {
			user, err := svc.CreateUser(&newUser)

			if err != nil {
				t.Error(err)
			}
			if user.Email != newUser.Email || user.FirstName != newUser.FirstName || user.LastName != newUser.LastName {
				t.Error("New user does not match created user")
			}
		}

	})
	t.Run("Read user with email", func(t *testing.T) {
		newUser := domain.User{
			FirstName: "Antony",
			LastName:  "Injila",
			Email:     "antony@gmail.com",
			Password:  "password",
		}
		user, err := svc.CreateUser(&newUser)
		if err != nil {
			t.Error(err)
		}
		user, err = svc.ReadUserWithEmail(user.Email)
		if err != nil {
			t.Error(err)
		}
		if user.Email != newUser.Email {
			t.Errorf("User with email %s is not same as %s ", user.Email, newUser.Email)
		}

	})
	t.Run("Read users", func(t *testing.T) {
		users, err := svc.ReadUsers()
		if err != nil {
			t.Error(err)
		}
		u := []*domain.User{}

		if !reflect.DeepEqual(reflect.TypeOf(users), reflect.TypeOf(u)) {
			t.Error(err)
		}
	})
	t.Run("Update user", func(t *testing.T) {
		newUser := domain.User{
			FirstName: "Antony",
			LastName:  "Injila",
			Email:     "antony@gmail.com",
			Password:  "password",
		}
		DBuser, err := svc.CreateUser(&newUser)
		if err != nil {
			t.Error(err)
		}
		DBuser.FirstName = "John"
		DBuser.LastName = "john@gmail.com"

		user, err := svc.UpdateUser(DBuser)
		if err != nil {
			t.Error(err)
		}
		if user.FirstName != DBuser.FirstName || user.LastName != DBuser.LastName {
			t.Error(err)
		}
	})
	t.Run("Delete user", func(t *testing.T) {
		newUser := domain.User{
			FirstName: "Antony",
			LastName:  "Injila",
			Email:     "antony@gmail.com",
			Password:  "password",
		}
		DBuser, err := svc.CreateUser(&newUser)
		if err != nil {
			t.Error(err)
		}

		err = svc.DeleteUser(DBuser.Email)
		if err != nil {
			t.Error(err)
		}

		_, err = svc.ReadUserWithEmail(DBuser.Email)
		if err == nil {
			t.Error(err)
		}

	})
	

	t.Run("Test create new project", func(t *testing.T) {
		newUser := domain.User{
			FirstName: "Antony",
			LastName:  "Injila",
			Email:     "antony@gmail.com",
			Password:  "password",
		}
		user, err := svc.CreateUser(&newUser)

		if err != nil {
			t.Error(err)
		}
		newProject := domain.Project{
			Title:  "Go gRPC for beginners",
			Body:   "This tutorial provides a basic Go programmer’s introduction to working with gRPC.",
			UserID: user.Id,
			Rate:   5,
		}

		project, err := svc.CreateProject(&newProject)
		if err != nil {
			t.Error(err)
		}

		if project.Title != newProject.Title || project.Body != newProject.Body || project.UserID != user.Id {
			t.Error("New user does not match created user")
		}

	})
	t.Run("Read project with id", func(t *testing.T) {
		newUser := domain.User{
			FirstName: "Antony",
			LastName:  "Injila",
			Email:     "antony@gmail.com",
			Password:  "password",
		}
		_, err := svc.ReadUserWithEmail(newUser.Email)

		if err != nil {
			user, err := svc.CreateUser(&newUser)
			if err != nil {
				t.Error(err)
			}
			newProject := domain.Project{
				Title:  "Go gRPC for beginners",
				Body:   "This tutorial provides a basic Go programmer’s introduction to working with gRPC.",
				UserID: user.Id,
				Rate:   5,
			}

			DBproject, err := svc.CreateProject(&newProject)
			if err != nil {
				t.Error(err)
			}

			project, err := svc.ReadProject(DBproject.Id)
			if err != nil {
				t.Error(err)
			}
			if project.Title != DBproject.Title || project.Body != DBproject.Body || project.UserID != user.Id {
				t.Error("New user does not match created user")
			}
		}

	})
	t.Run("Read projects", func(t *testing.T) {
		projects, err := svc.ReadProjects()
		if err != nil {
			t.Error(err)
		}
		p := []*domain.Project{}

		if !reflect.DeepEqual(reflect.TypeOf(projects), reflect.TypeOf(p)) {
			t.Error(err)
		}
	})
	t.Run("Update project", func(t *testing.T) {
		newUser := domain.User{
			FirstName: "Antony",
			LastName:  "Injila",
			Email:     "antony@gmail.com",
			Password:  "password",
		}
		user, err := svc.CreateUser(&newUser)
		if err != nil {
			t.Error(err)
		}
		newProject := domain.Project{
			Title:  "Go gRPC for beginners",
			Body:   "This tutorial provides a basic Go programmer’s introduction to working with gRPC.",
			UserID: user.Id,
			Rate:   5,
		}

		DBproject, err := svc.CreateProject(&newProject)
		if err != nil {
			t.Error(err)
		}
		DBproject.Title = "Master gRPC for beginners"
		DBproject.Body = "This tutorial provides a basic Go programmer’s introduction to working with gRPC.\nOur example is a simple route mapping application that lets clients get information about features on their route, create a summary of their route, and exchange route information such as traffic updates with the server and other clients."

		project, err := svc.UpdateProject(DBproject)
		if project.Title != DBproject.Title || project.Body != DBproject.Body {
			t.Error(err)
		}
	})
	t.Run("Delete Project", func(t *testing.T) {
		newUser := domain.User{
			FirstName: "Antony",
			LastName:  "Injila",
			Email:     "antony@gmail.com",
			Password:  "password",
		}
		user, err := svc.CreateUser(&newUser)
		if err != nil {
			t.Error(err)
		}
		newProject := domain.Project{
			Title:  "Go gRPC for beginners",
			Body:   "This tutorial provides a basic Go programmer’s introduction to working with gRPC.",
			UserID: user.Id,
			Rate:   5,
		}

		DBproject, err := svc.CreateProject(&newProject)
		if err != nil {
			t.Error(err)
		}
		err = svc.DeleteProject(DBproject.Id)
		if err != nil {
			t.Error(err)
		}

	})
	t.Run("Delete all test projects", func(t *testing.T) {
		projects, err := svc.ReadProjects()
		if err != nil {
			t.Error(err)
		}

		for _, project := range projects {

			err := svc.DeleteProject(project.Id)
			if err != nil {
				t.Error(err)
			}
		}
	})

}
