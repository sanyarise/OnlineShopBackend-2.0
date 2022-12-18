package cash

import (
	"OnlineShopBackend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var _ Cash = &RedisCash{}

type RedisCash struct {
	*redis.Client
	TTL    time.Duration
	logger *zap.Logger
}

type results struct {
	Responses []models.Item
}

// NewRedisCash initialize redis client
func NewRedisCash(host, port string, ttl time.Duration, logger *zap.Logger) (*RedisCash, error) {
	logger.Debug("Enter in NewRedisCash()")
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redis: %w", err)
	}
	logger.Debug("Redis Client ping success")
	cashTTL := ttl * time.Minute
	c := &RedisCash{
		Client: client,
		TTL:    cashTTL,
		logger: logger,
	}
	return c, nil
}

// Close close redis client
func (cash *RedisCash) Close() error {
	cash.logger.Debug("Enter in RedisCash Close()")
	return cash.Client.Close()
}

// CheckCash checks for data in the cache
func (cash *RedisCash) CheckCash(ctx context.Context, key string) bool {
	cash.logger.Debug("Enter in cash CheckCash()")
	check := cash.Exists(ctx, key)
	result, err := check.Result()
	if err != nil {
		cash.logger.Error(fmt.Errorf("error on check cash: %w", err).Error())
		return false
	}
	cash.logger.Debug(fmt.Sprintf("Check Cash with key: %s is %v", key, result))
	if result == 0 {
		cash.logger.Debug(fmt.Sprintf("Redis: get record %q not exist", key))
		return false
	} else {
		cash.logger.Debug(fmt.Sprintf("Key %q in cash found success", key))
		return true
	}
}

// CreateCash add data in the cash
func (cash *RedisCash) CreateItemsCash(ctx context.Context, res []models.Item, key string) error {
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

func (cash *RedisCash) CreateItemsQuantityCash(ctx context.Context, value int, key string) error {
	cash.logger.Debug("Enter in cash CreateItemsQuantityCash()")
	err := cash.Set(ctx, key, value, cash.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: set key %q: %w", key, err)
	}
	return nil
}

// GetCash retrieves data from the cache
func (cash *RedisCash) GetItemsCash(ctx context.Context, key string) ([]models.Item, error) {
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

func (cash *RedisCash) GetItemsQuantityCash(ctx context.Context, key string) (int, error) {
	cash.logger.Debug("Enter in cash GetItemsQuantityCash()")
	data, err := cash.Get(ctx, key).Int()
	if err != nil {
		return data, err
	}
	return data, nil
}

func (cash *RedisCash) ShutDown(timeout int) {
	cash.logger.Debug("Enter in cash ShutDown()")
	ctxWithTimiout, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	status := cash.Shutdown(ctxWithTimiout)
	_, err := status.Result()
	if err != nil {
		cash.logger.Warn(fmt.Sprintf("cash shutdown error: %v", err))
		return
	}
}
