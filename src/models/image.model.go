package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UUID      string             `json:"uuid" bson:"uuid"`
	Name      string             `json:"name" bson:"name" validate:"required"`
	FileID    string             `json:"fileID" bson:"fileID"`
	Tag       string             `json:"tag" bson:"tag" validate:"required,oneof=blog iptv ott"`
	FileName  string             `json:"fileName" bson:"FileName"`
	ImageURL  string             `json:"imageURL" bson:"imageURL"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"-" bson:"updated_at"`
	IsActive  bool               `json:"-" bson:"isActive" validate:"required"`
	IsDeleted bool               `json:"-" bson:"isDeleted"`
}
