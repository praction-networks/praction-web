package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckServiceAvailability(ctx context.Context, point *models.PointRequest) ([]models.FeatureProperties, error) {
	// Fetch service area polygons and their properties from the database
	features, err := GetArrayOfFeatures(ctx)
	if err != nil {
		if err.Error() == "no service areas found" {
			logger.Info("No service areas found in the database")
			return nil, fmt.Errorf("no service areas found")
		}

		logger.Error("Error fetching service area features", "error", err)
		return nil, fmt.Errorf("failed to fetch service area features: %w", err)
	}
	// Slice to hold the matching feature properties
	var matchingFeatures []models.FeatureProperties

	// Check if the point is inside any of the polygons
	for _, feature := range features {
		if PointInPolygon(point.Latitude, point.Longitude, feature.Geometry.Coordinates) {
			logger.Info(fmt.Sprintf("%f latitude and %f longitude is in service area", point.Latitude, point.Longitude))
			matchingFeatures = append(matchingFeatures, feature.Properties)
		}
	}

	// If no matching feature is found, return an error
	if len(matchingFeatures) == 0 {
		return nil, fmt.Errorf("the point with latitude %f and longitude %f is not within any service area", point.Latitude, point.Longitude)
	}

	// Return the matching feature properties
	return matchingFeatures, nil
}

func GetArrayOfFeatures(ctx context.Context) ([]models.Feature, error) {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("ServiceArea")

	// Slice to hold the features
	var featureCollections []models.FeatureCollection

	// Find all documents in the collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No documents found in the database.")
			return nil, fmt.Errorf("no service areas found")
		}
		logger.Error(fmt.Sprintf("Error retrieving documents from database: %v", err))
		return nil, fmt.Errorf("error fetching documents: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode each document into FeatureCollection
	for cursor.Next(ctx) {
		var featureCollection models.FeatureCollection
		if err := cursor.Decode(&featureCollection); err != nil {
			logger.Error("Error decoding document", "error", err)
			return nil, fmt.Errorf("error decoding document: %w", err)
		}
		featureCollections = append(featureCollections, featureCollection)
	}

	// Extract all features from collections
	var features []models.Feature
	for _, collection := range featureCollections {
		features = append(features, collection.Features...)
	}

	if len(features) == 0 {
		logger.Info("The collection is empty; no service areas found.")
		return nil, fmt.Errorf("no service areas found")
	}

	return features, nil

}

// PointInPolygon checks if a point (lat, lon) is inside a polygon
func PointInPolygon(lat, lon float64, polygons [][][]float64) bool {
	// Iterate through each polygon (each polygon is [][]float64)
	for _, polygon := range polygons {
		// Ray-casting algorithm to check if the point is inside the polygon
		oddNodes := false
		n := len(polygon)
		j := n - 1

		// Loop through each edge of the polygon
		for i := 0; i < n; i++ {
			if polygon[i][1] > lat != (polygon[j][1] > lat) {
				if lon < (polygon[j][0]-polygon[i][0])*(lat-polygon[i][1])/(polygon[j][1]-polygon[i][1])+polygon[i][0] {
					oddNodes = !oddNodes
				}
			}
			j = i
		}

		// If the point is inside this polygon, return true
		if oddNodes {
			return true
		}
	}

	// If the point is not inside any polygon, return false
	return false
}
