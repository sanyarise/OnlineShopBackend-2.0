package delivery

import (
	"OnlineShopBackend/internal/handlers"
	"log"
)

type Delivery struct {
	h *handlers.Handlers
}

func NewDelivery(handlers *handlers.Handlers) *Delivery {
	log.Println("Enter in NewDelivery()")
	return &Delivery{h: handlers}
}
