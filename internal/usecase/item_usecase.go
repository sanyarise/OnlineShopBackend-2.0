package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"OnlineShopBackend/internal/repository/cash"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ IItemUsecase = &ItemUsecase{}

const (
	itemsListKey     = "ItemsList"
	itemsQuantityKey = "ItemsQuantity"
)

type ItemUsecase struct {
	itemStore repository.ItemStore
	itemCash  cash.Cash
	logger    *zap.Logger
}

func NewItemUsecase(itemStore repository.ItemStore, itemCash cash.Cash, logger *zap.Logger) IItemUsecase {
	return &ItemUsecase{itemStore: itemStore, itemCash: itemCash, logger: logger}
}

// CreateItem call database method and returns id of created item or error
func (usecase *ItemUsecase) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	usecase.logger.Debug("Enter in usecase CreateItem()")
	id, err := usecase.itemStore.CreateItem(ctx, item)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create item: %w", err)
	}
	err = usecase.UpdateCash(ctx, id, "create")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cash: %v", err))
	} else {
		usecase.logger.Info("Update cash success")
	}
	return id, nil
}

// UpdateItem call database method to update item and returns error or nil
func (usecase *ItemUsecase) UpdateItem(ctx context.Context, item *models.Item) error {
	usecase.logger.Debug("Enter in usecase UpdateItem()")
	err := usecase.itemStore.UpdateItem(ctx, item)
	if err != nil {
		return fmt.Errorf("error on update item: %w", err)
	}
	err = usecase.UpdateCash(ctx, item.Id, "update")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cash: %v", err))
	} else {
		usecase.logger.Info("Update cash success")
	}
	return nil
}

// GetItem call database and returns *models.Item with given id or returns error
func (usecase *ItemUsecase) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	usecase.logger.Debug("Enter in usecase GetItem()")
	item, err := usecase.itemStore.GetItem(ctx, id)
	if err != nil {
		return &models.Item{}, fmt.Errorf("error on get item: %w", err)
	}
	return item, nil
}

// ItemsList call database method and returns slice with all models.Item or error
func (usecase *ItemUsecase) ItemsList(ctx context.Context, offset, limit int) ([]models.Item, error) {
	usecase.logger.Debug("Enter in usecase ItemsList()")
	if ok := usecase.itemCash.CheckCash(ctx, itemsListKey); !ok {
		itemIncomingChan, err := usecase.itemStore.ItemsList(ctx)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			usecase.logger.Debug(fmt.Sprintf("item from channel is: %v", item))
			items = append(items, item)
		}
		err = usecase.itemCash.CreateItemsCash(ctx, items, itemsListKey)
		if err != nil {
			return nil, fmt.Errorf("error on create items list cash: %w", err)
		}
		err = usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey)
		if err != nil {
			return nil, fmt.Errorf("error on create items quantity cash: %w", err)
		}
	}
	items, err := usecase.itemCash.GetItemsCash(ctx, itemsListKey)
	if err != nil {
		return nil, fmt.Errorf("error on get cash: %w", err)
	}
	if offset > len(items) {
		return nil, fmt.Errorf("error: offset bigger than lenght of items, offset: %d, lenght of items: %d", offset, len(items))
	}
	itemsWithLimit := make([]models.Item, 0, limit)
	var counter = 0
	for i := offset; i < len(items); i++ {
		if counter == limit {
			break
		}
		itemsWithLimit = append(itemsWithLimit, items[i])
		counter++
	}
	return itemsWithLimit, nil
}

func (usecase *ItemUsecase) ItemsQuantity(ctx context.Context) (int, error) {
	usecase.logger.Debug("Enter in usecase ItemsQuantity()")
	if ok := usecase.itemCash.CheckCash(ctx, itemsQuantityKey); !ok {
		if ok := usecase.itemCash.CheckCash(ctx, itemsListKey); !ok {
			_, err := usecase.ItemsList(ctx, 0, 1)
			if err != nil {
				return -1, fmt.Errorf("error on create items list: %w", err)
			}
		} else {
			items, err := usecase.itemCash.GetItemsCash(ctx, itemsListKey)
			if err != nil {
				return -1, fmt.Errorf("error on get items list cash: %w", err)
			}
			if items == nil {
				return -1, fmt.Errorf("items list is not exists")
			}
			err = usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey)
			if err != nil {
				return -1, fmt.Errorf("error on create items quantity cash: %w", err)
			}
		}
	}
	quantity, err := usecase.itemCash.GetItemsQuantityCash(ctx, itemsQuantityKey)
	return quantity, err
}

