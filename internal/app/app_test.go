package app

import (
	"context"
	"net/http"
	"os"
	"testing"
	"tezos-delegation-service/config"
	"tezos-delegation-service/internal/api/handlers"
	"tezos-delegation-service/internal/db"
	"tezos-delegation-service/internal/poller"
	"tezos-delegation-service/logger"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestNewApp(t *testing.T) {
	var cfg *config.Config
	var database db.Querier
	var srv *http.Server
	var handler *handlers.Handler
	var poll *poller.Poller
	var zeroLogger *logger.Logger

	app := fx.New(
		fx.Provide(
			config.LoadConfig,
			db.NewDatabase,
			func() *http.Server { return &http.Server{} }, // Assuming a basic HTTP server
			handlers.NewHandler,
			poller.NewPoller,
			func() zerolog.Logger {
				return zerolog.New(os.Stdout)
			},
			logger.NewLogger,
		),
		fx.Populate(&cfg, &database, &srv, &handler, &zeroLogger, &poll),
	)

	assert.NotNil(t, app)

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := app.Start(startCtx)
	assert.NoError(t, err)

	assert.NotNil(t, cfg)
	assert.NotNil(t, database)
	assert.NotNil(t, srv)
	assert.NotNil(t, handler)
	assert.NotNil(t, zeroLogger)
	assert.NotNil(t, poll)

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = app.Stop(stopCtx)
	assert.NoError(t, err)
}
