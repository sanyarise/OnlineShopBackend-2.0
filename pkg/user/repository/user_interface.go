package repository

import (
	"context"
	"online_shop_backend/pkg/models"

	"github.com/google/uuid"
)

type User interface {
	Get(ctx context.Context, email string) (models.User, error)
	Create(ctx context.Context, user models.User) (uuid.UUID, error)
}
