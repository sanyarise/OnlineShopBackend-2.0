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
		Id:   testId,
		Name: "test name",
	}
	emptyCategory = &models.Category{}
	testId        = uuid.New()
	testChan      = make(chan models.Category, 2)

	categories = []models.Category{*testModelCategoryWithId}
)

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cash := mocks.NewMockICategoriesCash(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cash, logger)

	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(testId, nil)
	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(false)
	res, err := usecase.CreateCategory(ctx, testModelCategory)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testId)

	emptyCategories := make([]models.Category, 0)
	categories := []models.Category{*testModelCategoryWithId}
	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(testId, nil)
	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cash.EXPECT().GetCategoriesListCash(ctx, categoriesListKey).Return(emptyCategories, nil)
	cash.EXPECT().CreateCategoriesListCash(ctx, categories, categoriesListKey).Return(nil)
	res, err = usecase.CreateCategory(ctx, testModelCategory)
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
	cash := mocks.NewMockICategoriesCash(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cash, logger)

	categoryRepo.EXPECT().UpdateCategory(ctx, testModelCategoryWithId).Return(nil)
	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(false)
	err := usecase.UpdateCategory(ctx, testModelCategoryWithId)
	require.NoError(t, err)

	categories := []models.Category{*testModelCategoryWithId}
	categoryRepo.EXPECT().UpdateCategory(ctx, testModelCategoryWithId).Return(nil)
	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cash.EXPECT().GetCategoriesListCash(ctx, categoriesListKey).Return(categories, nil)
	cash.EXPECT().CreateCategoriesListCash(ctx, categories, categoriesListKey).Return(nil)
	err = usecase.UpdateCategory(ctx, testModelCategoryWithId)
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
	cash := mocks.NewMockICategoriesCash(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cash, logger)

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
	cash := mocks.NewMockICategoriesCash(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cash, logger)

	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(true)
	cash.EXPECT().GetCategoriesListCash(ctx, categoriesListKey).Return(categories, nil)
	res, err := usecase.GetCategoryList(ctx)
	require.NoError(t, err)
	require.Equal(t, res, categories)

	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(true)
	cash.EXPECT().GetCategoriesListCash(ctx, categoriesListKey).Return(nil, fmt.Errorf("error"))
	res, err = usecase.GetCategoryList(ctx)
	require.Error(t, err)
	require.Nil(t, res)

	testChan <- *testModelCategoryWithId
	close(testChan)
	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(false)
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan, nil)
	cash.EXPECT().CreateCategoriesListCash(ctx, categories, categoriesListKey).Return(nil)
	cash.EXPECT().GetCategoriesListCash(ctx, categoriesListKey).Return(categories, nil)
	res, err = usecase.GetCategoryList(ctx)
	require.NoError(t, err)
	require.Equal(t, res, categories)

	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(false)
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan, fmt.Errorf("error"))
	res, err = usecase.GetCategoryList(ctx)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Category, 1)
	testChan2 <- *testModelCategoryWithId
	close(testChan2)
	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(false)
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan2, nil)
	cash.EXPECT().CreateCategoriesListCash(ctx, categories, categoriesListKey).Return(fmt.Errorf("error"))
	res, err = usecase.GetCategoryList(ctx)
	require.Error(t, err)
	require.Nil(t, res)

	testChan3 := make(chan models.Category, 1)
	testChan3 <- *testModelCategoryWithId
	close(testChan3)
	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(false)
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan3, nil)
	cash.EXPECT().CreateCategoriesListCash(ctx, categories, categoriesListKey).Return(nil)
	cash.EXPECT().GetCategoriesListCash(ctx, categoriesListKey).Return(nil, fmt.Errorf("error"))
	res, err = usecase.GetCategoryList(ctx)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestUpdateCategoryCash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cash := mocks.NewMockICategoriesCash(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cash, logger)

	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(false)
	err := usecase.UpdateCash(ctx, testId, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(nil, fmt.Errorf("error"))
	err = usecase.UpdateCash(ctx, testId, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cash.EXPECT().GetCategoriesListCash(ctx, categoriesListKey).Return(nil, fmt.Errorf("error"))
	err = usecase.UpdateCash(ctx, testId, "create")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cash.EXPECT().GetCategoriesListCash(ctx, categoriesListKey).Return(categories, nil)
	cash.EXPECT().CreateCategoriesListCash(ctx, categories, categoriesListKey).Return(fmt.Errorf("error"))
	err = usecase.UpdateCash(ctx, testId, "update")
	require.Error(t, err)

	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cash.EXPECT().GetCategoriesListCash(ctx, categoriesListKey).Return(categories, nil)
	cash.EXPECT().CreateCategoriesListCash(ctx, categories, categoriesListKey).Return(nil)
	err = usecase.UpdateCash(ctx, testId, "update")
	require.NoError(t, err)
}

func TestDeleteCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cash := mocks.NewMockICategoriesCash(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cash, logger)

	categoryRepo.EXPECT().DeleteCategory(ctx, testId).Return(fmt.Errorf("error"))
	err := usecase.DeleteCategory(ctx, testId)
	require.Error(t, err)

	categoryRepo.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cash.EXPECT().GetCategoriesListCash(ctx, categoriesListKey).Return(categories, nil)
	cash.EXPECT().CreateCategoriesListCash(ctx, []models.Category{}, categoriesListKey).Return(nil)
	err = usecase.DeleteCategory(ctx, testId)
	require.NoError(t, err)

	categoryRepo.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	cash.EXPECT().CheckCash(ctx, categoriesListKey).Return(false)
	err = usecase.DeleteCategory(ctx, testId)
	require.NoError(t, err)
}

func TestGetCategoryByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cash := mocks.NewMockICategoriesCash(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cash, logger)

	categoryRepo.EXPECT().GetCategoryByName(ctx, testModelCategoryWithId.Name).Return(nil, fmt.Errorf("error"))
	res, err := usecase.GetCategoryByName(ctx, testModelCategoryWithId.Name)
	require.Error(t, err)
	require.Nil(t, res)

	categoryRepo.EXPECT().GetCategoryByName(ctx, testModelCategoryWithId.Name).Return(testModelCategoryWithId, nil)
	res, err = usecase.GetCategoryByName(ctx, testModelCategoryWithId.Name)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testModelCategoryWithId)
}

func TestDeleteCategoryCash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cash := mocks.NewMockICategoriesCash(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cash, logger)

	cash.EXPECT().DeleteCash(ctx, "testNamenameasc").Return(err)
	err := usecase.DeleteCategoryCash(ctx, "testName")
	require.Error(t, err)

	cash.EXPECT().DeleteCash(ctx, "testNamenameasc").Return(nil)
	cash.EXPECT().DeleteCash(ctx, "testNamenamedesc").Return(nil)
	cash.EXPECT().DeleteCash(ctx, "testNamepriceasc").Return(nil)
	cash.EXPECT().DeleteCash(ctx, "testNamepricedesc").Return(nil)
	cash.EXPECT().DeleteCash(ctx, "testNameQuantity").Return(err)
	err = usecase.DeleteCategoryCash(ctx, "testName")
	require.Error(t, err)

	cash.EXPECT().DeleteCash(ctx, "testNamenameasc").Return(nil)
	cash.EXPECT().DeleteCash(ctx, "testNamenamedesc").Return(nil)
	cash.EXPECT().DeleteCash(ctx, "testNamepriceasc").Return(nil)
	cash.EXPECT().DeleteCash(ctx, "testNamepricedesc").Return(nil)
	cash.EXPECT().DeleteCash(ctx, "testNameQuantity").Return(nil)
	err = usecase.DeleteCategoryCash(ctx, "testName")
	require.NoError(t, err)
}