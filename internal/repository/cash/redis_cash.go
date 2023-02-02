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
	logger.Sugar().Debugf("Enter in NewRedisCash() with args: host: %s, port: %s, ttl: %v, logger", host, port, ttl)
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

// ShutDown is func for graceful shutdown redis connection
func (cash *RedisCash) ShutDown(timeout int) error{
	cash.logger.Sugar().Debugf("Enter in cash ShutDown() with args: timeout: %d", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	status := cash.Shutdown(ctx)
	result, err := status.Result()
	if err != nil {
		return err
	}
	cash.logger.Info(result)
	return nil
}