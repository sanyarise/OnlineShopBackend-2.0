package handlers

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)
// Category is struct for DTO
type Category struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateCategory transform Category to models.Category and call usecase CreateCategory
func (handlers *Handlers) CreateCategory(ctx context.Context, category Category) (uuid.UUID, error) {
	handlers.logger.Debug("Enter in handlers CreateCategory()")
	newCategory := &models.Category{
		Name:        category.Name,
		Description: category.Description,
	}
	id, err := handlers.repo.CreateCategory(ctx, newCategory)
	if err != nil {
		return id, err
	}
	return id, nil
}

// GetCategoryList returns list of all categories
func (handlers *Handlers) GetCategoryList(ctx context.Context) ([]Category, error) {
	handlers.logger.Debug("Enter in handlers GetCategoryList()")
	res := make([]Category, 0, 100)
	categories, err := handlers.repo.GetCategoryList(ctx)
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
