package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// isValidUUIDv4 checks if a string is a valid UUID v4
func IsValidUUIDv4(s string) bool {
	parsedUUID, err := uuid.Parse(s)
	return err == nil && parsedUUID.Version() == 4
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

	var image models.Image

	err := imageCollection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&image)
	if err != nil {
		logger.Error("Image not found for UUID", "UUID", uuid, "error", err)
		return "", fmt.Errorf("image not found for UUID: %s", uuid)
	}

	return image.ImageURL, nil
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

	if updateBlogData.MetaTitle != "" {
		updateFields["metaTitle"] = updateBlogData.BlogTitle
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
		logger.Info("Recived Blog Images is", "Image", updateBlogData.BlogImage)
		if IsValidUUIDv4(updateBlogData.BlogImage) {
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

		logger.Info("Passed Blog Images is", "Image", updateFields["blogImage"])
	}

	// Handling FeatureImage
	if updateBlogData.FeatureImage != "" {
		if IsValidUUIDv4(updateBlogData.FeatureImage) {
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

	if len(updateBlogData.EmbeddedMedia) > 0 { // ✅ Check if array is not empty
		var updatedMedia []string

		for _, media := range updateBlogData.EmbeddedMedia { // ✅ Loop over each media entry
			if IsValidUUIDv4(media) { // ✅ Check if it's a valid UUID
				imageURL, err := getImageURLFromDatabase(ctx, media)
				if err != nil {
					logger.Error("Failed to fetch image URL", "UUID", media, "error", err)
					return fmt.Errorf("failed to retrieve image URL for UUID: %w", err)
				}
				updatedMedia = append(updatedMedia, imageURL) // ✅ Store URL
			} else if isValidURL(media) { // ✅ If valid URL, store directly
				updatedMedia = append(updatedMedia, media)
			} else {
				logger.Error("Invalid embeddedMedia format", "embeddedMedia", media)
				return fmt.Errorf("embeddedMedia must be a valid UUID v4 or a URL")
			}
		}

		// ✅ Update the embeddedMedia field with the modified array
		updateFields["embeddedMedia"] = updatedMedia
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
