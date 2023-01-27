package user

import (
	"OnlineShopBackend/internal/delivery/user/jwtauth"
	"OnlineShopBackend/internal/models"
	"github.com/google/uuid"
)

type LoginResponseData struct {
	CartId uuid.UUID `json:"cartId"`
	Token jwtauth.Token `json:"token"`

}

type CreateUserData struct {
	ID        uuid.UUID          `json:"id"`
	Firstname string             `json:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty"`
	Password  string             `json:"password,omitempty"`
	Email     string             `json:"email,omitempty"`
	Address   models.UserAddress `json:"address,omitempty"`
	Rights    models.Rights      `json:"rights"`
}