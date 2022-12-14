package handlers

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type IItemHandlers interface {
	CreateItem(ctx context.Context, item Item) (uuid.UUID, error)
	UpdateItem(ctx context.Context, item Item) error
	GetItem(ctx context.Context, id string) (Item, error)
	ItemsList(ctx context.Context, offset, limit int) ([]Item, error)
	ItemsQuantity(ctx context.Context) (int, error)
	SearchLine(ctx context.Context, param string) ([]Item, error)
}

type ICategoryHandlers interface {
	CreateCategory(ctx context.Context, category Category) (uuid.UUID, error)
	GetCategoryList(ctx context.Context) ([]Category, error)
}

type IUserHandlers interface {
	CreateUser(ctx context.Context, user User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
