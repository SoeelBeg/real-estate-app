package routes

import (
	"real-estate-app/backend/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all the API routes
func SetupRoutes(r *gin.RouterGroup) {
	r.GET("/properties", controllers.GetProperties)
	r.POST("/properties", controllers.AddProperty)
	r.PUT("/properties/:id", controllers.UpdateProperty)
	r.DELETE("/properties/:id", controllers.DeleteProperty)
}
