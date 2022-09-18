package repository

import (
	"context"
	"fmt"
	"online_shop_backend/pkg/models"
	"online_shop_backend/pkg/storage"

	"github.com/google/uuid"
)

type userRepo struct {
	storage *storage.Storage
}

func New(store *storage.Storage) User {
	return &userRepo{
		storage: store,
	}
}

func (ur *userRepo) Get(ctx context.Context, email string) (models.User, error) {
	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("done with context")
	default:
		pool := ur.storage.Pool
		row := pool.QueryRow(ctx, `SELECT FROM users WHERE email=$1`, email)
		res := models.User{}
		err := row.Scan(&res)
		if err != nil {
			return models.User{}, fmt.Errorf("ca't get data from db: %w", err)
		}
		return res, nil
	}

}

func (ur *userRepo) Create(ctx context.Context, user models.User) (uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return uuid.Nil, fmt.Errorf("context done")
	default:
		pool := ur.storage.Pool
		row := pool.QueryRow(ctx,
			`INSERT INTO users (name, email, password, address, rights)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			user.Name, user.Email, user.Password, user.Address, user.Rights.ID)
		// TODO: think about rights
		var id uuid.UUID
		err := row.Scan(&id)
		if err != nil {
			return uuid.Nil, fmt.Errorf("can't read data from db %w", err)
		}
		return id, nil
	}
}
