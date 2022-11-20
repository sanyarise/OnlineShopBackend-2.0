package handlers

import "OnlineShopBackend/internal/usecase"

type Handlers struct {
	repo *usecase.Storage
}

func NewHandlers(repo *usecase.Storage) *Handlers {
	return &Handlers{repo: repo}
}
