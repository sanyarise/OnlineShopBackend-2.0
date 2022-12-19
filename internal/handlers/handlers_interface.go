package handlers

import (
	"context"

	"github.com/google/uuid"
)

type IItemHandlers interface {
	CreateItem(ctx context.Context, item Item) (uuid.UUID, error)
	UpdateItem(ctx context.Context, item Item) error
	GetItem(ctx context.Context, id string) (Item, error)
	ItemsList(ctx context.Context, offset, limit int) ([]Item, error)
	ItemsQuantity(ctx context.Context) (int, error)
	SearchLine(ctx context.Context, param string, offset, limit int) ([]Item, error)
	GetItemsByCategory(ctx context.Context, categoryName string, offset, limit int) ([]Item, error)
}

type ICategoryHandlers interface {
	CreateCategory(ctx context.Context, category Category) (uuid.UUID, error)
	UpdateCategory(ctx context.Context, category Category) error
	GetCategory(ctx context.Context, id string) (Category, error)
	GetCategoryList(ctx context.Context) ([]Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID)(string, error)
	GetCategoryByName(ctx context.Context, name string) (Category, error)
}
