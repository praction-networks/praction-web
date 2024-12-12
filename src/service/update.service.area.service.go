package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateAreaService(ctx context.Context, updateArea *models.UpdateFeture) error {

	logger.Info("Starting UpdateAreaService", "updateArea", updateArea)

	// MongoDB client and collection
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("ServiceArea")

	// Fetch the current document to validate duplicates
	var existingFeatureCollection models.FeatureCollection
	err := collection.FindOne(ctx, bson.M{}).Decode(&existingFeatureCollection)
	if err != nil && err != mongo.ErrNoDocuments {
		logger.Error("Error fetching existing features", "error", err)
		return fmt.Errorf("error fetching existing features: %w", err)
	}

	existingAreaNames := make(map[string]bool)
	for _, feature := range existingFeatureCollection.Features {
		existingAreaNames[feature.Properties.AreaName] = true
	}

	// Remove areas by names if specified
	if len(updateArea.RemoveArea) > 0 {
		logger.Info("Attempting to remove areas", "removeAreaCount", len(updateArea.RemoveArea))
		for _, areaName := range updateArea.RemoveArea {
			logger.Info("Attempting to remove area", "areaName", areaName)

			filter := bson.M{}
			update := bson.M{
				"$pull": bson.M{
					"features": bson.M{"properties.areaName": areaName},
				},
				"$setOnInsert": bson.M{
					"type":      "FeatureCollection",
					"updatedAt": time.Now(),
				},
			}
			result, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				logger.Error("Error removing area", "areaName", areaName, "error", err)
				return fmt.Errorf("failed to remove area '%s': %w", areaName, err)
			}
			if result.ModifiedCount == 0 {
				logger.Warn("No matching area found for removal", "areaName", areaName)
			}
			logger.Info("Area successfully removed", "areaName", areaName)
		}
	}

	// Add new areas if specified
	if len(updateArea.AddArea) > 0 {
		logger.Info("Attempting to add areas", "addAreaCount", len(updateArea.AddArea))

		// Validate and insert each feature
		for i := range updateArea.AddArea {
			feature := &updateArea.AddArea[i]
			// Generate UUID for the feature if not already set

			if existingAreaNames[feature.Properties.AreaName] {
				logger.Error("Duplicate areaName found", "areaName", feature.Properties.AreaName)
				return fmt.Errorf("duplicate areaName '%s' found", feature.Properties.AreaName)
			}
			if feature.UUID == "" {
				feature.UUID = uuid.New().String()
			} // Generate UUID for the feature

			logger.Info("Adding area", "areaName", feature.Properties.AreaName, "uuid", feature.UUID)

			filter := bson.M{}
			update := bson.M{
				"$push": bson.M{
					"features": feature,
				},
				"$setOnInsert": bson.M{
					"type":      "FeatureCollection",
					"updatedAt": time.Now(),
				},
			}
			_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
			if err != nil {
				if mongoErr, ok := err.(mongo.WriteException); ok {
					for _, writeErr := range mongoErr.WriteErrors {
						if writeErr.Code == 11000 { // Duplicate Key Error Code
							logger.Error("Duplicate feature found", "areaName", feature.Properties.AreaName, "uuid", feature.UUID)
							return fmt.Errorf("duplicate feature with areaName '%s' and UUID '%s'", feature.Properties.AreaName, feature.UUID)
						}
					}
				}
				logger.Error("Error adding area", "areaName", feature.Properties.AreaName, "error", err)
				return fmt.Errorf("failed to add area '%s': %w", feature.Properties.AreaName, err)
			}

			// Add the new areaName to the map to track during the same operation
			existingAreaNames[feature.Properties.AreaName] = true
			logger.Info("Area successfully added", "areaName", feature.Properties.AreaName)
		}
	}

	logger.Info("UpdateAreaService completed successfully")
	return nil
}
