package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetOneBlog(ctx context.Context, id string) (*models.Blog, error) {

	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	// Convert the string ID to a MongoDB ObjectID
	var filter bson.M
	objectID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		// Valid ObjectID
		filter = bson.M{
			"_id":        objectID,
			"isApproved": true,
			"isActive":   true,
			"isDeleted":  false,
			"status":     "published",
		}
	} else {
		// Treat as slug
		filter = bson.M{
			"slug":       id,
			"isApproved": true,
			"isActive":   true,
			"isDeleted":  false,
			"status":     "published",
		}
	}

	// Query the collection
	var blog models.Blog
	err = collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info(fmt.Sprintf("No blog found with ID: %s", id))
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving blog from the database: %v", err))
		return nil, fmt.Errorf("error fetching blog details: %w", err)
	}

	if len(blog.CommentsList) > 0 {
		commentCollection := client.Database("practionweb").Collection("BlogComments")
		commentFilter := bson.M{"_id": bson.M{"$in": blog.CommentsList}}

		// Fetch the comments by ObjectIDs from CommentsList
		commentCursor, err := commentCollection.Find(ctx, commentFilter)
		if err != nil {
			logger.Error(fmt.Sprintf("Error retrieving comments for blog %v: %v", blog.ID, err))
			return nil, fmt.Errorf("error retrieving comments: %w", err)
		}
		defer commentCursor.Close(ctx)

		var comments []models.Comments
		if err = commentCursor.All(ctx, &comments); err != nil {
			logger.Error(fmt.Sprintf("Error decoding comments: %v", err))
			return nil, fmt.Errorf("error decoding comments: %w", err)
		}

		logger.Info(fmt.Sprintf("Fetched Comments: %v", comments))

		// Add the comments to the blog's Comments field
		blog.Comments = comments
	}

	logger.Info(fmt.Sprintf("Retrieved blog with ID: %s successfully.", id))
	return &blog, nil
}
