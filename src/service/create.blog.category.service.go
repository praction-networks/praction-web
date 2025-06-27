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

// Predefined error for duplicate key
var ErrDuplicateKey = errors.New("duplicate key error")

// CreateBlogCategoryService handles the creation of a new blog category.
func CreateBlogCategoryService(ctx context.Context, blogCategory models.BlogCategory) error {
	// Assign essential metadata
	blogCategory.UUID = uuid.New().String()
	blogCategory.CreatedAt = time.Now()
	blogCategory.IsActive = true
	blogCategory.UpdatedAt = time.Now()

	// Insert the category into the database
	if err := insertBlogCategoryIntoDB(ctx, blogCategory); err != nil {
		if errors.Is(err, ErrDuplicateKey) {
			// Provide a clear and informative message about the duplicate field
			logger.Error(fmt.Sprintf("Duplicate entry for category: %s (Slug: %s)", blogCategory.Name, blogCategory.Slug))
			return fmt.Errorf("a category with the slug '%s' already exists", blogCategory.Slug)
		}
		logger.Error(fmt.Sprintf("Failed to insert category into DB: %v", err))
		return fmt.Errorf("failed to create blog category: %w", err)
	}

	logger.Info(fmt.Sprintf("Blog category '%s' (Slug: %s) created successfully.", blogCategory.Name, blogCategory.Slug))
	return nil
}

// insertBlogCategoryIntoDB inserts the blog category into the MongoDB database.
func insertBlogCategoryIntoDB(ctx context.Context, blogCategory models.BlogCategory) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("BlogCategory")

	// Attempt to insert the document
	_, err := collection.InsertOne(ctx, blogCategory)
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
