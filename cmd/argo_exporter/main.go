package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Systemnick/argo_exporter/internal/indications/repository/indications/mur_tcp"
	"github.com/Systemnick/argo_exporter/internal/indications/useCase/indications"
	"github.com/Systemnick/argo_exporter/internal/registrars/repository/registrars/static"
	"github.com/Systemnick/argo_exporter/internal/registrars/useCase/registrars"
	"github.com/Systemnick/argo_exporter/pkg/config"
	"github.com/Systemnick/argo_exporter/pkg/http/server"
)

var (
	ServiceName = "argo_exporter"
	Version     = "development"
)

func main() {
	logger := initLogger()

	cfg, _ := config.New(ServiceName, Version)

	logger.Info().Msgf("App %s, version: %s", cfg.ServiceName, cfg.Version)

	registrarRepo := static.NewRepository(logger)
	regs := registrars.NewUseCase(logger, registrarRepo)

	// indicationsRepoFabric := mur_tcp.NewRepositoryFabric()
	// indicationsUseCase := indications.NewUseCase(regs, indicationsRepoFabric)

	indicationsRepo := mur_tcp.NewRepository(logger)
	indicationsUseCase := indications.NewUseCase(logger, regs, indicationsRepo)

	actualIndications, cancel := indicationsUseCase.ScrapePeriodic(context.Background(), cfg.ScrapeInterval)
	defer cancel()

	server.Start(logger, cfg, actualIndications)
}

func initLogger() *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	zerolog.DurationFieldUnit = time.Nanosecond

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		With().Caller().Logger().
		With().Timestamp().Logger()

	return &logger
}

func initContext() context.Context {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	zerolog.DurationFieldUnit = time.Nanosecond

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		With().Caller().Logger().
		With().Timestamp().Logger()

	// logger := log.With().Str("component", "module").Logger()
	ctx := logger.WithContext(context.Background())

	return ctx
}
