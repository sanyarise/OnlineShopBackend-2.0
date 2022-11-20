package handlers

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Item struct {
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int32  `json:"price,omitempty"`
	Category    string `json:"category,omitempty"`
	Vendor      string `json:"vendor,omitempty"`
	Image       string `json:"image,omitempty"`
}

func (h *Handlers) CreateItem(ctx context.Context, item Item) (uuid.UUID, error) {
	log.Println("Enter in handlers CreateItem()")
	log.Println(item)
	cid, err := uuid.Parse(item.Category)
	if err != nil {
		return uuid.UUID{}, err
	}
	newItem := &models.Item{
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    cid,
		Vendor:      item.Vendor,
		Image:       item.Image,
	}
	id, err := h.repo.CreateItem(ctx, newItem)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (h *Handlers) UpdateItem(ctx context.Context, item Item) error {
	id, _ := uuid.Parse(item.Id)
	cid, _ := uuid.Parse(item.Category)
	updateItem := &models.Item{
		Id:          id,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    cid,
		Vendor:      item.Vendor,
		Image:       item.Image,
	}
	return h.repo.UpdateItem(ctx, updateItem)
}

func (h *Handlers) GetItem(ctx context.Context, id string) (Item, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return Item{}, fmt.Errorf("invalid uuid: %w", err)
	}
	item, err := h.repo.GetItem(ctx, uid)
	if err != nil {
		return Item{}, err
	}
	return Item{
		Id:          item.Id.String(),
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    item.Category.String(),
		Vendor:      item.Vendor,
		Image:       item.Image,
	}, nil
}

func (h *Handlers) ItemsList(ctx context.Context) ([]Item, error) {
	res := make([]Item, 0, 100)
	items, err := h.repo.ItemsList(ctx)
	if err != nil {
		return res, err
	}
	for {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case item, ok := <-items:
			if !ok {
				return res, nil
			}
			res = append(res, Item{
				Id:          item.Id.String(),
				Title:       item.Title,
				Description: item.Description,
				Price:       item.Price,
				Category:    item.Category.String(),
				Vendor:      item.Vendor,
			})
		}
		return res, nil
	}
}

func (h *Handlers) SearchLine(ctx context.Context, param string) ([]Item, error) {
	res := make([]Item, 0, 100)
	items, err := h.repo.SearchLine(ctx, param)
	if err != nil {
		return res, err
	}
	for {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case item, ok := <-items:
			if !ok {
				return res, nil
			}
			res = append(res, Item{
				Id:          item.Id.String(),
				Title:       item.Title,
				Description: item.Description,
				Price:       item.Price,
				Category:    item.Category.String(),
				Vendor:      item.Vendor,
			})
		}
		return res, nil
	}
}
