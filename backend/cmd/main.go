package main

import (
	"context"
	"log"
	"os"
	"real-estate-app/backend/config"
	"real-estate-app/backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to MongoDB
	client, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Create a new Gin router
	router := gin.Default()

	// Enable CORS for all origins
	router.Use(cors.Default())

	// ✅ Serve frontend static files correctly
	router.Static("/css", "./frontend/css")
	router.Static("/js", "./frontend/js")
	router.Static("/pages", "./frontend/pages")
	router.Static("/images", "./frontend/images")
	router.Static("/uploads", "./uploads")          // Serve uploaded images
	router.StaticFile("/", "./frontend/index.html") // Serve index.html

	// ✅ Set up API routes under /api
	api := router.Group("/api")
	routes.SetupRoutes(api)

	// ✅ Set the port from environment variables (for Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port for local testing
	}
	router.Run(":" + port)
}
