package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUserReferal(ctx context.Context, UserReferal models.UserRefrence) error {

	UserReferal.UUID = uuid.New().String()
	UserReferal.CreatedAt = time.Now()
	UserReferal.OTPExpireTime = time.Now().Add(30 * time.Minute)
	UserReferal.IsVerified = false
	// Insert the user into the MongoDB database
	err := insertUserReferlIntoDB(ctx, UserReferal)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert User Refereal into DB: %v", err))
		return fmt.Errorf("failed to create User refrel: %w", err)
	}
	logger.Info(fmt.Sprintf("Intrest Registred successfully for user %s created successfully.", UserReferal.ReferedBy.Name))
	return nil
}

// insertUserIntoDB inserts the new user into the MongoDB database
func insertUserReferlIntoDB(ctx context.Context, userIntrest models.UserRefrence) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("UserReferal")

	// Insert the user document into the collection
	_, err := collection.InsertOne(ctx, userIntrest)
	if err != nil {

		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log the duplicate key error
					// Identify the field that caused the duplicate key error
					duplicateField := identifyDuplicateField(writeErr.Message)
					if duplicateField != "" {
						// Return a specific error message
						return fmt.Errorf("duplicate %s detected. Please ensure the %s is unique", duplicateField, duplicateField)
					}
				}
			}
		}

		return fmt.Errorf("error inserting User Refrence into database: %w", err)

	}

	return nil
}

// identifyDuplicateField identifies the field (email or mobile) causing the duplicate key error
func identifyDuplicateField(message string) string {
	// Look for the "email" or "mobile" field in the error message
	if ContainsSubstring(message, "email") {
		return "email"
	} else if ContainsSubstring(message, "mobile") {
		return "mobile"
	}
	return ""
}

// containsSubstring checks if the given message contains the specified substring
func ContainsSubstring(message, substring string) bool {
	return len(message) > 0 && len(substring) > 0 && strings.Contains(message, substring)
}
