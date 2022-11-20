package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"log"

	"github.com/google/uuid"
)

func (r *Pgrepo) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	log.Println("Enter in repository CreateItem()")
	var id uuid.UUID
	_ = r.db.QueryRowContext(ctx, `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item.Title,
		item.Category,
		item.Description,
		item.Price,
		item.Vendor,
	).Scan(&id)

	log.Printf("id is %v\n", id)
	return id, nil
}

func (r *Pgrepo) UpdateItem(ctx context.Context, item *models.Item) error {return nil}

func (r *Pgrepo) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) { return nil, nil }

func (r *Pgrepo) ItemsList(ctx context.Context) (chan models.Item, error) { return nil, nil }

func (r *Pgrepo) SearchLine(ctx context.Context, param string) (chan models.Item, error) {
	return nil, nil
}
