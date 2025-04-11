package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PaginationParams represents pagination, filtering, and sorting parameters
func GetAllBlogService(ctx context.Context, params utils.PaginationParams) ([]models.Blog, error) {

	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("Blog")
	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.SortField == "" {
		params.SortField = "createdAt"
	}
	if params.SortOrder != 1 && params.SortOrder != -1 {
		params.SortOrder = -1
	}

	filter := buildBlogFilters(
		bson.M{
			"isApproved": true,
			"isActive":   true,
			"isDeleted":  false,
			"status":     "published",
		},
		params.Filters,
	)

	// Define find options for pagination and sorting
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: params.SortField, Value: params.SortOrder}})
	findOptions.SetSkip(int64((params.Page - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	// Query the collection
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No blogs found in the database to show.")
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving blogs from the database: %v", err))
		return nil, fmt.Errorf("error fetching blogs details: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice
	var blogs []models.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		logger.Error(fmt.Sprintf("Error decoding users: %v", err))
		return nil, fmt.Errorf("error decoding users: %w", err)
	}

	// Fetch comments for each blog based on the Comment IDs
	// Now we will retrieve the comments for each blog
	for i, blog := range blogs {
		// Ensure CommentsList is not empty
		if len(blog.CommentsList) > 0 {
			commentCollection := client.Database("uvfiberweb").Collection("BlogComments")
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

			// Add the comments to the blog's Comments field
			blogs[i].Comments = comments
		}
	}

	logger.Info(fmt.Sprintf("Retrieved %d blog(s) successfully with pagination and sorting and filters.", len(blogs)))
	return blogs, nil
}

func GetAdminAllBlogService(ctx context.Context, params utils.PaginationParams) ([]models.Blog, error) {

	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("Blog")
	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.SortField == "" {
		params.SortField = "createdAt"
	}
	if params.SortOrder != 1 && params.SortOrder != -1 {
		params.SortOrder = 1
	}

	// Build the query filters
	filter := bson.M{}

	// Add additional filters for category and tag if present
	if categories, ok := params.Filters["category"]; ok {
		// Filter blogs where category matches any of the given values
		filter["category"] = bson.M{"$in": categories}
	}

	if tags, ok := params.Filters["tag"]; ok {
		// Filter blogs where tag matches any of the given values
		filter["tag"] = bson.M{"$in": tags}
	}

	for key, value := range params.Filters {
		filter[key] = value
	}

	// Define find options for pagination and sorting
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: params.SortField, Value: params.SortOrder}})
	findOptions.SetSkip(int64((params.Page - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	// Query the collection
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No blogs found in the database to show.")
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving blogs from the database: %v", err))
		return nil, fmt.Errorf("error fetching blogs details: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice
	var blogs []models.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		logger.Error(fmt.Sprintf("Error decoding users: %v", err))
		return nil, fmt.Errorf("error decoding users: %w", err)
	}

	// Fetch comments for each blog based on the Comment IDs
	// Now we will retrieve the comments for each blog
	for i, blog := range blogs {
		// Ensure CommentsList is not empty
		if len(blog.CommentsList) > 0 {
			commentCollection := client.Database("uvfiberweb").Collection("BlogComments")
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

			// Add the comments to the blog's Comments field
			blogs[i].Comments = comments
		}
	}

	logger.Info(fmt.Sprintf("Retrieved %d blog(s) successfully with pagination and sorting and filters.", len(blogs)))
	return blogs, nil
}

func buildBlogFilters(base bson.M, filters map[string]interface{}) bson.M {
	for key, value := range filters {
		switch key {
		case "category", "tag":
			base[key] = bson.M{"$in": value}
		default:
			base[key] = value
		}
	}
	return base
}
