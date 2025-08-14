package app

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	httphandler "github.com/shanth1/authorization/internal/adapters/input/http"
	"github.com/shanth1/authorization/internal/adapters/output/auth/telegram"
	"github.com/shanth1/authorization/internal/adapters/output/repository/inmemory"
	"github.com/shanth1/authorization/internal/adapters/output/token"
	authcfg "github.com/shanth1/authorization/internal/config/auth"
	"github.com/shanth1/authorization/internal/core/ports"
	"github.com/shanth1/authorization/internal/core/service"
	"github.com/shanth1/gotools/log"
)

func Run(ctx, shutdownCtx context.Context, cfg *authcfg.Config) {
	logger := log.FromContext(ctx)

	telegramProvider, err := telegram.NewTelegramProvider(cfg.Telegram.BotToken, cfg.Telegram.BotName)
	if err != nil {
		logger.Fatal().Err(err).Msg("new telegram provider")
	}

	userRepo := inmemory.NewMemoryUserRepository()
	_ = inmemory.NewInMemoryCache()

	_ = token.NewJWTService(cfg.JWT.SecretKey)
	service.NewOIDCService(nil, nil, nil, nil)
	_ = service.NewAuthService(userRepo, nil, []ports.AuthProvider{telegramProvider})

	router := gin.New()
	router.Use(gin.Recovery())

	// TODO: cors config
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // TODO: config
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Internal-Secret"}
	router.Use(cors.New(config))

	httpHandler := httphandler.NewOIDCHandlers(nil, nil, nil, nil, nil)
	httpHandler.RegisterRoutes(router)

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
