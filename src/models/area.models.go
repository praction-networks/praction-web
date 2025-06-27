package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AreaCategory represents the structure for blog/area categories
type ServiceAreaPage struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AreaName         string             `json:"areaName" bson:"areaName" validate:"required,min=2,max=100"`
	UUID             string             `json:"uuid" bson:"uuid"`
	AreaImage        string             `json:"areaImage" bson:"areaImage" validate:"required,uuid4"`
	OTTs             []string           `json:"otts" bson:"otts" validate:"required,dive,required"`
	AreaBusinessName string             `json:"areaBusinessName" bson:"areaBusinessName" validate:"required"`
	AreaAddress      string             `json:"areaAddress" bson:"areaAddress" validate:"required"`
	AreaURL          string             `json:"areaURL" bson:"areaURL" validate:"required,url"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsActive         bool               `json:"-" bson:"isActive"`
	IsDeleted        bool               `json:"-" bson:"isDeleted"`
}
