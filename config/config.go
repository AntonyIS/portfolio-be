package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Env          string
	Port         string
	UsersTable   string
	ProjectTable string
	Region       string
}

func NewConfiguration(ENV string) *AppConfig {
	var (
		serverPort       = os.Getenv("SERVER_PORT")
		region           = os.Getenv("AWS_DEFAULT_REGION")
		env              = os.Getenv("ENV")
		userTablename    = ""
		projectTablename = ""
	)
	switch ENV {
	case "Development":
		userTablename = os.Getenv("TEST_USERS_TABLE")
		projectTablename = os.Getenv("TEST_PROJECTS_TABLE")

	case "Production":
		userTablename = os.Getenv("USERS_TABLE")
		projectTablename = os.Getenv("PROJECTS_TABLE")
	}

	return &AppConfig{
		Env:          env,
		Port:         serverPort,
		UsersTable:   userTablename,
		ProjectTable: projectTablename,
		Region:       region,
	}
}

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return nil
}
