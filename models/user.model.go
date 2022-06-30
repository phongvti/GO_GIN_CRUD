package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct{
	Id primitive.ObjectID 	`json:"id,omitempty" bson:"id,omitempty"`
	Name string 			`json:"name,omitempty" validate:"required" bson:"name,omitempty"`
	Location string 		`json:"location,omitempty" validate:"required" bson:"location,omitempty"`
	Title string 			`json:"title,omitempty" validate:"required" bson:"title,omitempty"`
}