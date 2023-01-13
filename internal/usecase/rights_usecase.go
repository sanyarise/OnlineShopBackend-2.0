package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ IRightsUsecase = usecaseRights{}

type usecaseRights struct {
	rightsRepo repository.RightsStore
	logger     *zap.Logger
}

func NewUsecaseRights(rigthsRepo repository.RightsStore, logger *zap.Logger) IRightsUsecase {
	return &usecaseRights{
		rightsRepo: rigthsRepo,
		logger:     logger,
	}
}

func (usecase usecaseRights) CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase CreateRights() with args: ctx, rights: %v", rights)

	id, err := usecase.rightsRepo.CreateRights(ctx, rights)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (usecase usecaseRights) UpdateRights(ctx context.Context, rights *models.Rights) error {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateRights() with args: ctx, rights: %v", rights)

	err := usecase.rightsRepo.UpdateRights(ctx, rights)
	if err != nil {
		return err
	}
	return nil
}

func (usecase usecaseRights) DeleteRights(ctx context.Context, id uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteRights() with args: ctx, id: %v", id)

	err := usecase.rightsRepo.DeleteRights(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (usecase usecaseRights) GetRights(ctx context.Context, id uuid.UUID) (*models.Rights, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetRights() with args: ctx, id: %v", id)

	rights, err := usecase.rightsRepo.GetRights(ctx, id)
	if err != nil {
		return nil, err
	}
	return rights, nil
}

func (usecase usecaseRights) RightsList(ctx context.Context) ([]models.Rights, error) {
	usecase.logger.Debug("Enter in usecase RightsList() with args: ctx")

	rightsList, err := usecase.rightsRepo.RightsList(ctx)
	if err != nil {
		return nil, err
	}
	return rightsList, nil
}
