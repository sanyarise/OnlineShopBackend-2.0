package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"OnlineShopBackend/internal/repository/cash"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ IItemUsecase = &ItemUsecase{}

// Keys for create and get cache
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

// AddFavouriteItem added item in list of favourites items
func (usecase *ItemUsecase) AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase AddFavouriteItem() with args: ctx, userId: %v, itemId: %v", userId, itemId)
	err := usecase.itemStore.AddFavouriteItem(ctx, userId, itemId)
	if err != nil {
		return err
	}
	usecase.UpdateFavouriteItemsCash(ctx, userId, itemId, "add")
	usecase.UpdateFavIdsCash(ctx, userId, itemId, "add")
	return nil
}

// DeleteFavouriteItem deleted item from list of favourites items
func (usecase *ItemUsecase) DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteFavouriteItem() with args: ctx, userId: %v, itemId: %v", userId, itemId)
	err := usecase.itemStore.DeleteFavouriteItem(ctx, userId, itemId)
	if err != nil {
		return err
	}
	usecase.UpdateFavouriteItemsCash(ctx, userId, itemId, "delete")
	usecase.UpdateFavIdsCash(ctx, userId, itemId, "delete")
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

// ItemsQuantity check cash and if cash not exists call database
// method and write in cash and returns quantity of all items
func (usecase *ItemUsecase) ItemsQuantity(ctx context.Context) (int, error) {
	usecase.logger.Debug("Enter in usecase ItemsQuantity() with args: ctx")
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	// Сheck the existence of a cache with the quantity of items
	if ok := usecase.itemCash.CheckCash(ctxT, itemsQuantityKey); !ok {
		// If a cache with the quantity of items does not exist,
		// check whether there is a cache with a list of items in the basic sorting
		if ok := usecase.itemCash.CheckCash(ctxT, itemsListKeyNameAsc); !ok {
			// If a cache with a list of items does not exist,
			// request a list of items from the database,
			// in this case, a cache is formed with a list of items
			limitOptions := map[string]int{"offset": 0, "limit": 1}
			sortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}
			_, err := usecase.ItemsList(ctxT, limitOptions, sortOptions)
			if err != nil {
				usecase.logger.Sugar().Warnf("error on create items list cache: %v", err)
			}
		} else {
			// If cache with list of items already exists
			// request an items list cache
			items, err := usecase.itemCash.GetItemsCash(ctxT, itemsListKeyNameAsc)
			if err != nil {
				usecase.logger.Sugar().Warnf("error on get items list cash with key: %s, error: %v", itemsListKeyNameAsc, err)
			}
			// If cache returns nil result, create an empty list of items
			if items == nil {
				items = make([]models.Item, 0)
			}
			// Create cache with items quantity
			err = usecase.itemCash.CreateItemsQuantityCash(ctxT, len(items), itemsQuantityKey)
			if err != nil {
				usecase.logger.Sugar().Warnf("error on create items quantity cash with key: %s, error: %v", itemsQuantityKey, err)
			}
		}
	}
	quantity, err := usecase.itemCash.GetItemsQuantityCash(ctxT, itemsQuantityKey)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get items quantity cash with key: %s, error: %v", itemsQuantityKey, err)
		// If get cache impossible get items from database
		itemsChan, err := usecase.itemStore.ItemsList(ctx)
		if err != nil {
			return -1, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemsChan {
			items = append(items, item)
		}
		quantity = len(items)
	}
	usecase.logger.Info("Get items quantity success")
	return quantity, nil
}

