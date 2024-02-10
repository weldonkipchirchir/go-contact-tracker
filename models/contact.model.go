package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `
	FirstName string             `json:"first_name" valiadte:"required,min=3"`
	LastName  string             `json:"last_name" validate:"required,min=3"`
	Twitter   string             `json:"twitter" validate:"required,min=3"`
	AvatarUrl string             `json:"avatar_url"`
	Notes     string             `json:"notes"`
}
