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

func UpdateUserOTP(ctx context.Context, userIntrest models.UserInterest) error {

	// Check if the user exists in the database
	err := updateOTPInDB(ctx, userIntrest)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to update OTP for User Interest: %v", err))
		return fmt.Errorf("failed to update OTP: %w", err)
	}

	logger.Info(fmt.Sprintf("OTP updated successfully for user %s.", userIntrest.Name))
	return nil
}

// updateOTPInDB updates the OTP and verification status for the user based on mobile or email
func updateOTPInDB(ctx context.Context, userIntrest models.UserInterest) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("UserIntrest")

	// Define the query for checking if the user exists based on email or mobile
	query := bson.M{
		"$or": []bson.M{
			{"email": userIntrest.Email},
			{"mobile": userIntrest.Mobile},
		},
	}

	// Check if the user already exists
	var existingUser models.UserInterest
	err := collection.FindOne(ctx, query).Decode(&existingUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No user found, return an error
			logger.Warn(fmt.Sprintf("No matching user found for mobile %s or email %s", userIntrest.Mobile, userIntrest.Email))
			return fmt.Errorf("no matching user found")
		}
		// Log any other errors
		logger.Error(fmt.Sprintf("Error checking user existence: %v", err))
		return fmt.Errorf("error checking user existence: %w", err)
	}

	// If user exists, update the OTP and set IsVerified as false (if not already verified)
	update := bson.M{
		"$set": bson.M{
			"updatedAt":  time.Now(),
			"IsVarified": false, // Keep it false since this is an OTP update
		},
	}

	// Update the OTP field for the matching user
	_, err = collection.UpdateOne(ctx, bson.M{"uuid": existingUser.UUID}, update)
	if err != nil {
		logger.Error(fmt.Sprintf("Error updating OTP for user %s: %v", userIntrest.Name, err))
		return fmt.Errorf("error updating OTP: %w", err)
	}

	logger.Info(fmt.Sprintf("Updated OTP for user %s successfully.", userIntrest.Name))
	return nil
}

func UpdateUserInterest(ctx context.Context, id string, userInterestUpdate *models.UserInterestUpdate) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("UserInterest")

	// Convert string ID to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Invalid User ID format", "userID", id, "error", err)
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Prepare dynamic update fields
	updateFields := bson.M{}

	// Append comments if provided
	if userInterestUpdate.Comments != nil {
		existingComments, ok := updateFields["comments"].([]string)
		if ok {
			updateFields["comments"] = append(existingComments, userInterestUpdate.Comments...)
		} else {
			updateFields["comments"] = userInterestUpdate.Comments
		}
	}

	// Interest Stage
	if userInterestUpdate.InterestStage != "" {
		updateFields["interestStage"] = userInterestUpdate.InterestStage
	}

	// Follow-up Date
	layout := "2006-01-02"
	// Validate and parse dates before updating
	if userInterestUpdate.FollowUpDate != "" {
		if _, err := time.Parse(layout, userInterestUpdate.FollowUpDate); err != nil {
			logger.Error("Invalid date format for FollowUpDate", "error", err)
			return fmt.Errorf("invalid date format for FollowUpDate, expected YYYY-MM-DD")
		}
		updateFields["followUpDate"] = userInterestUpdate.FollowUpDate
	}

	if userInterestUpdate.PreferredInstallationDate != "" {
		if _, err := time.Parse(layout, userInterestUpdate.PreferredInstallationDate); err != nil {
			logger.Error("Invalid date format for PreferredInstallationDate", "error", err)
			return fmt.Errorf("invalid date format for PreferredInstallationDate, expected YYYY-MM-DD")
		}
		updateFields["preferredInstallationDate"] = userInterestUpdate.PreferredInstallationDate
	}

	// Is Installation Agreed
	if userInterestUpdate.IsInstallationAgreed {
		updateFields["isInstallationAgreed"] = userInterestUpdate.IsInstallationAgreed
	}

	if userInterestUpdate.InstallationDate != "" {
		if _, err := time.Parse(layout, userInterestUpdate.InstallationDate); err != nil {
			logger.Error("Invalid date format for InstallationDate", "error", err)
			return fmt.Errorf("invalid date format for InstallationDate, expected YYYY-MM-DD")
		}
		updateFields["installationDate"] = userInterestUpdate.InstallationDate
	}

	// Installation Status
	if userInterestUpdate.InstallationStatus != "" {
		updateFields["installationStatus"] = userInterestUpdate.InstallationStatus
	}

	// Installation Notes
	if userInterestUpdate.InstallationNotes != nil {
		existingNotes, ok := updateFields["installationNotes"].([]string)
		if ok {
			updateFields["installationNotes"] = append(existingNotes, userInterestUpdate.InstallationNotes...)
		} else {
			updateFields["installationNotes"] = userInterestUpdate.InstallationNotes
		}
	}

	// Selected Plan
	if userInterestUpdate.SelectedPlan != "" {
		updateFields["selectedPlan"] = userInterestUpdate.SelectedPlan
	}

	// Update `UpdatedAt` field
	updateFields["updatedAt"] = time.Now()

	// If no update fields exist, return an error
	if len(updateFields) == 0 {
		logger.Warn("No valid fields provided for update", "userID", id)
		return fmt.Errorf("no valid fields provided for update")
	}

	// Update the document
	update := bson.M{"$set": updateFields}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		logger.Error("Error updating User Interest in the database", "userID", id, "error", err)
		return fmt.Errorf("error updating user interest: %w", err)
	}

	// Check if any document was modified
	if result.MatchedCount == 0 {
		logger.Warn("No User Interest found with the given ID for updating", "userID", id)
		return fmt.Errorf("no user interest found with the given ID: %s", id)
	}

	if result.ModifiedCount == 0 {
		logger.Warn("User Interest update request did not modify any fields", "userID", id)
		return fmt.Errorf("user interest update did not modify any fields")
	}

	logger.Info("User Interest updated successfully", "userID", id)
	return nil
}
