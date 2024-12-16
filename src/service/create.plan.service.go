package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// Define custom errors for specific failure cases
var (
	ErrDatabaseInsert   = errors.New("failed to insert into database")
	ErrDatabaseInternal = errors.New("database internal error")
)

func CreatePlan(ctx context.Context, plan models.Plan) error {
	// Generate UUID for the plan
	plan.UUID = uuid.New().String()

	// Generate unique PlanID for each PlanSpecific in PlanDetail
	for i := range plan.PlanDetail {
		plan.PlanDetail[i].PlanID = uuid.New().String()
	}

	// Insert the plan into the MongoDB database
	err := insertPlanIntoDB(ctx, plan)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create plan [Category: %s]: %v", plan.Category, err))
		// Return detailed error for the handler
		return fmt.Errorf("failed to create plan: %w", err)
	}

	logger.Info(fmt.Sprintf("Plan [Category: %s, UUID: %s] created successfully.", plan.Category, plan.UUID))
	return nil
}

func insertPlanIntoDB(ctx context.Context, plan models.Plan) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Plan")

	// Insert the plan document into the collection
	_, err := collection.InsertOne(ctx, plan)
	if err != nil {
		// Handle MongoDB write exceptions
		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log duplicate key error and return a specific error
					logger.Warn(fmt.Sprintf("Duplicate key error: %v", writeErr.Message))
					return fmt.Errorf("%w: %v", ErrDuplicateKey, writeErr.Message)
				}
			}
		}

		// Handle other MongoDB errors
		if errors.Is(err, mongo.ErrClientDisconnected) {
			logger.Error("Database client disconnected")
			return fmt.Errorf("%w: client disconnected", ErrDatabaseInternal)
		}

		logger.Error(fmt.Sprintf("Error inserting plan into database: %v", err))
		return fmt.Errorf("%w: %v", ErrDatabaseInsert, err)
	}

	logger.Info(fmt.Sprintf("Plan [UUID: %s] inserted into database successfully.", plan.UUID))
	return nil
}
