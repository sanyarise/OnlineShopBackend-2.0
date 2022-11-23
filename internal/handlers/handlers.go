package handlers

import (
	"OnlineShopBackend/internal/usecase"
	"log"

	"go.uber.org/zap"
)

type Handlers struct {
	repo *usecase.Storage
	l *zap.Logger
}

func NewHandlers(repo *usecase.Storage, logger *zap.Logger) *Handlers {
	log.Println("Enter in NewHandlers()")
	return &Handlers{repo: repo, l: logger}
}
