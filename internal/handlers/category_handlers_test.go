package handlers

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/mocks"
	"context"
	"fmt"
	"testing"
	"time"

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
		Id: testCatId,
		Name: "test",
		Description: "test",
		Image: "test",
	}
	testCategoryWithId = Category{
		Id: testCatId.String(),
		Name: "test",
		Description: "test",
		Image: "test",
	}
	testCategoryInvalidId = Category{
		Id: "invalid Id",
	}
	emptyCategory = Category{}
	emptyModelCategory = models.Category{}
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


	testChan := make(chan models.Category, 2)
	testChan <- testModelCategoryWithId
	close(testChan)
	testChan2 := make(chan models.Category, 1)

	expect := make([]Category, 0, 100)
	expect = append(expect, testCategoryWithId)
	usecase.EXPECT().GetCategoryList(ctx).Return(testChan, nil)
	res, err := handlers.GetCategoryList(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)

	usecase.EXPECT().GetCategoryList(ctx).Return(testChan2, fmt.Errorf("error on categories list query context"))
	res, err = handlers.GetCategoryList(ctx)
	expect2 := make([]Category, 0, 100)
	require.Error(t, err)
	require.Equal(t, res, expect2)
	ctx, cancel := context.WithDeadline(context.Background(), <-time.After(1*time.Microsecond))
	expect = make([]Category, 0, 100)
	usecase.EXPECT().GetCategoryList(ctx).Return(testChan2, nil)
	res, err = handlers.GetCategoryList(ctx)
	require.Error(t, err)
	require.Equal(t, res, expect)
	cancel()
}
