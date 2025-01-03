package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Property struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Location     string             `json:"location" bson:"location"`
	Price        float64            `json:"price" bson:"price"`
	PropertyType string             `json:"propertyType" bson:"propertyType"`
	Description  string             `json:"description" bson:"description"`
	Image        string             `json:"image" bson:"image"`
	UserID       string             `json:"user_id" bson:"user_id"`
}