// ItemsQuantityInCategory check cash and if cash not exists call database
// method and write in cash and returns quantity of items in category
func (usecase *ItemUsecase) ItemsQuantityInCategory(ctx context.Context, categoryName string) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsQuantityInCategory() with args: ctx, categoryName: %s", categoryName)
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	// Сheck the existence of a cache with the quantity of items in category
	if ok := usecase.itemCash.CheckCash(ctxT, categoryName+"Quantity"); !ok {
		// If a cache with the quantity of items in category does not exist,
		// check whether there is a cache with a list of items in category in the basic sorting
		if ok := usecase.itemCash.CheckCash(ctxT, categoryName+"nameasc"); !ok {
			// If a cache with a list of items in category does not exist,
			// request a list of items in category from the database,
			// in this case, a cache is formed with a list of items in category
			limitOptions := map[string]int{"offset": 0, "limit": 1}
			sortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}
			_, err := usecase.GetItemsByCategory(ctx, categoryName, limitOptions, sortOptions)
			if err != nil {
				usecase.logger.Sugar().Warnf("error on create items list in category cache: %v", err)
			}
		} else {
			// If cache with list of items in category already exists
			// request an items list in category cache
			items, err := usecase.itemCash.GetItemsCash(ctxT, categoryName+"nameasc")
			if err != nil {
				usecase.logger.Sugar().Warnf("error on get items list in category cash with key: %s, error: %v", categoryName+"nameasc", err)
			} else {
				usecase.logger.Sugar().Infof("Get items list in category with key: %s success", categoryName+"nameasc")
			}
			// If cache returns nil result, create an empty list of items in category
			if items == nil {
				items = make([]models.Item, 0)
			}
			// Create cache with quantity of items in category
			err = usecase.itemCash.CreateItemsQuantityCash(ctxT, len(items), categoryName+"Quantity")
			if err != nil {
				usecase.logger.Sugar().Warnf("error on create quantity of items in category cache with key: %s, error: %v", categoryName+"Quantity", err)
			} else {
				usecase.logger.Sugar().Infof("Create items quantity of items in category with key: %s success", categoryName+"Quantity")
			}
		}
	}
	quantity, err := usecase.itemCash.GetItemsQuantityCash(ctxT, categoryName+"Quantity")
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get quantity of items in category cache with key: %s, error: %v", categoryName+"Quantity", err)
		// If get cache impossible get items from database
		itemsChan, err := usecase.itemStore.GetItemsByCategory(ctx, categoryName)
		if err != nil {
			return -1, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemsChan {
			items = append(items, item)
		}
		quantity = len(items)
	}
	usecase.logger.Info("Get quantity of items in category success")
	return quantity, nil
}

// ItemsQuantityInSearch check cash and if cash not exists call database method and write
// in cash and returns quantity of items in search request
func (usecase *ItemUsecase) ItemsQuantityInSearch(ctx context.Context, search string) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsQuantityInCategory() with args: ctx, search: %s", search)
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	// Сheck the existence of a cache with the quantity of items in search
	if ok := usecase.itemCash.CheckCash(ctxT, search+"Quantity"); !ok {
		// If a cache with the quantity of items does not exist,
		// check whether there is a cache with a list of items in the basic sorting
		if ok := usecase.itemCash.CheckCash(ctxT, search+"nameasc"); !ok {
			// If a cache with a list of items in search does not exist,
			// request a list of items from the database, in this case,
			// a cache is formed with a list of items in search
			limitOptions := map[string]int{"offset": 0, "limit": 1}
			sortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}
			_, err := usecase.SearchLine(ctx, search, limitOptions, sortOptions)
			if err != nil {
				usecase.logger.Sugar().Warnf("error on create items list in search cache: %v", err)
			}
		} else {
			// If cache with list of items already exists
			// request an items list in search cache
			items, err := usecase.itemCash.GetItemsCash(ctxT, search+"nameasc")
			if err != nil {
				usecase.logger.Sugar().Warnf("error on get items list cash with key: %s, error: %v", search+"nameasc", err)
				if items == nil {
					// If cache returns nil result, create an empty list of items
					items = make([]models.Item, 0)
				}
				// Create cache with items quantity in search
				err = usecase.itemCash.CreateItemsQuantityCash(ctxT, len(items), search+"Quantity")
				if err != nil {
					usecase.logger.Sugar().Warnf("error on create items quantity in search cash with key: %s, error: %v", search+"Quantity", err)
				} else {
					usecase.logger.Sugar().Infof("Create items quantity in search cash with key: %s success", search+"Quantity")
				}
			}
		}
	}
	quantity, err := usecase.itemCash.GetItemsQuantityCash(ctxT, search+"Quantity")
	if err != nil {
		// If get cache impossible get items in search from database
		itemsChan, err := usecase.itemStore.SearchLine(ctx, search)
		if err != nil {
			return -1, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemsChan {
			items = append(items, item)
		}
		quantity = len(items)
	}
	return quantity, nil
}

