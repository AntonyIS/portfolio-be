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
	flag.StringVar(&env, "Environment", "Development", "The environment the application is running")
}

func main() {
	config := config.NewConfiguration(env)
	repo := repository.NewDynamoDBRepository(config)
	svc := services.NewPortfolioService(&repo)
	gin.InitGinRoutes(*svc)
}
