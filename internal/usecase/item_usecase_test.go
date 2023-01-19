package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository/mocks"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testItemId    = uuid.New()
	testModelItem = models.Item{
		Title:       "test",
		Description: "test",
		Category:    models.Category{},
	}
	testItemWithId = models.Item{
		Id:          testItemId,
		Title:       "test",
		Description: "test",
		Category:    models.Category{},
	}
	testItemWithId2 = models.Item{
		Id:       testItemId,
		Category: models.Category{},
	}
	emptyItem = models.Item{}
	param     = "test"
	items     = []models.Item{testItemWithId}
	items2    = []models.Item{testItemWithId, testItemWithId2}
	newItem   = &models.Item{
		Id:          testItemId,
		Title:       "test",
		Description: "test",
		Category:    models.Category{},
		Price:       0,
		Vendor:      "test",
	}
	cashItem = models.Item{
		Id:          testItemId,
		Title:       "test",
		Description: "test",
		Category:    models.Category{},
		Price:       0,
		Vendor:      "test",
	}
	testCategoryName          = "testName"
	err                       = errors.New("error")
	testLimitOptionsItemsList = map[string]int{
		"offset": 0,
		"limit":  1,
	}
	testLimitOptionsItemsList2 = map[string]int{
		"offset": 2,
		"limit":  1,
	}
	testSortOptionsItemsList = map[string]string{
		"sortType":  "name",
		"sortOrder": "asc",
	}
)

func TestCreateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	itemRepo.EXPECT().CreateItem(ctx, &testModelItem).Return(uuid.Nil, err)
	res, err := usecase.CreateItem(ctx, &testModelItem)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)

	itemRepo.EXPECT().CreateItem(ctx, &testModelItem).Return(testId, nil)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameDesc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyPriceAsc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyPriceDesc).Return(false)
	res, err = usecase.CreateItem(ctx, &testModelItem)
	require.NoError(t, err)
	require.Equal(t, res, testId)
}

func TestUpdateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	itemRepo.EXPECT().UpdateItem(ctx, &testModelItem).Return(err)
	err := usecase.UpdateItem(ctx, &testModelItem)
	require.Error(t, err)

	itemRepo.EXPECT().UpdateItem(ctx, &testModelItem).Return(nil)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameDesc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyPriceAsc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyPriceDesc).Return(false)
	err = usecase.UpdateItem(ctx, &testModelItem)
	require.NoError(t, err)
}

func TestGetItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(&testItemWithId, nil)
	res, err := usecase.GetItem(ctx, testItemId)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, &testItemWithId)

	err = fmt.Errorf("error on get item")
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(&emptyItem, err)
	res, err = usecase.GetItem(ctx, testItemId)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestItemsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()
	testItemChan := make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)

	cash.EXPECT().CheckCash(ctx, itemsListKey+"nameasc").Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testItemChan, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, itemsListKey+"nameasc").Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey).Return(nil)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey+"nameasc").Return(items, nil)

	res, err := usecase.ItemsList(ctx, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, itemsListKey+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey+"nameasc").Return(items, nil)
	res, err = usecase.ItemsList(ctx, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, itemsListKey+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey+"nameasc").Return(items2, nil)
	res, err = usecase.ItemsList(ctx, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, itemsListKey+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey+"nameasc").Return(items, nil)
	res, err = usecase.ItemsList(ctx, testLimitOptionsItemsList2, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error on itemslist()")
	cash.EXPECT().CheckCash(ctx, itemsListKey+"nameasc").Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testItemChan, err)
	res, err = usecase.ItemsList(ctx, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cash.EXPECT().CheckCash(ctx, itemsListKey+"nameasc").Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan2, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, itemsListKey+"nameasc").Return(err)
	res, err = usecase.ItemsList(ctx, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testItemWithId
	close(testChan3)
	cash.EXPECT().CheckCash(ctx, itemsListKey+"nameasc").Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan3, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, itemsListKey+"nameasc").Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey).Return(nil)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey+"nameasc").Return(nil, err)
	res, err = usecase.ItemsList(ctx, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan4 := make(chan models.Item, 1)
	testChan4 <- testItemWithId
	close(testChan4)
	err = fmt.Errorf("error on create items quantity cash")
	cash.EXPECT().CheckCash(ctx, itemsListKey+"nameasc").Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan4, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, itemsListKey+"nameasc").Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey).Return(err)
	res, err = usecase.ItemsList(ctx, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestSearchLine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	testItemChan := make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testItemChan, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(nil)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items, nil)
	res, err := usecase.SearchLine(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.SearchLine(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items2, nil)
	res, err = usecase.SearchLine(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.SearchLine(ctx, param, testLimitOptionsItemsList2, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error on search()")
	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testItemChan, err)
	res, err = usecase.SearchLine(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan2, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(err)
	res, err = usecase.SearchLine(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testItemWithId
	close(testChan3)
	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan3, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(nil)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(nil, err)
	res, err = usecase.SearchLine(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestGetItemByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	testItemChan := make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testItemChan, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, param+"Quantity")
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items, nil)
	res, err := usecase.GetItemsByCategory(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	testItemChan = make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)
	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testItemChan, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, param+"Quantity").Return(fmt.Errorf("error"))
	res, err = usecase.GetItemsByCategory(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.GetItemsByCategory(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items2, nil)
	res, err = usecase.GetItemsByCategory(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.GetItemsByCategory(ctx, param, testLimitOptionsItemsList2, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error on search()")
	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testItemChan, err)
	res, err = usecase.GetItemsByCategory(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testChan2, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(err)
	res, err = usecase.GetItemsByCategory(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testItemWithId
	close(testChan3)
	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testChan3, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, param+"Quantity")
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(nil, err)
	res, err = usecase.GetItemsByCategory(ctx, param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestItemsQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(true)
	cash.EXPECT().GetItemsQuantityCash(ctx, itemsQuantityKey).Return(1, nil)
	res, err := usecase.ItemsQuantity(ctx)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(true)
	cash.EXPECT().GetItemsQuantityCash(ctx, itemsQuantityKey).Return(-1, fmt.Errorf("error on get items quantity cash"))
	res, err = usecase.ItemsQuantity(ctx)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKeyNameAsc).Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, itemsQuantityKey).Return(nil)
	cash.EXPECT().GetItemsQuantityCash(ctx, itemsQuantityKey).Return(1, nil)
	res, err = usecase.ItemsQuantity(ctx)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKeyNameAsc).Return(nil, fmt.Errorf("error on get items list cash"))
	res, err = usecase.ItemsQuantity(ctx)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKeyNameAsc).Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, itemsQuantityKey).Return(fmt.Errorf("error on create items quantity cash"))
	res, err = usecase.ItemsQuantity(ctx)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKeyNameAsc).Return(items, nil)
	cash.EXPECT().GetItemsQuantityCash(ctx, itemsQuantityKey).Return(1, nil)
	res, err = usecase.ItemsQuantity(ctx)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(nil, fmt.Errorf("error"))
	res, err = usecase.ItemsQuantity(ctx)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKeyNameAsc).Return(nil, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 0, itemsQuantityKey).Return(nil)
	cash.EXPECT().GetItemsQuantityCash(ctx, itemsQuantityKey).Return(0, nil)
	res, err = usecase.ItemsQuantity(ctx)
	require.NoError(t, err)
	require.Equal(t, res, 0)
}

func TestItemsQuantityInCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	q := "Quantity"
	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(true)
	cash.EXPECT().GetItemsQuantityCash(ctx, testCategoryName+q).Return(1, nil)
	res, err := usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(true)
	cash.EXPECT().GetItemsQuantityCash(ctx, testCategoryName+q).Return(-1, fmt.Errorf("error on get items quantity cash"))
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testCategoryName+"nameasc").Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, testCategoryName+q).Return(nil)
	cash.EXPECT().GetItemsQuantityCash(ctx, testCategoryName+q).Return(1, nil)
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testCategoryName+"nameasc").Return(nil, fmt.Errorf("error on get items list cash"))
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testCategoryName+"nameasc").Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, testCategoryName+q).Return(fmt.Errorf("error on create items quantity cash"))
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName+"nameasc").Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testCategoryName+"nameasc")
	cash.EXPECT().GetItemsQuantityCash(ctx, testCategoryName+q).Return(1, nil)
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName+"nameasc").Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, testCategoryName).Return(nil, fmt.Errorf("error"))
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testCategoryName+"nameasc").Return(nil, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 0, testCategoryName+q).Return(nil)
	cash.EXPECT().GetItemsQuantityCash(ctx, testCategoryName+q).Return(0, nil)
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.NoError(t, err)
	require.Equal(t, res, 0)
}

