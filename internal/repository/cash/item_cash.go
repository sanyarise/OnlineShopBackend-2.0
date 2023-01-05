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

func NewItemsCash(cash *RedisCash, logger *zap.Logger) IItemsCash {
	logger.Debug("Enter in cash NewItemsCash")
	return &ItemsCash{cash, logger}
}

// CheckCash checks for data in the cache
func (cash *ItemsCash) CheckCash(ctx context.Context, key string) bool {
	cash.logger.Sugar().Debugf("Enter in cash CheckCash() with args: ctx, key: %s", key)
	check := cash.Exists(ctx, key)
	result, err := check.Result()
	if err != nil {
		cash.logger.Error(fmt.Errorf("error on check cash: %w", err).Error())
		return false
	}
	cash.logger.Debug(fmt.Sprintf("Check Cash with key: %s is %v", key, result))
	if result == 0 {
		cash.logger.Debug(fmt.Sprintf("Redis: get record %s not exist", key))
		return false
	} else {
		cash.logger.Debug(fmt.Sprintf("Key %s in cash found success", key))
		return true
	}
}

// CreateCash add data in the cash
func (cash *ItemsCash) CreateItemsCash(ctx context.Context, res []models.Item, key string) error {
	cash.logger.Sugar().Debugf("Enter in cash CreateItemsCash() with args: ctx, res, key: %s", key)
	in := results{
		Responses: res,
	}
	data, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal unknown item: %w", err)
	}

	err = cash.Set(ctx, key, data, cash.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: error on set key %s: %w", key, err)
	}
	cash.logger.Info(fmt.Sprintf("Cash with key: %s create success", key))
	return nil
}

// CreateItemsQuantityCash create cash for items quantity
func (cash *ItemsCash) CreateItemsQuantityCash(ctx context.Context, value int, key string) error {
	cash.logger.Sugar().Debugf("Enter in cash CreateItemsQuantityCash() with args: ctx, value: %d, key: %s", value, key)
	err := cash.Set(ctx, key, value, cash.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: error on set key %q: %w", key, err)
	}
	cash.logger.Info(fmt.Sprintf("Cash with key: %s create success", key))
	return nil
}

// GetItemsCash retrieves data from the cache
func (cash *ItemsCash) GetItemsCash(ctx context.Context, key string) ([]models.Item, error) {
	cash.logger.Sugar().Debugf("Enter in cash GetItemsCash() with args: ctx, key: %s", key)
	res := results{}
	data, err := cash.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		cash.logger.Debug("Success get nil result")
		return nil, nil
	} else if err != nil {
		cash.logger.Sugar().Errorf("Error on get cash: %v", err)
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		cash.logger.Sugar().Warnf("Can't json unmarshal data: %v", data)
		return nil, err
	}
	cash.logger.Debug("Get cash success")
	return res.Responses, nil
}

// GetItemsQuantityCash retrieves data from the cache
func (cash *ItemsCash) GetItemsQuantityCash(ctx context.Context, key string) (int, error) {
	cash.logger.Sugar().Debugf("Enter in cash GetItemsQuantityCash() with args: ctx, key: %s", key)
	data, err := cash.Get(ctx, key).Int()
	if err != nil {
		cash.logger.Sugar().Errorf("Error on get cash: %v", err)
		return data, err
	}
	cash.logger.Debug("Get cash success")
	return data, nil
}
