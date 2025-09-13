package router

import (
	"go_di_architecture/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures the complete routing structure for the application.
func SetupRouter(r *gin.Engine) {
	// Global middleware handlers
	r.Use(middleware.RequestIDHandler())
	r.Use(middleware.ExceptionHandler())
	// r.Use(middleware.LoggingHandler())

	// Versioned API routes
	v1 := r.Group("/api/v1")
	{
		// Module routes
		SetupModuleRoutes(v1)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
