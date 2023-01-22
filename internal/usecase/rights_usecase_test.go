package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository/mocks"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testRightsNoId = &models.Rights{
		Name: "Test",
	}
	testRightsWithId = &models.Rights{
		ID:   testRightsId,
		Name: "Test",
	}
	testRightsList = []models.Rights{*testRightsWithId}
	testRightsId   = uuid.New()
)

func TestCreateRights(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	rightsRepo := mocks.NewMockRightsStore(ctrl)
	usecase := NewRightsUsecase(rightsRepo, logger)
	ctx := context.Background()

	rightsRepo.EXPECT().CreateRights(ctx, testRightsNoId).Return(uuid.Nil, err)
	res, err := usecase.CreateRights(ctx, testRightsNoId)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)

	rightsRepo.EXPECT().CreateRights(ctx, testRightsNoId).Return(testRightsId, nil)
	res, err = usecase.CreateRights(ctx, testRightsNoId)
	require.NoError(t, err)
	require.Equal(t, res, testRightsId)
}

func TestUpdateRights(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	rightsRepo := mocks.NewMockRightsStore(ctrl)
	usecase := NewRightsUsecase(rightsRepo, logger)
	ctx := context.Background()

	rightsRepo.EXPECT().UpdateRights(ctx, testRightsWithId).Return(err)
	err := usecase.UpdateRights(ctx, testRightsWithId)
	require.Error(t, err)

	rightsRepo.EXPECT().UpdateRights(ctx, testRightsWithId).Return(nil)
	err = usecase.UpdateRights(ctx, testRightsWithId)
	require.NoError(t, err)
}

func TestDeleteRights(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	rightsRepo := mocks.NewMockRightsStore(ctrl)
	usecase := NewRightsUsecase(rightsRepo, logger)
	ctx := context.Background()

	rightsRepo.EXPECT().DeleteRights(ctx, testRightsId).Return(err)
	err := usecase.DeleteRights(ctx, testRightsId)
	require.Error(t, err)

	rightsRepo.EXPECT().DeleteRights(ctx, testRightsId).Return(nil)
	err = usecase.DeleteRights(ctx, testRightsId)
	require.NoError(t, err)
}

func TestGetRights(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	rightsRepo := mocks.NewMockRightsStore(ctrl)
	usecase := NewRightsUsecase(rightsRepo, logger)
	ctx := context.Background()

	rightsRepo.EXPECT().GetRights(ctx, testRightsId).Return(nil, err)
	res, err := usecase.GetRights(ctx, testRightsId)
	require.Error(t, err)
	require.Nil(t, res)

	rightsRepo.EXPECT().GetRights(ctx, testRightsId).Return(testRightsWithId, nil)
	res, err = usecase.GetRights(ctx, testRightsId)
	require.NoError(t, err)
	require.Equal(t, res, testRightsWithId)
}

func TestRightsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	rightsRepo := mocks.NewMockRightsStore(ctrl)
	usecase := NewRightsUsecase(rightsRepo, logger)
	ctx := context.Background()

	rightsRepo.EXPECT().RightsList(ctx).Return(nil, err)
	res, err := usecase.RightsList(ctx)
	require.Error(t, err)
	require.Nil(t, res)

	rightsRepo.EXPECT().RightsList(ctx).Return(testRightsList, nil)
	res, err = usecase.RightsList(ctx)
	require.NoError(t, err)
	require.Equal(t, res, testRightsList)
}

func TestRightsByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	rightsRepo := mocks.NewMockRightsStore(ctrl)
	usecase := NewRightsUsecase(rightsRepo, logger)
	ctx := context.Background()

	rightsRepo.EXPECT().GetRightsByName(ctx, "test").Return(nil, err)
	res, err := usecase.GetRightsByName(ctx, "test")
	require.Error(t, err)
	require.Nil(t, res)

	rightsRepo.EXPECT().GetRightsByName(ctx, "test").Return(testRightsWithId, nil)
	res, err = usecase.GetRightsByName(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, res, testRightsWithId)
}
