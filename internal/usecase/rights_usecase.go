package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ IRightsUsecase = rightsUsecase{}

type rightsUsecase struct {
	rightsRepo repository.RightsStore
	logger     *zap.Logger
}

func NewRightsUsecase(rigthsRepo repository.RightsStore, logger *zap.Logger) IRightsUsecase {
	return &rightsUsecase{
		rightsRepo: rigthsRepo,
		logger:     logger,
	}
}

func (usecase rightsUsecase) CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase CreateRights() with args: ctx, rights: %v", rights)

	id, err := usecase.rightsRepo.CreateRights(ctx, rights)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (usecase rightsUsecase) UpdateRights(ctx context.Context, rights *models.Rights) error {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateRights() with args: ctx, rights: %v", rights)

	err := usecase.rightsRepo.UpdateRights(ctx, rights)
	if err != nil {
		return err
	}
	return nil
}

func (usecase rightsUsecase) DeleteRights(ctx context.Context, id uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteRights() with args: ctx, id: %v", id)

	err := usecase.rightsRepo.DeleteRights(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (usecase rightsUsecase) GetRights(ctx context.Context, id uuid.UUID) (*models.Rights, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetRights() with args: ctx, id: %v", id)

	rights, err := usecase.rightsRepo.GetRights(ctx, id)
	if err != nil {
		return nil, err
	}
	return rights, nil
}

func (usecase rightsUsecase) RightsList(ctx context.Context) ([]models.Rights, error) {
	usecase.logger.Debug("Enter in usecase RightsList() with args: ctx")

	rightsList, err := usecase.rightsRepo.RightsList(ctx)
	if err != nil {
		return nil, err
	}
	return rightsList, nil
}

func (usecase rightsUsecase) GetRightsByName(ctx context.Context, name string) (*models.Rights, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetRightsByName() with args: ctx, name: %s", name)

	rights, err := usecase.rightsRepo.GetRightsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return rights, nil
}