// ItemsQuantityInFavourite check cash and if cash not exists call database
// method and write in cash and returns quantity of items in favourite
func (usecase *ItemUsecase) ItemsQuantityInFavourite(ctx context.Context, userId uuid.UUID) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetFavouriteQuantity() with args: ctx, userId: %v", userId)
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	// Сheck the existence of a cache with the quantity of items in favourite
	if ok := usecase.itemCash.CheckCash(ctxT, userId.String()+"Quantity"); !ok {
		// If a cache with the quantity of items in category does not exist,
		// check whether there is a cache with a list of items in favourite
		// in the basic sorting
		if ok := usecase.itemCash.CheckCash(ctxT, userId.String()+"nameasc"); !ok {
			// If a cache with a list of items in favourite does not exist,
			// request a list of items in favourite from the database,
			// in this case, a cache is formed with a list of items in favourite
			limitOptions := map[string]int{"offset": 0, "limit": 1}
			sortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}
			_, err := usecase.GetFavouriteItems(ctx, userId, limitOptions, sortOptions)
			if err != nil {
				usecase.logger.Sugar().Warnf("error on create items list in favourite cache: %v", err)
			}
		} else {
			// If cache with list of items in favourite already exists
			// request an items list in favourite cache
			items, err := usecase.itemCash.GetItemsCash(ctxT, userId.String()+"nameasc")
			if err != nil {
				usecase.logger.Sugar().Warnf("error on get items list in favourite cash with key: %s, error: %v", userId.String()+"nameasc", err)
			} else {
				usecase.logger.Sugar().Infof("Get items list in favourite cash with key: %s success", userId.String()+"nameasc")
			}
			// If cache returns nil result, create an empty list of items in favourite
			if items == nil {
				items = make([]models.Item, 0)
			}
			// Create cache with quantity of items in favourite
			err = usecase.itemCash.CreateItemsQuantityCash(ctxT, len(items), userId.String()+"Quantity")
			if err != nil {
				usecase.logger.Sugar().Warnf("error on create quantity of items in favourite cache with key: %s, error: %v", userId.String()+"Quantity", err)
			} else {
				usecase.logger.Sugar().Infof("Create items in favourite quantity cash with key: %s success", userId.String()+"Quantity")
			}
		}
	}
	quantity, err := usecase.itemCash.GetItemsQuantityCash(ctxT, userId.String()+"Quantity")
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get quantity of items in favourite cash with key: %s, error: %v", userId.String()+"Quantity", err)
		// If get cache impossible get items from database
		itemsChan, err := usecase.itemStore.GetFavouriteItems(ctx, userId)
		if err != nil {
			return -1, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemsChan {
			items = append(items, item)
		}
		quantity = len(items)
	}
	return quantity, nil
}

// ItemsList call database method and returns slice with all models.Item or error
func (usecase *ItemUsecase) ItemsList(ctx context.Context, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsList() with args: ctx, limitOptions: %v, sortOptions: %v", limitOptions, sortOptions)
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]
	// Check whether there is a cache with that name
	if ok := usecase.itemCash.CheckCash(ctxT, itemsListKey+sortType+sortOrder); !ok {
		// If the cache does not exist, request a list of items from the database
		itemIncomingChan, err := usecase.itemStore.ItemsList(ctx)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		// Sort the list of items based on the sorting parameters
		usecase.SortItems(items, sortType, sortOrder)
		// Create a cache with a sorted list of items
		err = usecase.itemCash.CreateItemsCash(ctxT, items, itemsListKey+sortType+sortOrder)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items list cash with key: %s, error: %v", itemsListKey+sortType+sortOrder, err)
		} else {
			usecase.logger.Sugar().Infof("Create items list cash with key: %s success", itemsListKey+sortType+sortOrder)
		}
		// Create a cache with a quantity of items in list
		err = usecase.itemCash.CreateItemsQuantityCash(ctxT, len(items), itemsQuantityKey)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items quantity cash with key: %s, error: %v", itemsQuantityKey, err)
		} else {
			usecase.logger.Sugar().Infof("Create items quantity cash with key: %s success", itemsQuantityKey)
		}
	}
	// Get items list from cache
	items, err := usecase.itemCash.GetItemsCash(ctxT, itemsListKey+sortType+sortOrder)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get cash with key: %s, err: %v", itemsListKey+sortType+sortOrder, err)
		// If error on get cache, request a list of items from the database
		itemIncomingChan, err := usecase.itemStore.ItemsList(ctx)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on itemStore.ItemsList: %v", err)
			return nil, err
		}
		dbItems := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			dbItems = append(dbItems, item)
		}
		// Sort the list of items based on the sorting parameters
		usecase.SortItems(dbItems, sortType, sortOrder)
		
		items = dbItems
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
		// Add items to the resulting list of items until the counter is equal to the limit
		itemsWithLimit = append(itemsWithLimit, items[i])
		counter++
	}
	return itemsWithLimit, nil
}

