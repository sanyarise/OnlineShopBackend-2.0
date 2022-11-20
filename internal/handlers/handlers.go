package handlers

import (
	"OnlineShopBackend/internal/usecase"
	"log"
)

type Handlers struct {
	repo *usecase.Storage
}

func NewHandlers(repo *usecase.Storage) *Handlers {
	log.Println("Enter in NewHandlers()")
	return &Handlers{repo: repo}
}
