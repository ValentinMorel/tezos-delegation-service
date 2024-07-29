package server

import (
	"fmt"
	"net/http"
	"time"

	"tezos-delegation-service/config"
	"tezos-delegation-service/internal/api/handlers"
	"tezos-delegation-service/internal/api/router"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Cfg    *config.Config
	engine *gin.Engine
}

func NewServer(cfg *config.Config, handler *handlers.Handler) *http.Server {
	engine := gin.New()
	router.RegisterRoutes(cfg, engine, handler)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Port),
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return srv
}
