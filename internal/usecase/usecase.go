package usecase

import "OnlineShopBackend/internal/repository"

type Storage struct {
	store repository.ItemStore
}

func NewStorage(store repository.ItemStore) *Storage {
	return &Storage{store: store}
}
