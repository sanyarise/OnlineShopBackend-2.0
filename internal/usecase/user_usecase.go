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
	logger    *zap.Logger
}

func NewUserUsecase(userStore repository.UserStore, logger *zap.Logger) IUserUsecase {
	return &UserUsecase{userStore: userStore, logger: logger}
}

/*type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Profile struct {
	Email     string  `json:"email,omitempty"`
	FirstName string  `json:"firstname,omitempty"`
	LastName  string  `json:"lastname,omitempty"`
	Address   Address `json:"address,omitempty"`
	Rights    Rights  `json:"rights,omitempty"`
}

type Address struct {
	Zipcode string `json:"zipcode,omitempty"`
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
	Street  string `json:"street,omitempty"`
}

type Rights struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Name  string    `json:"name,omitempty"`
	Rules []string  `json:"rules,omitempty"`
}*/

func (usecase *UserUsecase) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase CreateUser() with args: ctx, user: %v", user)

	rights, err := usecase.userStore.GetRightsId(ctx, "customer")
	if err != nil {
		return &models.User{}, err
	}

	usecaseUser := &models.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		Address: models.UserAddress{
			Zipcode: user.Address.Zipcode,
			Country: user.Address.Country,
			City:    user.Address.City,
			Street:  user.Address.Street,
		},
		Rights: models.Rights{
			ID:    rights.ID,
			Name:  rights.Name,
			Rules: rights.Rules,
		},
	}

	id, err := usecase.userStore.Create(ctx, usecaseUser)
	if err != nil {
		return &models.User{}, fmt.Errorf("error on create user: %w", err)
	}
	return id, nil
}

func (usecase *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetUserByEmail() with args: ctx, email: %s", email)

	var user *models.User
	user, err := usecase.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *UserUsecase) GetRightsId(ctx context.Context, name string) (*models.Rights, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase user GetRightsId() with args: ctx, name: %s", name)
	rights, err := usecase.userStore.GetRightsId(ctx, name)
	if err != nil {
		return nil, err
	}
	return &rights, nil
}

func (usecase *UserUsecase) UpdateUserData(ctx context.Context, user *models.User) (*models.User, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateUserData() with args: ctx, user: %v", user)
	user, err := usecase.userStore.UpdateUserData(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *UserUsecase) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetuserById() with args: ctx, id: %v", id)

	user, err := usecase.userStore.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *UserUsecase) GetUsersList(ctx context.Context) ([]models.User, error) {
	usecase.logger.Debug("Enter in usecase GetUsersList()")
	users, err := usecase.userStore.GetUsersList(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (usecase *UserUsecase) ChangeUserRole(ctx context.Context, userId uuid.UUID, rightsId uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase ChangeUserRole() with args: ctx, userId: %v, rightsId: %v", userId, rightsId)
	err := usecase.userStore.ChangeUserRole(ctx, userId, rightsId)
	if err != nil {
		return err
	}
	return nil
}

func (usecase *UserUsecase) ChangeUserPassword(ctx context.Context, userId uuid.UUID, newPassword string) error {
	usecase.logger.Sugar().Debugf("Enter in usecase ChangeUserPassword() with args: ctx, userId: %v, newPassword: %s", userId, newPassword)
	err := usecase.userStore.ChangeUserPassword(ctx, userId, newPassword)
	if err != nil {
		return err
	}
	return nil
}

func (usecase *UserUsecase) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteUser() with args: ctx, userId: %v", userId)
	err := usecase.userStore.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}