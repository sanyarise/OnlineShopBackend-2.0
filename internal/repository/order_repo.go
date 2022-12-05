package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type order struct {
	storage Storage
	logger  *zap.SugaredLogger
}

func (o *order) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("stopped with context")
	default:
		pool := o.storage.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			o.logger.Errorf("can't create transaction: %s", err)
			return nil, fmt.Errorf("can't create transaction: %w", err)
		}
		defer func() {
			if err != nil {
				o.logger.Info("rollback transaction")
				tx.Rollback(ctx)
			} else {
				o.logger.Info("commit transaction")
				tx.Commit(ctx)
			}
		}()
		row := tx.QueryRow(ctx, `INSERT INTO orders (shipment_time, user_id, status, address) 
		VALUES ($1, $2, $3, $4) RETURNING id`, order.ShipmentTime, order.User.ID, order.Status, order.Address)
		err = row.Scan(&order.ID)
		if err != nil {
			o.logger.Errorf("can't add new order: %w", err)
			return nil, fmt.Errorf("can't add new order: %w", err)
		}
		query := `INSERT INTO order_items (order_id, item_id) VALUES`
		itemsString := ""
		items := make([]interface{}, 0, len(order.Items))
		for ind, item := range order.Items {
			itemsString += fmt.Sprintf("($%d $$d),", 1, ind+2)
			items = append(items, item.Id.String())
		}
		itemsString = itemsString[:len(itemsString)-1]
		_, err = tx.Exec(ctx, fmt.Sprintf("%s %s", query, itemsString), items...)
		if err != nil {
			o.logger.Errorf("can't add items to order: %s", err)
			return nil, fmt.Errorf("can't add items to order: %w", err)
		}
		return &models.Order{}, nil
	}
}

func (o *order) DeleteOrder(ctx context.Context, order *models.Order) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := o.storage.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		defer func() {
			if err != nil {
				o.logger.Info("rollback transaction")
				tx.Rollback(ctx)
			} else {
				o.logger.Info("commit transaction")
				tx.Commit(ctx)
			}
		}()
		_, err = tx.Exec(ctx, `DELETE FROM order_items WHERE cart_id=$1`, order.ID)
		if err != nil {
			o.logger.Errorf("can't delete order items from order: %s", err)
			return fmt.Errorf("can't delete order items from order: %w", err)
		}
		_, err = tx.Exec(ctx, `DELETE FROM orders WHERE id=$1`, order.ID)
		if err != nil {
			o.logger.Errorf("can't delete order: %s", err)
			return fmt.Errorf("can't delete order: %w", err)
		}
		return nil
	}
}
func (o *order) ChangeAddress(ctx context.Context, order *models.Order, address models.UserAddress) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := o.storage.GetPool()
		_, err := pool.Exec(ctx, `UPDATE orders SET address=$1 WHERE id=%2`, address, order.ID)
		if err != nil {
			o.logger.Errorf("can't update address: %s", err)
			return fmt.Errorf("can't update address: %w", err)
		}
		return nil
	}
}
func (o *order) ChangeStatus(ctx context.Context, order *models.Order, status models.Status) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := o.storage.GetPool()
		_, err := pool.Exec(ctx, `UPDATE orders SET status=$1 WHERE id=%2`, status, order.ID)
		if err != nil {
			o.logger.Errorf("can't update status: %s", err)
			return fmt.Errorf("can't update status: %w", err)
		}
		return nil
	}
}
func (o *order) GetOrderByID(ctx context.Context, id uuid.UUID) (models.Order, error) {
	select {
	case <-ctx.Done():
		o.logger.Errorf("context closed")
		return models.Order{}, fmt.Errorf("context closed")
	default:
		pool := o.storage.GetPool()
		ordr := models.Order{
			Items: make([]models.Item, 0),
		}
		rows, err := pool.Query(ctx, `SELECT items.id, items.name, categories.id, categories.name, categories.description,
				items.description, items.price, items.vendor, orders.id, orders.shipment_time,
				orders.status, orders.address from items INNER JOIN categories ON categories.id=category  INNER JOIN order_items ON
				items.id=order_items.item_id INNER JOIN orders ON orders.id=order_items.order_id and orders.id = $1 ORDER BY order_id ASC`, id)
		if err != nil {
			o.logger.Errorf("can't get order from db: %s", err)
			return ordr, fmt.Errorf("can't get order from db: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			item := models.Item{}
			if err := rows.Scan(&item.Id, &item.Title, &item.Category.Id, &item.Category.Name, &item.Category.Description,
				&item.Description, &item.Price, &item.Vendor, &ordr.ID, &ordr.ShipmentTime, &ordr.Status, &ordr.Address); err != nil {
				o.logger.Errorf("can't scan data to order object: %w", err)
				return models.Order{}, err
			}
			ordr.Items = append(ordr.Items, item)
		}
		return ordr, nil
	}

}

func (o *order) GetOrdersForUser(ctx context.Context, user *models.User) (chan models.Order, error) {
	select {
	case <-ctx.Done():
		o.logger.Errorf("context closed")
		return nil, fmt.Errorf("context closed")
	default:
		pool := o.storage.GetPool()
		resChan := make(chan models.Order, 1)
		go func() {
			defer close(resChan)
			rows, err := pool.Query(ctx, `SELECT items.id, items.name, categories.id, categories.name, categories.description,
			items.description, items.price, items.vendor, orders.id, orders.shipment_time,
			orders.status, orders.address from items INNER JOIN categories ON categories.id=category  INNER JOIN order_items ON
			items.id=order_items.item_id INNER JOIN orders ON orders.id=order_items.order_id and orders.user_id = $1 ORDER BY order_id ASC`, user.ID)
			if err != nil {
				o.logger.Errorf("can't get order from db: %s", err)
				return
			}
			defer rows.Close()
			prevOrder := models.Order{
				Items: make([]models.Item, 0),
			}
			for rows.Next() {
				item := models.Item{}
				order := models.Order{}
				if err := rows.Scan(&item.Id, &item.Title, &item.Category.Id, &item.Category.Name, &item.Category.Description,
					&item.Description, &item.Price, &item.Vendor, &order.ID, &order.ShipmentTime, &order.Status, &order.Address); err != nil {
					o.logger.Errorf("can't scan data to order object: %w", err)
					return
				}
				if order.ID != prevOrder.ID {
					prevOrder = order
					resChan <- prevOrder
				} else {
					prevOrder.Items = append(prevOrder.Items, item)
				}
			}
			resChan <- prevOrder
		}()
		return resChan, nil
	}
}
