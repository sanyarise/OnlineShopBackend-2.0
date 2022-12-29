package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type categoryRepo struct {
	storage *PGres
	logger  *zap.SugaredLogger
}

type Category struct {
	Id          uuid.UUID
	Name        string
	Description string
	Image       string
	DeletedAt   *time.Time
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

	repoCategory := &Category{
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
	}
	var id uuid.UUID
	pool := repo.storage.GetPool()

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("can't create transaction: %s", err)
		return uuid.Nil, fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("can't rollback %s", err)
			}

		} else {
			repo.logger.Info("transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("can't commit %s", err)
			}
		}
	}()
	row := tx.QueryRow(ctx, `INSERT INTO categories(name, description, picture, deleted_at)
	values ($1, $2, $3, $4) RETURNING id`,
		repoCategory.Name,
		repoCategory.Description,
		repoCategory.Image,
		nil,
	)
	if err := row.Scan(&id); err != nil {
		repo.logger.Errorf("can't scan %s", err)
		return uuid.Nil, fmt.Errorf("can't scan %w", err)
	}

	repo.logger.Debugf("id is %v\n", id)
	return id, nil
}

func (repo *categoryRepo) UpdateCategory(ctx context.Context, category *models.Category) error {
	repo.logger.Debug("Enter in repository UpdateCategory()")
	pool := repo.storage.GetPool()
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("can't create transaction: %s", err)
		return fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("can't rollback %s", err)
			}

		} else {
			repo.logger.Info("transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("can't commit %s", err)
			}
		}
	}()

	_, err = tx.Exec(ctx, `UPDATE categories SET name=$1, description=$2, picture=$3 WHERE id=$4`,
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

func (repo *categoryRepo) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	repo.logger.Debug("Enter in repository GetCategory()")
	category := models.Category{}
	pool := repo.storage.GetPool()
	row := pool.QueryRow(ctx,
		`SELECT id, name, description, picture FROM categories WHERE deleted_at is null AND id = $1`, id)
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

func (repo *categoryRepo) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	repo.logger.Debug("Enter in repository GetCategoryByName()")
	category := models.Category{}
	pool := repo.storage.GetPool()
	row := pool.QueryRow(ctx,
		`SELECT id, name, description, picture FROM categories WHERE deleted_at is null AND name = $1`, name)
	err := row.Scan(
		&category.Id,
		&category.Name,
		&category.Description,
		&category.Image,
	)
	if err != nil {
		repo.logger.Errorf("error in rows scan get category by name: %s", err)
		return &models.Category{}, fmt.Errorf("error in rows scan get category by name: %w", err)
	}
	return &category, nil
}

func (repo *categoryRepo) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	repo.logger.Debug("Enter in repository GetCategoryList()")
	categoryChan := make(chan models.Category, 100)
	go func() {
		defer close(categoryChan)
		category := &models.Category{}

		pool := repo.storage.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT id, name, description, picture FROM categories WHERE deleted_at is null`)
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

func (repo *categoryRepo) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	repo.logger.Debug("Enter in repository DeleteCategory()")
	pool := repo.storage.GetPool()
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("can't create transaction: %s", err)
		return fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("can't rollback %s", err)
			}

		} else {
			repo.logger.Info("transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("can't commit %s", err)
			}
		}
	}()
	_, err = tx.Exec(ctx, `UPDATE categories SET deleted_at=$1 WHERE id=$2`,
		time.Now(), id)
	if err != nil {
		repo.logger.Errorf("error on delete category %s: %s", id, err)
		return fmt.Errorf("error on delete category %s: %w", id, err)
	}
	repo.logger.Infof("Category %s successfully deleted from database", id)
	return nil
}