func TestUpdateCash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	cash.EXPECT().CheckCash(ctx, "ItemsList"+"nameasc").Return(false)
	cash.EXPECT().CheckCash(ctx, "ItemsList"+"namedesc").Return(false)
	cash.EXPECT().CheckCash(ctx, "ItemsList"+"priceasc").Return(false)
	cash.EXPECT().CheckCash(ctx, "ItemsList"+"pricedesc").Return(false)
	err := usecase.UpdateCash(ctx, uuid.New(), "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKeyNameAsc).Return(nil, err)
	err = usecase.UpdateCash(ctx, testId, "create")
	require.Error(t, err)

	cashResults := make([]models.Item, 0, 1)
	cashResults = append(cashResults, cashItem)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKeyNameAsc).Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testId).Return(nil, err)
	err = usecase.UpdateCash(ctx, testId, "create")
	require.Error(t, err)

	cashResults = make([]models.Item, 0, 1)
	cashResults = append(cashResults, cashItem)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKeyNameAsc).Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testId).Return(&testItemWithId, nil)
	cash.EXPECT().CreateItemsCash(ctx, cashResults, itemsListKeyNameAsc).Return(err)
	err = usecase.UpdateCash(ctx, testId, "update")
	require.Error(t, err)

	cashResults = make([]models.Item, 0, 1)
	cashResults = append(cashResults, cashItem)
	updateResults := append(cashResults, cashItem)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKeyNameAsc).Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testId).Return(&testItemWithId, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(updateResults), itemsQuantityKey).Return(err)
	err = usecase.UpdateCash(ctx, testId, "create")
	require.Error(t, err)
}

func TestUpdateItemsInCategoryCash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	cashResults := make([]models.Item, 0, 1)
	cashResults = append(cashResults, cashItem)
	updateResults := make([]models.Item, 0, 2)
	updateResults = append(updateResults, *newItem)
	updateResults = append(updateResults, *newItem)

	cash.EXPECT().CheckCash(ctx, newItem.Category.Name+"nameasc").Return(false)
	cash.EXPECT().CheckCash(ctx, newItem.Category.Name+"namedesc").Return(false)
	cash.EXPECT().CheckCash(ctx, newItem.Category.Name+"priceasc").Return(false)
	cash.EXPECT().CheckCash(ctx, newItem.Category.Name+"pricedesc").Return(false)
	err := usecase.UpdateItemsInCategoryCash(ctx, newItem, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, newItem.Category.Name+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"nameasc").Return(nil, fmt.Errorf("error"))
	err = usecase.UpdateItemsInCategoryCash(ctx, newItem, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, newItem.Category.Name+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"nameasc").Return(cashResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(updateResults), newItem.Category.Name+"Quantity").Return(fmt.Errorf("error"))
	err = usecase.UpdateItemsInCategoryCash(ctx, newItem, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, newItem.Category.Name+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"nameasc").Return(cashResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(updateResults), newItem.Category.Name+"Quantity").Return(nil)
	cash.EXPECT().CreateItemsCash(ctx, updateResults, newItem.Category.Name+"nameasc").Return(fmt.Errorf("error"))
	err = usecase.UpdateItemsInCategoryCash(ctx, newItem, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, newItem.Category.Name+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"nameasc").Return(cashResults, nil)
	cash.EXPECT().CreateItemsCash(ctx, cashResults, newItem.Category.Name+"nameasc").Return(nil)

	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"namedesc").Return(cashResults, nil)
	cash.EXPECT().CreateItemsCash(ctx, cashResults, newItem.Category.Name+"namedesc").Return(nil)

	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"priceasc").Return(cashResults, nil)
	cash.EXPECT().CreateItemsCash(ctx, cashResults, newItem.Category.Name+"priceasc").Return(nil)

	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"pricedesc").Return(cashResults, nil)
	cash.EXPECT().CreateItemsCash(ctx, cashResults, newItem.Category.Name+"pricedesc").Return(nil)

	err = usecase.UpdateItemsInCategoryCash(ctx, newItem, "update")
	require.NoError(t, err)

	deletedResults := []models.Item{testItemWithId}
	deleteResults := []models.Item{}
	cash.EXPECT().CheckCash(ctx, newItem.Category.Name+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"nameasc").Return(deletedResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(deleteResults), testItemWithId.Category.Name+"Quantity").Return(err)
	cash.EXPECT().CreateItemsCash(ctx, deleteResults, newItem.Category.Name+"nameasc").Return(nil)

	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"namedesc").Return(deletedResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(deleteResults), testItemWithId.Category.Name+"Quantity").Return(err)
	cash.EXPECT().CreateItemsCash(ctx, deleteResults, newItem.Category.Name+"namedesc").Return(nil)

	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"priceasc").Return(deletedResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(deleteResults), testItemWithId.Category.Name+"Quantity").Return(err)
	cash.EXPECT().CreateItemsCash(ctx, deleteResults, newItem.Category.Name+"priceasc").Return(nil)

	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name+"pricedesc").Return(deletedResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(deleteResults), testItemWithId.Category.Name+"Quantity").Return(err)
	cash.EXPECT().CreateItemsCash(ctx, deleteResults, newItem.Category.Name+"pricedesc").Return(nil)
	err = usecase.UpdateItemsInCategoryCash(ctx, newItem, "delete")
	require.NoError(t, err)
}

func TestDeleteItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	itemRepo.EXPECT().DeleteItem(ctx, testId).Return(err)
	err := usecase.DeleteItem(ctx, testId)
	require.Error(t, err)

	itemRepo.EXPECT().DeleteItem(ctx, testId).Return(nil)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameAsc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyNameDesc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyPriceAsc).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKeyPriceDesc).Return(false)
	err = usecase.DeleteItem(ctx, testId)
	require.NoError(t, err)
}

func TestAddFavouriteItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	itemRepo.EXPECT().AddFavouriteItem(ctx, testId, testItemId).Return(err)
	err := usecase.AddFavouriteItem(ctx, testId, testItemId)
	require.Error(t, err)

	itemRepo.EXPECT().AddFavouriteItem(ctx, testId, testItemId).Return(nil)
	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"namedesc").Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"priceasc").Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"pricedesc").Return(false)
	err = usecase.AddFavouriteItem(ctx, testId, testItemId)
	require.NoError(t, err)
}

func TestDeleteFavouriteItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	itemRepo.EXPECT().DeleteFavouriteItem(ctx, testId, testItemId).Return(err)
	err := usecase.DeleteFavouriteItem(ctx, testId, testItemId)
	require.Error(t, err)

	itemRepo.EXPECT().DeleteFavouriteItem(ctx, testId, testItemId).Return(nil)
	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"namedesc").Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"priceasc").Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"pricedesc").Return(false)
	err = usecase.DeleteFavouriteItem(ctx, testId, testItemId)
	require.NoError(t, err)
}

func TestGetFavouriteItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()
	param = testId.String()
	paramns := testId
	testItemChan := make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testItemChan, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, param+"Quantity")
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items, nil)
	res, err := usecase.GetFavouriteItems(ctx, paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	testItemChan = make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)
	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testItemChan, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, param+"Quantity").Return(fmt.Errorf("error"))
	res, err = usecase.GetFavouriteItems(ctx, paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.GetFavouriteItems(ctx, paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items2, nil)
	res, err = usecase.GetFavouriteItems(ctx, paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.GetFavouriteItems(ctx, paramns, testLimitOptionsItemsList2, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error on search()")
	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testItemChan, err)
	res, err = usecase.GetFavouriteItems(ctx, paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testChan2, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(err)
	res, err = usecase.GetFavouriteItems(ctx, paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testItemWithId
	close(testChan3)
	cash.EXPECT().CheckCash(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testChan3, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param+"nameasc").Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, param+"Quantity")
	cash.EXPECT().GetItemsCash(ctx, param+"nameasc").Return(nil, err)
	res, err = usecase.GetFavouriteItems(ctx, paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestItemsQuantityInFavourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	q := "Quantity"
	cash.EXPECT().CheckCash(ctx, testId.String()+q).Return(true)
	cash.EXPECT().GetItemsQuantityCash(ctx, testId.String()+q).Return(1, nil)
	res, err := usecase.ItemsQuantityInFavourite(ctx, testId)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, testId.String()+q).Return(true)
	cash.EXPECT().GetItemsQuantityCash(ctx, testId.String()+q).Return(-1, fmt.Errorf("error on get items quantity cash"))
	res, err = usecase.ItemsQuantityInFavourite(ctx, testId)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testId.String()+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, testId.String()+q).Return(nil)
	cash.EXPECT().GetItemsQuantityCash(ctx, testId.String()+q).Return(1, nil)
	res, err = usecase.ItemsQuantityInFavourite(ctx, testId)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, testId.String()+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return(nil, fmt.Errorf("error on get items list cash"))
	res, err = usecase.ItemsQuantityInFavourite(ctx, testId)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testId.String()+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, testId.String()+q).Return(fmt.Errorf("error on create items quantity cash"))
	res, err = usecase.ItemsQuantityInFavourite(ctx, testId)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testId.String()+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc")
	cash.EXPECT().GetItemsQuantityCash(ctx, testId.String()+q).Return(1, nil)
	res, err = usecase.ItemsQuantityInFavourite(ctx, testId)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, testId.String()+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(false)
	itemRepo.EXPECT().GetFavouriteItems(ctx, testId).Return(nil, fmt.Errorf("error"))
	res, err = usecase.ItemsQuantityInFavourite(ctx, testId)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testId.String()+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return(nil, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 0, testId.String()+q).Return(nil)
	cash.EXPECT().GetItemsQuantityCash(ctx, testId.String()+q).Return(0, nil)
	res, err = usecase.ItemsQuantityInFavourite(ctx, testId)
	require.NoError(t, err)
	require.Equal(t, res, 0)
}

