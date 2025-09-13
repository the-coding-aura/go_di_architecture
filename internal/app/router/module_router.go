package router

import (
	"go_di_architecture/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

// SetupModuleRoutes configures all routes related to module resources.
func SetupModuleRoutes(api *gin.RouterGroup) {
	// Create a dedicated group for module endpoints
	modules := api.Group("/modules")
	{
		handler := handlers.NewModuleHandler()

		// Collection endpoints
		modules.POST("", handler.CreateModule) // POST /api/v1/modules

		// Resource endpoints
		modules.GET("/:id", handler.GetModuleById) // GET /api/v1/modules/{id}
	}
}
