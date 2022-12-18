package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type cart struct {
	storage *PGres
	logger  *zap.SugaredLogger
}

var _ CartStore = (*cart)(nil)

func NewCartStore(storage *PGres, logger *zap.SugaredLogger) CartStore {
	return &cart{
		storage: storage,
		logger:  logger,
	}
}

// Create Shall we add items at the moment we create cart
func (c *cart) Create(ctx context.Context, cart *models.Cart) (*models.Cart, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		row := pool.QueryRow(ctx, `INSERT INTO carts (user_id, expire_at) VALUES ($1, $2) RETURNING id`,
			cart.UserID, cart.ExpireAt)
		err := row.Scan(&cart.ID)
		if err != nil {
			c.logger.Error(err)
			return nil, fmt.Errorf("can't create cart object: %w", err)
		}
		return cart, nil
	}
}

// AddItemToCart Maybe add to item
func (c *cart) AddItemToCart(ctx context.Context, cart *models.Cart, item *models.Item) error {
	select {
	case <-ctx.Done():
		c.logger.Error("context closed")
		return fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		_, err := pool.Exec(ctx, `INSERT INTO cart_items (cart_id, item_id) VALUES ($1, $2)`, cart.ID, item.Id)
		if err != nil {
			c.logger.Errorf("can't add item to cart: %s", err)
			return fmt.Errorf("can't add item to cart: %w", err)
		}
		return nil
	}
}

func (c *cart) DeleteCart(ctx context.Context, cart *models.Cart) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		defer func() {
			if err != nil {
				c.logger.Errorf("transaction rolled back")
				if err = tx.Rollback(ctx); err != nil {
					c.logger.Errorf("can't rollback %s", err)
				}

			} else {
				c.logger.Info("transaction commited")
				if err != tx.Commit(ctx) {
					c.logger.Errorf("can't commit %s", err)
				}
			}
		}()
		_, err = tx.Exec(ctx, `DELETE FROM cart_items WHERE cart_id=$1`, cart.ID)
		if err != nil {
			c.logger.Errorf("can't delete cart items from cart: %s", err)
			return fmt.Errorf("can't delete cart items from cart: %w", err)
		}
		_, err = tx.Exec(ctx, `DELETE FROM carts WHERE id=$1`, cart.ID)
		if err != nil {
			c.logger.Errorf("can't delete cart: %s", err)
			return fmt.Errorf("can't delete cart: %w", err)
		}
		return nil
	}
}
func (c *cart) DeleteItemFromCart(ctx context.Context, cart *models.Cart, item *models.Item) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		_, err := pool.Exec(ctx, `DELETE FROM cart_items WHERE item_id=$1 AND cart_id=$2`, item.Id, cart.ID)
		if err != nil {
			c.logger.Errorf("can't delete item from cart: %s", err)
			return fmt.Errorf("can't delete item from cart: %w", err)
		}
		return nil
	}
}

func (c *cart) SelectItemsFromCart(ctx context.Context, cart *models.Cart) ([]*models.Item, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		rows, err := pool.Query(ctx, `
				SELECT 	i.id, i.name, i.category, i.description, i.price, i.vendor, i.pictures
				FROM cart_items c, items i
				WHERE c.cart_id=$1 and i.id = c.item_id`, cart.ID)
		if err != nil {
			c.logger.Errorf("can't select items from cart: %s", err)
			return nil, fmt.Errorf("can't select items from cart: %w", err)
		}
		items := make([]*models.Item, 0)
		for rows.Next() {
			v, err := rows.Values()
			if err != nil {
				c.logger.Errorf("can't select items from cart: %s", err)
				return nil, fmt.Errorf("can't select items from cart: %w", err)
			}
			i := models.Item{
				Id:          v[0].(uuid.UUID),
				Title:       v[1].(string),
				Description: v[3].(string),
				Price:       v[4].(int32),
				Category: models.Category{
					Id:          v[2].(uuid.UUID),
					Name:        "",
					Description: "",
				},
				Vendor: v[5].(string),
				Images: v[6].([]string),
			}
			items = append(items, &i)
		}
		return items, nil
	}
}
