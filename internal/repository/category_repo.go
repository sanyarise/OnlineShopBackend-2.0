package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type categoryRepo struct {
	storage *PGres
	logger  *zap.SugaredLogger
}

var _ CategoryStore = (*categoryRepo)(nil)

func NewCategoryRepo(store *PGres, log *zap.SugaredLogger) CategoryStore {
	return &categoryRepo{
		storage: store,
		logger:  log,
	}
}

func (repo *categoryRepo) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	repo.logger.Debug("Enter in repository CreateCategory()")
	select {
	case <-ctx.Done():
		return uuid.Nil, fmt.Errorf("context closed")
	default:
		var id uuid.UUID
		pool := repo.storage.GetPool()
		row := pool.QueryRow(ctx, `INSERT INTO categories(name, description, picture)
	values ($1, $2, $3) RETURNING id`,
			category.Name,
			category.Description,
			category.Image,
		)
		if err := row.Scan(&id); err != nil {
			repo.logger.Errorf("can't scan %s", err)
			return uuid.Nil, fmt.Errorf("can't scan %w", err)
		}

		repo.logger.Debugf("id is %v\n", id)
		return id, nil
	}
}

func (repo *categoryRepo) UpdateCategory(ctx context.Context, category *models.Category) error {
	repo.logger.Debug("Enter in repository UpdateCategory()")
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := repo.storage.GetPool()
		_, err := pool.Exec(ctx, `UPDATE categories SET name=$1, description=$2, picture=$3 WHERE id=$4`,
			category.Name,
			category.Description,
			category.Image,
			category.Id)
		if err != nil {
			repo.logger.Errorf("error on update category %s: %s", category.Id, err)
			return fmt.Errorf("error on update category %s: %w", category.Id, err)
		}
		repo.logger.Infof("category %s successfully updated", category.Id)
		return nil
	}
}

func (repo *categoryRepo) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	repo.logger.Debug("Enter in repository GetCategory()")
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context closed")
	default:
		category := models.Category{}
		pool := repo.storage.GetPool()
		row := pool.QueryRow(ctx,
			`SELECT id, name, description, picture FROM categories WHERE id = $1`, id)
		err := row.Scan(
			&category.Id,
			&category.Name,
			&category.Description,
			&category.Image,
		)
		if err != nil {
			repo.logger.Errorf("error in rows scan get category by id: %s", err)
			return &models.Category{}, fmt.Errorf("error in rows scan get category by id: %w", err)
		}
		return &category, nil
	}
}

func (repo *categoryRepo) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	repo.logger.Debug("Enter in repository GetCategoryList()")
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context closed")
	default:
		categoryChan := make(chan models.Category, 100)
		go func() {
			defer close(categoryChan)
			category := &models.Category{}

			pool := repo.storage.GetPool()
			rows, err := pool.Query(ctx, `
		SELECT id, name, description, picture FROM categories`)
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
					&category.Image,
				); err != nil {
					repo.logger.Error(err.Error())
					return
				}
				categoryChan <- *category
			}
		}()

		return categoryChan, nil
	}
}
