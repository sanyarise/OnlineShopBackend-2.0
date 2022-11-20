package delivery

import "OnlineShopBackend/internal/handlers"

type Delivery struct {
	h *handlers.Handlers
}

func NewDelivery(handlers *handlers.Handlers) *Delivery {
	return &Delivery{h: handlers}
}
