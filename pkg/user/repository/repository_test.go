package repository_test

import (
	"context"
	"online_shop_backend/pkg/config"
	"online_shop_backend/pkg/models"
	"online_shop_backend/pkg/storage"
	"online_shop_backend/pkg/user/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var conf = config.Config{
	DSN: "postgres://shopteam:123@localhost:5432/shop?sslmode=disable",
}

func TestAddUser(t *testing.T) {
	strg := storage.New()
	err := strg.Start(context.Background(), conf)
	require.NoError(t, err)
	rep := repository.New(strg)
	id, err := rep.Create(context.Background(), models.User{
		Name:     "testName",
		Password: "testPassword",
		Email:    "testEmail",
	})
	require.NoError(t, err)
	defer func() {
		strg.Pool.Exec(context.Background(), "delete from users")
		strg.ShutDown(context.TODO())
	}()
	assert.NotEqual(t, uuid.Nil, id)
}

func TestGetUser(t *testing.T) {
	strg := storage.New()
	err := strg.Start(context.Background(), conf)
	require.NoError(t, err)
	testModel := models.User{
		Name:     "testName",
		Password: "testPassword",
		Email:    "testEmail",
	}
	row := strg.Pool.QueryRow(context.Background(), `INSERT INTO users (name, email, rights)
	VALUES ($1, $2, (SELECT id FROM rights WHERE name='user')) RETURNING id`,
		testModel.Name, testModel.Email)
	var id uuid.UUID
	err = row.Scan(&id)
	require.NoError(t, err)
	defer func() {
		strg.Pool.Exec(context.Background(), "delete from users")
		strg.ShutDown(context.TODO())
	}()
	repo := repository.New(strg)
	res, err := repo.Get(context.Background(), testModel.Email)
	require.NoError(t, err)
	assert.Equal(t, testModel.Name, res.Name)
	assert.Equal(t, id, res.ID)

}
