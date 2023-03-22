package main

import (
	"flag"
	"fmt"

	"github.com/AntonyIS/portfolio-be/config"
	"github.com/AntonyIS/portfolio-be/internal/adapters/http/gin"
	"github.com/AntonyIS/portfolio-be/internal/adapters/repository"
	"github.com/AntonyIS/portfolio-be/internal/core/services"
)

var fg string

func init() {
	config.LoadEnv()
	flag.StringVar(&fg, "flagname", "Production", "The Environment the application is running")
}

func main() {
	// Load application configuration
	fmt.Println(fg)
	config := config.NewConfiguration(fg)
	// DynamoDB repository
	repo := repository.NewDynamoDBRepository(config)
	// Portifolio service
	svc := services.NewPortfolioService(&repo)
	// Gin server
	gin.InitGinRoutes(*svc)
}
