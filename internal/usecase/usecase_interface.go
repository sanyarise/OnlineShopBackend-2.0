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
	ItemsList(ctx context.Context, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error)
	ItemsQuantity(ctx context.Context) (int, error)
	ItemsQuantityInCategory(ctx context.Context, categoryName string) (int, error)
	SearchLine(ctx context.Context, param string, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error)
	GetItemsByCategory(ctx context.Context, categoryName string, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error)
	UpdateCash(ctx context.Context, id uuid.UUID, op string) error
	UpdateItemsInCategoryCash(ctx context.Context, newItem *models.Item, op string) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
	AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
	DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
	GetFavouriteItems(ctx context.Context, userId uuid.UUID, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error)
	ItemsQuantityInFavourite(ctx context.Context, userId uuid.UUID) (int, error)
	UpdateFavouriteItemsCash(ctx context.Context, userId uuid.UUID, itemId uuid.UUID, op string)
	SortItems(items []models.Item, sortType string, sortOrder string)
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

type ICartUsecase interface {
	GetCart(ctx context.Context, cartId uuid.UUID) (*models.Cart, error)
	DeleteItemFromCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error
	Create(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
	AddItemToCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error
	DeleteCart(ctx context.Context, cartId uuid.UUID) error
	GetCartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error)
}

type IUserUsecase interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetRightsId(ctx context.Context, name string) (*models.Rights, error)
	UpdateUserData(ctx context.Context, id uuid.UUID, user *models.User) (*models.User, error)
}
