package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ IUserUsecase = &UserUsecase{}

type UserUsecase struct {
	userStore repository.UserStore
	logger *zap.Logger
}

func NewUserUsecase(userStore repository.UserStore, logger *zap.Logger) IUserUsecase {
	return &UserUsecase{userStore: userStore, logger: logger}
}

func (usecase *UserUsecase) CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error) {
	usecase.logger.Debug("Enter in usecase CreateUser()")

	id, err := usecase.userStore.Create(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create user: %w", err)
	}
	return id.ID, nil
}

func (usecase *UserUsecase) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	return user, nil
}

func (usecase *UserUsecase) GetRightsId(ctx context.Context, name string) (*models.Rights, error) {
	//var rights models.Rights

	rights, err := usecase.userStore.GetRightsId(ctx, name)
	if err != nil {
		return nil, err
	}

	return &rights, nil
}




