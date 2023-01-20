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

//type JwtRealisation struct {
//	ttl int64
//}



type Credentials struct {
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
}

func (usecase *UserUsecase) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	usecase.logger.Debug("Enter in usecase CreateUser()")

	rights, err := usecase.userStore.GetRightsId(ctx, "Customer")
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
	var user *models.User

	user, err := usecase.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return &models.User{}, err
	}

	//user.Password = ""

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

func (usecase *UserUsecase) UpdateUserData(ctx context.Context, id uuid.UUID, user *models.User) (*models.User, error) {
	user, err := usecase.userStore.UpdateUserData(ctx, id, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
