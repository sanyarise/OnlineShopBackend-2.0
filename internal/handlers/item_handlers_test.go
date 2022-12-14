package handlers

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/mocks"
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testModelCategory = models.Category{
		Id:          uuid.New(),
		Name:        "TestName",
		Description: "TestDescription",
	}
	testHandlersCategory = Category{
		Id:          testModelCategory.Id.String(),
		Name:        "TestName",
		Description: "TestDescription",
	}
	testHandlersCategoryWithInvalidId = Category{
		Id:          "InvalidId",
		Name:        "TestName",
		Description: "TestDescription",
	}
	testHandlersItem = Item{
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testHandlersCategory,
		Price:       1,
		Vendor:      "TestVendor",
	}
	testHandlersItemWithInvalidCategoryId = Item{
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testHandlersCategoryWithInvalidId,
		Price:       1,
		Vendor:      "TestVendor",
	}
	testHandlersItemWithIdWithInvalidCategoryId = Item{
		Id:       testNewId.String(),
		Category: testHandlersCategoryWithInvalidId,
	}
	testHandlersItemWithInvalidID = Item{
		Id: "invalid id",
	}
	testModelItem = &models.Item{
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testModelCategory,
		Price:       1,
		Vendor:      "TestVendor",
	}
	testNewId              = uuid.New()
	testHandlersItemWithId = Item{
		Id:          testNewId.String(),
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testHandlersCategory,
		Price:       1,
		Vendor:      "TestVendor",
	}
	testModelItemWithId = &models.Item{
		Id:          testNewId,
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testModelCategory,
		Price:       1,
		Vendor:      "TestVendor",
	}
	testParam = "est"
	testSlice = []models.Item{*testModelItemWithId}
	testSlice2 = []Item{testHandlersItemWithId}
	emptyHandlersSlice = []Item{}
	emptyModelsSlice = []models.Item{}
)

func TestCreateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockIItemUsecase(ctrl)
	handlers := NewItemHandlers(usecase, logger)

	res, err := handlers.CreateItem(ctx, testHandlersItemWithInvalidCategoryId)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)

	usecase.EXPECT().CreateItem(ctx, testModelItem).Return(testNewId, nil)
	res, err = handlers.CreateItem(ctx, testHandlersItem)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testNewId)

	err = fmt.Errorf("error on create item")
	usecase.EXPECT().CreateItem(ctx, testModelItem).Return(uuid.Nil, err)
	res, err = handlers.CreateItem(ctx, testHandlersItem)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)
}

func TestUpdateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockIItemUsecase(ctrl)
	handlers := NewItemHandlers(usecase, logger)

	err := handlers.UpdateItem(ctx, testHandlersItemWithInvalidID)
	require.Error(t, err)

	err = handlers.UpdateItem(ctx, testHandlersItemWithIdWithInvalidCategoryId)
	require.Error(t, err)

	usecase.EXPECT().UpdateItem(ctx, testModelItemWithId).Return(nil)
	err = handlers.UpdateItem(ctx, testHandlersItemWithId)
	require.NoError(t, err)

	err = fmt.Errorf("error on update item")
	usecase.EXPECT().UpdateItem(ctx, testModelItemWithId).Return(err)
	err = handlers.UpdateItem(ctx, testHandlersItemWithId)
	require.Error(t, err)
}

func TestGetItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockIItemUsecase(ctrl)
	handlers := NewItemHandlers(usecase, logger)

	usecase.EXPECT().GetItem(ctx, testNewId).Return(testModelItemWithId, nil)
	res, err := handlers.GetItem(ctx, testNewId.String())
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testHandlersItemWithId)

	res, err = handlers.GetItem(ctx, "invalidId")
	require.Error(t, err)
	require.Equal(t, res, Item{})

	err = fmt.Errorf("error on get item")
	usecase.EXPECT().GetItem(ctx, testNewId).Return(&models.Item{}, err)
	res, err = handlers.GetItem(ctx, testNewId.String())
	require.Error(t, err)
	require.Equal(t, res, Item{})
}

func TestItemsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockIItemUsecase(ctrl)
	handlers := NewItemHandlers(usecase, logger)

	testSlice := make([]Item, 0, 1)
	testSlice = append(testSlice, testHandlersItemWithId)
	testModelSlice := make([]models.Item, 0, 1)
	testModelSlice = append(testModelSlice, *testModelItemWithId)

	usecase.EXPECT().ItemsList(ctx, 0, 1).Return(testModelSlice, nil)
	res, err := handlers.ItemsList(ctx, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testSlice)

	err = fmt.Errorf("error on itemslist()")
	testSlice2 := make([]Item, 0, 1)
	testModelSlice2 := make([]models.Item, 0, 1)
	usecase.EXPECT().ItemsList(ctx, 0, 1).Return(testModelSlice2, err)
	res, err = handlers.ItemsList(ctx, 0, 1)
	require.Error(t, err)
	require.Equal(t, res, testSlice2)
}

func TestItemsQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockIItemUsecase(ctrl)
	handlers := NewItemHandlers(usecase, logger)

	usecase.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	res, err := handlers.ItemsQuantity(ctx)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	usecase.EXPECT().ItemsQuantity(ctx).Return(-1, fmt.Errorf("error on get items quantity"))
	res, err = handlers.ItemsQuantity(ctx)
	require.Error(t, err)
	require.Equal(t, res, -1)
}

func TestSearchLine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockIItemUsecase(ctrl)
	handlers := NewItemHandlers(usecase, logger)

	usecase.EXPECT().SearchLine(ctx, testParam, 0, 1).Return(testSlice, nil)
	res, err := handlers.SearchLine(ctx, testParam, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testSlice2)

	err = fmt.Errorf("error on search line()")
	usecase.EXPECT().SearchLine(ctx, testParam, 0, 1).Return(emptyModelsSlice, err)
	res, err = handlers.SearchLine(ctx, testParam, 0, 1)
	require.Error(t, err)
	require.Equal(t, res, emptyHandlersSlice)
}

func TestGetItemsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockIItemUsecase(ctrl)
	handlers := NewItemHandlers(usecase, logger)

	usecase.EXPECT().GetItemsByCategory(ctx, testParam, 0, 1).Return(testSlice, nil)
	res, err := handlers.GetItemsByCategory(ctx, testParam, 0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testSlice2)

	err = fmt.Errorf("error")
	usecase.EXPECT().GetItemsByCategory(ctx, testParam, 0, 1).Return(emptyModelsSlice, err)
	res, err = handlers.GetItemsByCategory(ctx, testParam, 0, 1)
	require.Error(t, err)
	require.Equal(t, res, emptyHandlersSlice)
}