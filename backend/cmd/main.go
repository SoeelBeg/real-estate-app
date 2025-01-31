package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"real-estate-app/backend/config" //with backend

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
	router.Static("../css", "../frontend/css")
	router.Static("../js", "../frontend/js")
	router.Static("../pages", "../frontend/pages")
	router.Static("./uploads", "./uploads") // Serve files from the uploads folder
	router.StaticFile("../", "../frontend/index.html")
	router.StaticFile("/favicon.ico", "./frontend/favicon.ico")
	router.Static("../images", "../frontend/images")
	router.StaticFile("/index.html", "../frontend/index.html")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port agar Render ne assign nahi kiya
	}

	// Server start karein
	fmt.Println("🚀 Server running on port:", port)
	log.Fatal(router.Run(":" + port))

}
