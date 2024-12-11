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

// CheckServiceAvailability checks if a point (latitude, longitude) is within any service area
func CheckServiceAvailability(ctx context.Context, point *models.PointRequest) error {
	// Fetch service area polygons from the database
	polygon, err := GetArrayofCoordinates(ctx)
	if err != nil {
		logger.Error("Error fetching service area polygons", "error", err)
		// Return an internal server error (500) if there's a database issue
		return fmt.Errorf("failed to fetch service area polygons: %w", err)
	}

	// Check if the point is inside any of the polygons
	if PointInPolygon(point.Latitude, point.Longitude, polygon) {
		logger.Info(fmt.Sprintf("%f latitude and %f longitude is in service area", point.Latitude, point.Longitude))
		return nil // The point is within a valid service area
	}

	// If the point is not inside any service area, return an error with 404 or 403
	return fmt.Errorf("the point with latitude %f and longitude %f is not within any service area", point.Latitude, point.Longitude)
}

// GetArrayofCoordinates retrieves all service area polygons from the database
func GetArrayofCoordinates(ctx context.Context) ([][][]float64, error) {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("ServiceArea")

	// Slice to hold all coordinates from the service area features
	var coordinatesList [][][]float64

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

	// Iterate through the cursor to decode each document and extract the coordinates
	for cursor.Next(ctx) {
		var serviceArea models.FeatureCollection
		if err := cursor.Decode(&serviceArea); err != nil {
			logger.Error(fmt.Sprintf("Error decoding document: %v", err))
			return nil, fmt.Errorf("error decoding document: %w", err)
		}

		// Iterate through features to collect coordinates
		for _, feature := range serviceArea.Features {
			coordinates := feature.Geometry.Coordinates

			// Check if coordinates is not empty before appending
			if len(coordinates) > 0 {
				coordinatesList = append(coordinatesList, coordinates...)
			}
		}
	}

	// Check for any errors during iteration
	if err := cursor.Err(); err != nil {
		logger.Error(fmt.Sprintf("Error iterating over cursor: %v", err))
		return nil, fmt.Errorf("error iterating over cursor: %w", err)
	}

	return coordinatesList, nil
}

// PointInPolygon checks if a point (lat, lon) is inside any polygon in a list of polygons
func PointInPolygon(lat, lon float64, polygons [][][]float64) bool {
	// Iterate through each polygon (each polygon is [][]float64)
	for _, polygon := range polygons {
		// Ray-casting algorithm to check if the point is inside the polygon
		oddNodes := false
		n := len(polygon)
		j := n - 1

		// Loop through each edge of the polygon
		for i := 0; i < n; i++ {
			// Check if the point is on the edge of the polygon
			if polygon[i][1] > lat && polygon[j][1] <= lat || polygon[j][1] > lat && polygon[i][1] <= lat {
				if polygon[i][0]+(lat-polygon[i][1])/(polygon[j][1]-polygon[i][1])*(polygon[j][0]-polygon[i][0]) < lon {
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
