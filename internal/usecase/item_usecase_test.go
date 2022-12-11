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

func TestCreateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()
	testCategoryId, _ := uuid.Parse("feb77bbc-1b8a-4739-bd68-d3b052af9a80")
	testCategory := models.Category{
		Id:          testCategoryId,
		Name:        "TestCategoryName",
		Description: "TestCategoryDescription",
	}
	testModelItem := &models.Item{
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testCategory,
		Price:       1,
		Vendor:      "TestVendor",
	}
	expect, _ := uuid.Parse("13574b3b-0c44-4864-89de-a086ad68ec4b")
	itemRepo.EXPECT().CreateItem(ctx, testModelItem).Return(expect, nil)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	res, err := usecase.CreateItem(ctx, testModelItem)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)

	err = fmt.Errorf("error on create item")
	itemRepo.EXPECT().CreateItem(ctx, testModelItem).Return(uuid.Nil, err)
	res, err = usecase.CreateItem(ctx, testModelItem)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)
}

func TestUpdateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	itemId, _ := uuid.Parse("feb77bbc-1b8a-4739-bd68-d3b052af9a80")
	testCategoryId, _ := uuid.Parse("b02c1542-dba1-46d2-ac71-e770c13d0d50")
	testCategory := models.Category{
		Id:          testCategoryId,
		Name:        "TestCategoryName",
		Description: "TestCategoryDescription",
	}
	testModelItem := &models.Item{
		Id:          itemId,
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testCategory,
		Price:       1,
		Vendor:      "TestVendor",
	}
	itemRepo.EXPECT().UpdateItem(ctx, testModelItem).Return(nil)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	err := usecase.UpdateItem(ctx, testModelItem)
	require.NoError(t, err)

	err = fmt.Errorf("error on update item")
	itemRepo.EXPECT().UpdateItem(ctx, testModelItem).Return(err)
	err = usecase.UpdateItem(ctx, testModelItem)
	require.Error(t, err)
}

func TestGetItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()
	id := "feb77bbc-1b8a-4739-bd68-d3b052af9a80"
	uid, _ := uuid.Parse(id)
	testCategory := models.Category{
		Id:          uuid.New(),
		Name:        "TestCategoryName",
		Description: "TestCategoryDescription",
	}
	testModelItem := &models.Item{
		Id:          uid,
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testCategory,
		Price:       1,
		Vendor:      "TestVendor",
	}
	itemRepo.EXPECT().GetItem(ctx, uid).Return(testModelItem, nil)
	res, err := usecase.GetItem(ctx, uid)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testModelItem)

	err = fmt.Errorf("error on get item")
	itemRepo.EXPECT().GetItem(ctx, uid).Return(&models.Item{}, err)
	res, err = usecase.GetItem(ctx, uid)
	require.Error(t, err)
	require.Equal(t, res, &models.Item{})
}

func TestItemsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()
	testCategory := models.Category{
		Id:          uuid.New(),
		Name:        "TestCategoryName",
		Description: "TestCategoryDescription",
	}
	id := "feb77bbc-1b8a-4739-bd68-d3b052af9a80"
	uid, _ := uuid.Parse(id)
	testModelItem := models.Item{
		Id:          uid,
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testCategory,
		Price:       1,
		Vendor:      "TestVendor",
	}
	testChan := make(chan models.Item, 1)
	testChan <- testModelItem
	close(testChan)
	expect := make([]models.Item, 0, 100)
	expect = append(expect, testModelItem)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan, nil)
	cash.EXPECT().CreateItemsListCash(ctx, expect, itemsListKey).Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(expect), itemsQuantityKey).Return(nil)
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey).Return(expect, nil)
	res, err := usecase.ItemsList(ctx, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey).Return(expect, nil)
	res, err = usecase.ItemsList(ctx, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)

	err = fmt.Errorf("error on itemslist()")
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan, err)
	res, err = usecase.ItemsList(ctx, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testModelItem
	close(testChan2)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan2, nil)
	cash.EXPECT().CreateItemsListCash(ctx, expect, itemsListKey).Return(err)
	res, err = usecase.ItemsList(ctx, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testModelItem
	close(testChan3)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan3, nil)
	cash.EXPECT().CreateItemsListCash(ctx, expect, itemsListKey).Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(expect), itemsQuantityKey).Return(nil)
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey).Return(nil, err)
	res, err = usecase.ItemsList(ctx, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)

	testChan4 := make(chan models.Item, 1)
	testChan4 <- testModelItem
	close(testChan4)
	err = fmt.Errorf("error on create items quantity cash")
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan4, nil)
	cash.EXPECT().CreateItemsListCash(ctx, expect, itemsListKey).Return(nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(expect), itemsQuantityKey).Return(err)
	res, err = usecase.ItemsList(ctx, 0, 1)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestSearchLine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()
	testCategory := models.Category{
		Id:          uuid.New(),
		Name:        "TestCategoryName",
		Description: "TestCategoryDescription",
	}
	id := "feb77bbc-1b8a-4739-bd68-d3b052af9a80"
	uid, _ := uuid.Parse(id)
	testModelItem := models.Item{
		Id:          uid,
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testCategory,
		Price:       1,
		Vendor:      "TestVendor",
	}
	param := "est"
	testChan := make(chan models.Item, 1)
	testChan <- testModelItem
	close(testChan)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan, nil)
	res, err := usecase.SearchLine(ctx, param)
	require.NoError(t, err)
	require.NotNil(t, res)

	err = fmt.Errorf("error on search line()")
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan, err)
	res, err = usecase.SearchLine(ctx, param)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestItemsQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
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
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey).Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, itemsQuantityKey).Return(nil)
	cash.EXPECT().GetItemsQuantityCash(ctx, itemsQuantityKey).Return(1, nil)
	res, err = usecase.ItemsQuantity(ctx)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey).Return(nil, fmt.Errorf("error on get items list cash"))
	res, err = usecase.ItemsQuantity(ctx)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey).Return([]models.Item{{Title: "testTitle"}}, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, 1, itemsQuantityKey).Return(fmt.Errorf("error on create items quantity cash"))
	res, err = usecase.ItemsQuantity(ctx)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cash.EXPECT().CheckCash(ctx, itemsQuantityKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey)
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
}

func TestUpdateCash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewItemUsecase(itemRepo, cash, logger)
	ctx := context.Background()

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(false)
	err := usecase.UpdateCash(ctx, uuid.New(), "create")
	require.Error(t, err)

	item := &models.Item{}
	id := uuid.New()
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, id).Return(item, fmt.Errorf("error on get item"))
	err = usecase.UpdateCash(ctx, id, "create")
	require.Error(t, err)

	testCategory := models.Category{
		Id:          uuid.New(),
		Name:        "TestName",
		Description: "TestDescription",
	}

	newItem := &models.Item{
		Id:          id,
		Title:       "test",
		Description: "test",
		Category:    testCategory,
		Price:       0,
		Vendor:      "test",
	}
	cashItem := models.Item{
		Id:          id,
		Title:       "test",
		Description: "test",
		Category:    testCategory,
		Price:       0,
		Vendor:      "test",
	}
	cashResults := make([]models.Item, 0, 1)
	cashResults = append(cashResults, cashItem)
	updateResults := make([]models.Item, 0, 1)
	updateResults = append(updateResults, *newItem)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, id).Return(newItem, nil)
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey).Return(nil, fmt.Errorf("error on get cash"))
	err = usecase.UpdateCash(ctx, id, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, id).Return(newItem, nil)
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey).Return(cashResults, nil)
	cash.EXPECT().CreateItemsListCash(ctx, updateResults, itemsListKey).Return(fmt.Errorf("error on create cash"))
	err = usecase.UpdateCash(ctx, id, "update")
	require.Error(t, err)

	updateResults = append(cashResults, *newItem)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, id).Return(newItem, nil)
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey).Return(cashResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(updateResults), itemsQuantityKey).Return(nil)
	cash.EXPECT().CreateItemsListCash(ctx, updateResults, itemsListKey).Return(nil)
	err = usecase.UpdateCash(ctx, id, "create")
	require.NoError(t, err)

	updateResults = append(cashResults, *newItem)
	cash.EXPECT().CheckCash(ctx, itemsListKey).Return(true)
	itemRepo.EXPECT().GetItem(ctx, id).Return(newItem, nil)
	cash.EXPECT().GetItemsListCash(ctx, itemsListKey).Return(cashResults, nil)
	cash.EXPECT().CreateItemsQuantityCash(ctx, len(updateResults), itemsQuantityKey).Return(fmt.Errorf("error"))
	err = usecase.UpdateCash(ctx, id, "create")
	require.Error(t, err)
}
