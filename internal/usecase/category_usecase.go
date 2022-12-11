package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ ICategoryUsecase = &CategoryUsecase{}

type CategoryUsecase struct {
	categoryStore repository.CategoryStore
	logger        *zap.Logger
}

func NewCategoryUsecase(store repository.CategoryStore, logger *zap.Logger) ICategoryUsecase {
	return &CategoryUsecase{categoryStore: store, logger: logger}
}

// / CreateCategory call database method and returns id of created category or error
func (usecase *CategoryUsecase) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	usecase.logger.Debug("Enter in usecase CreateCategory()")
	id, err := usecase.categoryStore.CreateCategory(ctx, category)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create category: %w", err)
	}
	return id, nil
}

// UpdateCategory call database method to update category and returns error or nil
func (usecase *CategoryUsecase) UpdateCategory(ctx context.Context, category *models.Category) error {
	usecase.logger.Debug("Enter in usecase UpdateCategory()")
	err := usecase.categoryStore.UpdateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("error on update category: %w", err)
	}
	return nil
}

// GetCategory call database and returns *models.Category with given id or returns error
func (usecase *CategoryUsecase) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	usecase.logger.Debug("Enter in usecase GetCategory()")
	category, err := usecase.categoryStore.GetCategory(ctx, id)
	if err != nil {
		return &models.Category{}, fmt.Errorf("error on get category: %w", err)
	}
	return category, nil
}

// GetCategoryList call database method and returns chan with all models.Category or error
func (usecase *CategoryUsecase) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	usecase.logger.Debug("Enter in usecase GetCategoryList()")
	categoryIncomingChan, err := usecase.categoryStore.GetCategoryList(ctx)
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
