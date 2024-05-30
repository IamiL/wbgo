package orderNatsStreaming

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log/slog"
	"time"
	orderNatsStreaming "wbnats/internal/controller/nutsServer/order/models"
	orderService "wbnats/internal/services/order"
	"wbnats/internal/services/order/models"
)

type OrderService interface {
	NewOrder(order *models.Order)
}

func NewOrderSaverHandler(log *slog.Logger, orderSaver *orderService.Order) func(*stan.Msg) {
	return func(m *stan.Msg) {

		newOrder := orderNatsStreaming.Order{}

		err := json.Unmarshal(m.Data, &newOrder)
		if err != nil {
			log.Info("Deserialization error", err)
			return
		}
		if newOrder.DateCreated == "" {
			log.Info("Deserialization error")
			return
		}
		dateCreated, err := time.Parse("2006-01-02T15:04:05Z", newOrder.DateCreated)

		its := []models.Item{}

		for _, item := range newOrder.Items {
			its = append(its, models.Item{
				ChrtID:      item.ChrtID,
				TrackNumber: item.TrackNumber,
				Price:       item.Price,
				RID:         item.RID,
				Name:        item.Name,
				Sale:        item.Sale,
				Size:        item.Size,
				TotalPrice:  item.TotalPrice,
				NmID:        item.NmID,
				Brand:       item.Brand,
				Status:      item.Status,
			})
		}
		(*orderSaver).NewOrder(&models.Order{
			UID:         newOrder.UID,
			TrackNumber: newOrder.TrackNumber,
			Entry:       newOrder.Entry,
			Delivery: models.Delivery{
				Name:    newOrder.Delivery.Name,
				Phone:   newOrder.Delivery.Phone,
				Zip:     newOrder.Delivery.Zip,
				City:    newOrder.Delivery.City,
				Address: newOrder.Delivery.Address,
				Region:  newOrder.Delivery.Region,
				Email:   newOrder.Delivery.Email,
			},
			Payment: models.Payment{
				Transaction:  newOrder.Payment.Transaction,
				RequestID:    newOrder.Payment.RequestID,
				Currency:     newOrder.Payment.Currency,
				Provider:     newOrder.Payment.Provider,
				Amount:       newOrder.Payment.Amount,
				PaymentDT:    newOrder.Payment.PaymentDT,
				Bank:         newOrder.Payment.Bank,
				DeliveryCost: newOrder.Payment.DeliveryCost,
				GoodsTotal:   newOrder.Payment.GoodsTotal,
				CustomFee:    newOrder.Payment.CustomFee,
			},
			Items:             its,
			Locale:            newOrder.Locale,
			InternalSignature: newOrder.InternalSignature,
			CustomerID:        newOrder.CustomerID,
			DeliveryService:   newOrder.DeliveryService,
			Shardkey:          newOrder.Shardkey,
			SmID:              newOrder.SmID,
			DateCreated:       dateCreated,
			OofShard:          newOrder.OofShard,
		})

	}
}
