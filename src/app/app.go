package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/config"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/router"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	router http.Handler
	client *mongo.Client
}

func New() (*App, error) {

	// Initialize the logger
	if err := logger.SetupLogger(); err != nil {
		log.Fatalf("Unable to set up logger: %v", err)
	}

	logger.Info("Logger initialized successfully")
	// Initialize MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize MongoDB
	err := database.InitializeMongo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	err = service.CreateUser(ctx, "praction", "Lcmanager123", "akshay.kumar@praction.in", "Akshay", "Chauhan", "admin")

	if err != nil {
		logger.Warn("Failed to Create user or user already exist")
	}

	fmt.Println("User created successfully")

	app := &App{
		router: router.LoadRoutes(),
		client: database.GetClient(),
	}

	return app, nil
}

func (a *App) Start(ctx context.Context) error {

	ChiConfig, err := config.ChiEnvGet()

	port := ChiConfig.Port

	if err != nil {
		logger.Warn("Unable to get Chi Env Config, Setting Default Port as 3000 Error : %v", err)
		port = "3000"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: a.router,
	}

	defer func() {
		database.CloseClient(ctx)
	}()

	logger.Info("Starting server on", "Port:", port)

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}

func (a *App) Stop(ctx context.Context) error {
	if a.client != nil {
		if err := a.client.Disconnect(ctx); err != nil {
			return fmt.Errorf("error disconnecting MongoDB client: %v", err)
		}
		logger.Info("MongoDB connection closed successfully")
	}
	return nil
}
