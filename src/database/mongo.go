package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/config"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitializeMongo(ctx context.Context) error {

	var err error

	cfg, err := config.MongoEnvGet()
	if err != nil {
		return fmt.Errorf("failed to initialize Mongo env config: %w", err)
	}

	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)).
		SetAuth(options.Credential{
			Username:      cfg.DBUser,
			Password:      cfg.DBPassword,
			AuthSource:    cfg.DBName,
			AuthMechanism: "SCRAM-SHA-256",
		}).

		// Add pool size settings (optional)
		SetMaxPoolSize(100).
		SetMinPoolSize(10)

	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("error creating MongoDB client: %v", err)
	}

	// Adding a timeout to the ping to avoid hanging
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check the connection
	err = client.Ping(ctxPing, nil)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	//Creating Indexing for DomainUser
	err = CreateIndexesForUser(ctx, cfg.DBName, "User")

	if err != nil {
		return fmt.Errorf("error while creating indexes: %v", err)
	}

	err = CreateIndexesForUserIntrest(ctx, cfg.DBName, "UserIntrest")

	if err != nil {
		return fmt.Errorf("error while creating indexes: %v", err)
	}

	err = CreateIndexesForPlan(ctx, cfg.DBName, "Plan")

	if err != nil {
		return fmt.Errorf("error while creating indexes: %v", err)
	}

	err = CreateIndexesForBlogCategory(ctx, cfg.DBName, "BlogCategory")

	if err != nil {
		return fmt.Errorf("error while creating indexes: %v", err)
	}

	err = CreateIndexesForBlogTag(ctx, cfg.DBName, "BlogTag")

	if err != nil {
		return fmt.Errorf("error while creating indexes: %v", err)
	}

	err = CreateIndexesForServiceArea(ctx, cfg.DBName, "ServiceArea")

	if err != nil {
		return fmt.Errorf("error while creating indexes: %v", err)
	}
	err = CreateIndexesForBlogComments(ctx, cfg.DBName, "BlogComments")

	if err != nil {
		return fmt.Errorf("error while creating indexes: %v", err)
	}
	err = createUniqueIndexesForUserReferenceCollection(ctx, cfg.DBName, "UserReferal")

	if err != nil {
		return fmt.Errorf("error while creating indexes: %v", err)
	}
	logger.Info("Index created for User collection successfully.")
	return nil
}

// CreateIndexes creates indexes on the required fields
func CreateIndexesForUser(ctx context.Context, dbName, collectionName string) error {
	collection := client.Database(dbName).Collection(collectionName)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1}, // Unique index on email
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"mobile": 1}, // Unique index on mobile
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"UserName": 1}, // Unique index on mobile
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"uuid": 1}, // Unique index on mobile
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}
	// Log success
	logger.Info(fmt.Sprintf("Indexes created for collection: %s", collectionName))

	return nil
}

// CreateIndexes creates indexes on the required fields
func CreateIndexesForUserIntrest(ctx context.Context, dbName, collectionName string) error {
	collection := client.Database(dbName).Collection(collectionName)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1}, // Unique index on email
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"mobile": 1}, // Unique index on mobile
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"uuid": 1}, // Unique index on mobile
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}
	// Log success
	logger.Info(fmt.Sprintf("Indexes created for collection: %s", collectionName))

	return nil
}

func CreateIndexesForPlan(ctx context.Context, dbName, collectionName string) error {
	collection := client.Database(dbName).Collection(collectionName)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"category": 1}, // Unique index on category
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}
	// Log success
	logger.Info(fmt.Sprintf("Indexes created for collection: %s", collectionName))

	return nil
}

func CreateIndexesForBlogCategory(ctx context.Context, dbName, collectionName string) error {
	collection := client.Database(dbName).Collection(collectionName)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1}, // Unique index on category
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"uuid": 1}, // Unique index on mobile
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}
	// Log success
	logger.Info(fmt.Sprintf("Indexes created for collection: %s", collectionName))

	return nil
}

func CreateIndexesForServiceArea(ctx context.Context, dbName, collectionName string) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Define multiple indexes
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "features.uuid", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "features.properties.areaName", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "uuid", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	// Create indexes in bulk
	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}

	// Log success
	logger.Info(fmt.Sprintf("Indexes created for collection: %s", collectionName))

	return nil
}

func CreateIndexesForBlogTag(ctx context.Context, dbName, collectionName string) error {
	collection := client.Database(dbName).Collection(collectionName)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1}, // Unique index on category
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"uuid": 1}, // Unique index on mobile
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}
	// Log success
	logger.Info(fmt.Sprintf("Indexes created for collection: %s", collectionName))

	return nil
}

func CreateIndexesForBlogComments(ctx context.Context, dbName, collectionName string) error {
	collection := client.Database(dbName).Collection(collectionName)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"description": 1}, // Unique index on category
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"uuid": 1}, // Unique index on mobile
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}
	// Log success
	logger.Info(fmt.Sprintf("Indexes created for collection: %s", collectionName))

	return nil
}

func createUniqueIndexesForUserReferenceCollection(ctx context.Context, dbName, collectionName string) error {
	collection := client.Database(dbName).Collection(collectionName)

	// Create unique index for Referrels.Email
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "referrels.email", Value: 1}}, // Unique for nested emails
		Options: options.Index().SetUnique(true).SetName("unique_email"),
	}

	// Create unique index for Referrels.Mobile
	mobileIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "referrels.mobile", Value: 1}}, // Unique for nested mobiles
		Options: options.Index().SetUnique(true).SetName("unique_mobile"),
	}

	// Apply the indexes
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{emailIndex, mobileIndex})
	if err != nil {
		return err
	}

	log.Println("Unique indexes for Email and Mobile created successfully.")
	return nil
}

// GetClient returns the MongoDB client instance// GetClient returns the MongoDB client instance
func GetClient() *mongo.Client {
	return client
}

// CloseClient closes the MongoDB client connection
func CloseClient(ctx context.Context) {
	if client != nil {
		if err := client.Disconnect(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error disconnecting MongoDB client: %v", err))
		} else {
			logger.Info("MongoDB connection closed successfully.")
		}
	}
}
