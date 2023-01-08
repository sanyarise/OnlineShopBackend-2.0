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
	ItemsQuantityInCategory(ctx context.Context, categoryName string) (int, error)
	SearchLine(ctx context.Context, param string, offset, limit int) ([]models.Item, error)
	GetItemsByCategory(ctx context.Context, categoryName string, offset, limit int) ([]models.Item, error)
	UpdateCash(ctx context.Context, id uuid.UUID, op string) error
	UpdateItemsInCategoryCash(ctx context.Context, newItem *models.Item, op string) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
}

type ICategoryUsecase interface {
	CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetCategoryList(ctx context.Context) ([]models.Category, error)
	UpdateCash(ctx context.Context, id uuid.UUID, op string) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
	DeleteCategoryCash(ctx context.Context, name string) error
}

type IOrderUsecase interface {
	PlaceOrder(ctx context.Context, cart *models. Cart, user *models.User) (*models.Order, error)
	ChangeStatus(ctx context.Context, order *models.Order, newStatus models.Status) error
	GetOrdersForUser(ctx context.Context, user *models.User) ([]models.Order, error)
	DeleteOrder(ctx context.Context, order *models.Order) error
	ChangeAddress(ctx context.Context, order *models.Order, newAddress models.UserAddress) error
	GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error)
}