package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllPlans(ctx context.Context) ([]models.Plan, error) {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Plan")

	// Define a slice to store retrieved plans
	var plans []models.Plan

	// Query to fetch all plans
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No plans found in the database.")
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving plans from database: %v", err))
		return nil, fmt.Errorf("error fetching plans: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode all documents into the slice
	if err = cursor.All(ctx, &plans); err != nil {
		logger.Error(fmt.Sprintf("Error decoding plans: %v", err))
		return nil, fmt.Errorf("error decoding plans: %w", err)
	}

	logger.Info(fmt.Sprintf("Retrieved %d plans successfully.", len(plans)))
	return plans, nil
}
