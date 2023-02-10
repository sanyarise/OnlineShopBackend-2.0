package cash

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type IItemsCash interface {
	CheckCash(ctx context.Context, key string) bool
	CreateItemsCash(ctx context.Context, res []models.Item, key string) error
	CreateItemsQuantityCash(ctx context.Context, value int, key string) error
	GetItemsCash(ctx context.Context, key string) ([]models.Item, error)
	GetItemsQuantityCash(ctx context.Context, key string) (int, error)
	CreateFavouriteItemsIdCash(ctx context.Context, res map[uuid.UUID]uuid.UUID, key string) error
	GetFavouriteItemsIdCash(ctx context.Context, key string) (*map[uuid.UUID]uuid.UUID, error)
}

type ICategoriesCash interface {
	CheckCash(ctx context.Context, key string) bool
	CreateCategoriesListCash(ctx context.Context, categories []models.Category, key string) error
	GetCategoriesListCash(ctx context.Context, key string) ([]models.Category, error)
	DeleteCash(ctx context.Context, key string) error
}
