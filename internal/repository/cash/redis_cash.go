package cash

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisCash struct {
	*redis.Client
	TTL    time.Duration
	logger *zap.Logger
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

// ShutDown is func for graceful shutdown redis connection
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
