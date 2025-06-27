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

// GetAllBlogService retrieves public blogs with filtering, pagination, and sorting
func GetAllBlogService(ctx context.Context, params utils.PaginationParams) ([]models.Blog, error) {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	// Base filters for public blogs
	baseFilter := bson.M{
		"isApproved": true,
		"isActive":   true,
		"isDeleted":  false,
		"status":     "published",
	}
	filter := buildBlogFilters(baseFilter, params.Filters)

	// Options for pagination and sorting
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: params.SortField, Value: params.SortOrder}})
	findOptions.SetSkip(int64((params.Page - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	logger.Info("Querying public blogs", "filter", filter, "sortField", params.SortField, "sortOrder", params.SortOrder)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No public blogs found")
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving blogs: %v", err))
		return nil, fmt.Errorf("error fetching blogs: %w", err)
	}
	defer cursor.Close(ctx)

	var blogs []models.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		logger.Error(fmt.Sprintf("Error decoding blogs: %v", err))
		return nil, fmt.Errorf("error decoding blogs: %w", err)
	}

	// Hydrate comments
	if err := hydrateBlogComments(ctx, client, blogs); err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Retrieved %d public blog(s)", len(blogs)))
	return blogs, nil
}

// GetAdminAllBlogService retrieves all blogs for admin view
func GetAdminAllBlogService(ctx context.Context, params utils.PaginationParams) ([]models.Blog, error) {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	filter := buildBlogFilters(bson.M{}, params.Filters)

	// Options
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: params.SortField, Value: params.SortOrder}})
	findOptions.SetSkip(int64((params.Page - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	logger.Info("Querying admin blogs", "filter", filter, "sortField", params.SortField, "sortOrder", params.SortOrder)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No admin blogs found")
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving blogs: %v", err))
		return nil, fmt.Errorf("error fetching blogs: %w", err)
	}
	defer cursor.Close(ctx)

	var blogs []models.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		logger.Error(fmt.Sprintf("Error decoding blogs: %v", err))
		return nil, fmt.Errorf("error decoding blogs: %w", err)
	}

	// Hydrate comments
	if err := hydrateBlogComments(ctx, client, blogs); err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Retrieved %d admin blog(s)", len(blogs)))
	return blogs, nil
}

// buildBlogFilters constructs filters supporting both direct match and $in
func buildBlogFilters(base bson.M, filters map[string]interface{}) bson.M {
	for key, value := range filters {
		switch v := value.(type) {
		case []string:
			base[key] = bson.M{"$in": v}
		case string:
			base[key] = v
		default:
			base[key] = v
		}
	}
	return base
}

// hydrateBlogComments attaches comments to each blog
func hydrateBlogComments(ctx context.Context, client *mongo.Client, blogs []models.Blog) error {
	commentCollection := client.Database("practionweb").Collection("BlogComments")

	for i, blog := range blogs {
		if len(blog.CommentsList) == 0 {
			continue
		}

		commentFilter := bson.M{"_id": bson.M{"$in": blog.CommentsList}}
		commentCursor, err := commentCollection.Find(ctx, commentFilter)
		if err != nil {
			logger.Error(fmt.Sprintf("Error retrieving comments for blog %v: %v", blog.ID, err))
			return fmt.Errorf("error retrieving comments: %w", err)
		}
		defer commentCursor.Close(ctx)

		var comments []models.Comments
		if err = commentCursor.All(ctx, &comments); err != nil {
			logger.Error(fmt.Sprintf("Error decoding comments: %v", err))
			return fmt.Errorf("error decoding comments: %w", err)
		}
		blogs[i].Comments = comments
	}
	return nil
}
