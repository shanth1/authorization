package main

import (
	"flag"
	"time"

	"github.com/shanth1/authorization/internal/app"
	"github.com/shanth1/authorization/internal/config"
	"github.com/shanth1/gotools/conf"
	"github.com/shanth1/gotools/ctx"
	"github.com/shanth1/gotools/env"
	"github.com/shanth1/gotools/flags"
	"github.com/shanth1/gotools/log"
)

type startCfg struct {
	EnvPath    string `flag:"env-path" usage:"[OPTIONAL] Path to env file"`
	ConfigPath string `flag:"config-path" usage:"[OPTIONAL] Path to the YAML config file"`
}

func main() {
	logger := log.New(log.WithService("auth"))

	ctx, shutdownCtx, cancel, shutdownCancel := ctx.WithGracefulShutdown(5 * time.Second)
	defer cancel()
	defer shutdownCancel()

	ctx = log.NewContext(ctx, logger)

	var startCfg startCfg
	if err := flags.RegisterFromStruct(&startCfg); err != nil {
		logger.Fatal().Err(err).Msg("Register flags from struct")
	}
	flag.Parse()

	cfg := &config.Config{}
	if err := env.LoadIntoStruct(startCfg.EnvPath, cfg); err != nil {
		logger.Fatal().Err(err).Msg("Load env into struct")
	}

	if err := conf.Load(startCfg.ConfigPath, cfg); err != nil {
		logger.Fatal().Err(err).Msg("Load config")
	}

	app.Run(ctx, shutdownCtx, cfg)
}
