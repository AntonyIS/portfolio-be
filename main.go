package main

import (
	"flag"

	"github.com/AntonyIS/portfolio-be/config"
	"github.com/AntonyIS/portfolio-be/internal/adapters/http/gin"
	"github.com/AntonyIS/portfolio-be/internal/adapters/repository"
	"github.com/AntonyIS/portfolio-be/internal/core/services"
)

var env string

func init() {
	config.LoadEnv()
	flag.StringVar(&env, "Environment", "Production", "The Environment the application is running")
}

func main() {
	// Load application configuration
	config := config.NewConfiguration(env)
	// DynamoDB repository
	repo := repository.NewDynamoDBRepository(config)
	// Portifolio service
	svc := services.NewPortfolioService(&repo)
	// Gin server
	gin.InitGinRoutes(*svc)
}
