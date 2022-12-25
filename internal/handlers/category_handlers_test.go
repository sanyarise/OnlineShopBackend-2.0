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
	testCatId        = uuid.New()
	testCategoryNoId = Category{
		Name:        "test",
		Description: "test",
	}
	testModelCategoryNoId = models.Category{
		Name:        "test",
		Description: "test",
	}
	testModelCategoryWithId = models.Category{
		Id:          testCatId,
		Name:        "test",
		Description: "test",
		Image:       "test",
	}
	testCategoryWithId = Category{
		Id:          testCatId.String(),
		Name:        "test",
		Description: "test",
		Image:       "test",
	}
	testCategoryInvalidId = Category{
		Id: "invalid Id",
	}
	emptyCategory      = Category{}
	emptyModelCategory = models.Category{}
	emptyCategories    = make([]Category, 0, 100)
	categories         = []models.Category{testModelCategoryWithId}
	resCategories      = []Category{testCategoryWithId}
)

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockICategoryUsecase(ctrl)
	handlers := NewCategoryHandlers(usecase, logger)

	usecase.EXPECT().CreateCategory(ctx, &testModelCategoryNoId).Return(testCatId, nil)
	res, err := handlers.CreateCategory(ctx, testCategoryNoId)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testCatId)

	err = fmt.Errorf("error on create category")
	usecase.EXPECT().CreateCategory(ctx, &testModelCategoryNoId).Return(uuid.Nil, err)
	res, err = handlers.CreateCategory(ctx, testCategoryNoId)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)
}

func TestUpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockICategoryUsecase(ctrl)
	handlers := NewCategoryHandlers(usecase, logger)

	err := handlers.UpdateCategory(ctx, testCategoryInvalidId)
	require.Error(t, err)

	usecase.EXPECT().UpdateCategory(ctx, &testModelCategoryWithId).Return(nil)
	err = handlers.UpdateCategory(ctx, testCategoryWithId)
	require.NoError(t, err)

	usecase.EXPECT().UpdateCategory(ctx, &testModelCategoryWithId).Return(fmt.Errorf("error"))
	err = handlers.UpdateCategory(ctx, testCategoryWithId)
	require.Error(t, err)
}

func TestGetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockICategoryUsecase(ctrl)
	handlers := NewCategoryHandlers(usecase, logger)

	res, err := handlers.GetCategory(ctx, "invalidId")
	require.Error(t, err)
	require.Equal(t, res, emptyCategory)

	usecase.EXPECT().GetCategory(ctx, testCatId).Return(&testModelCategoryWithId, nil)
	res, err = handlers.GetCategory(ctx, testCatId.String())
	require.NoError(t, err)
	require.Equal(t, res, testCategoryWithId)

	usecase.EXPECT().GetCategory(ctx, testCatId).Return(&emptyModelCategory, fmt.Errorf("error"))
	res, err = handlers.GetCategory(ctx, testCatId.String())
	require.Error(t, err)
	require.Equal(t, res, emptyCategory)
}

func TestGetCategoryList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockICategoryUsecase(ctrl)
	handlers := NewCategoryHandlers(usecase, logger)

	usecase.EXPECT().GetCategoryList(ctx).Return(nil, fmt.Errorf("error"))
	res, err := handlers.GetCategoryList(ctx)
	require.Error(t, err)
	require.Equal(t, res, emptyCategories)

	usecase.EXPECT().GetCategoryList(ctx).Return(categories, nil)
	res, err = handlers.GetCategoryList(ctx)
	require.NoError(t, err)
	require.Equal(t, res, resCategories)
}

func TestDeleteCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockICategoryUsecase(ctrl)
	handlers := NewCategoryHandlers(usecase, logger)

	usecase.EXPECT().DeleteCategory(ctx, testCatId).Return(fmt.Errorf("error"))
	err := handlers.DeleteCategory(ctx, testCatId)
	require.Error(t, err)

	usecase.EXPECT().DeleteCategory(ctx, testCatId).Return(nil)
	err = handlers.DeleteCategory(ctx, testCatId)
	require.NoError(t, err)
}

func TestGetCategoryByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	usecase := mocks.NewMockICategoryUsecase(ctrl)
	handlers := NewCategoryHandlers(usecase, logger)

	usecase.EXPECT().GetCategoryByName(ctx, testModelCategoryWithId.Name).Return(nil, fmt.Errorf("error"))
	res, err := handlers.GetCategoryByName(ctx, testModelCategoryWithId.Name)
	require.Error(t, err)
	require.Equal(t, res, emptyCategory)

	usecase.EXPECT().GetCategoryByName(ctx, testModelCategoryWithId.Name).Return(&testModelCategoryWithId, nil)
	res, err = handlers.GetCategoryByName(ctx, testModelCategoryWithId.Name)
	require.NoError(t, err)
	require.Equal(t, res, testCategoryWithId)
}
