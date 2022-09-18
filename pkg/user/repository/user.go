package repository

import (
	"context"
	"fmt"
	"online_shop_backend/pkg/models"
	"online_shop_backend/pkg/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
		row := pool.QueryRow(ctx, `SELECT id, name, password FROM users WHERE email=$1`, email)
		res := models.User{}
		var password pgtype.Text
		err := row.Scan(&res.ID, &res.Name, &password)
		if password.Valid {
			res.Password = password.String
		}
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
			`INSERT INTO users (name, email, password, rights)
			VALUES ($1, $2, $3, (SELECT id FROM rights WHERE name='user')) RETURNING id`,
			user.Name, user.Email, user.Password)
		// TODO: think about rights
		var id uuid.UUID
		err := row.Scan(&id)
		if err != nil {
			return uuid.Nil, fmt.Errorf("can't read data from db %w", err)
		}
		return id, nil
	}
}
