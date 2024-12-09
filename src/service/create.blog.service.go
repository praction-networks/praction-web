package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	apperror "github.com/praction-networks/quantum-ISP365/webapp/src/appError"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateBlog(ctx context.Context, blog models.Blog) error {

	// Check if Blog Title is already used
	err := getBlogByName(ctx, blog.BlogTitle)
	if err != nil {
		if strings.Contains(err.Error(), "database error") {
			// Handle database error
			logger.Error("Database error: ", "Error", err)
			return fmt.Errorf("failed to connect with database")
		} else if err.Error() == "blog with this title already exists" {
			// Handle the case where the blog already exists
			logger.Info("Blog exists, cannot create new one")
			return fmt.Errorf("blog exists, cannot create new one with same blog title")
		}
	}
	// Check if Blog Image is available to use
	logger.Info("Checking if Image UUID is available to use")

	blogImageURL, featureBlogImageURL, err := getImageByUUID(ctx, blog.BlogImage, blog.FeatureImage)
	if err != nil {
		if errors.Is(err, apperror.ErrBlogImageNotFound) {
			logger.Error("Blog Image not found")
			return fmt.Errorf("blog image not found")
		} else if errors.Is(err, apperror.ErrFeatureImageNotFound) {
			logger.Error("Feature Blog Image not found")
			return fmt.Errorf("feature blog image not found")
		} else {
			// Handle other errors like database connection issues
			return fmt.Errorf("error fetching images: %v", err)
		}
	}

	blog.BlogImage = blogImageURL
	blog.FeatureImage = featureBlogImageURL
	// Check if Assigned Titles are in Database

	err = checkBlogCategory(ctx, blog.Category)

	if err != nil {
		logger.Error("Can not find blog category inside the databased", "error", err)
		return err
	}

	//Check if Assigne Tag are in database

	err = checkBlogTag(ctx, blog.Tag)

	if err != nil {
		logger.Error("Can not find blog tag inside the databased", "error", err)
		return err
	}

	blog.UUID = uuid.New().String()
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	blog.Status = "draft"
	blog.IsApproved = false
	blog.IsActive = false
	blog.IsDeleted = false
	blog.View = 0
	blog.CommentsCount = 0
	blog.Shares = 0

	err = insertBlogIntoDB(ctx, blog)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert Blog into DB: %v", err))
		return fmt.Errorf("failed to create Blog: %w", err)
	}

	logger.Info(fmt.Sprintf("Plan %s created successfully.", blog.BlogTitle))
	return nil
}

// insertUserIntoDB inserts the new user into the MongoDB database
func insertBlogIntoDB(ctx context.Context, blog models.Blog) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	// Insert the user document into the collection
	_, err := collection.InsertOne(ctx, blog)
	if err != nil {

		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log the duplicate key error
					logger.Info(fmt.Sprintf("Duplicate key error: %v", writeErr.Message))
				}

			}
		}

		return fmt.Errorf("error inserting blog into database: %w", err)

	}

	return nil
}

func getBlogByName(ctx context.Context, blogTitle string) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	filter := bson.D{{Key: "blogTitle", Value: blogTitle}}

	var blog models.Blog
	err := collection.FindOne(ctx, filter).Decode(&blog)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No blog found with this title, so it can be created
			return nil
		}
		// Handle other errors
		logger.Error("Error fetching blog:", "error", err)
		return err
	}

	// If a blog is found with the given title, return an error to prevent creation
	return fmt.Errorf("blog with this title already exists")
}

func getImageByUUID(ctx context.Context, blogImageUUID, featureImageUUID string) (string, string, error) {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("BlogImage")

	// Find BlogImage using blogImageUUID
	filterBlogImage := bson.D{{Key: "uuid", Value: blogImageUUID}}
	var blogImage models.Image

	err := collection.FindOne(ctx, filterBlogImage).Decode(&blogImage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return error for blog image not found
			logger.Error("Blog image not found", "UUID", blogImageUUID)
			return "", "", apperror.ErrBlogImageNotFound
		}
		// Handle other errors
		logger.Error("Error fetching blog image", "error", err)
		return "", "", apperror.ErrFetchingImage
	}

	// Find FeatureImage using featureImageUUID
	filterFeatureImage := bson.D{{Key: "uuid", Value: featureImageUUID}}
	var featureImage models.Image

	err = collection.FindOne(ctx, filterFeatureImage).Decode(&featureImage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return error for feature image not found
			logger.Error("Feature image not found", "UUID", featureImageUUID)
			return blogImage.ImageURL, "", apperror.ErrFeatureImageNotFound
		}
		// Handle other errors
		logger.Error("Error fetching feature image", "error", err)
		return blogImage.ImageURL, "", apperror.ErrFetchingImage
	}

	// Return both image URLs if both are found
	return blogImage.ImageURL, featureImage.ImageURL, nil
}

func checkBlogCategory(ctx context.Context, categorys []string) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("BlogCategory")

	// Create a filter for matching any of the provided UUIDs
	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$in", Value: categorys}}}}

	// Find documents matching the UUIDs
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		// Handle any error while querying the database
		logger.Error("Error fetching blog categories", "error", err)
		return fmt.Errorf("error fetching blog categories: %v", err)
	}
	defer cursor.Close(ctx)

	// Track how many UUIDs were found in the database
	foundCount := 0
	for cursor.Next(ctx) {
		foundCount++
	}

	// If the number of found categories is less than the number of provided UUIDs, return an error
	if foundCount < len(categorys) {
		return fmt.Errorf("one or more blog categories not found")
	}

	// If all categories are found, return nil (no error)
	return nil
}

func checkBlogTag(ctx context.Context, tags []string) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("BlogTag")

	// Create a filter for matching any of the provided UUIDs
	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$in", Value: tags}}}}

	// Find documents matching the UUIDs
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		// Handle any error while querying the database
		logger.Error("Error fetching blog tag", "error", err)
		return fmt.Errorf("error fetching blog tag: %v", err)
	}
	defer cursor.Close(ctx)

	// Track how many UUIDs were found in the database
	foundCount := 0
	for cursor.Next(ctx) {
		foundCount++
	}

	// If the number of found categories is less than the number of provided UUIDs, return an error
	if foundCount < len(tags) {
		return fmt.Errorf("one or more blog tag not found")
	}

	// If all categories are found, return nil (no error)
	return nil
}
