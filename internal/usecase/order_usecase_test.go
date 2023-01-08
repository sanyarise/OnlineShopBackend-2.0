package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testUser = models.User{
		Firstname: "TestName",
		Lastname:  "TestLastName",
		Password:  "123",
		Email:     "testeamil@123.ru",
		Address: models.UserAddress{
			Zipcode: "123",
			Country: "Israel",
			City:    "Haifa",
			Street:  "דניל 4",
		},
		Rights: models.Rights{
			Name:  "pacan",
			Rules: []string{"buy"},
		},
	}

	testCategory = models.Category{
		Name:        "Electronics",
		Description: "Electric stuff",
		Image:       "image.url",
	}

	testItem1 = models.Item{
		Title:       "testItem1",
		Description: "Awesome chinese item",
		Price:       300,
		Category:    testCategory,
		Vendor:      "chinese factory",
		Images:      []string{},
	}
	testItem2 = models.Item{
		Title:       "testItem2",
		Description: "Awesome chinese item",
		Price:       500,
		Category:    testCategory,
		Vendor:      "russian factory",
		Images:      []string{},
	}

	testOrder = models.Order{
		ShipmentTime: time.Now().Add(models.StandardShipmentPeriod),
		User:         testUser,
		Address:      testUser.Address,
		Status:       models.StatusCreated,
		Items: []models.Item{
			testItem1, testItem2,
		},
	}

	lgr = zap.NewExample().Sugar()
)

type orderRepoMock struct {
	err error
}

var _ repository.OrderStore = (*orderRepoMock)(nil)

func (orMock *orderRepoMock) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	order.ID, _ = uuid.NewRandom()
	return order, orMock.err
}
func (orMock *orderRepoMock) DeleteOrder(ctx context.Context, order *models.Order) error {
	return orMock.err
}
func (orMock *orderRepoMock) ChangeAddress(ctx context.Context, order *models.Order, address models.UserAddress) error {
	order.Address = address
	return orMock.err
}
func (orMock *orderRepoMock) ChangeStatus(ctx context.Context, order *models.Order, status models.Status) error {
	order.Status = status
	return orMock.err
}

func (orMock *orderRepoMock) GetOrderByID(ctx context.Context, id uuid.UUID) (models.Order, error) {
	orderID, _ := uuid.NewRandom()
	userID, _ := uuid.NewRandom()
	itemID1, _ := uuid.NewRandom()
	itemID2, _ := uuid.NewRandom()
	order := testOrder
	order.ID = orderID
	order.User.ID = userID
	order.Items[0].Id = itemID1
	order.Items[1].Id = itemID2
	return order, orMock.err
}

func (orMock *orderRepoMock) GetOrdersForUser(ctx context.Context, user *models.User) (chan models.Order, error) {
	res := make(chan models.Order, 1)
	return res, orMock.err
}

func TestPlaceOrder(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{}, lgr)
	cartID, _ := uuid.NewRandom()
	userID, _ := uuid.NewRandom()
	cart := models.Cart{
		ID:     cartID,
		UserID: userID,
		Items: []models.Item{
			testItem1, testItem2,
		},
		ExpireAt: time.Now().Add(2 * time.Hour),
	}
	res, err := uscs.PlaceOrder(context.Background(), &cart, &testUser)
	require.NoError(t, err)
	assert.Equal(t, testUser.Address, res.Address)
	assert.Equal(t, cart.Items, res.Items)
}

func TestPlaceOrderDBError(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{err: fmt.Errorf("test error")}, lgr)
	cartID, _ := uuid.NewRandom()
	userID, _ := uuid.NewRandom()
	cart := models.Cart{
		ID:     cartID,
		UserID: userID,
		Items: []models.Item{
			testItem1, testItem2,
		},
		ExpireAt: time.Now().Add(2 * time.Hour),
	}
	res, err := uscs.PlaceOrder(context.Background(), &cart, &testUser)
	require.Error(t, err)
	assert.Nil(t, res)
}

func TestChangeStatus(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{}, lgr)
	err := uscs.ChangeStatus(context.Background(), &testOrder, models.StatusProcessed)
	defer func() {
		testOrder.Status = models.StatusCreated
	}()
	require.NoError(t, err)
	assert.Equal(t, models.StatusProcessed, testOrder.Status)

}

func TestChangeStatusError(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{err: fmt.Errorf("test error")}, lgr)
	err := uscs.ChangeStatus(context.Background(), &testOrder, models.StatusProcessed)
	defer func() {
		testOrder.Status = models.StatusCreated
	}()
	require.Error(t, err)
}

func TestChangeAddress(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{}, lgr)
	err := uscs.ChangeAddress(context.Background(), &testOrder, models.UserAddress{
		Street:  "הלל 49",
		City:    "חיפה",
		Zipcode: "313455",
		Country: "Israel",
	})
	defer func() {
		testOrder.Status = models.StatusCreated
	}()
	require.NoError(t, err)
}

func TestChangeAddressError(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{err: fmt.Errorf("test error")}, lgr)
	err := uscs.ChangeAddress(context.Background(), &testOrder, models.UserAddress{
		Street:  "הלל 49",
		City:    "חיפה",
		Zipcode: "313455",
		Country: "Israel",
	})
	defer func() {
		testOrder.Status = models.StatusCreated
	}()
	require.Error(t, err)
}
