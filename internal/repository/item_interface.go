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
}
