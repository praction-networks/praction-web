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

func GetOnePlanByUUID(ctx context.Context, planID string) (models.PlanSpecific, error) {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Plan")

	// Query to fetch all plans
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No plans found in the database.")
			return models.PlanSpecific{}, fmt.Errorf("no plans found in database: %w", err)
		}
		logger.Error(fmt.Sprintf("Error retrieving plans from database: %v", err))
		return models.PlanSpecific{}, fmt.Errorf("error fetching plans from database: %w", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error closing cursor: %v", err))
		}
	}()

	// Iterate through all plans and find the matching PlanID in PlanDetail
	for cursor.Next(ctx) {
		var plan models.Plan
		if err = cursor.Decode(&plan); err != nil {
			logger.Error(fmt.Sprintf("Error decoding plan: %v", err))
			return models.PlanSpecific{}, fmt.Errorf("error decoding plan from database: %w", err)
		}

		// Check for the matching PlanID in PlanDetail
		for _, planSpecific := range plan.PlanDetail {
			if planSpecific.PlanID == planID {
				logger.Info(fmt.Sprintf("Found matching PlanSpecific for PlanID: %s", planID))
				return planSpecific, nil
			}
		}
	}

	logger.Info(fmt.Sprintf("No PlanSpecific found with PlanID: %s", planID))
	return models.PlanSpecific{}, fmt.Errorf("planID %s not found", planID)
}
