package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CheckIfImageExists checks whether the image with the given hash already exists in the database.
func CheckIfImageExists(imageHash string) (bool, error) {
	// Get the MongoDB client
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("BlogImage")

	// Check if an image with the same hash exists in the database
	var image models.Image
	err := collection.FindOne(context.TODO(), bson.M{"filehash": imageHash}).Decode(&image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If no document is found, it's not a duplicate
			return false, nil
		}
		// Other errors (e.g., connection issues)
		return false, fmt.Errorf("error checking for duplicate image: %v", err)
	}

	// If an image is found, it's a duplicate
	return true, nil
}
