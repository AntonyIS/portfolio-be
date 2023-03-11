package main

import (
	"github.com/AntonyIS/portfolio-be/internal/adapters/http"
)

func main() {
	// Gin server
	http.InitGinRoutes()
}
