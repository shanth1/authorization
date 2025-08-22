package app

import (
	"context"
	"net/http"

	"github.com/shanth1/authorization/internal/adapters/input/httphandler"
	"github.com/shanth1/authorization/internal/config"
	"github.com/shanth1/gotools/log"
)

func Run(ctx, shutdownCtx context.Context, cfg *config.Config) {
	logger := log.FromContext(ctx)
	logger = logger.With(log.Str("env", cfg.Env))

	httpHandler := httphandler.New(cfg)
	router := httpHandler.SetupRouter()

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
