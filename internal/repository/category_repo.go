package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type categoryRepo struct {
	storage Storage
	logger  *zap.SugaredLogger
}

var _ CategoryStore = (*categoryRepo)(nil)

func (repo *categoryRepo) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	repo.logger.Debug("Enter in repository CreateCategory()")
	var id uuid.UUID
	pool := repo.storage.GetPool()
	row := pool.QueryRow(ctx, `INSERT INTO categories(name, description)
	values ($1, $2) RETURNING id`,
		category.Name,
		category.Description,
	)
	row.Scan(&id)

	repo.logger.Debugf("id is %v\n", id)
	return id, nil
}

func (repo *categoryRepo) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	repo.logger.Debug("Enter in repository GetCategoryList()")
	categoryChan := make(chan models.Category, 100)
	go func() {
		defer close(categoryChan)
		category := &models.Category{}

		pool := repo.storage.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT id, name, description FROM categories`)
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
