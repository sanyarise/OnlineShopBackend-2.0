package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type ItemRepo struct {
	db *sql.DB
}

func NewItemRepo(db *sql.DB) *ItemRepo {
	return &ItemRepo{db: db}
}

func (ir *ItemRepo) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error)

func (ir *ItemRepo) UpdateItem(ctx context.Context, item *models.Item) error

func (ir *ItemRepo) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error)

func (ir *ItemRepo) ItemsList(ctx context.Context, s string) ([]*models.Item, error)
