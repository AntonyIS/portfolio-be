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
		userTablename    = ""
		projectTablename = ""
	)
	switch ENV {
	case "dev":
		userTablename = os.Getenv("TEST_USERS_TABLE")
		projectTablename = os.Getenv("TEST_PROJECTS_TABLE")

	case "pro":
		userTablename = os.Getenv("USERS_TABLE")
		projectTablename = os.Getenv("PROJECTS_TABLE")
	}

	return &AppConfig{
		Env:          ENV,
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
