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
	logger.Debug("Enter in cash NewCategoriesCash()")
	return &CategoriesCash{cash, logger}
}

// CheckCash checks for data in the cache
func (cash *CategoriesCash) CheckCash(ctx context.Context, key string) bool {
	cash.logger.Sugar().Debugf("Enter in cash CheckCash() with args: ctx, key: %s", key)
	check := cash.Exists(ctx, key)
	result, err := check.Result()
	if err != nil {
		cash.logger.Error(fmt.Errorf("error on check cash: %w", err).Error())
		return false
	}
	cash.logger.Debug(fmt.Sprintf("Check cash with key: %s is %v", key, result))
	if result == 0 {
		cash.logger.Debug(fmt.Sprintf("Redis: key %s not exist", key))
		return false
	} else {
		cash.logger.Debug(fmt.Sprintf("Key %s in cash found success", key))
		return true
	}
}

// CreateCategoriesListCash creates cash of categories list
func (cash *CategoriesCash) CreateCategoriesListCash(ctx context.Context, categories []models.Category, key string) error {
	cash.logger.Sugar().Debugf("Enter in CategoriesCash CreateCash() with args: ctx, categories []models.Category, key: %s", key)
	in := categoriesData{
		Categories: categories,
	}
	bytesData, err := json.Marshal(in)
	if err != nil {
		cash.logger.Sugar().Warnf("Error on json marshal data: %v", in)
		return fmt.Errorf("marshal unknown category: %w", err)
	}

	err = cash.Set(ctx, key, bytesData, cash.TTL).Err()
	if err != nil {
		cash.logger.Sugar().Warnf("Error on set cash with key: %s, error: %v", key, err)
		return fmt.Errorf("error on set cash with key: %v, error: %w", key, err)
	}
	cash.logger.Debug(fmt.Sprintf("Cash with key %s write in redis success", key))
	return nil
}

// GetCategoriesListCash retrieves data from the cash
func (cash *CategoriesCash) GetCategoriesListCash(ctx context.Context, key string) ([]models.Category, error) {
	cash.logger.Sugar().Debugf("Enter in cash GetCategoriesListCash() with args: ctx, key: %s", key)
	categories := categoriesData{}
	bytesData, err := cash.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		cash.logger.Debug("Success get nil result")
		return nil, nil
	} else if err != nil {
		cash.logger.Sugar().Errorf("Error on get cash: %v", err)
		return nil, err
	}
	err = json.Unmarshal(bytesData, &categories)
	if err != nil {
		cash.logger.Sugar().Warnf("Can't json unmarshal data: %v", bytesData)
		return nil, err
	}
	cash.logger.Debug("Get cash success")
	return categories.Categories, nil
}

// DeleteCash deleted cash by key
func (cash *CategoriesCash) DeleteCash(ctx context.Context, key string) error {
	cash.logger.Debug(fmt.Sprintf("Enter in cash DeleteCash with args: ctx, key: %s", key))
	err := cash.Del(ctx, key).Err()
	if err != nil {
		cash.logger.Sugar().Warnf("Error on delete cash with key: %s", key)
		return err
	}
	cash.logger.Sugar().Infof("Delete cash with key: %s success", key)
	return nil
}
