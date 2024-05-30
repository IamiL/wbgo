package app

import (
	"log/slog"
	HTTPApp "wbnats/internal/app/HTTPServer"
	natsStreamingApp "wbnats/internal/app/natsStreaming"
	"wbnats/internal/config"
	"wbnats/internal/repository/postgres"
	"wbnats/internal/services/order"
)

type App struct {
	NatsStreaming *natsStreamingApp.App
	HTTPServer    *HTTPApp.App
}

func New(
	log *slog.Logger,
	clusterID string,
	clientID string,
	dbConfig config.PostgresConfig,
	HTTPConfig config.HTTPServer,
) *App {
	storage, err := postgres.New(dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.User, dbConfig.Pass)
	if err != nil {
		panic(err)
	}

	order := orderService.New(log, storage, storage)

	nutsApp := natsStreamingApp.New(log, clusterID, clientID, order)

	httpApp := HTTPApp.New(log, HTTPConfig.Port, HTTPConfig.Timeout, order)

	storage.RestoreCache()

	return &App{
		NatsStreaming: nutsApp,
		HTTPServer:    httpApp,
	}
}
