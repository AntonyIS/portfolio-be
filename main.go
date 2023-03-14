package main

import (
	"github.com/AntonyIS/portfolio-be/config"
	"github.com/AntonyIS/portfolio-be/internal/adapters/http"
)

func init() {
	config.LoadEnv()
}
func main() {
	// Gin server
	http.InitGinRoutes()
}
