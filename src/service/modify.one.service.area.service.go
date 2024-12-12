package service

import (
	"context"
	"fmt"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ModifyServiceArea updates the properties and geometry of a specific feature within a service area document.
func ModifyServiceArea(ctx context.Context, objectID string, updatedFeature *models.Feature) error {
	logger.Info("Starting ModifyServiceArea", "objectID", objectID, "area", updatedFeature.Properties.AreaName)

	// MongoDB client and collection
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("ServiceArea")

	// Convert objectID string to ObjectID type
	objID, err := primitive.ObjectIDFromHex(objectID)
	if err != nil {
		logger.Error("Invalid ObjectID format", "objectID", objectID, "error", err)
		return fmt.Errorf("invalid ObjectID: %w", err)
	}

	// Define filter to find the document by ObjectID and the feature by UUID
	filter := bson.M{
		"_id":           objID,
		"features.uuid": updatedFeature.UUID,
	}

	// Fetch the existing feature collection
	var existingFeatureCollection models.FeatureCollection
	err = collection.FindOne(ctx, filter).Decode(&existingFeatureCollection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Warn("No matching feature found for update", "objectID", objectID)
			return fmt.Errorf("no matching feature found with UUID '%s' in document '%s'", updatedFeature.UUID, objectID)
		}
		logger.Error("Error fetching existing feature", "error", err)
		return fmt.Errorf("error fetching existing feature: %w", err)
	}

	// Find the specific feature to update
	var existingProperties *models.FeatureProperties
	var existingGeometry *models.GeoJSONPolygon

	for _, feature := range existingFeatureCollection.Features {
		if feature.UUID == updatedFeature.UUID {
			existingProperties = &feature.Properties
			existingGeometry = &feature.Geometry
			break
		}
	}

	if existingProperties == nil || existingGeometry == nil {
		logger.Warn("Matching feature not found within the document", "uuid", updatedFeature.UUID)
		return fmt.Errorf("no matching feature found with UUID '%s'", updatedFeature.UUID)
	}

	// Prepare the fields to update
	updateFields := bson.M{}

	// Compare and update properties
	if updatedFeature.Properties.AreaName != "" && updatedFeature.Properties.AreaName != existingProperties.AreaName {
		updateFields["features.$.properties.areaName"] = updatedFeature.Properties.AreaName
	}
	if len(updatedFeature.Properties.AvailableService) > 0 && !areStringSlicesEqual(updatedFeature.Properties.AvailableService, existingProperties.AvailableService) {
		updateFields["features.$.properties.availableService"] = updatedFeature.Properties.AvailableService
	}
	if updatedFeature.Properties.Pincode != "" && updatedFeature.Properties.Pincode != existingProperties.Pincode {
		updateFields["features.$.properties.pincode"] = updatedFeature.Properties.Pincode
	}
	if len(updatedFeature.Properties.SubArea) > 0 && !areStringSlicesEqual(updatedFeature.Properties.SubArea, existingProperties.SubArea) {
		updateFields["features.$.properties.subArea"] = updatedFeature.Properties.SubArea
	}

	// Compare and update geometry
	if !isGeometryEqual(updatedFeature.Geometry, *existingGeometry) {
		updateFields["features.$.geometry"] = updatedFeature.Geometry
	}

	// Update only if there are changes
	if len(updateFields) == 0 {
		logger.Info("No updates detected for the feature", "uuid", updatedFeature.UUID)
		return nil
	}

	// Add updatedAt timestamp
	updateFields["updatedAt"] = time.Now()

	// Execute the update operation
	update := bson.M{
		"$set": updateFields,
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Error updating service area", "error", err)
		return fmt.Errorf("error updating service area: %w", err)
	}

	// Check if any document was modified
	if result.MatchedCount == 0 {
		logger.Warn("No document matched for update", "objectID", objectID, "uuid", updatedFeature.UUID)
		return fmt.Errorf("no matching feature found with UUID '%s' in document '%s'", updatedFeature.UUID, objectID)
	}

	logger.Info("Feature updated successfully", "objectID", objectID, "uuid", updatedFeature.UUID)
	return nil
}

// Helper function to compare string slices
func areStringSlicesEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	stringMap := make(map[string]bool, len(slice1))
	for _, val := range slice1 {
		stringMap[val] = true
	}
	for _, val := range slice2 {
		if !stringMap[val] {
			return false
		}
	}
	return true
}

// Helper function to compare geometries
func isGeometryEqual(geometry1, geometry2 models.GeoJSONPolygon) bool {
	if geometry1.Type != geometry2.Type {
		return false
	}
	if len(geometry1.Coordinates) != len(geometry2.Coordinates) {
		return false
	}
	for i := range geometry1.Coordinates {
		if len(geometry1.Coordinates[i]) != len(geometry2.Coordinates[i]) {
			return false
		}
		for j := range geometry1.Coordinates[i] {
			if geometry1.Coordinates[i][j][0] != geometry2.Coordinates[i][j][0] || geometry1.Coordinates[i][j][1] != geometry2.Coordinates[i][j][1] {
				return false
			}
		}
	}
	return true
}
