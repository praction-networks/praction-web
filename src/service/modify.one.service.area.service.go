package service

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ModifyServiceArea updates the properties and geometry of a specific feature within a service area document.
func ModifyServiceArea(ctx context.Context, objectID string, updatedFeature *models.UpdateOneArea) error {
	logger.Info("Starting ModifyServiceArea", "objectID", objectID, "area", updatedFeature.UpdateArea.AreaName)

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

	// Prepare the update fields for properties using $mergeObjects
	propertiesUpdate := bson.M{}
	val := reflect.ValueOf(updatedFeature.UpdateArea)
	typ := reflect.TypeOf(updatedFeature.UpdateArea)

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldTag := typ.Field(i).Tag.Get("bson")
		if fieldValue.Kind() == reflect.Slice && fieldValue.Len() == 0 {
			continue // Skip empty slices
		}
		if !reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
			propertiesUpdate[fieldTag] = fieldValue.Interface()
		}
	}

	if len(propertiesUpdate) == 0 {
		logger.Info("No updates detected for properties", "uuid", updatedFeature.UUID)
		return nil
	}

	// Define the MongoDB update operation
	update := bson.M{
		"$set": bson.M{
			"features.$.properties": bson.M{
				"$mergeObjects": bson.M{
					"properties": propertiesUpdate,
				},
			},
			"updatedAt": time.Now(), // Update the FeatureCollection's updatedAt field
		},
	}

	// Execute the update operation
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
