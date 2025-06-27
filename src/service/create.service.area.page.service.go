package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperror "github.com/praction-networks/quantum-ISP365/webapp/src/appError"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateAreaPage creates a new Area Page entry
func CreateAreaPage(ctx context.Context, area models.ServiceAreaPage) error {
	logger.Info("Starting Area Page creation for:", "areaName", area.AreaName)

	imageURL, err := getPageImageByUUID(ctx, area.AreaImage)
	if err != nil {
		if errors.Is(err, apperror.ErrPageImageNotFound) {
			logger.Error("Area image not found for UUID", "uuid", area.AreaImage)
			return fmt.Errorf("area image not found for UUID: %s", area.AreaImage)
		}
		logger.Error("Error fetching area image", "error", err)
		return fmt.Errorf("failed to fetch area image: %w", err)
	}

	area.AreaImage = imageURL
	area.UUID = uuid.New().String()
	area.CreatedAt = time.Now()
	area.UpdatedAt = time.Now()
	area.IsActive = true
	area.IsDeleted = false

	return insertAreaPageIntoDB(ctx, area)
}

func insertAreaPageIntoDB(ctx context.Context, area models.ServiceAreaPage) error {
	collection := database.GetClient().Database("practionweb").Collection("AreaPage")

	_, err := collection.InsertOne(ctx, area)
	if err != nil {
		if writeErr, ok := err.(mongo.WriteException); ok {
			for _, e := range writeErr.WriteErrors {
				if e.Code == 11000 {
					logger.Warn("Duplicate area entry detected", "message", e.Message)
				}
			}
		}
		return fmt.Errorf("error inserting area into database: %w", err)
	}
	return nil
}

func getPageImageByUUID(ctx context.Context, uuid string) (string, error) {
	collection := database.GetClient().Database("practionweb").Collection("Image")

	var image models.Image
	err := collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Warn("Image not found for UUID", "uuid", uuid)
			return "", apperror.ErrPageImageNotFound
		}
		logger.Error("Error querying image collection", "error", err)
		return "", apperror.ErrFetchingImage
	}

	return image.ImageURL, nil
}

// GetAreaPage returns one AreaPage by UUID
func GetAreaPage(ctx context.Context, uuid string) (*models.ServiceAreaPage, error) {
	collection := database.GetClient().Database("practionweb").Collection("AreaPage")
	var area models.ServiceAreaPage
	err := collection.FindOne(ctx, bson.M{"uuid": uuid, "isDeleted": false}).Decode(&area)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get area page: %w", err)
	}
	return &area, nil
}

// GetAllBlogService retrieves public blogs with filtering, pagination, and sorting
func GetAllPageAreaService(ctx context.Context, params utils.PaginationParams) ([]models.ServiceAreaPage, error) {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("AreaPage")

	// Base filters for public blogs
	baseFilter := bson.M{
		"isActive":  true,
		"isDeleted": false,
	}
	filter := buildBlogFilters(baseFilter, params.Filters)

	// Options for pagination and sorting
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: params.SortField, Value: params.SortOrder}})
	findOptions.SetSkip(int64((params.Page - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	logger.Info("Querying page meta", "filter", filter, "sortField", params.SortField, "sortOrder", params.SortOrder)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No page meta found")
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving page meta: %v", err))
		return nil, fmt.Errorf("error fetching page meta: %w", err)
	}
	defer cursor.Close(ctx)

	var pageMeta []models.ServiceAreaPage
	if err = cursor.All(ctx, &pageMeta); err != nil {
		logger.Error(fmt.Sprintf("Error decoding page meta: %v", err))
		return nil, fmt.Errorf("error decoding page meta: %w", err)
	}

	logger.Info(fmt.Sprintf("Retrieved %d public blog(s)", len(pageMeta)))
	return pageMeta, nil
}

// DeleteAreaPage marks an AreaPage as deleted
func DeleteAreaPage(ctx context.Context, id primitive.ObjectID) error {
	collection := database.GetClient().Database("practionweb").Collection("AreaPage")
	update := bson.M{
		"$set": bson.M{
			"isDeleted": true,
			"updatedAt": time.Now(),
		},
	}

	filter := bson.M{"_id": id, "isDeleted": false}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete area page: %w", err)
	}
	if res.MatchedCount == 0 {
		return apperror.ErrNotFound
	}
	return nil
}
