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
	ItemsList(ctx context.Context, number int) (chan models.Item, error)
	SearchLine(ctx context.Context, param string, number int) (chan models.Item, error)
}

type CategoryStore interface {
	CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error)
	GetCategoryList(ctx context.Context) (chan models.Category, error)
}

type UserStore interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type CartStore interface {
	Create(ctx context.Context, cart *models.Cart) (*models.Cart, error)
	//? Maybe add to item
	AddItemToCart(ctx context.Context, cart *models.Cart, item *models.Item) error
	DeleteCart(cxt context.Context, cart *models.Cart) error
	DeleteItemFromCart(ctx context.Context, cart *models.Cart, item *models.Item) error
}

type OrderStore interface {
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
	DeleteOrder(ctx context.Context, order *models.Order) error
	ChangeAddress(ctx context.Context, order *models.Order, address models.UserAddress) error
	ChangeStatus(ctx context.Context, order *models.Order, status models.Status) error
	GetOrderByID(ctx context.Context, id uuid.UUID) (models.Order, error)
	GetOrdersForUser(ctx context.Context, user *models.User) (chan models.Order, error)
}
