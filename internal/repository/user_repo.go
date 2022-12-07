package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type user struct {
	storage Storage
	logger  *zap.SugaredLogger
}

var _ UserStore = (*user)(nil)

func NewUser(storage Storage, logger *zap.SugaredLogger) UserStore {
	return &user{
		storage: storage,
		logger:  logger,
	}
}

func (u *user) Create(ctx context.Context, user *models.User) (*models.User, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		// we create rights and address somewhere in usecase or get them from user
		row := pool.QueryRow(ctx, `INSERT INTO users 
		(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
			user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
			user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)
		var id uuid.UUID
		err := row.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("can't create user: %w", err)
		}
		user.ID = id
		return user, nil
	}
}

func (u *user) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		row := pool.QueryRow(ctx, `SELECT users.id, users.name, lastname, password, email, rights.id, zipcode, country, city, street,
		rights.name, rights.rules
		FROM users INNER JOIN rights ON email=$1 and rights.id=users.rights`, email)
		var user = models.User{}
		err := row.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Password, &user.Email, &user.Rights.ID,
			&user.Address.Zipcode, &user.Address.Country, &user.Address.City, &user.Address.Street, &user.Rights.Name, &user.Rights.Rules)
		if err != nil {
			return models.User{}, fmt.Errorf("can't get user from database: %w", err)
		}
		return user, err
	}
}
