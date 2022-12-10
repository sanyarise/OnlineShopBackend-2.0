package handlers

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ IItemHandlers = &ItemHandlers{}

type ItemHandlers struct {
	usecase usecase.IItemUsecase
	logger  *zap.Logger
}

func NewItemHandlers(usecase usecase.IItemUsecase, logger *zap.Logger) *ItemHandlers {
	return &ItemHandlers{usecase: usecase, logger: logger}
}

// Item is struct for DTO
type Item struct {
	Id          string   `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       int32    `json:"price,omitempty"`
	Category    Category `json:"category,omitempty"`
	Vendor      string   `json:"vendor,omitempty"`
	Images      []string `json:"image,omitempty"`
}

// CreateItem transform Item to models.Item and call usecase CreateItem
func (handlers *ItemHandlers) CreateItem(ctx context.Context, item Item) (uuid.UUID, error) {
	handlers.logger.Debug("Enter in handlers CreateItem()")
	categoryId, err := uuid.Parse(item.Category.Id)
	if err != nil {
		return uuid.Nil, err
	}
	itemCategory := models.Category{
		Id:          categoryId,
		Name:        item.Category.Name,
		Description: item.Category.Description,
	}
	newItem := &models.Item{
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    itemCategory,
		Vendor:      item.Vendor,
		Images:      item.Images,
	}
	id, err := handlers.usecase.CreateItem(ctx, newItem)
	if err != nil {
		return id, err
	}
	return id, nil
}

// UpdateItem transform Item to models.Item and call usecase UpdateItem
func (handlers *ItemHandlers) UpdateItem(ctx context.Context, item Item) error {
	handlers.logger.Debug("Enter in handlers UpdateItem()")
	id, err := uuid.Parse(item.Id)
	if err != nil {
		return fmt.Errorf("invalid item uuid: %w", err)
	}
	categoryId, err := uuid.Parse(item.Category.Id)
	if err != nil {
		return fmt.Errorf("invalid category uuid: %w", err)
	}
	itemCategory := models.Category{
		Id:          categoryId,
		Name:        item.Category.Name,
		Description: item.Category.Description,
	}

	updateItem := &models.Item{
		Id:          id,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    itemCategory,
		Vendor:      item.Vendor,
		Images:      item.Images,
	}
	return handlers.usecase.UpdateItem(ctx, updateItem)
}

// GetItem returns Item on id
func (handlers *ItemHandlers) GetItem(ctx context.Context, id string) (Item, error) {
	handlers.logger.Debug("Enter in handlers GetItem()")
	uid, err := uuid.Parse(id)
	if err != nil {
		return Item{}, fmt.Errorf("invalid item uuid: %w", err)
	}
	item, err := handlers.usecase.GetItem(ctx, uid)
	if err != nil {
		return Item{}, err
	}

	handlersCategory := Category{
		Id:          item.Category.Id.String(),
		Name:        item.Category.Name,
		Description: item.Category.Description,
	}

	return Item{
		Id:          item.Id.String(),
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    handlersCategory,
		Vendor:      item.Vendor,
		Images:      item.Images,
	}, nil
}

// ItemsList returns list of all Items
func (handlers *ItemHandlers) ItemsList(ctx context.Context, offset, limit int) ([]Item, error) {
	handlers.logger.Debug("Enter in handlers ItemsList()")
	res := make([]Item, 0, limit)
	items, err := handlers.usecase.ItemsList(ctx, offset, limit)
	if err != nil {
		return res, err
	}
	for _, item := range items {
		handlersCategory := Category{
			Id:          item.Category.Id.String(),
			Name:        item.Category.Name,
			Description: item.Category.Description,
		}

		res = append(res, Item{
			Id:          item.Id.String(),
			Title:       item.Title,
			Description: item.Description,
			Price:       item.Price,
			Category:    handlersCategory,
			Vendor:      item.Vendor,
		})
	}
	return res, nil
}

func (handlers *ItemHandlers) ItemsQuantity(ctx context.Context) (int, error) {
	handlers.logger.Debug("Enter in handlers ItemsQuantity()")
	quantity, err := handlers.usecase.ItemsQuantity(ctx)
	if err != nil {
		return quantity, err
	}
	return quantity, nil
}

// SearchLine returns list of Items with parameters
func (handlers *ItemHandlers) SearchLine(ctx context.Context, param string) ([]Item, error) {
	handlers.logger.Debug("Enter in handlers SearchLine()")
	res := make([]Item, 0, 100)
	items, err := handlers.usecase.SearchLine(ctx, param)
	if err != nil {
		return res, err
	}
	for {
		select {
		case <-ctx.Done():
			handlers.logger.Debug("handlers SearchLine() ctx.Done recieved")
			return res, ctx.Err()
		case item, ok := <-items:
			if !ok {
				return res, nil
			}
			handlersCategory := Category{
				Id:          item.Category.Id.String(),
				Name:        item.Category.Name,
				Description: item.Category.Description,
			}
			res = append(res, Item{
				Id:          item.Id.String(),
				Title:       item.Title,
				Description: item.Description,
				Price:       item.Price,
				Category:    handlersCategory,
				Vendor:      item.Vendor,
			})
		}
	}
}
