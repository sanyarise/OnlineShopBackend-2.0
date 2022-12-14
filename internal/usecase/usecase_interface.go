package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type IItemUsecase interface {
	CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error)
	UpdateItem(ctx context.Context, item *models.Item) error
	GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error)
	ItemsList(ctx context.Context, offset, limit int) ([]models.Item, error)
	ItemsQuantity(ctx context.Context) (int, error)
	SearchLine(ctx context.Context, param string, offset, limit int) ([]models.Item, error)
	GetItemsByCategory(ctx context.Context, categoryName string, offset, limit int) ([]models.Item, error)
	UpdateCash(ctx context.Context, id uuid.UUID, op string) error
}

type ICategoryUsecase interface {
	CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetCategoryList(ctx context.Context) (chan models.Category, error)
}
