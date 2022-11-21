package delivery

import (
	"OnlineShopBackend/internal/handlers"
	"log"

	"go.uber.org/zap"
)

type Delivery struct {
	h *handlers.Handlers
	l *zap.Logger
}

func NewDelivery(handlers *handlers.Handlers, logger *zap.Logger) *Delivery {
	log.Println("Enter in NewDelivery()")
	return &Delivery{h: handlers, l: logger}
}
