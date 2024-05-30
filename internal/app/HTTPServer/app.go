package HTTPApp

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
	timeoutMiddleware "wbnats/internal/controller/http-server/middleware/timeout"
	orderHTTPHandler "wbnats/internal/controller/http-server/order"
	orderService "wbnats/internal/services/order"
)

type App struct {
	log    *slog.Logger
	engine *gin.Engine
	port   string
}

func New(
	log *slog.Logger,
	port string,
	Timeout time.Duration,
	orderService *orderService.Order) *App {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(timeoutMiddleware.New(Timeout))

	r.GET("/orders/:id", orderHTTPHandler.NewOrderHandler(log, orderService))

	return &App{
		log:    log,
		engine: r,
		port:   port,
	}
}

func (a *App) Run() error {
	err := a.engine.Run(a.port)
	if err != nil {
		return err
	}
	return nil
}
