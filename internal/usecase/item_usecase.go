package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"OnlineShopBackend/internal/repository/cash"
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ IItemUsecase = &ItemUsecase{}

const (
	itemsListKey          = "ItemsList"
	itemsListKeyNameAsc   = "ItemsListnameasc"
	itemsListKeyNameDesc  = "ItemsListnamedesc"
	itemsListKeyPriceAsc  = "ItemsListpriceasc"
	itemsListKeyPriceDesc = "ItemsListpricedesc"
	itemsQuantityKey      = "ItemsQuantity"
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
		usecase.logger.Debug(err.Error())
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
		usecase.logger.Debug(err.Error())
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
func (usecase *ItemUsecase) ItemsList(ctx context.Context, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsList() with args: ctx, limitOptions: %v, sortOptions: %v", limitOptions, sortOptions)

	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]

	if ok := usecase.itemCash.CheckCash(ctx, itemsListKey+sortType+sortOrder); !ok {
		itemIncomingChan, err := usecase.itemStore.ItemsList(ctx)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		usecase.SortItems(items, sortType, sortOrder)
		err = usecase.itemCash.CreateItemsCash(ctx, items, itemsListKey+sortType+sortOrder)
		if err != nil {
			return nil, fmt.Errorf("error on create items list cash: %w", err)
		}
		err = usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey)
		if err != nil {
			return nil, fmt.Errorf("error on create items quantity cash: %w", err)
		}
	}
	items, err := usecase.itemCash.GetItemsCash(ctx, itemsListKey+sortType+sortOrder)
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
		if ok := usecase.itemCash.CheckCash(ctx, itemsListKeyNameAsc); !ok {
			limitOptions := map[string]int{"offset": 0, "limit": 1}
			sortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}
			_, err := usecase.ItemsList(ctx, limitOptions, sortOptions)
			if err != nil {
				return -1, fmt.Errorf("error on create items list: %w", err)
			}
		} else {
			items, err := usecase.itemCash.GetItemsCash(ctx, itemsListKeyNameAsc)
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
		if ok := usecase.itemCash.CheckCash(ctx, categoryName+"nameasc"); !ok {
			limitOptions := map[string]int{"offset": 0, "limit": 1}
			sortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}
			_, err := usecase.GetItemsByCategory(ctx, categoryName, limitOptions, sortOptions)
			if err != nil {
				return -1, fmt.Errorf("error on create items list: %w", err)
			}
		} else {
			items, err := usecase.itemCash.GetItemsCash(ctx, categoryName+"nameasc")
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
func (usecase *ItemUsecase) SearchLine(ctx context.Context, param string, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase SearchLine() with args: ctx, param: %s, limitOptions: %v, sortOptions: %v", param, limitOptions, sortOptions)

	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]

	if ok := usecase.itemCash.CheckCash(ctx, param+sortType+sortOrder); !ok {
		itemIncomingChan, err := usecase.itemStore.SearchLine(ctx, param)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		usecase.SortItems(items, sortType, sortOrder)
		err = usecase.itemCash.CreateItemsCash(ctx, items, param+sortType+sortOrder)
		if err != nil {
			return nil, fmt.Errorf("error on create search list cash: %w", err)
		}
	}

	items, err := usecase.itemCash.GetItemsCash(ctx, param+sortType+sortOrder)
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
func (usecase *ItemUsecase) GetItemsByCategory(ctx context.Context, categoryName string, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetItemsByCategory() with args: ctx, categoryName: %s, limitOptions: %v, sortOptions: %v", categoryName, limitOptions, sortOptions)

	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]

	if ok := usecase.itemCash.CheckCash(ctx, categoryName+sortType+sortOrder); !ok {
		itemIncomingChan, err := usecase.itemStore.GetItemsByCategory(ctx, categoryName)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}

		usecase.SortItems(items, sortType, sortOrder)

		err = usecase.itemCash.CreateItemsCash(ctx, items, categoryName+sortType+sortOrder)
		if err != nil {
			return nil, fmt.Errorf("error on create get items by category cash: %w", err)
		}
		err = usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), categoryName+"Quantity")
		if err != nil {
			return nil, fmt.Errorf("error on create items quantity in category cash: %w", err)
		}
	}

	items, err := usecase.itemCash.GetItemsCash(ctx, categoryName+sortType+sortOrder)
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

	if !usecase.itemCash.CheckCash(ctx, itemsListKeyNameAsc) && !usecase.itemCash.CheckCash(ctx, itemsListKeyNameDesc) && !usecase.itemCash.CheckCash(ctx, itemsListKeyPriceAsc) && !usecase.itemCash.CheckCash(ctx, itemsListKeyPriceDesc) {
		return fmt.Errorf("cash is not exists")
	}
	newItem := &models.Item{}
	cashKeys := []string{itemsListKeyNameAsc, itemsListKeyNameDesc, itemsListKeyPriceAsc, itemsListKeyPriceDesc}

	for _, key := range cashKeys {
		items, err := usecase.itemCash.GetItemsCash(ctx, key)
		if err != nil {
			return fmt.Errorf("error on get cash: %w", err)
		}
		if op == "delete" {
			newItem.Id = id
		} else {
			newItem, err = usecase.itemStore.GetItem(ctx, id)
			if err != nil {
				usecase.logger.Sugar().Errorf("error on get item: %v", err)
				return err
			}
		}

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

		switch {
		case key == itemsListKeyNameAsc:
			usecase.SortItems(items, "name", "asc")
		case key == itemsListKeyNameDesc:
			usecase.SortItems(items, "name", "desc")
		case key == itemsListKeyPriceAsc:
			usecase.SortItems(items, "price", "asc")
		case key == itemsListKeyPriceDesc:
			usecase.SortItems(items, "price", "desc")
		}
		err = usecase.itemCash.CreateItemsCash(ctx, items, key)
		if err != nil {
			return err
		}
		usecase.logger.Sugar().Infof("Cash of items list with key: %s update success", key)
	}
	err := usecase.UpdateItemsInCategoryCash(ctx, newItem, op)
	if err != nil {
		usecase.logger.Error(err.Error())
	}
	return nil
}

