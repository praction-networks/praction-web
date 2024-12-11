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
	Properties FeatureProperties `bson:"properties" json:"properties"`                       // Properties like area_name and available_service
	Geometry   GeoJSONPolygon    `bson:"geometry" json:"geometry" validate:"required"`       // GeoJSON Polygon (required)
}

// FeatureProperties represents the properties of a feature, like area_name and available_service
type FeatureProperties struct {
	AreaName         string `bson:"areaName" json:"areaName" validate:"required"`                 // Area name (required)
	AvailableService string `bson:"availableService" json:"availableService" validate:"required"` // Available services (required)
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

// PointRequest represents the request body for validating if a point (latitude, longitude) is inside a service area.
type PointRequest struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}