// GetItemsByCategory call database method and returns chan with all models.Item in category or error
func (usecase *ItemUsecase) GetItemsByCategory(ctx context.Context, categoryName string, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetItemsByCategory() with args: ctx, categoryName: %s, limitOptions: %v, sortOptions: %v", categoryName, limitOptions, sortOptions)
	
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]

	// Check whether there is a cache of items in category
	if ok := usecase.itemCash.CheckCash(ctxT, categoryName+sortType+sortOrder); !ok {
		// If the cache does not exist, request a list of items in
		// category from the database
		itemIncomingChan, err := usecase.itemStore.GetItemsByCategory(ctx, categoryName)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		// Sort the list of items in category based on the sorting parameters
		usecase.SortItems(items, sortType, sortOrder)
		// Create a cache with a sorted list of items in category
		err = usecase.itemCash.CreateItemsCash(ctxT, items, categoryName+sortType+sortOrder)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items cash with key: %s, error: %v", categoryName+sortType+sortOrder, err)
		} else {
			usecase.logger.Sugar().Infof("Create items cash with key: %s success", categoryName+sortType+sortOrder)
		}
		// Create a cache with a quantity of items in category
		err = usecase.itemCash.CreateItemsQuantityCash(ctxT, len(items), categoryName+"Quantity")
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items in category quantity cash with key: %s, error: %v", categoryName+"Quantity", err)
		} else {
			usecase.logger.Sugar().Infof("Create items in category quantity cash with key: %s success", categoryName+"Quantity")
		}
	}
	// Get items list from cache
	items, err := usecase.itemCash.GetItemsCash(ctxT, categoryName+sortType+sortOrder)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get cache with key: %s, error: %v", categoryName+sortType+sortOrder, err)
		// If error on get cache, request a list of items from the database
		itemIncomingChan, err := usecase.itemStore.GetItemsByCategory(ctx, categoryName)
		if err != nil {
			return nil, err
		}
		dbItems := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			dbItems = append(dbItems, item)
		}
		// Sort the list of items based on the sorting parameters
		usecase.SortItems(dbItems, sortType, sortOrder)
		items = dbItems
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
		// Add items to the resulting list of items until the counter is equal to the limit
		itemsWithLimit = append(itemsWithLimit, items[i])
		counter++
	}
	return itemsWithLimit, nil
}

// SearchLine call database method and returns chan with all models.Item with given params or error
func (usecase *ItemUsecase) SearchLine(ctx context.Context, param string, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase SearchLine() with args: ctx, param: %s, limitOptions: %v, sortOptions: %v", param, limitOptions, sortOptions)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]

	// Check whether there is a cache of this search request
	if ok := usecase.itemCash.CheckCash(ctxT, param+sortType+sortOrder); !ok {
		// If the cache does not exist, request a list of items by
		// search request from the database
		itemIncomingChan, err := usecase.itemStore.SearchLine(ctx, param)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		// Create a cache with a quantity of items in list by search request
		err = usecase.itemCash.CreateItemsQuantityCash(ctxT, len(items), param+"Quantity")
		if err != nil {
			usecase.logger.Warn("can't create items quantity cash: %v", zap.Error(err))
		} else {
			usecase.logger.Info("Items quantity cash create success")
		}
		// Sort the list of items in search request based on the sorting parameters
		usecase.SortItems(items, sortType, sortOrder)
		// Create a cache with a sorted list of items in search request
		err = usecase.itemCash.CreateItemsCash(ctxT, items, param+sortType+sortOrder)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create cash of items in search with key: %s, error: %v", param+sortType+sortOrder, err)
		} else {
			usecase.logger.Sugar().Infof("Create cash of items in search with key: %s success", param+sortType+sortOrder)
		}
	}
	// Get items list from cache
	items, err := usecase.itemCash.GetItemsCash(ctxT, param+sortType+sortOrder)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get cache with key: %s, error: %v", param+sortType+sortOrder, err)
		// If error on get cache, request a list of items from the database
		itemIncomingChan, err := usecase.itemStore.SearchLine(ctx, param)
		if err != nil {
			return nil, err
		}
		dbItems := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			dbItems = append(dbItems, item)
		}
		// Sort the list of items based on the sorting parameters
		usecase.SortItems(dbItems, sortType, sortOrder)
		items = dbItems
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
		// Add items to the resulting list of items until the counter is equal to the limit
		itemsWithLimit = append(itemsWithLimit, items[i])
		counter++
	}
	return itemsWithLimit, nil
}

