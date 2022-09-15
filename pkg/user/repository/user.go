package repository

import (
	"online_shop_backend/pkg/models"
	"online_shop_backend/pkg/storage"
)

type User interface {
	Get(email string) models.User
	Create(user models.User) error
}

type userRepo struct {
	storage *storage.Storage
}

func New(store *storage.Storage) User {
	return &userRepo{
		storage: store,
	}
}

func (ur *userRepo) Get(email string) models.User {
	return models.User{}
}

func (ur *userRepo) Create(user models.User) error {
	return nil
}
