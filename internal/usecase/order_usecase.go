package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type order struct {
	orderStore repository.OrderStore
	logger     *zap.SugaredLogger
}

var _ IOrderUsecase = (*order)(nil)

func NewOrderUsecase(orderStore repository.OrderStore, logger *zap.SugaredLogger) IOrderUsecase {
	return &order{
		orderStore: orderStore,
		logger:     logger,
	}
}

func (o *order) PlaceOrder(ctx context.Context, cart *models.Cart, user *models.User) (*models.Order, error) {
	select {
	case <-ctx.Done():
		o.logger.Error("context closed")
		return nil, fmt.Errorf("context closed")
	default:
		ordr := models.Order{
			User:         *user,
			Address:      user.Address,
			Status:       models.StatusCreated,
			ShipmentTime: time.Now().Add(models.ProlongedShipmentPeriod),
			Items:        append([]models.Item{}[:0:0], cart.Items...),
		}
		res, err := o.orderStore.Create(ctx, &ordr)
		if err != nil {
			o.logger.Errorf("can't add order to db %s", err)
			return nil, fmt.Errorf("can't place order to db : %w", err)
		}
		return res, nil
	}
}

func (o *order) ChangeStatus(ctx context.Context, order *models.Order, newStatus models.Status) error {
	return nil
}
func (o *order) GetOrdersForUser(ctx context.Context, user *models.User) ([]models.Order, error) {
	return []models.Order{}, nil
}
func (o *order) DeleteOrder(ctx context.Context, order *models.Order) error {
	return nil
}
func (o *order) ChangeAddress(ctx context.Context, order *models.Order, newAddress models.UserAddress) error {
	return nil
}
func (o *order) GetOrder(ctx context.Context, id uuid.UUID) (models.Order, error) {
	return models.Order{}, nil
}
