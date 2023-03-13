package main

import (
	config "github.com/AntonyIS/portfolio-be/config"
	"github.com/AntonyIS/portfolio-be/internal/adapters/http"
)

var (
	port = "8080"
)

func main() {
	config := config.Config(port)
	// Gin server
	http.InitGinRoutes(config.Port)
}
