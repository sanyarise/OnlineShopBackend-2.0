package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
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

func (r *Pgrepo) UpdateItem(ctx context.Context, item *models.Item) error {
	_, err := r.db.ExecContext(ctx, `UPDATE items SET name=$1, category=$2, description=$3, price=$4, vendor=$5 WHERE id=$6`,
		item.Title,
		item.Category,
		item.Description,
		item.Price,
		item.Vendor,
		item.Id)
	if err != nil {
		return fmt.Errorf("error on update item: %w)", err)
	}
	return nil
}

func (r *Pgrepo) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	item := models.Item{}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, category, description, price, vendor FROM items WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("error on get item: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.Category,
			&item.Description,
			&item.Price,
			&item.Vendor,
		); err != nil {
			return nil, fmt.Errorf("error in rows scan get item by id: %w", err)
		}
	}
	if uuid.UUID.String(item.Id) == "00000000-0000-0000-0000-000000000000" {
		err = fmt.Errorf("id not found")
		return nil, err
	}

	return &item, nil
}

func (r *Pgrepo) ItemsList(ctx context.Context) (chan models.Item, error) {
	chout := make(chan models.Item, 100)

	go func() {
		defer close(chout)
		item := &models.Item{}

		rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, category, description, price, vendor FROM items`)
		if err != nil {
			msg := fmt.Errorf("error on items list query context: %w", err)
			r.l.Error(msg.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&item.Id,
				&item.Title,
				&item.Category,
				&item.Description,
				&item.Price,
				&item.Vendor,
			); err != nil {
				r.l.Error(err.Error())
				return
			}
			fmt.Println(item)
			chout <- *item
		}
	}()

	return chout, nil
}

func (r *Pgrepo) SearchLine(ctx context.Context, param string) (chan models.Item, error) {
	chout := make(chan models.Item, 100)

	go func() {
		defer close(chout)
		item := &models.Item{}

		rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, category, description, price, vendor FROM items WHERE name LIKE $1 OR description LIKE $1 OR vendor LIKE $1`,
			"%"+param+"%")
		if err != nil {
			msg := fmt.Errorf("error on search line query context: %w", err)
			r.l.Error(msg.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&item.Id,
				&item.Title,
				&item.Category,
				&item.Description,
				&item.Price,
				&item.Vendor,
			); err != nil {
				r.l.Error(err.Error())
				return
			}
			fmt.Println(item)
			chout <- *item
		}
	}()

	return chout, nil
}
