package repository

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

func (r *Pgrepo) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error)

func (r *Pgrepo) UpdateItem(ctx context.Context, item *models.Item) error

func (r *Pgrepo) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error)

func (r *Pgrepo) ItemsList(ctx context.Context) (chan models.Item, error)

func (r *Pgrepo) SearchLine(ctx context.Context, param string) (chan models.Item, error)
