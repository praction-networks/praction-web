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

func GetAllJobs(ctx context.Context) ([]models.Job, error) {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("Job")

	// Define a slice to store retrieved plans
	var jobs []models.Job

	// Query to fetch all plans
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No jobs found in the database.")
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving jobs from database: %v", err))
		return nil, fmt.Errorf("error fetching jobs: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode all documents into the slice
	if err = cursor.All(ctx, &jobs); err != nil {
		logger.Error(fmt.Sprintf("Error decoding plans: %v", err))
		return nil, fmt.Errorf("error decoding plans: %w", err)
	}

	logger.Info(fmt.Sprintf("Retrieved %d plans successfully.", len(jobs)))
	return jobs, nil
}
