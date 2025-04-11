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
func GetAllUnavailableAreaUserService(ctx context.Context, params utils.PaginationParams, db string) ([]models.UnAvailableArea, error) {

	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection(db)
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
			logger.Info("No users found in the database who showed interest in our services.")
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving user data from the database: %v", err))
		return nil, fmt.Errorf("error fetching user details: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice
	var users []models.UnAvailableArea
	if err = cursor.All(ctx, &users); err != nil {
		logger.Error(fmt.Sprintf("Error decoding users: %v", err))
		return nil, fmt.Errorf("error decoding users: %w", err)
	}

	logger.Info(fmt.Sprintf("Retrieved %d user(s) successfully with pagination and sorting.", len(users)))
	return users, nil
}
