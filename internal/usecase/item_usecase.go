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
	itemCash  cash.IItemsCash
	logger    *zap.Logger
}

func NewItemUsecase(itemStore repository.ItemStore, itemCash cash.IItemsCash, logger *zap.Logger) IItemUsecase {
	logger.Debug("Enter in usecase NewItemUsecase()")
	return &ItemUsecase{itemStore: itemStore, itemCash: itemCash, logger: logger}
}

// CreateItem call database method and returns id of created item or error
func (usecase *ItemUsecase) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase CreateItem() with args: ctx, item: %v", item)
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
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateItem() with args: ctx, item: %v", item)
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
	usecase.logger.Sugar().Debugf("Enter in usecase GetItem() with args: ctx, id: %v", id)
	item, err := usecase.itemStore.GetItem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error on get item: %w", err)
	}
	return item, nil
}

// ItemsList call database method and returns slice with all models.Item or error
func (usecase *ItemUsecase) ItemsList(ctx context.Context, offset, limit int) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsList() with args: ctx, offset: %d, limit: %d", offset, limit)
	if ok := usecase.itemCash.CheckCash(ctx, itemsListKey); !ok {
		itemIncomingChan, err := usecase.itemStore.ItemsList(ctx)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
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

// ItemsQuantity check cash and if cash not exists call database
// method and write in cash and returns quantity of all items
func (usecase *ItemUsecase) ItemsQuantity(ctx context.Context) (int, error) {
	usecase.logger.Debug("Enter in usecase ItemsQuantity() with args: ctx")
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
				items = make([]models.Item, 0)
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

func (usecase *ItemUsecase) ItemsQuantityInCategory(ctx context.Context, categoryName string) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsQuantityInCategory() with args: ctx, categoryName: %s", categoryName)
	if ok := usecase.itemCash.CheckCash(ctx, categoryName+"Quantity"); !ok {
		if ok := usecase.itemCash.CheckCash(ctx, categoryName); !ok {
			_, err := usecase.GetItemsByCategory(ctx, categoryName, 0, 1)
			if err != nil {
				return -1, fmt.Errorf("error on create items list: %w", err)
			}
		} else {
			items, err := usecase.itemCash.GetItemsCash(ctx, categoryName)
			if err != nil {
				return -1, fmt.Errorf("error on get items list cash: %w", err)
			}
			if items == nil {
				items = make([]models.Item, 0)
			}
			err = usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), categoryName+"Quantity")
			if err != nil {
				return -1, fmt.Errorf("error on create items quantity cash: %w", err)
			}
		}
	}
	quantity, err := usecase.itemCash.GetItemsQuantityCash(ctx, categoryName+"Quantity")
	return quantity, err
}

// SearchLine call database method and returns chan with all models.Item with given params or error
func (usecase *ItemUsecase) SearchLine(ctx context.Context, param string, offset, limit int) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase SearchLine() with args: ctx, param: %s, offset: %d, limit: %d", param, offset, limit)
	if ok := usecase.itemCash.CheckCash(ctx, param); !ok {
		itemIncomingChan, err := usecase.itemStore.SearchLine(ctx, param)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
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
	usecase.logger.Sugar().Debugf("Enter in usecase GetItemsByCategory() with args: ctx, categoryName: %s, offset: %d, limit: %d", categoryName, offset, limit)
	if ok := usecase.itemCash.CheckCash(ctx, categoryName); !ok {
		itemIncomingChan, err := usecase.itemStore.GetItemsByCategory(ctx, categoryName)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		err = usecase.itemCash.CreateItemsCash(ctx, items, categoryName)
		if err != nil {
			return nil, fmt.Errorf("error on create get items by category cash: %w", err)
		}
		err = usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), categoryName+"Quantity")
		if err != nil {
			return nil, fmt.Errorf("error on create items quantity in category cash: %w", err)
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

// UpdateCash updating cash when creating or updating item
func (usecase *ItemUsecase) UpdateCash(ctx context.Context, id uuid.UUID, op string) error {
	usecase.logger.Sugar().Debugf("Enter in itemUsecase UpdateCash() with args: ctx, id: %v, op: %s", id, op)
	if !usecase.itemCash.CheckCash(ctx, itemsListKey) {
		return fmt.Errorf("cash is not exists")
	}
	items, err := usecase.itemCash.GetItemsCash(ctx, itemsListKey)
	if err != nil {
		return fmt.Errorf("error on get cash: %w", err)
	}
	newItem := &models.Item{}
	if op == "delete" {
		newItem.Id = id
	} else {
		newItem, err = usecase.itemStore.GetItem(ctx, id)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get item: %v", err)
			return err
		}
		err = usecase.UpdateItemsInCategoryCash(ctx, newItem, op)
		if err != nil {
			usecase.logger.Error(err.Error())
		}
	}
	fmt.Printf("newItem is: %v\n", newItem)

	if op == "update" {
		for i, item := range items {
			if item.Id == newItem.Id {
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
	if op == "delete" {
		for i, item := range items {
			if item.Id == newItem.Id {
				items = append(items[:i], items[i+1:]...)
				err := usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey)
				if err != nil {
					return fmt.Errorf("error on create items quantity cash: %w", err)
				}
				break
			}
		}
	}
	err = usecase.itemCash.CreateItemsCash(ctx, items, itemsListKey)
	if err != nil {
		return err
	}
	usecase.logger.Info("Cash of items list update success")
	return nil
}

// UpdateItemsInCategoryCash update cash items from category
func (usecase *ItemUsecase) UpdateItemsInCategoryCash(ctx context.Context, newItem *models.Item, op string) error {
	usecase.logger.Debug(fmt.Sprintf("Enter in usecase UpdateItemsInCategoryCash() with args: ctx, newItem: %v, op: %s", newItem, op))
	categoryItemsKey := newItem.Category.Name
	categoryItemsQuantityKey := categoryItemsKey + "Quantity"

	if !usecase.itemCash.CheckCash(ctx, categoryItemsKey) {
		return fmt.Errorf("cash with key: %s is not exist", categoryItemsKey)
	}
	items, err := usecase.itemCash.GetItemsCash(ctx, categoryItemsKey)
	if err != nil {
		return fmt.Errorf("error on get cash: %w", err)
	}
	usecase.logger.Debug(fmt.Sprintf("items after get items cash: %v", items))
	if op == "update" {
		for i, item := range items {
			if item.Id == newItem.Id {
				items[i] = *newItem
				break
			}
		}
	}
	if op == "create" {
		items = append(items, *newItem)
		err := usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), categoryItemsQuantityKey)
		if err != nil {
			return fmt.Errorf("error on create items quantity cash: %w", err)
		}
	}
	if op == "delete" {
		for i, item := range items {
			if item.Id == newItem.Id {
				items = append(items[:i], items[i+1:]...)
				usecase.logger.Debug(fmt.Sprintf("items after delete item: %v", items))
				err := usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), categoryItemsQuantityKey)
				if err != nil {
					return fmt.Errorf("error on create items quantity cash: %w", err)
				}
				break
			}
		}
	}
	err = usecase.itemCash.CreateItemsCash(ctx, items, categoryItemsKey)
	if err != nil {
		return err
	}
	usecase.logger.Info("Delete category list cash success")
	return nil
}

// DeleteItem call database method for deleting item
func (usecase *ItemUsecase) DeleteItem(ctx context.Context, id uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteItem() with args: ctx, id: %v", id)
	err := usecase.itemStore.DeleteItem(ctx, id)
	if err != nil {
		return err
	}
	err = usecase.UpdateCash(ctx, id, "delete")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cash: %v", err))
	}
	return nil
}