// GetFavouriteItems call database method and returns chan with models.Item from list of favourites item or error
func (usecase *ItemUsecase) GetFavouriteItems(ctx context.Context, userId uuid.UUID, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetFavouriteItems() with args: ctx, userId: %v", userId)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]
	// Check whether there is a cache of items in favourites
	if ok := usecase.itemCash.CheckCash(ctxT, userId.String()+sortType+sortOrder); !ok {
		// If the cache does not exist, request a list of items in
		// favourites from the database
		itemIncomingChan, err := usecase.itemStore.GetFavouriteItems(ctx, userId)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		// Sort the list of items in favourites
		// based on the sorting parameters
		usecase.SortItems(items, sortType, sortOrder)
		// Create a cache with a sorted list of items in favourites
		err = usecase.itemCash.CreateItemsCash(ctxT, items, userId.String()+sortType+sortOrder)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create favourite items cash with key: %s, error: %v", userId.String()+sortType+sortOrder, err)
		} else {
			usecase.logger.Sugar().Infof("Create favourite items cash with key: %s success", userId.String()+sortType+sortOrder)
		}
		// Create a cache with a quantity of items in favourites
		err = usecase.itemCash.CreateItemsQuantityCash(ctxT, len(items), userId.String()+"Quantity")
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items in favourites quantity cash with key: %s, error: %v", userId.String()+"Quantity", err)
		} else {
			usecase.logger.Sugar().Infof("Create items in favourites quantity cash with key: %s success", userId.String()+"Quantity")
		}
	}
	// Get items list from cache
	items, err := usecase.itemCash.GetItemsCash(ctxT, userId.String()+sortType+sortOrder)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get items in favourite cash with key: %s, error: %v", userId.String()+sortType+sortOrder, err)
		// If error on get cache, request a list of items in favourite from the database
		itemIncomingChan, err := usecase.itemStore.GetFavouriteItems(ctx, userId)
		if err != nil {
			return nil, err
		}
		dbItems := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			dbItems = append(dbItems, item)
		}
		// Sort the list of items based on the sorting parameters
		usecase.SortItems(dbItems, sortType, sortOrder)
		items = dbItems
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
		// Add items to the resulting list of items until the counter is equal to the limit
		itemsWithLimit = append(itemsWithLimit, items[i])
		counter++
	}
	return itemsWithLimit, nil
}

// GetFavouriteItemsId calls database method and returns map with identificators of favourite items of user or error
func (usecase *ItemUsecase) GetFavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetFavouriteItemsId() with args: ctx, userId: %v", userId)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	// Check whether there is a cache of identificators of favourite items
	if !usecase.itemCash.CheckCash(ctxT, userId.String()+"Fav") {
		// If the cache does not exist, request a quantity of
		// favourite items
		quantity, err := usecase.ItemsQuantityInFavourite(ctx, userId)
		if err != nil && quantity == -1 {
			usecase.logger.Warn(err.Error())
			return nil, err
		}
		if quantity == 0 {
			return nil, models.ErrorNotFound{}
		}
		// If quantity > 0 request a map with identificators of
		// favourite items from database
		favUids, err := usecase.itemStore.GetFavouriteItemsId(ctx, userId)
		if err != nil && errors.Is(err, models.ErrorNotFound{}) {
			return nil, models.ErrorNotFound{}
		}
		if err != nil {
			return nil, err
		}
		// Create cache with favourite items identificators
		err = usecase.itemCash.CreateFavouriteItemsIdCash(ctxT, *favUids, userId.String()+"Fav")
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create favourite items id cash with key: %s, error: %v", userId.String()+"Fav", err)
		} else {
			usecase.logger.Sugar().Infof("Create favourite items id cash with key: %s success", userId.String()+"Fav")
		}
	}
	// Get favourite items identificators from cache
	favUids, err := usecase.itemCash.GetFavouriteItemsIdCash(ctxT, userId.String()+"Fav")
	if err != nil {
		usecase.logger.Sugar().Errorf("error on get favourite items id cash: %v", err)
		// If error on get cache, request a map of favourite items identificators
		// from the database
		dbFavUids, err := usecase.itemStore.GetFavouriteItemsId(ctx, userId)
		if err != nil && errors.Is(err, models.ErrorNotFound{}) {
			return nil, models.ErrorNotFound{}
		}
		if err != nil {
			return nil, err
		}
		favUids = dbFavUids
	}
	return favUids, nil
}

