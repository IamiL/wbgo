package orderHTTPHandler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	orderService "wbnats/internal/services/order"
	"wbnats/internal/services/order/models"
)

type OrderService interface {
	Order(ctx context.Context, uid string) (*models.Order, error)
}

func NewOrderHandler(log *slog.Logger, order *orderService.Order) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid := c.Param("id")
		ord, err := (*order).Order(c, uid)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
			return
		}
		c.JSON(http.StatusOK, ord)
		return
	}
}
