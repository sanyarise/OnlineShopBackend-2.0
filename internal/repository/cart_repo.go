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
func (c *cart) Create(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	c.logger.Debugf("Enter in repository cart Create() with args: ctx, userId: %v", userId)
	select {
	case <-ctx.Done():
		return uuid.Nil, fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		var cartId uuid.UUID
		row := pool.QueryRow(ctx, `INSERT INTO carts (user_id) VALUES ($1) RETURNING id`,
			userId)
		err := row.Scan(&cartId)
		if err != nil {
			c.logger.Error(err)
			return uuid.Nil, fmt.Errorf("can't create cart object: %w", err)
		}
		c.logger.Info("Create cart success")
		return cartId, nil
	}
}

// AddItemToCart Maybe add to item
func (c *cart) AddItemToCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error {
	c.logger.Debugf("Enter in repository cart AddItemToCart() with args: ctx, cartId: %v, itemId: %v", cartId, itemId)
	select {
	case <-ctx.Done():
		c.logger.Error("context closed")
		return fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		_, err := pool.Exec(ctx, `INSERT INTO cart_items (cart_id, item_id) VALUES ($1, $2)`, cartId, itemId)
		if err != nil {
			c.logger.Errorf("can't add item to cart: %s", err)
			return fmt.Errorf("can't add item to cart: %w", err)
		}
		return nil
	}
}

func (c *cart) DeleteCart(ctx context.Context, cartId uuid.UUID) error {
	c.logger.Debug("Enter in repository cart DeleteCart() with args: ctx, cartId: %v", cartId)
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
		_, err = tx.Exec(ctx, `DELETE FROM cart_items WHERE cart_id=$1`, cartId)
		if err != nil {
			c.logger.Errorf("can't delete cart items from cart: %s", err)
			return fmt.Errorf("can't delete cart items from cart: %w", err)
		}
		_, err = tx.Exec(ctx, `DELETE FROM carts WHERE id=$1`, cartId)
		if err != nil {
			c.logger.Errorf("can't delete cart: %s", err)
			return fmt.Errorf("can't delete cart: %w", err)
		}
		c.logger.Info("Delete cart with id: %v from database success", cartId)
		return nil
	}
}
func (c *cart) DeleteItemFromCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error {
	c.logger.Debug("Enter in repository cart DeleteItemFromCart() with args: ctx, cartId: %v, itemId: %v", cartId, itemId)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		_, err := pool.Exec(ctx, `DELETE FROM cart_items WHERE item_id=$1 AND cart_id=$2`, itemId, cartId)
		if err != nil {
			c.logger.Errorf("can't delete item from cart: %s", err)
			return fmt.Errorf("can't delete item from cart: %w", err)
		}
		c.logger.Info("Delete item from cart success")
		return nil
	}
}

func (c *cart) GetCart(ctx context.Context, cartId uuid.UUID) (*models.Cart, error) {
	c.logger.Debug("Enter in repository cart SelectItemsFromCart() with args: ctx, cartId: %v", cartId)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		var userId uuid.UUID
		row := pool.QueryRow(ctx, `SELECT user_id FROM carts WHERE id = $1`, cartId)
		err := row.Scan(&userId)
		if err != nil {
			c.logger.Error(err)
			return nil, fmt.Errorf("can't read user id: %w", err)
		}
		c.logger.Debug("read user id success: %v", userId)
		item := models.Item{}
		rows, err := pool.Query(ctx, `
		SELECT 	i.id, i.name, i.price, i.pictures
		FROM cart_items c, items i
		WHERE c.cart_id=$1 and i.id = c.item_id`, cartId)
		if err != nil {
			c.logger.Errorf("can't select items from cart: %s", err)
			return nil, fmt.Errorf("can't select items from cart: %w", err)
		}
		defer rows.Close()
		c.logger.Debug("read info from db in pool.Query success")
		items := make([]models.Item, 0, 100)
		for rows.Next() {
			if err := rows.Scan(
				&item.Id,
				&item.Title,
				&item.Price,
				&item.Images,
			); err != nil {
				c.logger.Error(err.Error())
				return nil, err
			}

			items = append(items, item)
		}
		c.logger.Info("Select items from cart success")
		c.logger.Info("Get cart success")
		return &models.Cart{
			Id:     cartId,
			UserId: userId,
			Items:  items,
		}, nil
	}
}
