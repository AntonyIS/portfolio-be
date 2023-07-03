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
	Env                string
	Port               string
	UsersTable         string
	ProjectTable       string
	AWSDefaultRegion   string
	AWSAccessKeyID     string
	AWSAccessSecretKey string
	Testing            bool
}

func NewConfiguration(Env string) *AppConfig {

	var (
		serverPort         = os.Getenv("SERVER_PORT")
		AWSDefaultRegion   = os.Getenv("AWS_DEFAULT_REGION")
		AWSAccessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
		AWSAccessSecretKey = os.Getenv("AWS_ACCESS_SECRET_KEY")
		userTablename      = os.Getenv("USER_TABLE")
		projectTablename   = os.Getenv("PROJECT_TABLE")
		testing            = false
	)

	switch Env {

	case "dev":
		testing = true

	case "pro":
		testing = false

	}
	return &AppConfig{
		Env:                Env,
		Port:               serverPort,
		UsersTable:         userTablename,
		ProjectTable:       projectTablename,
		AWSDefaultRegion:   AWSDefaultRegion,
		AWSAccessKeyID:     AWSAccessKeyID,
		AWSAccessSecretKey: AWSAccessSecretKey,
		Testing:            testing,
	}
}

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return nil
}
