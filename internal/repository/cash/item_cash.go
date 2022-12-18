package cash

import (
	"OnlineShopBackend/internal/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var _ IItemsCash = &ItemsCash{}

type ItemsCash struct {
	*RedisCash
	logger *zap.Logger
}

type results struct {
	Responses []models.Item
}

func NewItemsCash(cash *RedisCash, logger *zap.Logger) *ItemsCash {
	return &ItemsCash{cash, logger}
}

// CreateCash add data in the cash
func (cash *ItemsCash) CreateItemsCash(ctx context.Context, res []models.Item, key string) error {
	cash.logger.Debug("Enter in cash CreateItemsListCash()")
	in := results{
		Responses: res,
	}
	data, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal unknown item: %w", err)
	}

	err = cash.Set(ctx, key, data, cash.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: set key %q: %w", key, err)
	}
	return nil
}

// CreateItemsQuantityCash create cash for items quantity
func (cash *ItemsCash) CreateItemsQuantityCash(ctx context.Context, value int, key string) error {
	cash.logger.Debug("Enter in cash CreateItemsQuantityCash()")
	err := cash.Set(ctx, key, value, cash.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: set key %q: %w", key, err)
	}
	return nil
}

// GetCash retrieves data from the cache
func (cash *ItemsCash) GetItemsCash(ctx context.Context, key string) ([]models.Item, error) {
	cash.logger.Debug("Enter in cash GetItemsListCash()")
	res := results{}
	data, err := cash.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res.Responses, nil
}

func (cash *ItemsCash) GetItemsQuantityCash(ctx context.Context, key string) (int, error) {
	cash.logger.Debug("Enter in cash GetItemsQuantityCash()")
	data, err := cash.Get(ctx, key).Int()
	if err != nil {
		return data, err
	}
	return data, nil
}
