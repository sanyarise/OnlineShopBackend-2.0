//go:build integration

package repository_test

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	store  = repository.NewPgxStorage(&zap.SugaredLogger{})
	logger = zap.NewExample().Sugar()
)

func TestConnection(t *testing.T) {
	store, err := store.InitStorage(context.Background(), "postgresql://localhost:5432/shop?user=shopteam&password=123")
	assert.NoError(t, err)
	assert.NotNil(t, store)

}

func TestCategoryCreate(t *testing.T) {
	store, err := store.InitStorage(context.Background(), "postgresql://localhost:5432/shop?user=shopteam&password=123")
	assert.NoError(t, err)
	assert.NotNil(t, store)
	cat := repository.NewCategoryRepo(store, logger)
	id, err := cat.CreateCategory(context.Background(), &models.Category{
		Name:        "testCat",
		Description: "Description",
	})
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories WHERE id=$1`, id)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestGetCatList(t *testing.T) {
	store, err := store.InitStorage(context.Background(), "postgresql://localhost:5432/shop?user=shopteam&password=123")
	assert.NoError(t, err)
	assert.NotNil(t, store)
	cat := repository.NewCategoryRepo(store, logger)
	store.GetPool().Exec(context.Background(), `INSERT INTO categories (name, description) VALUES
	('1', '1des'), ('2', '2des'), ('3', '3des')`)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM categories`)
	res, err := cat.GetCategoryList(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, res)
	for c := range res {
		require.Contains(t, []string{"1", "2", "3"}, c.Name)
		require.Contains(t, []string{"1des", "2des", "3des"}, c.Description)
	}
}

func TestUserCreate(t *testing.T) {
	store, err := store.InitStorage(context.Background(), "postgresql://localhost:5432/shop?user=shopteam&password=123")
	assert.NoError(t, err)
	assert.NotNil(t, store)
	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ($1, $2) RETURNING id`,
		user.Rights.Name, user.Rights.Rules)
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)
	u := repository.NewUser(store, logger)
	res, err := u.Create(context.Background(), &user)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, res.ID)
}

func TestGetUSer(t *testing.T) {
	store, err := store.InitStorage(context.Background(), "postgresql://localhost:5432/shop?user=shopteam&password=123")
	assert.NoError(t, err)
	assert.NotNil(t, store)
	user := models.User{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Password:  "123",
		Email:     "123@mail.ru",
		Address: models.UserAddress{
			Zipcode: "123455",
			Country: "Russia",
			City:    "Moscow",
			Street:  "Polyanka, 10",
		},
		Rights: models.Rights{
			Name:  models.Admin,
			Rules: []string{"Do everything"},
		},
	}
	row := store.GetPool().QueryRow(context.Background(), `INSERT INTO rights (name, rules) VALUES ('admin', $1) RETURNING id`, []string{})
	err = row.Scan(&user.Rights.ID)
	defer store.GetPool().Exec(context.TODO(), `DELETE FROM rights`)
	assert.NoError(t, err)

	_, err = store.GetPool().Exec(context.Background(), `INSERT INTO users 
	(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
		user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)
	defer store.GetPool().Exec(context.Background(), `DELETE FROM users`)
	assert.NoError(t, err)
	u := repository.NewUser(store, logger)
	res, err := u.GetUserByEmail(context.Background(), "123@mail.ru")
	assert.NoError(t, err)
	require.Equal(t, user.Firstname, res.Firstname)
	require.Equal(t, user.Password, res.Password)
	require.Equal(t, user.Email, res.Email)
}
