package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateJob(ctx context.Context, job models.Job) error {

	job.UUID = uuid.New().String()
	job.JobID = generateJobID()
	// Insert the user into the MongoDB database
	err := insertJobIntoDB(ctx, job)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert Jon into DB: %v", err))
		return fmt.Errorf("failed to create Job: %w", err)
	}

	logger.Info(fmt.Sprintf("Job %s created successfully.", job.Title))
	return nil
}

// insertUserIntoDB inserts the new user into the MongoDB database
func insertJobIntoDB(ctx context.Context, job models.Job) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("Job")

	// Insert the user document into the collection
	_, err := collection.InsertOne(ctx, job)
	if err != nil {

		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log the duplicate key error
					logger.Info(fmt.Sprintf("Duplicate key error: %v", writeErr.Message))
				}

			}
		}

		return fmt.Errorf("error inserting Job into database: %w", err)

	}

	return nil
}

func generateJobID() string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random 6-digit number (between 100000 and 999999)
	randomNumber := rand.Intn(900000) + 100000

	// Format the JobID as "PR" followed by the random number
	jobID := "PR" + strconv.Itoa(randomNumber)

	return jobID
}
