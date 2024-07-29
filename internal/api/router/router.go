package router

import (
	"tezos-delegation-service/config"
	"tezos-delegation-service/internal/api"
	"tezos-delegation-service/internal/api/handlers"
	"tezos-delegation-service/internal/middleware"
	"tezos-delegation-service/logger"

	pagination "github.com/webstradev/gin-pagination"

	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(cfg *config.Config, r *gin.Engine, handler *handlers.Handler) {
	// middlewares
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(logger.NewLogger()))
	r.Use(middleware.NewCors().Middleware())
	r.Use(middleware.NewRateLimiter(cfg.RateLimit).Middleware())

	paginator := pagination.New("page", "rowsPerPage", "1", "15", 5, 150)
	r.Use(paginator)

	// openapi specs
	swagger, _ := api.GetSwagger()

	// validator
	r.Use(oapiMiddleware.OapiRequestValidator(swagger))

	api.RegisterHandlers(r, handler)
}
