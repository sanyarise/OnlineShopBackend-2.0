package handlers

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type Item struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int32  `json:"price,omitempty"`
	Category    string `json:"category,omitempty"`
	Image       string `json:"image,omitempty"`
}

func (h *Handlers) CreateItem(ctx context.Context, item Item) (uuid.UUID, error) {
	newItem := &models.Item{
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    item.Category,
		Image:       item.Image,
	}
	id, err := h.repo.CreateItem(ctx, newItem)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (h *Handlers) UpdateItem(ctx context.Context, item Item) error {
	updateItem := &models.Item{
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    item.Category,
		Image:       item.Image,
	}
	return h.repo.UpdateItem(ctx, updateItem)
}
