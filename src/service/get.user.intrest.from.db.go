package service

import (
	"context"
	"fmt"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserDetailsFromDB(userDara *models.UserOTPResend) (models.UserInterest, error) {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("UserIntrest")

	// Context with timeout to avoid indefinite hanging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find user by email or mobile
	query := bson.M{
		"$or": []bson.M{
			{"email": userDara.Email},
			{"mobile": userDara.Mobile},
		},
	}

	var existingUser models.UserInterest
	err := collection.FindOne(ctx, query).Decode(&existingUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No matching user found
			logger.Warn("No matching user found in the database", "Query", query)
			return models.UserInterest{}, fmt.Errorf("user not found in database")
		}

		// Handle unexpected MongoDB errors
		logger.Error("Error querying the database", "Error", err)
		return models.UserInterest{}, fmt.Errorf("database error: unable to fetch user data")
	}

	return existingUser, nil

}
