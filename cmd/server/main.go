package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tezos-delegation-service/internal/app"
	"tezos-delegation-service/logger"
)

func main() {
	// Run the main application logic and handle any errors
	if err := run(); err != nil {
		// Exit with status code 1 on error
		os.Exit(1)
	}
}

// run initializes the application, starts it, and handles graceful shutdown.
func run() error {
	// Initialize the logger
	log := logger.NewLogger()

	// Create a new application instance
	app := app.NewApp()

	// Create a context that listens for OS interrupts and SIGTERM signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the application and handle any startup errors
	if err := app.Start(ctx); err != nil {
		log.Log.Fatal().Err(err).Msg("Failed to start the application")
		return err
	}

	// Wait for the context to be done (e.g., signal received)
	<-ctx.Done()
	log.Log.Info().Msg("Shutting down gracefully, press Ctrl+C again to force")

	// Create a new context with a timeout for graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop the application and handle any shutdown errors
	if err := app.Stop(ctx); err != nil {
		log.Log.Fatal().Err(err).Msg("Failed to stop the application")
		return err
	}

	return nil
}