// UpdateItemsInCategoryCash update cash items from category
func (usecase *ItemUsecase) UpdateItemsInCategoryCash(ctx context.Context, newItem *models.Item, op string) error {
	usecase.logger.Debug(fmt.Sprintf("Enter in usecase UpdateItemsInCategoryCash() with args: ctx, newItem: %v, op: %s", newItem, op))
	categoryItemsKeyNameAsc := newItem.Category.Name + "nameasc"
	categoryItemsKeyNameDesc := newItem.Category.Name + "namedesc"
	categoryItemsKeyPriceAsc := newItem.Category.Name + "priceasc"
	categoryItemsKeyPriceDesc := newItem.Category.Name + "pricedesc"
	categoryItemsQuantityKey := newItem.Category.Name + "Quantity"

	keys := []string{categoryItemsKeyNameAsc, categoryItemsKeyNameDesc, categoryItemsKeyPriceAsc, categoryItemsKeyPriceDesc}

	if !usecase.itemCash.CheckCash(ctx, categoryItemsKeyNameAsc) &&
		!usecase.itemCash.CheckCash(ctx, categoryItemsKeyNameDesc) &&
		!usecase.itemCash.CheckCash(ctx, categoryItemsKeyPriceAsc) &&
		!usecase.itemCash.CheckCash(ctx, categoryItemsKeyPriceDesc) {
		return fmt.Errorf("cash is not exist")
	}
	for _, key := range keys {
		items, err := usecase.itemCash.GetItemsCash(ctx, key)
		if err != nil {
			return fmt.Errorf("error on get cash: %w", err)
		}
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
					err := usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), categoryItemsQuantityKey)
					if err != nil {
						return fmt.Errorf("error on create items quantity cash: %w", err)
					}
					break
				}
			}
		}
		switch {
		case key == categoryItemsKeyNameAsc:
			usecase.SortItems(items, "name", "asc")
		case key == categoryItemsKeyNameDesc:
			usecase.SortItems(items, "name", "desc")
		case key == categoryItemsKeyPriceAsc:
			usecase.SortItems(items, "price", "asc")
		case key == categoryItemsKeyPriceDesc:
			usecase.SortItems(items, "price", "desc")
		}
		err = usecase.itemCash.CreateItemsCash(ctx, items, key)
		if err != nil {
			return err
		}
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

// SortItems sorts items
func (usecase *ItemUsecase) SortItems(items []models.Item, sortType string, sortOrder string) {
	usecase.logger.Sugar().Debugf("Enter in usecase SortItems() with args: items []models.Item, sortType: %s, sortOrder: %s", sortType, sortOrder)
	sortType = strings.ToLower(sortType)
	sortOrder = strings.ToLower(sortOrder)
	switch {
	case sortType == "name" && sortOrder == "asc":
		sort.Slice(items, func(i, j int) bool { return items[i].Title < items[j].Title })
		return
	case sortType == "name" && sortOrder == "desc":
		sort.Slice(items, func(i, j int) bool { return items[i].Title > items[j].Title })
		return
	case sortType == "price" && sortOrder == "asc":
		sort.Slice(items, func(i, j int) bool { return items[i].Price < items[j].Price })
		return
	case sortType == "price" && sortOrder == "desc":
		sort.Slice(items, func(i, j int) bool { return items[i].Price > items[j].Price })
		return
	default:
		usecase.logger.Sugar().Errorf("unknown type of sort: %v", sortType)
	}
}

func (usecase *ItemUsecase) AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase AddFavouriteItem() with args: ctx, userId: %v, itemId: %v", userId, itemId)
	err := usecase.itemStore.AddFavouriteItem(ctx, userId, itemId)
	if err != nil {
		return err
	}
	return nil
}

func (usecase *ItemUsecase) DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteFavouriteItem() with args: ctx, userId: %v, itemId: %v", userId, itemId)
	err := usecase.itemStore.DeleteFavouriteItem(ctx, userId, itemId)
	if err != nil {
		return err
	}
	return nil
}

func (usecase *ItemUsecase) GetFavouriteItems(ctx context.Context, userId uuid.UUID) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetFavouriteItems() with args: ctx, userId: %v", userId)
	items, err := usecase.itemStore.GetFavouriteItems(ctx, userId)
	if err != nil {
		return nil, err
	}
	return items, nil
}
