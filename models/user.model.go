package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `
	Name        string             `json:"name" validate:"required,min=3"`
	PhoneNumber string             `json:"phoneNumber" validate:"required,phoneNumber"`
	Email       string             `json:"email" validate:"email,required"`
	Password    string             `json:"password" validate:"email,required"`
}

type UserResponse struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `
	Name        string             `json:"name" validate:"required,min=3"`
	PhoneNumber string             `json:"phoneNumber" validate:"required,phoneNumber"`
	Email       string             `json:"email" validate:"email,required"`
}
