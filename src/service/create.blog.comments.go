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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoBlogFound       = errors.New("no blog found")
	ErrDatabaseOperation = errors.New("database error")
)

func CreateBlogComments(ctx context.Context, blogComments models.Comments) error {

	//Check If Blog is Valid or Not

	err := isBlogAvailable(ctx, blogComments.BlogID)
	if err != nil {
		if errors.Is(err, ErrNoBlogFound) {
			logger.Info(fmt.Sprintf("No blog found with ID: %v", blogComments.BlogID))
			return fmt.Errorf("blog with ID %v does not exist", blogComments.BlogID)
		}
		if errors.Is(err, ErrDatabaseOperation) {
			logger.Error(fmt.Sprintf("Database error while checking blog availability: %v", err))
			return fmt.Errorf("internal database error occurred")
		}

		// Handle other unexpected errors
		logger.Error(fmt.Sprintf("Unexpected error: %v", err))
		return fmt.Errorf("unexpected error occurred: %v", err)
	}

	blogComments.UUID = uuid.New().String()
	blogComments.CreatedAt = time.Now()
	blogComments.IsActive = true
	blogComments.UpdatedAt = time.Now()
	blogComments.IsDeleted = false

	// Insert the user into the MongoDB database
	CommentID, err := insertBlogCommentsIntoDB(ctx, blogComments)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert User Intrst into DB: %v", err))
		return fmt.Errorf("failed to create plan: %w", err)
	}

	// Append Blog COmments ID to Blog Comments Section

	err = appendBlogCommentsToBlog(ctx, blogComments.BlogID, CommentID)

	if err != nil {
		logger.Error("Failed to Update Blog With Comments")
		return err
	}

	logger.Info(fmt.Sprintf("Comments Created successfully for blog ID %s created successfully.", blogComments.BlogID))
	return nil
}

// insertBlogCommentsIntoDB inserts a new comment into the BlogComments collection and returns the comment ID
func insertBlogCommentsIntoDB(ctx context.Context, blogComments models.Comments) (primitive.ObjectID, error) {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("BlogComments")

	// Insert the comment document into the collection
	result, err := collection.InsertOne(ctx, blogComments)
	if err != nil {
		// Handle duplicate key error if any
		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log the duplicate key error
					logger.Info(fmt.Sprintf("Duplicate key error: %v", writeErr.Message))
				}
			}
		}
		return primitive.NilObjectID, fmt.Errorf("error inserting Blog Comments: %w", err)
	}

	// Extract the inserted ID and cast it to primitive.ObjectID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("failed to cast inserted ID to ObjectID")
	}

	return insertedID, nil
}

func isBlogAvailable(ctx context.Context, blogID primitive.ObjectID) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	// Create a filter to search for the blog by its ID
	filter := bson.D{{Key: "_id", Value: blogID}}

	// Attempt to find the blog
	err := collection.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info(fmt.Sprintf("No blog found with ID: %v", blogID))
			return ErrNoBlogFound
		}
		logger.Error(fmt.Sprintf("Error fetching blog with ID %v: %v", blogID, err))
		return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	return nil
}
func appendBlogCommentsToBlog(ctx context.Context, blogID primitive.ObjectID, commentID primitive.ObjectID) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	// Filter to locate the blog by its ID
	filter := bson.M{"_id": blogID}

	// Update operation to add the CommentID to the CommentsList if it doesn't already exist
	update := bson.M{
		"$addToSet": bson.M{"comments": commentID}, // Prevents duplicates
	}

	// Execute the update
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error(fmt.Sprintf("Error appending CommentID %v to blog with ID %v: %v", commentID, blogID, err))
		return fmt.Errorf("failed to append comment to blog: %w", err)
	}

	logger.Info(fmt.Sprintf("Successfully appended CommentID %v to blog with ID %v.", commentID, blogID))
	return nil
}
