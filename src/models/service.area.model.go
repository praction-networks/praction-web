package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GeoJSONPolygon represents a GeoJSON Polygon
type GeoJSONPolygon struct {
	Type        string        `bson:"type" json:"type" validate:"required,oneof=Polygon"`                  // "Polygon" (required)
	Coordinates [][][]float64 `bson:"coordinates" json:"coordinates" validate:"required,coordinatesRange"` // [[[lng, lat], ...]] (required)
}

// Feature represents a GeoJSON Feature, which could either be a Polygon or Point
type Feature struct {
	Type       string            `bson:"type" json:"type" validate:"required,oneof=Feature"` // "Feature" (required)
	UUID       string            `bson:"uuid" json:"uuid"`
	Properties FeatureProperties `bson:"properties" json:"properties"`                 // Properties like area_name and available_service
	Geometry   GeoJSONPolygon    `bson:"geometry" json:"geometry" validate:"required"` // GeoJSON Polygon (required)
}

// FeatureProperties represents the properties of a feature, including additional location details
type FeatureProperties struct {
	AreaName         string   `bson:"areaName" json:"areaName" validate:"required"`                                // Area name (required)
	AvailableService []string `bson:"availableService" json:"availableService" validate:"required,min=1,dive"`     // At least one service required           // Available services (required)
	Pincode          string   `bson:"pincode" json:"pincode" validate:"required,len=6"`                            // Optional pincode, expected to be 6 digits
	SubArea          []string `bson:"subArea" json:"subArea" validate:"required,min=1,dive"`                       // Optional sub-area or locality name
	Zone             string   `bson:"zone,omitempty" json:"zone,omitempty" validate:"omitempty"`                   // Optional zone or region
	Landmark         string   `bson:"landmark,omitempty" json:"landmark,omitempty" validate:"omitempty"`           // Optional nearby landmark
	Taluk            string   `bson:"taluk,omitempty" json:"taluk,omitempty" validate:"omitempty"`                 // Optional Taluk or sub-district
	Division         string   `bson:"division,omitempty" json:"division,omitempty" validate:"omitempty"`           // Optional division within a district or state
	District         string   `bson:"district,omitempty" json:"district,omitempty" validate:"omitempty"`           // Optional district
	Region           string   `bson:"region,omitempty" json:"region,omitempty" validate:"omitempty"`               // Optional larger region
	Circle           string   `bson:"circle,omitempty" json:"circle,omitempty" validate:"omitempty"`               // Optional administrative circle
	State            string   `bson:"state,omitempty" json:"state,omitempty" validate:"omitempty"`                 // Optional state or province
	Country          string   `bson:"country,omitempty" json:"country,omitempty" validate:"omitempty"`             // Optional country
	Stroke           string   `bson:"stroke,omitempty" json:"stroke,omitempty" validate:"omitempty"`               // Stroke color for visualization
	StrokeWidth      int      `bson:"strokeWidth,omitempty" json:"strokeWidth,omitempty" validate:"omitempty"`     // Stroke width for visualization
	StrokeOpacity    float64  `bson:"strokeOpacity,omitempty" json:"strokeOpacity,omitempty" validate:"omitempty"` // Stroke opacity for visualization
	Fill             string   `bson:"fill,omitempty" json:"fill,omitempty" validate:"omitempty"`                   // Fill color for visualization
	FillOpacity      float64  `bson:"fillOpacity,omitempty" json:"fillOpacity,omitempty" validate:"omitempty"`     // Fill opacity for visualization
	Notes            string   `bson:"notes,omitempty" json:"notes,omitempty" validate:"omitempty,max=500"`         // Optional additional notes, up to 500 characters
}

// FeatureCollection represents a collection of features (GeoJSON FeatureCollection)
type FeatureCollection struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UUID      string             `bson:"uuid" json:"uuid"`
	Type      string             `bson:"type" json:"type" validate:"required,oneof=FeatureCollection"` // "FeatureCollection" (required)
	Features  []Feature          `bson:"features" json:"features" validate:"required,dive"`            // List of features (required)
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`                                   // Creation time
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`                                   // Last updated time
}

type UpdateFeture struct {
	AddArea    []Feature `bson:"addArea" json:"addArea" validate:"dive"`
	RemoveArea []string  `bson:"removeArea" json:"removeArea" validate:"dive"`
}

type UpdateOneArea struct {
	UUID       string            `bson:"-" json:"uuid" validate:"required,uuid4"`
	UpdateArea FeatureProperties `bson:"-" json:"updateArea" validate:"required"`
}
type PointRequest struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type ServiceCheck struct {
	Pincode string `bson:"pincode,omitempty" json:"pincode,omitempty" validate:"len=6"`
}
