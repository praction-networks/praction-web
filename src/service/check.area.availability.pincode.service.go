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

// CheckServiceByPinCode checks service availability by pincode
func CheckServiceByPinCode(ctx context.Context, pincode string) ([]models.FeatureProperties, error) {
	// Fetch service area properties by pincode
	featureProperties, err := GetArrayofPinCode(ctx, pincode)
	if err != nil {
		logger.Error("Error fetching service area properties", "error", err)
		return nil, fmt.Errorf("failed to fetch service area properties: %w", err)
	}

	// Return the matched FeatureProperties
	if len(featureProperties) == 0 {
		logger.Info(fmt.Sprintf("No service areas found for pincode: %s", pincode))
		return nil, fmt.Errorf("no service areas found for pincode %s", pincode)
	}

	return featureProperties, nil
}

// GetArrayofCoordinates retrieves all service area polygons from the database
func GetArrayofPinCode(ctx context.Context, pincode string) ([]models.FeatureProperties, error) {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("ServiceArea")

	// Slice to hold all coordinates from the service area features
	var featurePropertiesList []models.FeatureProperties

	filter := bson.M{
		"features.properties.pincode": pincode,
	}
	// Find all documents in the collection
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No documents found in the database for the given pincode.")
			return nil, fmt.Errorf("no service areas found for pincode %s", pincode)
		}
		logger.Error(fmt.Sprintf("Error retrieving documents from database: %v", err))
		return nil, fmt.Errorf("error fetching documents: %w", err)
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor to decode each document and extract the coordinates
	// Iterate through the cursor to decode each document and extract FeatureProperties
	for cursor.Next(ctx) {
		var serviceArea models.FeatureCollection
		if err := cursor.Decode(&serviceArea); err != nil {
			logger.Error(fmt.Sprintf("Error decoding document: %v", err))
			return nil, fmt.Errorf("error decoding document: %w", err)
		}

		// Iterate through features to collect properties with matching pincode
		for _, feature := range serviceArea.Features {
			if feature.Properties.Pincode == pincode {
				featurePropertiesList = append(featurePropertiesList, feature.Properties)
			}
		}
	}

	// Check for any errors during iteration
	if err := cursor.Err(); err != nil {
		logger.Error(fmt.Sprintf("Error iterating over cursor: %v", err))
		return nil, fmt.Errorf("error iterating over cursor: %w", err)
	}

	// Return the list of FeatureProperties
	return featurePropertiesList, nil
}
