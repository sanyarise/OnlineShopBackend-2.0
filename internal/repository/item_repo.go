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

type itemRepo struct {
	storage *PGres
	logger  *zap.SugaredLogger
}

func NewItemRepo(storage *PGres, logger *zap.SugaredLogger) ItemStore {
	return &itemRepo{
		storage: storage,
		logger:  logger,
	}
}

var _ ItemStore = (*itemRepo)(nil)

func (repo *itemRepo) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	repo.logger.Debug("Enter in repository CreateItem()")
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
	row := tx.QueryRow(ctx, `INSERT INTO items(name, category, description, price, vendor, pictures, deleted_at)
	values ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		item.Title,
		item.Category.Id,
		item.Description,
		item.Price,
		item.Vendor,
		item.Images,
		nil,
	)
	err = row.Scan(&id)
	if err != nil {
		repo.logger.Errorf("can't create item %s", err)
		return uuid.Nil, fmt.Errorf("can't create item %w", err)
	}
	repo.logger.Debugf("id is %v\n", id)
	return id, nil
}

func (repo *itemRepo) UpdateItem(ctx context.Context, item *models.Item) error {
	repo.logger.Debug("Enter in repository UpdateItem()")
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

	_, err = tx.Exec(ctx, `UPDATE items SET name=$1, category=$2, description=$3, price=$4, vendor=$5, pictures = $6 WHERE id=$7`,
		item.Title,
		item.Category.Id,
		item.Description,
		item.Price,
		item.Vendor,
		item.Images,
		item.Id)
	if err != nil {
		repo.logger.Errorf("error on update item %s: %s", item.Id, err)
		return fmt.Errorf("error on update item %s: %w", item.Id, err)
	}
	repo.logger.Infof("item %s successfully updated", item.Id)
	return nil
}

func (repo *itemRepo) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	repo.logger.Debug("Enter in repository GetItem()")
	item := models.Item{}
	pool := repo.storage.GetPool()
	row := pool.QueryRow(ctx,
		`SELECT items.id, items.name, category, categories.name, categories.description, categories.picture, items.description, price, vendor, pictures FROM items INNER JOIN categories ON category=categories.id and items.id = $1 WHERE items.deleted_at is null AND categories.deleted_at is null`, id)
	err := row.Scan(
		&item.Id,
		&item.Title,
		&item.Category.Id,
		&item.Category.Name,
		&item.Category.Description,
		&item.Category.Image,
		&item.Description,
		&item.Price,
		&item.Vendor,
		&item.Images,
	)
	if err != nil {
		repo.logger.Errorf("error in rows scan get item by id: %s", err)
		return &models.Item{}, fmt.Errorf("error in rows scan get item by id: %w", err)
	}
	return &item, nil
}

func (repo *itemRepo) ItemsList(ctx context.Context) (chan models.Item, error) {
	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		item := &models.Item{}

		pool := repo.storage.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT items.id, items.name, category, categories.name, categories.description, categories.picture, items.description, price, vendor, pictures FROM items INNER JOIN categories ON category=categories.id WHERE items.deleted_at is null AND categories.deleted_at is null`)
		if err != nil {
			msg := fmt.Errorf("error on items list query context: %w", err)
			repo.logger.Error(msg.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&item.Id,
				&item.Title,
				&item.Category.Id,
				&item.Category.Name,
				&item.Category.Description,
				&item.Category.Image,
				&item.Description,
				&item.Price,
				&item.Vendor,
				&item.Images,
			); err != nil {
				repo.logger.Error(err.Error())
				return
			}
			itemChan <- *item
		}
	}()

	return itemChan, nil
}

func (repo *itemRepo) SearchLine(ctx context.Context, param string) (chan models.Item, error) {
	repo.logger.Debug("Enter in repository SearchLine()")
	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		item := &models.Item{}
		pool := repo.storage.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT items.id, items.name, category, categories.name, categories.description,categories.picture, items.description, price, vendor, pictures FROM items INNER JOIN categories ON category=categories.id WHERE items.name LIKE $1 OR items.description LIKE $1 OR vendor LIKE $1 OR categories.name LIKE $1 WHERE items.deleted_at is null AND categories.deleted_at is null`,
			"%"+param+"%")
		if err != nil {
			msg := fmt.Errorf("error on search line query context: %w", err)
			repo.logger.Error(msg.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&item.Id,
				&item.Title,
				&item.Category.Id,
				&item.Category.Name,
				&item.Category.Description,
				&item.Category.Image,
				&item.Description,
				&item.Price,
				&item.Vendor,
				&item.Images,
			); err != nil {
				repo.logger.Error(err.Error())
				return
			}
			repo.logger.Info(fmt.Sprintf("find item: %v", item))
			itemChan <- *item
		}
	}()
	return itemChan, nil
}

func (repo *itemRepo) GetItemsByCategory(ctx context.Context, categoryName string) (chan models.Item, error) {
	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		item := &models.Item{}
		pool := repo.storage.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT items.id, items.name, category, categories.name, categories.description,categories.picture, items.description, price, vendor, pictures FROM items INNER JOIN categories ON category=categories.id WHERE WHERE items.deleted_at is null AND categories.deleted_at is null AND categories.name=$1`,
			categoryName)
		if err != nil {
			msg := fmt.Errorf("error on get items by category query context: %w", err)
			repo.logger.Error(msg.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&item.Id,
				&item.Title,
				&item.Category.Id,
				&item.Category.Name,
				&item.Category.Description,
				&item.Category.Image,
				&item.Description,
				&item.Price,
				&item.Vendor,
				&item.Images,
			); err != nil {
				repo.logger.Error(err.Error())
				return
			}
			repo.logger.Info(fmt.Sprintf("find item: %v", item))
			itemChan <- *item
		}
	}()
	return itemChan, nil
}

func (repo *itemRepo) DeleteItem(ctx context.Context, id uuid.UUID) error {
	repo.logger.Debug("Enter in repository DeleteItem()")
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
	_, err = tx.Exec(ctx, `UPDATE items SET deleted_at=$1 WHERE id=$2`,
		time.Now(), id)
	if err != nil {
		repo.logger.Errorf("error on delete item %s: %s", id, err)
		return fmt.Errorf("error on delete item %s: %w", id, err)
	}
	repo.logger.Infof("Item %s successfully deleted from database", id)
	return nil
}
