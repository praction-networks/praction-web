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
	imageCollection := client.Database("practionweb").Collection("Image")

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

	// Enrich plans with OTT details
	for i := range plans {
		for j := range plans[i].PlanDetail {
			if len(plans[i].PlanDetail[j].OTTs) > 0 {
				// Fetch related OTT details from the Image collection
				var ottDetails []models.Image
				filter := bson.M{"_id": bson.M{"$in": plans[i].PlanDetail[j].OTTs}}
				ottCursor, err := imageCollection.Find(ctx, filter)
				if err != nil {
					logger.Error(fmt.Sprintf("Error retrieving OTTs for plan detail %v: %v", plans[i].PlanDetail[j].PlanID, err))
					return nil, fmt.Errorf("error retrieving OTTs for plan detail %v: %w", plans[i].PlanDetail[j].PlanID, err)
				}
				defer ottCursor.Close(ctx)

				if err := ottCursor.All(ctx, &ottDetails); err != nil {
					logger.Error(fmt.Sprintf("Error decoding OTTs for plan detail %v: %v", plans[i].PlanDetail[j].PlanID, err))
					return nil, fmt.Errorf("error decoding OTTs for plan detail %v: %w", plans[i].PlanDetail[j].PlanID, err)
				}

				// Attach OTT details to the PlanSpecific item
				plans[i].PlanDetail[j].OttDetails = ottDetails

			}
			plans[i].PlanDetail[j].OTTs = nil
		}
	}

	logger.Info(fmt.Sprintf("Retrieved %d plans successfully.", len(plans)))
	return plans, nil
}
