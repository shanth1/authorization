package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	httphandler "github.com/shanth1/authorization/internal/adapters/input/http"
	"github.com/shanth1/authorization/internal/adapters/output/repository/inmemory"
	"github.com/shanth1/authorization/internal/adapters/output/token"
	authcfg "github.com/shanth1/authorization/internal/config/auth"
	"github.com/shanth1/authorization/internal/core/service"
	"github.com/shanth1/gotools/log"
)

func Run(ctx, shutdownCtx context.Context, cfg *authcfg.Config) {
	logger := log.FromContext(ctx)

	userRepo := inmemory.NewInMemoryUserRepository()
	cache := inmemory.NewInMemoryCache()
	tokenSvc := token.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	authSvc := service.NewAuthService(userRepo, cache, cache, tokenSvc, cfg.Telegram.BotName)

	router := gin.New()
	router.Use(gin.Recovery())
	httpHandler := httphandler.NewHandler(authSvc)
	httpHandler.InitRoutes(router, cfg.HTTPServer.APIToken)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("failed to start http server")
		}
	}()
	logger.Info().Str("address", cfg.HTTPServer.Address).Msg("server started")

	<-ctx.Done()

	logger.Info().Msg("shutting down server...")

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal().Err(err).Msg("server shutdown failed")
	}

	logger.Info().Msg("server exited gracefully")
}
