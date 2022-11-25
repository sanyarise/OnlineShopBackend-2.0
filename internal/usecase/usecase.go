package usecase

import (
	"OnlineShopBackend/internal/repository"

	"go.uber.org/zap"
)

type Storage struct {
	itemStore     repository.ItemStore
	categoryStore repository.CategoryStore
	logger        *zap.Logger
}

func NewStorage(itemStore repository.ItemStore, categoryStore repository.CategoryStore, logger *zap.Logger) *Storage {
	return &Storage{itemStore: itemStore, categoryStore: categoryStore, logger: logger}
}
