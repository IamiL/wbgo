package orderService

import (
	"context"
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

func (o *Order) NewOrder(order *models.Order) {
	err := o.ordSaver.SaveOrder(order)
	if err != nil {
		o.log.Error("asd", err)
		return
	}
}

func (o *Order) Order(ctx context.Context, uid string) (models.Order, error) {

	order, err := o.ordProvider.Order(ctx, uid)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}
