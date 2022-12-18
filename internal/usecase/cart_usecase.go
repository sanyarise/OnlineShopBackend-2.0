package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var Cart *CartUseCase

type CartUseCase struct {
	store  repository.CartStore
	logger *zap.Logger
}

func NewCartUseCase(store *repository.CartStore, logger *zap.Logger) *CartUseCase {
	Cart = &CartUseCase{store: *store, logger: logger}
	return Cart
}

func (c *CartUseCase) GetCart(ctx context.Context, cartID uuid.UUID) (*models.Cart, error) {
	cart := &models.Cart{ID: cartID}
	items, err := c.store.SelectItemsFromCart(ctx, cart)
	if err != nil {
		return nil, err
	}
	cart.Items = items
	return cart, nil
}

func (c *CartUseCase) DeleteItemFromCart(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID) error {
	cart := &models.Cart{ID: cartID}
	item := &models.Item{Id: itemID}
	err := c.store.DeleteItemFromCart(ctx, cart, item)
	if err != nil {
		return err
	}
	return nil
}

func (c *CartUseCase) Create(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	cart := &models.Cart{UserID: userID}
	cart, err := c.store.Create(ctx, cart)
	if err != nil {
		return cart.ID, err
	}
	return cart.ID, nil
}

func (c *CartUseCase) AddItemToCart(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID) error {
	cart := &models.Cart{ID: cartID}
	item := &models.Item{Id: itemID}
	err := c.store.AddItemToCart(ctx, cart, item)
	if err != nil {
		return err
	}
	return nil
}

func (c *CartUseCase) DeleteCart(ctx context.Context, cartID uuid.UUID) error {
	cart := &models.Cart{ID: cartID}
	err := c.store.DeleteCart(ctx, cart)
	if err != nil {
		return err
	}
	return nil
}
