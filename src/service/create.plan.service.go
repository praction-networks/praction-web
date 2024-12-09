package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePlan(ctx context.Context, plan models.Plan) error {

	plan.UUID = uuid.New().String()
	// Insert the user into the MongoDB database
	err := insertPlanIntoDB(ctx, plan)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert plan into DB: %v", err))
		return fmt.Errorf("failed to create plan: %w", err)
	}

	logger.Info(fmt.Sprintf("Plan %s created successfully.", plan.Category))
	return nil
}

// insertUserIntoDB inserts the new user into the MongoDB database
func insertPlanIntoDB(ctx context.Context, plan models.Plan) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Plan")

	// Insert the user document into the collection
	_, err := collection.InsertOne(ctx, plan)
	if err != nil {

		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log the duplicate key error
					logger.Info(fmt.Sprintf("Duplicate key error: %v", writeErr.Message))
				}

			}
		}

		return fmt.Errorf("error inserting plan into database: %w", err)

	}

	return nil
}
