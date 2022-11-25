package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Pgrepo) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	r.l.Debug("Enter in repository CreateCategory()")
	var id uuid.UUID
	_ = r.db.QueryRowContext(ctx, `INSERT INTO categories(name, description)
	values ($1, $2) RETURNING id`,
		category.Name,
		category.Description,
	).Scan(&id)

	r.l.Debug(fmt.Sprintf("id is %v\n", id))
	return id, nil
}

func (r *Pgrepo) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	chout := make(chan models.Category, 100)

	go func() {
		defer close(chout)
		category := &models.Category{}

		rows, err := r.db.QueryContext(ctx, `
		SELECT id, name,description FROM categories`)
		if err != nil {
			msg := fmt.Errorf("error on categories list query context: %w", err)
			r.l.Error(msg.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&category.Id,
				&category.Name,
				&category.Description,
			); err != nil {
				r.l.Error(err.Error())
				return
			}
			fmt.Println(category)
			chout <- *category
		}
	}()

	return chout, nil
}
