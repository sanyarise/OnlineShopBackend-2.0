package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"github.com/google/uuid"
)

func (repo *Pgrepo) Create(ctx context.Context, user models.User) (uuid.UUID, error) {
		var id uuid.UUID
		return id, nil
}

func (repo *Pgrepo) Get(ctx context.Context, email string) (models.User, error){
	var user models.User
	return user, nil
}