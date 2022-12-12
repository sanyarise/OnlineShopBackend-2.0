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
	testModelCategory = &models.Category{
		Name: "test name",
	}
	testModelCategoryWithId = &models.Category{
		Id: testId,
	}
	emptyCategory = &models.Category{}
	testId   = uuid.New()
	testChan = make(chan models.Category, 2)
)

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, logger)

	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(testId, nil)
	res, err := usecase.CreateCategory(ctx, testModelCategory)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testId)
	err = fmt.Errorf("error on create category")
	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(uuid.Nil, err)
	res, err = usecase.CreateCategory(ctx, testModelCategory)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)
}

func TestUpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, logger)
	categoryRepo.EXPECT().UpdateCategory(ctx, testModelCategoryWithId).Return(nil)
	err := usecase.UpdateCategory(ctx, testModelCategoryWithId)
	require.NoError(t, err)
	categoryRepo.EXPECT().UpdateCategory(ctx, testModelCategoryWithId).Return(fmt.Errorf("error on update"))
	err = usecase.UpdateCategory(ctx, testModelCategoryWithId)
	require.Error(t, err)
}

func TestGetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, logger)

	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	res, err := usecase.GetCategory(ctx, testId)
	require.NoError(t, err)
	require.Equal(t, res, testModelCategoryWithId)

	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(emptyCategory, fmt.Errorf("error on get category"))
	res, err = usecase.GetCategory(ctx, testId)
	require.Error(t, err)
	require.Equal(t, res, emptyCategory)
}

func TestGetCategoryList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, logger)
	testChan <- *testModelCategoryWithId
	close(testChan)
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan, nil)
	res, err := usecase.GetCategoryList(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)

	categoryRepo.EXPECT().GetCategoryList(ctx).Return(nil, fmt.Errorf("error on categories list query context"))
	res, err = usecase.GetCategoryList(ctx)
	require.Error(t, err)
	require.Nil(t, res)
}
