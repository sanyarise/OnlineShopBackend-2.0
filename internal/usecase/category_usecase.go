package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateCategory call database method and returns id of created category or error
func (storage *Storage) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	storage.logger.Debug("Enter in usecase CreateCategory()")
	id, err := storage.categoryStore.CreateCategory(ctx, category)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create category: %w", err)
	}
	return id, nil
}

// GetCategoryList call database method and returns chan with all models.Category or error
func (storage *Storage) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	storage.logger.Debug("Enter in usecase GetCategoryList()")
	categoryIncomingChan, err := storage.categoryStore.GetCategoryList(ctx)
	if err != nil {
		return nil, err
	}
	categoryOutChan := make(chan models.Category, 100)
	go func() {
		defer close(categoryOutChan)
		for {
			select {
			case <-ctx.Done():
				return
			case category, ok := <-categoryIncomingChan:
				if !ok {
					return
				}
				categoryOutChan <- category
			}
		}
	}()
	return categoryOutChan, nil

}
