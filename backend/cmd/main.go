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
	router.Static("/css", "./frontend/css")                     // Serve CSS from frontend folder
	router.Static("/js", "./frontend/js")                       // Serve JS files from frontend folder
	router.Static("/pages", "./frontend/pages")                 // Serve pages
	router.Static("/uploads", "./uploads")                      // Serve files from uploads folder
	router.StaticFile("/favicon.ico", "./frontend/favicon.ico") // Serve favicon
	router.StaticFile("/", "./frontend/index.html")             // Serve homepage
	router.Static("/images", "./frontend/images")               // Serve images

	// Set up API routes under /api
	api := router.Group("/api")
	routes.SetupRoutes(api)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback for local testing
	}
	router.Run(":" + port)
}
