package main

import (
	"github.com/AntonyIS/portfolio-be/config"
	"github.com/AntonyIS/portfolio-be/internal/adapters/http"
	"github.com/AntonyIS/portfolio-be/internal/adapters/repository"
	"github.com/AntonyIS/portfolio-be/internal/core/services"
)

func init() {
	config.LoadEnv()
}
func main() {
	// Load application configuration
	config := config.NewConfiguration()
	// DynamoDB repository
	repo := repository.NewDynamoDBRepository(config)
	// Portifolio service
	svc := services.NewPortfolioService(&repo)
	// Gin server
	http.InitGinRoutes(*svc)
}