// UpdateCash updating cash when creating or updating item
func (usecase *ItemUsecase) UpdateCash(ctx context.Context, id uuid.UUID, op string) error {
	usecase.logger.Sugar().Debugf("Enter in itemUsecase UpdateCash() with args: ctx, id: %v, op: %s", id, op)
	// Check the presence of a cache with all possible keys
	if !usecase.itemCash.CheckCash(ctx, itemsListKeyNameAsc) &&
		!usecase.itemCash.CheckCash(ctx, itemsListKeyNameDesc) &&
		!usecase.itemCash.CheckCash(ctx, itemsListKeyPriceAsc) &&
		!usecase.itemCash.CheckCash(ctx, itemsListKeyPriceDesc) {
		// If the cache with any of the keys does not return the error
		return fmt.Errorf("cash is not exists")
	}
	newItem := &models.Item{}
	cashKeys := []string{itemsListKeyNameAsc, itemsListKeyNameDesc, itemsListKeyPriceAsc, itemsListKeyPriceDesc}
	// Sort through all possible keys
	for _, key := range cashKeys {
		// For each key get a cache
		items, err := usecase.itemCash.GetItemsCash(ctx, key)
		if err != nil {
			return fmt.Errorf("error on get cash: %w", err)
		}
		// If the renewal of the cache is associated with
		// removal of item, we use
		// empty item with ID from parameters
		// method
		if op == "delete" {
			newItem.Id = id
		} else {
			// Otherwise, we get item from the database
			newItem, err = usecase.itemStore.GetItem(ctx, id)
			if err != nil {
				usecase.logger.Sugar().Errorf("error on get item: %v", err)
				return err
			}
		}
		// Сhange the list of items in accordance with the operation
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

		// Sort the list of items
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
		// Record the updated cache
		err = usecase.itemCash.CreateItemsCash(ctx, items, key)
		if err != nil {
			return err
		}
		usecase.logger.Sugar().Infof("Cash of items list with key: %s update success", key)
	}
	// Update the cache of the item list in the category
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
	// Check the presence of a cache with all possible keys
	if !usecase.itemCash.CheckCash(ctx, categoryItemsKeyNameAsc) &&
		!usecase.itemCash.CheckCash(ctx, categoryItemsKeyNameDesc) &&
		!usecase.itemCash.CheckCash(ctx, categoryItemsKeyPriceAsc) &&
		!usecase.itemCash.CheckCash(ctx, categoryItemsKeyPriceDesc) {
		// If the cache with any of the keys does not return the error
		return fmt.Errorf("cash is not exist")
	}
	// Sort through all possible keys
	for _, key := range keys {
		// For each key get a cache
		items, err := usecase.itemCash.GetItemsCash(ctx, key)
		if err != nil {
			return fmt.Errorf("error on get cash: %w", err)
		}
		// Сhange the list of items in accordance with the operation
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
		// Sort the list of items
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
		// Record the updated cache
		err = usecase.itemCash.CreateItemsCash(ctx, items, key)
		if err != nil {
			return err
		}
	}
	usecase.logger.Info("Update category list cash success")
	return nil
}

