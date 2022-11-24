package handlers

import (
	"OnlineShopBackend/internal/models"
	"context"
	"log"

	"github.com/google/uuid"
)

type Category struct {
	Id string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`
}

// CreateCategory transform Category to models.Category and call usecase CreateCategory
func (h *Handlers) CreateCategory(ctx context.Context, category Category) (uuid.UUID, error) {
	log.Println("Enter in handlers CreateCategory()")
	newCategory := &models.Category{
		Name:        category.Name,
		Description: category.Description,
	}
	id, err := h.repo.CreateCategory(ctx, newCategory)
	if err != nil {
		return id, err
	}
	return id, nil
}

// GetCategoryList returns list of all categories
func (h *Handlers) GetCategoryList(ctx context.Context) ([]Category, error) {
	res := make([]Category, 0, 100)
	items, err := h.repo.GetCategoryList(ctx)
	if err != nil {
		return res, err
	}
	for {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case category, ok := <-items:
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
