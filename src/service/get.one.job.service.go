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

func GetOneJob(ctx context.Context, id string) (*models.Job, error) {

	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Job")

	// Convert the string ID to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid job ID format: %v", err))
		return nil, fmt.Errorf("invalid job ID format: %w", err)
	}
	// Build the query filter
	filter := bson.M{
		"_id":                      objectID,
		"applicationDeadline.time": bson.M{"$gt": time.Now()},
		"status":                   "Open",
	}

	// Query the collection
	var job models.Job
	err = collection.FindOne(ctx, filter).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info(fmt.Sprintf("No job found with ID: %s", id))
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving job from the database: %v", err))
		return nil, fmt.Errorf("error fetching job details: %w", err)
	}

	logger.Info(fmt.Sprintf("Retrieved job with ID: %s successfully.", id))
	return &job, nil
}
