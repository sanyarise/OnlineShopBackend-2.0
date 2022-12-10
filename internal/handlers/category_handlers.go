package handlers

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ ICategoryHandlers = &CategoryHandlers{}

type CategoryHandlers struct {
	usecase usecase.ICategoryUsecase
	logger  *zap.Logger
}

func NewCategoryHandlers(usecase usecase.ICategoryUsecase, logger *zap.Logger) *CategoryHandlers {
	return &CategoryHandlers{usecase: usecase, logger: logger}
}

// Category is struct for DTO
type Category struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateCategory transform Category to models.Category and call usecase CreateCategory
func (handlers *CategoryHandlers) CreateCategory(ctx context.Context, category Category) (uuid.UUID, error) {
	handlers.logger.Debug("Enter in handlers CreateCategory()")
	newCategory := &models.Category{
		Name:        category.Name,
		Description: category.Description,
	}
	id, err := handlers.usecase.CreateCategory(ctx, newCategory)
	if err != nil {
		return id, err
	}
	return id, nil
}

// GetCategoryList returns list of all categories
func (handlers *CategoryHandlers) GetCategoryList(ctx context.Context) ([]Category, error) {
	handlers.logger.Debug("Enter in handlers GetCategoryList()")
	res := make([]Category, 0, 100)
	categories, err := handlers.usecase.GetCategoryList(ctx)
	if err != nil {
		return res, err
	}
	for {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case category, ok := <-categories:
			if !ok {
				return res, nil
			}
			res = append(res, Category{
				Id:          category.Id.String(),
				Name:        category.Name,
				Description: category.Description,
			})
		}
	}
}
