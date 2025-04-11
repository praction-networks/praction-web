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

// CreateUnAvailableAreaRequest creates a new UnAvailableArea request in the database
func CreateUnAvailableAreaRequest(ctx context.Context, unAvailableArea models.UnAvailableArea) error {
	logger.Info("Creating a new UnAvailableArea request")

	// Assign a new UUID if not already set
	if unAvailableArea.UUID == "" {
		unAvailableArea.UUID = uuid.New().String()
	}

	// Set creation and update timestamps
	now := time.Now()
	unAvailableArea.CreatedAt = now
	unAvailableArea.UpdatedAt = now

	// Insert the request into the MongoDB database
	err := insertUnAvailableAreaIntoDB(ctx, unAvailableArea)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert UnAvailableArea request into DB: %v", err))
		return fmt.Errorf("failed to create unavailable area request: %w", err)
	}

	logger.Info(fmt.Sprintf("UnAvailableArea request with UUID %s created successfully.", unAvailableArea.UUID))
	return nil
}

// insertUnAvailableAreaIntoDB inserts the UnAvailableArea request into the MongoDB database
func insertUnAvailableAreaIntoDB(ctx context.Context, unAvailableArea models.UnAvailableArea) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("unavailableAreaRequest")

	// Insert the request document into the collection
	_, err := collection.InsertOne(ctx, unAvailableArea)
	if err != nil {
		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log the duplicate key error
					logger.Info(fmt.Sprintf("Duplicate key error: %v", writeErr.Message))
				}
			}
		}
		return fmt.Errorf("error inserting unavailable area request: %w", err)
	}

	return nil
}
