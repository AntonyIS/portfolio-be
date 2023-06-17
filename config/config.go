package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrNotFound       = errors.New("item not found")
	ErrInvalidItem    = errors.New("invalid item")
	ErrInternalServer = errors.New("internal server error")
)

type AppConfig struct {
	Env          string
	Port         string
	UsersTable   string
	ProjectTable string
	Region       string
	Testing      bool
}

func NewConfiguration(Env string) *AppConfig {

	var (
		serverPort       = os.Getenv("SERVER_PORT")
		region           = os.Getenv("AWS_DEFAULT_REGION")
		userTablename    = ""
		projectTablename = ""
		testing          = false
	)

	switch Env {

	case "dev":
		userTablename = os.Getenv("USERS_TABLE")
		projectTablename = os.Getenv("PROJECTS_TABLE")
		testing = true

	case "pro":
		userTablename = os.Getenv("USERS_TABLE")
		projectTablename = os.Getenv("PROJECTS_TABLE")
		testing = false

	default:
		userTablename = os.Getenv("USERS_TABLE")
		projectTablename = os.Getenv("PROJECTS_TABLE")
	}
	return &AppConfig{
		Env:          Env,
		Port:         serverPort,
		UsersTable:   userTablename,
		ProjectTable: projectTablename,
		Region:       region,
		Testing:      testing,
	}
}

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return nil
}
