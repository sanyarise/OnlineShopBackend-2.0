package handlers

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type User struct {
	Id        string `json:"id,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	Cellphone string `json:"cellphone,omitempty"`
	Zipcode   int64  `json:"zipcode,omitempty"`
	Country   string `json:"country,omitempty"`
	City      string `json:"city,omitempty"`
	Street    string `json:"street,omitempty"`
	Rights    Rights `json:"rights"`
}

type Rights struct {
	Id    string   `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Rules []string `json:"rules,omitempty"`
}

type UserHandlers struct {
	usecase usecase.IUserUsecase
	logger  *zap.Logger
}

func NewUserHandlers(usecase usecase.IUserUsecase, logger *zap.Logger) *UserHandlers {
	return &UserHandlers{usecase: usecase, logger: logger}
}

func (handlers *UserHandlers) CreateUser(ctx context.Context, user User) (uuid.UUID, error) {
	handlers.logger.Debug("Enter in handlers CreateItem()")
	rights, err := handlers.usecase.GetRightsId(ctx, models.Customer)
	newUser := &models.User{
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Rights: models.Rights{
			ID:    rights.ID,
			Name:  rights.Name,
			Rules: rights.Rules,
		},
	}

	id, err := handlers.usecase.CreateUser(ctx, newUser)
	if err != nil {
		return id, err
	}
	return id, nil

}
