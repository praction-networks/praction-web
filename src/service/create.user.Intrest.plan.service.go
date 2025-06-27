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

func CreateUserIntrestForPlan(ctx context.Context, userIntrestPlan models.AvailableUserRequest) error {

	userIntrestPlan.UUID = uuid.New().String()
	userIntrestPlan.CreatedAt = time.Now()
	userIntrestPlan.OTPExpireTime = time.Now().Add(30 * time.Minute)
	userIntrestPlan.IsVerified = false
	// Insert the user into the MongoDB database
	err := insertUserPlanIntrestIntoDB(ctx, userIntrestPlan)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert User Intrst into DB: %v", err))
		return fmt.Errorf("failed to create plan: %w", err)
	}

	logger.Info(fmt.Sprintf("Intrest Registred successfully for user %s created successfully.", userIntrestPlan.FirstName))
	return nil
}

// insertUserIntoDB inserts the new user into the MongoDB database
func insertUserPlanIntrestIntoDB(ctx context.Context, userIntrest models.AvailableUserRequest) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("UserPlanIntrest")

	// Insert the user document into the collection
	_, err := collection.InsertOne(ctx, userIntrest)
	if err != nil {

		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log the duplicate key error
					logger.Info(fmt.Sprintf("Duplicate key error: %v", writeErr.Message))
				}

			}
		}

		return fmt.Errorf("error inserting user data into User Plan Intrest: %w", err)

	}

	return nil
}
