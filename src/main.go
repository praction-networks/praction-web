package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/app"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
)

// @title Praction Web APIs Server
// @version 1.0.1
// @description This API Server will be used to serve Praction backend API Server
// @BasePath /

func main() {
	// Create the application
	app, err := app.New()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to create application: %v", err))
	}

	// Handle OS signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the application in a goroutine
	go func() {
		if err := app.Start(ctx); err != nil {
			logger.Fatal(fmt.Sprintf("Application start error: %v", err))
		}
	}()

	// Wait for termination signal
	<-ctx.Done()
	logger.Info("Shutting down...")

	// Perform graceful shutdown with a timeout
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctxShutdown); err != nil {
		logger.Error(fmt.Sprintf("Error during shutdown: %v", err))
	} else if ctxShutdown.Err() == context.DeadlineExceeded {
		logger.Warn("Shutdown process exceeded timeout")
	} else {
		logger.Info("Application stopped")
	}
}