// SearchLine call database method and returns chan with all models.Item with given params or error
func (usecase *ItemUsecase) SearchLine(ctx context.Context, param string, offset, limit int) ([]models.Item, error) {
	usecase.logger.Debug("Enter in usecase SearchLine()")
	if ok := usecase.itemCash.CheckCash(ctx, param); !ok {
		itemIncomingChan, err := usecase.itemStore.SearchLine(ctx, param)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			usecase.logger.Debug(fmt.Sprintf("item from channel is: %v", item))
			items = append(items, item)
		}
		err = usecase.itemCash.CreateItemsCash(ctx, items, param)
		if err != nil {
			return nil, fmt.Errorf("error on create search list cash: %w", err)
		}
	}

	items, err := usecase.itemCash.GetItemsCash(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("error on get cash: %w", err)
	}
	if offset > len(items) {
		return nil, fmt.Errorf("error: offset bigger than lenght of items, offset: %d, lenght of items: %d", offset, len(items))
	}
	itemsWithLimit := make([]models.Item, 0, limit)
	var counter = 0
	for i := offset; i < len(items); i++ {
		if counter == limit {
			break
		}
		itemsWithLimit = append(itemsWithLimit, items[i])
		counter++
	}
	return itemsWithLimit, nil
}

// GetItemsByCategory call database method and returns chan with all models.Item in category or error
func (usecase *ItemUsecase) GetItemsByCategory(ctx context.Context, categoryName string, offset, limit int) ([]models.Item, error) {
	usecase.logger.Debug("Enter in usecase GetItemsByCategory()")
	if ok := usecase.itemCash.CheckCash(ctx, categoryName); !ok {
		itemIncomingChan, err := usecase.itemStore.GetItemsByCategory(ctx, categoryName)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			usecase.logger.Debug(fmt.Sprintf("item from channel is: %v", item))
			items = append(items, item)
		}
		err = usecase.itemCash.CreateItemsCash(ctx, items, categoryName)
		if err != nil {
			return nil, fmt.Errorf("error on create search list cash: %w", err)
		}
	}

	items, err := usecase.itemCash.GetItemsCash(ctx, categoryName)
	if err != nil {
		return nil, fmt.Errorf("error on get cash: %w", err)
	}
	if offset > len(items) {
		return nil, fmt.Errorf("error: offset bigger than lenght of items, offset: %d, lenght of items: %d", offset, len(items))
	}
	itemsWithLimit := make([]models.Item, 0, limit)
	var counter = 0
	for i := offset; i < len(items); i++ {
		if counter == limit {
			break
		}
		itemsWithLimit = append(itemsWithLimit, items[i])
		counter++
	}
	return itemsWithLimit, nil
}

// updateCash updating cash when creating or updating item
func (usecase *ItemUsecase) UpdateCash(ctx context.Context, id uuid.UUID, op string) error {
	usecase.logger.Debug("Enter in usecase UpdateCash()")
	if !usecase.itemCash.CheckCash(ctx, itemsListKey) {
		return fmt.Errorf("cash is not exists")
	}
	newItem, err := usecase.itemStore.GetItem(ctx, id)
	if err != nil {
		return fmt.Errorf("error on get item: %w", err)
	}
	items, err := usecase.itemCash.GetItemsCash(ctx, itemsListKey)
	if err != nil {
		return fmt.Errorf("error on get cash: %w", err)
	}
	if op == "update" {
		for i, item := range items {
			if item.Id == id {
				items[i] = *newItem
				break
			}
		}
	}
	if op == "create" {
		items = append(items, *newItem)
		err := usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey)
		if err != nil {
			return fmt.Errorf("error on create items quantity cash: %w", err)
		}
	}

	return usecase.itemCash.CreateItemsCash(ctx, items, itemsListKey)
}
