package service

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// isValidUUIDv4 checks if a string is a valid UUID v4
func isValidUUIDv4(s string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(s)
}

// isValidURL checks if a string is a valid URL
func isValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

// getImageURLFromDatabase searches for an image by UUID and returns its URL
func getImageURLFromDatabase(ctx context.Context, uuid string) (string, error) {
	client := database.GetClient()
	imageCollection := client.Database("practionweb").Collection("Image")

	var imageDoc struct {
		URL string `bson:"url"`
	}

	err := imageCollection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&imageDoc)
	if err != nil {
		logger.Error("Image not found for UUID", "UUID", uuid, "error", err)
		return "", fmt.Errorf("image not found for UUID: %s", uuid)
	}

	return imageDoc.URL, nil
}

func UpdateOneBlog(ctx context.Context, id string, updateBlogData *models.BlogUpdate) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	// Convert string ID to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Invalid blog ID format", "blogID", id, "error", err)
		return fmt.Errorf("invalid blog ID format: %w", err)
	}

	// Prepare dynamic update fields
	updateFields := bson.M{}

	// Blog Title
	if updateBlogData.BlogTitle != "" {
		updateFields["blogTitle"] = updateBlogData.BlogTitle
	}

	// Slug
	if updateBlogData.Slug != "" {
		updateFields["slug"] = updateBlogData.Slug
	}

	// Blog Description
	if updateBlogData.BlogDescription != "" {
		updateFields["blogDescription"] = updateBlogData.BlogDescription
	}

	// Meta Description
	if updateBlogData.MetaDescription != "" {
		updateFields["metaDescription"] = updateBlogData.MetaDescription
	}

	// Meta Keywords
	if updateBlogData.MetaKeywords != nil {
		updateFields["metaKeywords"] = updateBlogData.MetaKeywords
	}

	// Embedded Media
	if updateBlogData.EmbeddedMedia != nil {
		updateFields["embeddedMedia"] = updateBlogData.EmbeddedMedia
	}

	// Summary
	if updateBlogData.Summary != "" {
		updateFields["summary"] = updateBlogData.Summary
	}

	// Category
	if updateBlogData.Category != nil {
		updateFields["category"] = updateBlogData.Category
	}

	// Tag
	if updateBlogData.Tag != nil {
		updateFields["tag"] = updateBlogData.Tag
	}

	// Status
	if updateBlogData.Status != "" {
		updateFields["status"] = updateBlogData.Status
	}

	// Views, Shares, Comments
	if updateBlogData.View >= 0 {
		updateFields["view"] = updateBlogData.View
	}
	if updateBlogData.Shares >= 0 {
		updateFields["shares"] = updateBlogData.Shares
	}
	if updateBlogData.CommentsCount >= 0 {
		updateFields["commentsCount"] = updateBlogData.CommentsCount
	}
	if updateBlogData.CommentsList != nil {
		updateFields["commentsList"] = updateBlogData.CommentsList
	}
	if updateBlogData.Comments != nil {
		updateFields["comments"] = updateBlogData.Comments
	}

	// Handling BlogImage
	if updateBlogData.BlogImage != "" {
		if isValidUUIDv4(updateBlogData.BlogImage) {
			imageURL, err := getImageURLFromDatabase(ctx, updateBlogData.BlogImage)
			if err != nil {
				logger.Error("Failed to fetch image URL", "UUID", updateBlogData.BlogImage, "error", err)
				return fmt.Errorf("failed to retrieve image URL for UUID: %w", err)
			}
			updateFields["blogImage"] = imageURL
		} else if isValidURL(updateBlogData.BlogImage) {
			updateFields["blogImage"] = updateBlogData.BlogImage
		} else {
			logger.Error("Invalid BlogImage format", "BlogImage", updateBlogData.BlogImage)
			return fmt.Errorf("BlogImage must be a valid UUID v4 or a URL")
		}
	}

	// Handling FeatureImage
	if updateBlogData.FeatureImage != "" {
		if isValidUUIDv4(updateBlogData.FeatureImage) {
			imageURL, err := getImageURLFromDatabase(ctx, updateBlogData.FeatureImage)
			if err != nil {
				logger.Error("Failed to fetch image URL", "UUID", updateBlogData.FeatureImage, "error", err)
				return fmt.Errorf("failed to retrieve image URL for UUID: %w", err)
			}
			updateFields["featureImage"] = imageURL
		} else if isValidURL(updateBlogData.FeatureImage) {
			updateFields["featureImage"] = updateBlogData.FeatureImage
		} else {
			logger.Error("Invalid FeatureImage format", "FeatureImage", updateBlogData.FeatureImage)
			return fmt.Errorf("FeatureImage must be a valid UUID v4 or a URL")
		}
	}

	// Update `UpdatedAt` field
	updateFields["updatedAt"] = time.Now()

	// If no update fields exist, return an error
	if len(updateFields) == 0 {
		logger.Warn("No valid fields provided for update", "blogID", id)
		return fmt.Errorf("no valid fields provided for update")
	}

	// Update the document
	update := bson.M{"$set": updateFields}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		logger.Error("Error updating blog in the database", "blogID", id, "error", err)
		return fmt.Errorf("error updating blog: %w", err)
	}

	// Check if any document was modified
	if result.MatchedCount == 0 {
		logger.Warn("No blog found with the given ID for updating", "blogID", id)
		return fmt.Errorf("no blog found with the given ID: %s", id)
	}

	if result.ModifiedCount == 0 {
		logger.Warn("Blog update request did not modify any fields", "blogID", id)
		return fmt.Errorf("blog update did not modify any fields")
	}

	logger.Info("Blog updated successfully", "blogID", id)
	return nil
}
