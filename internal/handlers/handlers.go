package handlers

import (
	"OnlineShopBackend/internal/usecase"

	"go.uber.org/zap"
)

type Handlers struct {
	repo   *usecase.Storage
	logger *zap.Logger
}

func NewHandlers(repo *usecase.Storage, logger *zap.Logger) *Handlers {
	logger.Debug("Enter in NewHandlers()")
	return &Handlers{repo: repo, logger: logger}
}
