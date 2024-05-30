package natsStreamingApp

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"log/slog"
	orderNatsStreaming "wbnats/internal/controller/nutsServer/order"
	orderService "wbnats/internal/services/order"
)

type App struct {
	log               *slog.Logger
	natsStreamConnect *stan.Conn
	orderService      *orderService.Order
	sub               *stan.Subscription
}
type Order interface {
}

func New(
	log *slog.Logger,
	clusterID string,
	clientID string,
	orderService *orderService.Order,
) *App {
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		{
			panic("nats streaming server connection error")
		}
	}

	return &App{
		log:               log,
		natsStreamConnect: &sc,
		orderService:      orderService,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const op = "natsStreamingApp.Run"

	sub, err := (*a.natsStreamConnect).Subscribe("foo", orderNatsStreaming.NewOrderSaverHandler(a.log, a.orderService))
	a.sub = &sub
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("nats streaming server started")
	return nil
}

func (a *App) Stop() {
	const op = "natsStreamingApp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping nats streaming server")

	(*a.sub).Unsubscribe()
	(*a.natsStreamConnect).Close()
}
