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

func SaveblogImage(ImageName, ImageID, ImageURL string, name string, tag string) (models.Image, error) {

	var Image models.Image
	Image.UUID = uuid.New().String()
	Image.FileName = ImageName
	Image.Name = name
	Image.Tag = tag
	Image.FileID = ImageID
	Image.ImageURL = ImageURL
	Image.CreatedAt = time.Now()
	Image.UpdatedAt = time.Now()
	Image.IsActive = true
	Image.IsDeleted = false

	ctx := context.TODO()
	err := insertGoogleImageIntoDB(ctx, Image)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert Blog Image Metadata into DB: %v", err))
		return models.Image{}, fmt.Errorf("failed to create plan: %w", err)
	}

	logger.Info(fmt.Sprintf("Google Image Metadata with uuid %s created successfully.", Image.UUID))
	return Image, nil
}

// insertUserIntoDB inserts the new user into the MongoDB database
func insertGoogleImageIntoDB(ctx context.Context, image models.Image) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("Image")

	// Insert the user document into the collection

	logger.Info("Image", "ImageData", image)

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
