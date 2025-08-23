package app

import (
	"context"
	"net/http"

	"github.com/shanth1/authorization/client/internal/config"
	"github.com/shanth1/authorization/client/internal/server"
	"github.com/shanth1/gotools/log"
)

func Run(ctx, shutdownCtx context.Context, cfg *config.Config) {
	logger := log.FromContext(ctx)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: server.NewHandler(cfg),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	logger.Info().Str("address", cfg.Address).Msg("server started")

	<-ctx.Done()

	logger.Info().Msg("shutting down server...")

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal().Err(err).Msg("server shutdown failed")
	}

	logger.Info().Msg("server exited gracefully")
}
