package app

import (
	"context"
	"net/http"
	"tezos-delegation-service/config"
	"tezos-delegation-service/internal/api/handlers"
	"tezos-delegation-service/internal/api/server"
	"tezos-delegation-service/internal/db"
	"tezos-delegation-service/internal/poller"
	"tezos-delegation-service/logger"

	"go.uber.org/fx"
)

// NewApp creates and returns a new fx.App instance with all necessary components
// provided and lifecycle hooks registered.
func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			config.LoadConfig,   // Provide configuration
			db.NewDatabase,      // Provide database connection
			server.NewServer,    // Provide HTTP server
			handlers.NewHandler, // Provide HTTP handlers
			poller.NewPoller,    // Provide polling service
			logger.NewLogger,    // Provide logger
		),
		fx.Invoke(
			registerHooks, // Register lifecycle hooks for server and poller
			startPoller,
		),
	)
}

// registerHooks registers the hooks for starting and stopping the HTTP server.
func registerHooks(lifecycle fx.Lifecycle, server *http.Server, zeroLogger *logger.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			zeroLogger.Log.Info().Msg("Starting HTTP server")

			go func() {
				// Start the HTTP server in a separate goroutine
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					zeroLogger.Log.Fatal().Err(err).Msg("Failed to start HTTP server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			zeroLogger.Log.Info().Msg("Stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})
}

// startPoller registers the hooks for starting and stopping the poller service.
func startPoller(p *poller.Poller, lifecycle fx.Lifecycle, zeroLogger *logger.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			zeroLogger.Log.Info().Msg("Starting Poller")

			go p.StartPolling(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			zeroLogger.Log.Info().Msg("Stopping Poller")
			// No specific stop logic needed for Poller in this example
			return nil
		},
	})
}
