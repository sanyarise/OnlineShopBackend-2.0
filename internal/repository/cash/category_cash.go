package cash

import (
	"OnlineShopBackend/internal/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var _ ICategoriesCash = &CategoriesCash{}

type CategoriesCash struct {
	*RedisCash
	logger *zap.Logger
}

type categoriesData struct {
	Categories []models.Category `json:"categories"`
}

func NewCategoriesCash(cash *RedisCash, logger *zap.Logger) ICategoriesCash {
	return &CategoriesCash{cash, logger}
}

// CheckCash checks for data in the cache
func (cash *CategoriesCash) CheckCash(ctx context.Context, key string) bool {
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

func (cash *CategoriesCash) CreateCategoriesListCash(ctx context.Context, categories []models.Category, key string) error {
	cash.logger.Debug("Enter in CategoriesCash CreateCash()")
	in := categoriesData{
		Categories: categories,
	}
	bytesData, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal unknown category: %w", err)
	}

	err = cash.Set(ctx, key, bytesData, cash.TTL).Err()
	if err != nil {
		return fmt.Errorf("error on set cash with key: %v, error: %w", key, err)
	}
	cash.logger.Debug(fmt.Sprintf("Cash with key %s write in redis success", key))
	return nil
}

// GetCategoriesListCash retrieves data from the cash
func (cash *CategoriesCash) GetCategoriesListCash(ctx context.Context, key string) ([]models.Category, error) {
	cash.logger.Debug("Enter in cash GetCategoriesListCash()")
	categories := categoriesData{}
	bytesData, err := cash.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytesData, &categories)
	if err != nil {
		return nil, err
	}
	return categories.Categories, nil
}

// DeleteCash deleted cash by key
func (cash *CategoriesCash) DeleteCash(ctx context.Context, key string) error {
	cash.logger.Debug(fmt.Sprintf("Enter in cash DeleteCash with args: ctx, key: %s", key))
	return cash.Del(ctx, key).Err()
}
