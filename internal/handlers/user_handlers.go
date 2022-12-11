package handlers

import (
	"OnlineShopBackend/internal/models"
	"context"
	"github.com/google/uuid"
)

type User struct {
	Id          string `json:"id,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName 	string `json:"last_name,omitempty"`
	Email       string  `json:"email,omitempty"`
	Password    string  `json:"password,omitempty"`
	Cellphone   string `json:"cellphone,omitempty"`
	Zipcode 	int64  `json:"zipcode,omitempty"`
	Country 	string `json:"country,omitempty"`
	City   		string `json:"city,omitempty"`
	Street 	 	string `json:"street,omitempty"`
}

func (handlers *Handlers) CreateUser(ctx context.Context, user User) (uuid.UUID, error) {
	handlers.logger.Debug("Enter in handlers CreateItem()")
	newUser := &models.User{
		Firstname: user.FirstName,
		Lastname: user.LastName,
		Email: user.Email,
		Password: user.Password,
		Zipcode: user.Zipcode,
		Country: user.Country,
		City: user.City,
		Street: user.Street,
	}

	id, err := handlers.repo.CreateUser(ctx, newUser)
	if err != nil {
		return id, err
	}
	return id, nil

}