func (usecase *ItemUsecase) UpdateFavouriteItemsCash(ctx context.Context, userId uuid.UUID, itemId uuid.UUID, op string) {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateFavouriteItemsCash() with args: ctx, userId: %v, itemId: %v, op: %s", userId, itemId, op)
	favouriteItemsKeyNameAsc := userId.String() + "nameasc"
	favouriteItemsKeyNameDesc := userId.String() + "namedesc"
	favouriteItemsKeyPriceAsc := userId.String() + "priceasc"
	favouriteItemsKeyPriceDesc := userId.String() + "pricedesc"
	favouriteItemsQuantityKey := userId.String() + "Quantity"

	keys := []string{favouriteItemsKeyNameAsc, favouriteItemsKeyNameDesc, favouriteItemsKeyPriceAsc, favouriteItemsKeyPriceDesc}
	// Check the presence of a cache with all possible keys
	if !usecase.itemCash.CheckCash(ctx, favouriteItemsKeyNameAsc) &&
		!usecase.itemCash.CheckCash(ctx, favouriteItemsKeyNameDesc) &&
		!usecase.itemCash.CheckCash(ctx, favouriteItemsKeyPriceAsc) &&
		!usecase.itemCash.CheckCash(ctx, favouriteItemsKeyPriceDesc) {
		// If the cache with any of the keys does not return the error
		usecase.logger.Error("cash is not exist")
		return
	}
	// Sort through all possible keys
	for _, key := range keys {
		// For each key get a cache
		items, err := usecase.itemCash.GetItemsCash(ctx, key)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get cash: %v", err)
			return
		}
		// Сhange the list of items in accordance with the operation
		if op == "add" {
			item, err := usecase.itemStore.GetItem(ctx, itemId)
			if err != nil {
				usecase.logger.Sugar().Errorf("error on get item: %v", err)
				return
			}
			items = append(items, *item)
			err = usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), favouriteItemsQuantityKey)
			if err != nil {
				usecase.logger.Sugar().Errorf("error on create items quantity cash: %w", err)
				return
			}
		}
		if op == "delete" {
			for i, item := range items {
				if item.Id == itemId {
					items = append(items[:i], items[i+1:]...)
					err := usecase.itemCash.CreateItemsQuantityCash(ctx, len(items), favouriteItemsQuantityKey)
					if err != nil {
						usecase.logger.Sugar().Errorf("error on create items quantity cash: %w", err)
						return
					}
					break
				}
			}
		}
		// Sort the list of items
		switch {
		case key == favouriteItemsKeyNameAsc:
			usecase.SortItems(items, "name", "asc")
		case key == favouriteItemsKeyNameDesc:
			usecase.SortItems(items, "name", "desc")
		case key == favouriteItemsKeyPriceAsc:
			usecase.SortItems(items, "price", "asc")
		case key == favouriteItemsKeyPriceDesc:
			usecase.SortItems(items, "price", "desc")
		}
		// Record the updated cache
		err = usecase.itemCash.CreateItemsCash(ctx, items, key)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on create favourite items cash: %v", err)
			return
		}
	}
	usecase.logger.Info("Update favourite items list cash success")
}

// UpdateFavIdsCash updates cash with favourite items identificators
func (usecase *ItemUsecase) UpdateFavIdsCash(ctx context.Context, userId, itemId uuid.UUID, op string) {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateFavIdsCash() with args userId: %v, itemId: %v", userId, itemId)
	// Check the presence of a cache with key
	if !usecase.itemCash.CheckCash(ctx, userId.String()+"Fav") {
		// If cache doesn't exists create it
		favMap := make(map[uuid.UUID]uuid.UUID)
		// Add itemId in map of favourite
		// item's identificators
		favMap[itemId] = userId

		// Record the cache with favourite items identificators
		err := usecase.itemCash.CreateFavouriteItemsIdCash(ctx, favMap, userId.String()+"Fav")
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create favourite items id cash: %v", err)
			return
		}
		usecase.logger.Info("create favourite items id cash success")
		return
	}
	// If cache exists get it
	favMapLink, err := usecase.itemCash.GetFavouriteItemsIdCash(ctx, userId.String()+"Fav")
	if err != nil {
		usecase.logger.Sugar().Warn("error on get favourite items id cash with key: %v, err: %v", userId.String()+"Fav", err)
		return
	}
	// Сhange the map of favourite items identificators
	// in accordance with the operation
	favMap := *favMapLink
	if op == "add" {
		favMap[itemId] = userId
	}
	if op == "delete" {
		delete(favMap, itemId)
	}
	// Record the updated cache
	err = usecase.itemCash.CreateFavouriteItemsIdCash(ctx, favMap, userId.String()+"Fav")
	if err != nil {
		usecase.logger.Sugar().Warn("error on create favourite items id cash: %v", err)
		return
	}
	usecase.logger.Info("Create favourite items id cash success")
}

// SortItems sorts list of items by sort parameters
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
