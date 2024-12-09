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

func CreateGoogleImage(ImageName, ImageID, ImageURL, mimeType string) (models.Image, error) {

	var googleImage models.Image
	googleImage.UUID = uuid.New().String()
	googleImage.FileName = ImageName
	googleImage.FileID = ImageID
	googleImage.ImageURL = ImageURL
	googleImage.MimeType = mimeType
	googleImage.CreatedAt = time.Now()
	googleImage.UpdatedAt = time.Now()
	googleImage.IsActive = true
	googleImage.IsDeleted = false

	ctx := context.TODO()
	err := insertGoogleImageIntoDB(ctx, googleImage)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert Google Image Metadata into DB: %v", err))
		return models.Image{}, fmt.Errorf("failed to create plan: %w", err)
	}

	logger.Info(fmt.Sprintf("Google Image Metadata with uuid %s created successfully.", googleImage.UUID))
	return googleImage, nil
}

// insertUserIntoDB inserts the new user into the MongoDB database
func insertGoogleImageIntoDB(ctx context.Context, image models.Image) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("BlogImage")

	// Insert the user document into the collection
	_, err := collection.InsertOne(ctx, image)
	if err != nil {

		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log the duplicate key error
					logger.Info(fmt.Sprintf("Duplicate key error: %v", writeErr.Message))
				}

			}
		}

		return fmt.Errorf("error inserting plan into database: %w", err)

	}

	return nil
}
