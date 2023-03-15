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

func NewConfiguration() *AppConfig {
	var (
		userTablename    = os.Getenv("USERS_TABLE")
		projectTablename = os.Getenv("PROJECTS_TABLE")
		serverPort       = os.Getenv("SERVER_PORT")
		region           = os.Getenv("AWS_DEFAULT_REGION")
		Env              = os.Getenv("ENV")
	)

	return &AppConfig{
		Env:          Env,
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
