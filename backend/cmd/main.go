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

	// Serve frontend static files
	router.Static("/css", "./frontend/css")
	router.Static("/js", "./frontend/js")
	router.Static("/pages", "./frontend/pages")
	router.Static("/uploads", "./uploads")
	router.Static("/images", "./frontend/images")
	router.StaticFile("/", "./frontend/index.html")
	router.StaticFile("/favicon.ico", "./frontend/favicon.ico")

	// Set up API routes under /api
	api := router.Group("/api")
	routes.SetupRoutes(api)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port if not specified
	}
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)

}
