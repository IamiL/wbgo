package orderService

import (
	"context"
	"fmt"
	"log/slog"
	"wbnats/internal/services/order/models"
)

type Order struct {
	log         *slog.Logger
	ordSaver    OrderSaver
	ordProvider OrderProvider
}

type OrderSaver interface {
	SaveOrder(order *models.Order) (err error)
}

type OrderProvider interface {
	Order(ctx context.Context, email string) (models.Order, error)
}

func New(
	log *slog.Logger,
	ordSaver OrderSaver,
	ordProvider OrderProvider,
) *Order {
	return &Order{
		log:         log,
		ordSaver:    ordSaver,
		ordProvider: ordProvider,
	}
}

func (o *Order) NewOrder(order *models.Order) error {
	const op = "Order.NewOrder"

	log := o.log.With(
		slog.String("op", op),
		slog.String("orderUID", order.UID),
	)

	log.Info("processing a new order")

	err := o.ordSaver.SaveOrder(order)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (o *Order) Order(ctx context.Context, uid string) (models.Order, error) {
	const op = "Order.Order"

	log := o.log.With(
		slog.String("op", op),
		slog.String("orderUID", uid),
	)

	log.Info("getting order information")
	order, err := o.ordProvider.Order(ctx, uid)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: %w", op, err)
	}
	return order, nil
}
