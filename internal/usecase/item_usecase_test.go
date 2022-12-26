package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository/mocks"
	"context"
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
	testCategoryName = "testName"
)

func TestCreateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	itemRepo.EXPECT().CreateItem(ctx, &testModelItem).Return(testItemId, nil)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	res, err := usecase.CreateItem(ctx, &testModelItem)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testItemId)

	items3 := append(items, testItemWithId)
	itemRepo.EXPECT().CreateItem(ctx, &testModelItem).Return(testItemId, nil)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(&testItemWithId, nil)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(items, nil)
	cash.EXPECT().CheckCash(ctx, testItemWithId.Category.Name).Return(false)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(items2), itemsQuantityKey).Return(nil)
	cash.EXPECT().CreateItemsCash(ctx, items3, itemsListKey).Return(nil)
	res, err = usecase.CreateItem(ctx, &testModelItem)
	require.NoError(t, err)
	require.NotNil(t, res)

	err = fmt.Errorf("test error")
	itemRepo.EXPECT().CreateItem(ctx, &testModelItem).Return(uuid.Nil, err)
	res, err = usecase.CreateItem(ctx, &testModelItem)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)
}

func TestUpdateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockIItemsCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	itemRepo.EXPECT().UpdateItem(ctx, &testItemWithId).Return(nil)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	err := usecase.UpdateItem(ctx, &testItemWithId)
	require.NoError(t, err)

	itemRepo.EXPECT().UpdateItem(ctx, &testItemWithId).Return(nil)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(&testItemWithId, nil)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(items, nil)
	cash.EXPECT().CheckCash(ctx, testItemWithId.Category.Name).Return(false)
	cash.EXPECT().CreateItemsCash(ctx, items, itemsListKey).Return(nil)
	err = usecase.UpdateItem(ctx, &testItemWithId)
	require.NoError(t, err)

	err = fmt.Errorf("error on update item")
	itemRepo.EXPECT().UpdateItem(ctx, &testItemWithId).Return(err)
	err = usecase.UpdateItem(ctx, &testItemWithId)
	require.Error(t, err)
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
	require.Equal(t, res, &emptyItem)
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

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testItemChan, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, itemsListKey).Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey).Return(nil)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(items, nil)
	res, err := usecase.ItemsList(ctx, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(items, nil)
	res, err = usecase.ItemsList(ctx, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(items2, nil)
	res, err = usecase.ItemsList(ctx, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(items, nil)
	res, err = usecase.ItemsList(ctx, 2, 1)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error on itemslist()")
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testItemChan, err)
	res, err = usecase.ItemsList(ctx, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan2, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, itemsListKey).Return(err)
	res, err = usecase.ItemsList(ctx, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testItemWithId
	close(testChan3)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan3, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, itemsListKey).Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey).Return(nil)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(nil, err)
	res, err = usecase.ItemsList(ctx, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	testChan4 := make(chan models.Item, 1)
	testChan4 <- testItemWithId
	close(testChan4)
	err = fmt.Errorf("error on create items quantity cash")
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan4, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, itemsListKey).Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(items), itemsQuantityKey).Return(err)
	res, err = usecase.ItemsList(ctx, 0, 1)
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

	cash.EXPECT().CheckCash(ctx, param).Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testItemChan, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param).Return(nil)
	cash.EXPECT().GetItemsCash(ctx, param).Return(items, nil)
	res, err := usecase.SearchLine(ctx, param, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param).Return(true)
	cash.EXPECT().GetItemsCash(ctx, param).Return(items, nil)
	res, err = usecase.SearchLine(ctx, param, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param).Return(true)
	cash.EXPECT().GetItemsCash(ctx, param).Return(items2, nil)
	res, err = usecase.SearchLine(ctx, param, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param).Return(true)
	cash.EXPECT().GetItemsCash(ctx, param).Return(items, nil)
	res, err = usecase.SearchLine(ctx, param, 2, 1)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error on search()")
	cash.EXPECT().CheckCash(ctx, param).Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testItemChan, err)
	res, err = usecase.SearchLine(ctx, param, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cash.EXPECT().CheckCash(ctx, param).Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan2, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param).Return(err)
	res, err = usecase.SearchLine(ctx, param, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testItemWithId
	close(testChan3)
	cash.EXPECT().CheckCash(ctx, param).Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan3, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param).Return(nil)
	cash.EXPECT().GetItemsCash(ctx, param).Return(nil, err)
	res, err = usecase.SearchLine(ctx, param, 0, 1)
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

	cash.EXPECT().CheckCash(ctx, param).Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testItemChan, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param).Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, param+"Quantity")
	cash.EXPECT().GetItemsCash(ctx, param).Return(items, nil)
	res, err := usecase.GetItemsByCategory(ctx, param, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	testItemChan = make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)
	cash.EXPECT().CheckCash(ctx, param).Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testItemChan, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param).Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, param+"Quantity").Return(fmt.Errorf("error"))
	res, err = usecase.GetItemsByCategory(ctx, param, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	cash.EXPECT().CheckCash(ctx, param).Return(true)
	cash.EXPECT().GetItemsCash(ctx, param).Return(items, nil)
	res, err = usecase.GetItemsByCategory(ctx, param, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param).Return(true)
	cash.EXPECT().GetItemsCash(ctx, param).Return(items2, nil)
	res, err = usecase.GetItemsByCategory(ctx, param, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cash.EXPECT().CheckCash(ctx, param).Return(true)
	cash.EXPECT().GetItemsCash(ctx, param).Return(items, nil)
	res, err = usecase.GetItemsByCategory(ctx, param, 2, 1)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error on search()")
	cash.EXPECT().CheckCash(ctx, param).Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testItemChan, err)
	res, err = usecase.GetItemsByCategory(ctx, param, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cash.EXPECT().CheckCash(ctx, param).Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testChan2, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param).Return(err)
	res, err = usecase.GetItemsByCategory(ctx, param, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testItemWithId
	close(testChan3)
	cash.EXPECT().CheckCash(ctx, param).Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testChan3, nil)
	cash.EXPECT().CreateItemsCash(ctx, items, param).Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, param+"Quantity")
	cash.EXPECT().GetItemsCash(ctx, param).Return(nil, err)
	res, err = usecase.GetItemsByCategory(ctx, param, 0, 1)
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
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, itemsQuantityKey).Return(nil)
	cash.EXPECT().GetItemsQuantityCash(ctx, itemsQuantityKey).Return(1, nil)
	res, err = usecase.ItemsQuantity(ctx)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(nil, fmt.Errorf("error on get items list cash"))
	res, err = usecase.ItemsQuantity(ctx)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, itemsQuantityKey).Return(fmt.Errorf("error on create items quantity cash"))
	res, err = usecase.ItemsQuantity(ctx)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey)
	cash.EXPECT().GetItemsQuantityCash(ctx, itemsQuantityKey).Return(1, nil)
	res, err = usecase.ItemsQuantity(ctx)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(nil, fmt.Errorf("error"))
	res, err = usecase.ItemsQuantity(ctx)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(nil, nil)
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
	cash.EXPECT().CheckCash(ctx, testCategoryName).Return(true)
	cash.EXPECT().GetItemsCash(ctx, testCategoryName).Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, testCategoryName+q).Return(nil)
	cash.EXPECT().GetItemsQuantityCash(ctx, testCategoryName+q).Return(1, nil)
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName).Return(true)
	cash.EXPECT().GetItemsCash(ctx, testCategoryName).Return(nil, fmt.Errorf("error on get items list cash"))
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName).Return(true)
	cash.EXPECT().GetItemsCash(ctx, testCategoryName).Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, testCategoryName+q).Return(fmt.Errorf("error on create items quantity cash"))
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName).Return(true)
	cash.EXPECT().GetItemsCash(ctx, testCategoryName)
	cash.EXPECT().GetItemsQuantityCash(ctx, testCategoryName+q).Return(1, nil)
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName).Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, testCategoryName).Return(nil, fmt.Errorf("error"))
	res, err = usecase.ItemsQuantityInCategory(ctx, testCategoryName)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, testCategoryName+q).Return(false)
	cash.EXPECT().CheckCash(ctx, testCategoryName).Return(true)
	cash.EXPECT().GetItemsCash(ctx, testCategoryName).Return(nil, nil)
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

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	err := usecase.UpdateCash(ctx, uuid.New(), "create")
	require.Error(t, err)

	cashResults := make([]models.Item, 0, 1)
	cashResults = append(cashResults, cashItem)
	updateResults := make([]models.Item, 0, 1)
	updateResults = append(updateResults, *newItem)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(nil, fmt.Errorf("error on get cash"))
	err = usecase.UpdateCash(ctx, testItemId, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, err)
	err = usecase.UpdateCash(ctx, testItemId, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(cashResults, nil)
	cash.EXPECT().CheckCash(ctx, newItem.Category.Name).Return(false)
	cash.EXPECT().CreateItemsCash(ctx, updateResults, itemsListKey).Return(fmt.Errorf("error on create cash"))
	err = usecase.UpdateCash(ctx, testItemId, "update")
	require.Error(t, err)

	updateResults = append(cashResults, *newItem)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(cashResults, nil)
	cash.EXPECT().CheckCash(ctx, newItem.Category.Name).Return(false)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(updateResults), itemsQuantityKey).Return(nil)
	cash.EXPECT().CreateItemsCash(ctx, updateResults, itemsListKey).Return(nil)
	err = usecase.UpdateCash(ctx, testItemId, "create")
	require.NoError(t, err)

	updateResults = append(cashResults, *newItem)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cash.EXPECT().GetItemsCash(ctx, itemsListKey).Return(cashResults, nil)
	cash.EXPECT().CheckCash(ctx, newItem.Category.Name).Return(false)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(updateResults), itemsQuantityKey).Return(fmt.Errorf("error"))
	err = usecase.UpdateCash(ctx, testItemId, "create")
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

	cash.EXPECT().CheckCash(ctx, newItem.Category.Name).Return(false)
	err := usecase.UpdateItemsInCategoryCash(ctx, newItem, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, newItem.Category.Name).Return(true)
	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name).Return(nil, fmt.Errorf("error"))
	err = usecase.UpdateItemsInCategoryCash(ctx, newItem, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, newItem.Category.Name).Return(true)
	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name).Return(cashResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(updateResults), newItem.Category.Name+"Quantity").Return(fmt.Errorf("error"))
	err = usecase.UpdateItemsInCategoryCash(ctx, newItem, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, newItem.Category.Name).Return(true)
	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name).Return(cashResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(updateResults), newItem.Category.Name+"Quantity").Return(nil)
	cash.EXPECT().CreateItemsCash(ctx, updateResults, newItem.Category.Name).Return(fmt.Errorf("error"))
	err = usecase.UpdateItemsInCategoryCash(ctx, newItem, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, newItem.Category.Name).Return(true)
	cash.EXPECT().GetItemsCash(ctx, newItem.Category.Name).Return(cashResults, nil)
	cash.EXPECT().CreateItemsCash(ctx, cashResults, newItem.Category.Name).Return(nil)
	err = usecase.UpdateItemsInCategoryCash(ctx, newItem, "update")
	require.NoError(t, err)
}
