package cash

import (
	"OnlineShopBackend/internal/models"
	"context"
)

type Cash interface {
	CheckCash(ctx context.Context, key string) bool
	CreateItemsCash(ctx context.Context, res []models.Item, key string) error
	CreateItemsQuantityCash(ctx context.Context, value int, key string) error
	GetItemsCash(ctx context.Context, key string) ([]models.Item, error)
	GetItemsQuantityCash(ctx context.Context, key string) (int, error)
}
