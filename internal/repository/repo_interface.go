package repository

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type ItemStore interface {
	CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error)
	UpdateItem(ctx context.Context, item *models.Item) error
	GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error)
	ItemsList(ctx context.Context) (chan models.Item, error)
	SearchLine(ctx context.Context, param string) (chan models.Item, error)
	GetItemsByCategory(ctx context.Context, categoryName string) (chan models.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
}

type CategoryStore interface {
	CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetCategoryList(ctx context.Context) (chan models.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
}

type UserStore interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type CartStore interface {
	Create(ctx context.Context, cart *models.Cart) (*models.Cart, error)
	AddItemToCart(ctx context.Context, cart *models.Cart, item *models.Item) error
	DeleteCart(ctx context.Context, cart *models.Cart) error
	DeleteItemFromCart(ctx context.Context, cart *models.Cart, item *models.Item) error
	SelectItemsFromCart(ctx context.Context, cart *models.Cart) ([]*models.Item, error)
}

type OrderStore interface {
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
	DeleteOrder(ctx context.Context, order *models.Order) error
	ChangeAddress(ctx context.Context, order *models.Order, address models.UserAddress) error
	ChangeStatus(ctx context.Context, order *models.Order, status models.Status) error
	GetOrderByID(ctx context.Context, id uuid.UUID) (models.Order, error)
	GetOrdersForUser(ctx context.Context, user *models.User) (chan models.Order, error)
}
