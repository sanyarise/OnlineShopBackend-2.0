package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"strings"
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
	repo.logger.Debugf("Enter in repository CreateCategory() with args: ctx, category: %v", category)
	repoCategory := &Category{
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
	}
	pool := repo.storage.GetPool()

	if id, ok := repo.isDeletedCategory(ctx, repoCategory.Name); ok {
		_, err := pool.Exec(ctx, `UPDATE categories SET description=$1, picture=$2, deleted_at=null WHERE name=$3`,
			repoCategory.Description,
			repoCategory.Image,
			repoCategory.Name)
		if err != nil {
			repo.logger.Debug(err.Error())
			return uuid.Nil, err
		}
		return id, nil
	}
	var id uuid.UUID
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("Can't create transaction: %s", err)
		return uuid.Nil, fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("Can't rollback %s", err)
			}

		} else {
			repo.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("Can't commit %s", err)
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
		repo.logger.Errorf("Can't scan %s", err)
		return uuid.Nil, fmt.Errorf("can't scan %w", err)
	}
	repo.logger.Debug("Category created success")
	repo.logger.Debugf("Id is %v\n", id)
	return id, nil
}

func (repo *categoryRepo) isDeletedCategory(ctx context.Context, name string) (uuid.UUID, bool) {
	repo.logger.Debug("Enter in repository is DeletedCategory() with args: ctx, name: %s", name)
	pool := repo.storage.GetPool()
	category := models.Category{}
	row := pool.QueryRow(ctx,
		`SELECT id FROM categories WHERE deleted_at is not null AND name = $1`, name)
	err := row.Scan(
		&category.Id,
	)
	if err == nil && category.Id != uuid.Nil {
		return category.Id, true
	}
	repo.logger.Error(err.Error())
	return uuid.Nil, false
}

func (repo *categoryRepo) UpdateCategory(ctx context.Context, category *models.Category) error {
	repo.logger.Debugf("Enter in repository UpdateCategory() with args: ctx, category: %v", category)
	pool := repo.storage.GetPool()
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("Can't create transaction: %s", err)
		return fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("Can't rollback %s", err)
			}

		} else {
			repo.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("Can't commit %s", err)
			}
		}
	}()

	_, err = tx.Exec(ctx, `UPDATE categories SET name=$1, description=$2, picture=$3 WHERE id=$4`,
		category.Name,
		category.Description,
		category.Image,
		category.Id)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("Error on update category %s: %s", category.Id, err)
		return models.ErrorNotFound{Message: fmt.Sprintf("category with id: %v not found", category.Id)}
	}
	if err != nil {
		repo.logger.Errorf("Error on update category %s: %s", category.Id, err)
		return fmt.Errorf("error on update category %s: %w", category.Id, err)
	}
	repo.logger.Infof("Category with id %s successfully updated", category.Id)
	return nil
}

func (repo *categoryRepo) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	repo.logger.Debugf("Enter in repository GetCategory() with args: ctx, id: %v", id)
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
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("Error in rows scan get category by id: %s", err)
		return &models.Category{}, models.ErrorNotFound{Message: fmt.Sprintf("category with id: %v not found", id)}
	}
	if err != nil {
		repo.logger.Errorf("Error in rows scan get category by id: %s", err)
		return &models.Category{}, fmt.Errorf("error in rows scan get category by id: %w", err)
	}
	repo.logger.Info("Get category success")
	return &category, nil
}

func (repo *categoryRepo) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	repo.logger.Debugf("Enter in repository GetCategoryByName() with args: ctx, name: %s", name)
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
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("Error in rows scan get category by name: %s", err)
		return &models.Category{}, models.ErrorNotFound{Message: fmt.Sprintf("category with name: %s not found", name)}
	}
	if err != nil {
		repo.logger.Errorf("Error in rows scan get category by name: %s", err)
		return &models.Category{}, fmt.Errorf("error in rows scan get category by name: %w", err)
	}
	repo.logger.Info("Get category by name success")
	return &category, nil
}

func (repo *categoryRepo) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	repo.logger.Debug("Enter in repository GetCategoryList() with args: ctx")
	categoryChan := make(chan models.Category, 100)
	go func() {
		defer close(categoryChan)
		category := &models.Category{}

		pool := repo.storage.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT id, name, description, picture FROM categories WHERE deleted_at is null`)
		if err != nil {
			repo.logger.Error(fmt.Errorf("error on categories list query context: %w", err).Error())
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
	repo.logger.Debugf("Enter in repository DeleteCategory() with args: ctx, id: %v", id)
	pool := repo.storage.GetPool()
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("Can't create transaction: %s", err)
		return fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("Can't rollback %s", err)
			}

		} else {
			repo.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("Can't commit %s", err)
			}
		}
	}()
	_, err = tx.Exec(ctx, `UPDATE categories SET deleted_at=$1 WHERE id=$2`,
		time.Now(), id)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("Error on delete category %s: %s", id, err)
		return models.ErrorNotFound{Message: fmt.Sprintf("category with id: %v not found", id)}
	}
	if err != nil {
		repo.logger.Errorf("Error on delete category %s: %s", id, err)
		return fmt.Errorf("error on delete category %s: %w", id, err)
	}
	repo.logger.Infof("Category with id: %s successfully deleted from database", id)
	return nil
}
