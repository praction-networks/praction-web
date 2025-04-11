package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateBlogCategoryService handles the creation of a new blog category.
func CreateBlogTagService(ctx context.Context, blogTag models.BlogTag) error {
	// Assign essential metadata
	blogTag.UUID = uuid.New().String()
	blogTag.IsActive = true
	blogTag.CreatedAt = time.Now()
	blogTag.UpdatedAt = time.Now()

	// Insert the category into the database
	if err := insertBlogTagIntoDB(ctx, blogTag); err != nil {
		if errors.Is(err, ErrDuplicateKey) {
			// Provide a clear and informative message about the duplicate field
			logger.Error(fmt.Sprintf("Duplicate entry for tag: %s (Slug: %s)", blogTag.Name, blogTag.Slug))
			return fmt.Errorf("a tag with the slug '%s' already exists", blogTag.Slug)
		}
		logger.Error(fmt.Sprintf("Failed to insert blog tag into DB: %v", err))
		return fmt.Errorf("failed to create blog tag: %w", err)
	}

	logger.Info(fmt.Sprintf("Blog category '%s' (Slug: %s) created successfully.", blogTag.Name, blogTag.Slug))
	return nil
}

// insertBlogCategoryIntoDB inserts the blog category into the MongoDB database.
func insertBlogTagIntoDB(ctx context.Context, blogTag models.BlogTag) error {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("BlogTag")

	// Attempt to insert the document
	_, err := collection.InsertOne(ctx, blogTag)
	if err != nil {
		// Handle MongoDB-specific errors
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			// Look for duplicate key error and log it with details
			for _, we := range writeErr.WriteErrors {
				if we.Code == 11000 { // Duplicate key error
					// Extract the duplicate key details
					if we.Message != "" {
						// Return the specific error message related to the field
						return fmt.Errorf("duplicate key error: %s", we.Message)
					}
					// Default error message if no specific message is available
					return ErrDuplicateKey
				}
			}
		}
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}
