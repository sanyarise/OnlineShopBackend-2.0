package delivery

import (
	"OnlineShopBackend/internal/filestorage"
	"OnlineShopBackend/internal/handlers"
	"log"

	"go.uber.org/zap"
)

type Delivery struct {
	h  *handlers.Handlers
	l  *zap.Logger
	fs filestorage.FileStorager
}

func NewDelivery(handlers *handlers.Handlers, logger *zap.Logger, fs filestorage.FileStorager) *Delivery {
	log.Println("Enter in NewDelivery()")
	return &Delivery{h: handlers, l: logger, fs: fs}
}
