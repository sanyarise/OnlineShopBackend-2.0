package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (repo *Pgrepo) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	repo.logger.Debug("Enter in repository CreateCategory()")
	var id uuid.UUID
	_ = repo.db.QueryRowContext(ctx, `INSERT INTO categories(name, description)
	values ($1, $2) RETURNING id`,
		category.Name,
		category.Description,
	).Scan(&id)

	repo.logger.Debug(fmt.Sprintf("id is %v\n", id))
	return id, nil
}

func (repo *Pgrepo) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	repo.logger.Debug("Enter in repository GetCategoryList()")
	categoryChan := make(chan models.Category, 100)
	go func() {
		defer close(categoryChan)
		category := &models.Category{}

		rows, err := repo.db.QueryContext(ctx, `
		SELECT id, name,description FROM categories`)
		if err != nil {
			msg := fmt.Errorf("error on categories list query context: %w", err)
			repo.logger.Error(msg.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&category.Id,
				&category.Name,
				&category.Description,
			); err != nil {
				repo.logger.Error(err.Error())
				return
			}
			categoryChan <- *category
		}
	}()

	return categoryChan, nil
}
