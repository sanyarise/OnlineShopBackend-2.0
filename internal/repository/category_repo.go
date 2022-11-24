package repository

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

func (r *Pgrepo) CreateCategory(ctx context.Context, item *models.Category) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (r *Pgrepo) GetCategoryList(ctx context.Context) (chan models.Category, error) { return nil, nil }
