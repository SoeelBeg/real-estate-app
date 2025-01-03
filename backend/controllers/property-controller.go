package controllers

import (
	"context"
	"fmt"
	"net/http"
	"real-estate-app/backend/config"
	"real-estate-app/backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetProperties retrieves all properties from the database.
func GetProperties(c *gin.Context) {
	collection := config.GetCollection("realestate", "properties")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var properties []models.Property
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &properties); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, properties)
}

func AddProperty(c *gin.Context) {
	// Form data parsing for non-file fields
	title := c.DefaultPostForm("title", "")
	location := c.DefaultPostForm("location", "")
	priceStr := c.DefaultPostForm("price", "") // price as string
	propertyType := c.DefaultPostForm("propertyType", "")
	description := c.DefaultPostForm("description", "")
	userID := c.DefaultPostForm("user_id", "") // Add user_id here

	// Convert price to float64
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
		return
	}

	// Handle image file upload
	file, _ := c.FormFile("image")
	if file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image uploaded"})
		return
	}

	// Save the file to the "uploads" directory
	filePath := fmt.Sprintf("/uploads/%s", file.Filename) // Updated path for public access
	if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Ensure user_id is not empty
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Create property object
	property := models.Property{
		ID:           primitive.NewObjectID(),
		Title:        title,
		Location:     location,
		Price:        price,
		PropertyType: propertyType,
		Description:  description,
		Image:        filePath,
		UserID:       userID, // Here we pass user_id
	}

	collection := config.GetCollection("realestate", "properties")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, property)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add property"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Property added successfully", "property": property})
}

// delete property ----------
func DeleteProperty(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	collection := config.GetCollection("realestate", "properties")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete property"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Property deleted successfully"})
}

// UpdateProperty updates the details of a property
func UpdateProperty(c *gin.Context) {
	id := c.Param("id") // Retrieve the property ID from URL parameters
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID format"})
		return
	}

	collection := config.GetCollection("realestate", "properties")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingProperty models.Property
	// Fetch the existing property details
	if err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingProperty); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Fetch the new values from the form
	title := c.DefaultPostForm("title", existingProperty.Title)
	location := c.DefaultPostForm("location", existingProperty.Location)
	priceStr := c.DefaultPostForm("price", fmt.Sprintf("%f", existingProperty.Price))
	propertyType := c.DefaultPostForm("propertyType", existingProperty.PropertyType)
	description := c.DefaultPostForm("description", existingProperty.Description)

	// Convert price to float
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
		return
	}

	// Handle the image update
	imagePath := existingProperty.Image
	file, _ := c.FormFile("image")
	if file != nil { // If a new file is uploaded
		filePath := fmt.Sprintf("/uploads/%s", file.Filename)
		if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		imagePath = filePath
	}

	// Create the updated fields
	updatedFields := bson.M{
		"title":        title,
		"location":     location,
		"price":        price,
		"propertyType": propertyType,
		"description":  description,
		"image":        imagePath,
	}

	// Update the property in the database
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updatedFields})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update property"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Property updated successfully"})
}

// Helper function to get a property by its ID
func getPropertyByID(objectID primitive.ObjectID) (*models.Property, error) {
	collection := config.GetCollection("realestate", "properties")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var property models.Property
	err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&property)
	if err != nil {
		return nil, err
	}

	return &property, nil
}
