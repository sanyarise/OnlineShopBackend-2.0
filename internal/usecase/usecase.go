package usecase

import (
	"OnlineShopBackend/internal/repository"

	"go.uber.org/zap"
)

type Storage struct {
	store repository.ItemStore
	l *zap.Logger
}

func NewStorage(store repository.ItemStore, logger *zap.Logger) *Storage {
	return &Storage{store: store, l: logger}
}