func TestUpdateFavouriteItemsCash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	cashResults := make([]models.Item, 0, 1)
	cashResults = append(cashResults, cashItem)
	updateResults := make([]models.Item, 0, 2)
	updateResults = append(updateResults, *newItem)
	updateResults = append(updateResults, *newItem)

	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return(nil, err)
	usecase.UpdateFavouriteItemsCash(ctx, testId, testItemId, "add")

	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(nil, err)
	usecase.UpdateFavouriteItemsCash(ctx, testId, testItemId, "add")

	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(nil, err)
	usecase.UpdateFavouriteItemsCash(ctx, testId, testItemId, "add")

	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(&testItem1, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 2, testId.String()+"Quantity").Return(err)
	usecase.UpdateFavouriteItemsCash(ctx, testId, testItemId, "add")

	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 2, testId.String()+"Quantity").Return(nil)
	cash.EXPECT().CreateItemsCash(ctx, updateResults, testId.String()+"nameasc").Return(err)
	usecase.UpdateFavouriteItemsCash(ctx, testId, testItemId, "add")

	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 2, testId.String()+"Quantity").Return(nil)
	cash.EXPECT().CreateItemsCash(ctx, updateResults, testId.String()+"nameasc").Return(nil)

	cash.EXPECT().GetItemsCash(ctx, testId.String()+"namedesc").Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 2, testId.String()+"Quantity").Return(nil)
	cash.EXPECT().CreateItemsCash(ctx, updateResults, testId.String()+"namedesc").Return(nil)

	cash.EXPECT().GetItemsCash(ctx, testId.String()+"priceasc").Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 2, testId.String()+"Quantity").Return(nil)
	cash.EXPECT().CreateItemsCash(ctx, updateResults, testId.String()+"priceasc").Return(nil)

	cash.EXPECT().GetItemsCash(ctx, testId.String()+"pricedesc").Return(cashResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 2, testId.String()+"Quantity").Return(nil)
	cash.EXPECT().CreateItemsCash(ctx, updateResults, testId.String()+"pricedesc").Return(nil)
	usecase.UpdateFavouriteItemsCash(ctx, testId, testItemId, "add")

	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return(updateResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, testId.String()+"Quantity").Return(nil)
	cash.EXPECT().CreateItemsCash(ctx, cashResults, testId.String()+"nameasc").Return(err)
	usecase.UpdateFavouriteItemsCash(ctx, testId, testItemId, "delete")

	cash.EXPECT().CheckCash(ctx, testId.String()+"nameasc").Return(true)
	cash.EXPECT().GetItemsCash(ctx, testId.String()+"nameasc").Return(updateResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, testId.String()+"Quantity").Return(err)
	usecase.UpdateFavouriteItemsCash(ctx, testId, testItemId, "delete")
}

func TestSortItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)

	testItems := []models.Item{
		{Title: "A"},
		{Title: "C"},
		{Title: "B"},
	}
	testItems2 := []models.Item{
		{Price: 10},
		{Price: 30},
		{Price: 20},
	}

	usecase.SortItems(testItems, "name", "asc")
	require.Equal(t, testItems, []models.Item{
		{Title: "A"},
		{Title: "B"},
		{Title: "C"},
	})
	usecase.SortItems(testItems, "name", "desc")
	require.Equal(t, testItems, []models.Item{
		{Title: "C"},
		{Title: "B"},
		{Title: "A"},
	})
	usecase.SortItems(testItems2, "price", "asc")
	require.Equal(t, testItems2, []models.Item{
		{Price: 10},
		{Price: 20},
		{Price: 30},
	})
	usecase.SortItems(testItems2, "price", "desc")
	require.Equal(t, testItems2, []models.Item{
		{Price: 30},
		{Price: 20},
		{Price: 10},
	})
	usecase.SortItems(testItems, "pricee", "desc")
}
