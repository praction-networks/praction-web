package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// SaveServiceArea saves a new ServiceArea to the database
func CreateServiceArea(ctx context.Context, serviceArea *models.FeatureCollection) error {

	// Add created and updated timestamps
	serviceArea.CreatedAt = time.Now()
	serviceArea.UpdatedAt = time.Now()
	serviceArea.UUID = uuid.New().String()

	// Generate a UUID for each feature in the FeatureCollection
	for i := range serviceArea.Features {
		serviceArea.Features[i].UUID = uuid.New().String()
		logger.Info("Generated UUID for feature", "uuid", serviceArea.Features[i].UUID, "areaName", serviceArea.Features[i].Properties.AreaName)
	}

	// Insert into MongoDB
	err := insertServiceAreaIntoDB(ctx, serviceArea)
	if err != nil {
		logger.Error("Error inserting service area", "error", err)
		return fmt.Errorf("failed to save service area: %w", err)
	}

	// Return the saved service area
	return nil
}

func insertServiceAreaIntoDB(ctx context.Context, area *models.FeatureCollection) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("ServiceArea")

	// Insert the user document into the collection
	_, err := collection.InsertOne(ctx, area)
	if err != nil {

		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log the duplicate key error
					logger.Info(fmt.Sprintf("Duplicate key error: %v", writeErr.Message))
				}

			}
		}

		return fmt.Errorf("error inserting Service Area into Database: %w", err)

	}

	return nil
}
