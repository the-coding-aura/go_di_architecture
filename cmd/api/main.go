package main

import (
	"go_di_architecture/internal/app/router"

	"github.com/gin-gonic/gin"
)

// @title Module API
// @version 1.0
// @description API for managing module entities in the system
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
//
// @info.title Module API
// @info.description This API provides operations for managing module entities.
// @info.termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
//
// @x-logo {"url": "https://example.com/logo.png", "backgroundColor": "#FFFFFF"}

func main() {
	r := gin.Default()

	// Setup routes
	router.SetupRouter(r)

	// Run the server
	r.Run(":8080")
}
