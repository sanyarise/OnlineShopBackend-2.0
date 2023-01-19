package delivery

import (
	"OnlineShopBackend/internal/delivery/cart"
	"OnlineShopBackend/internal/delivery/order"
	fs "OnlineShopBackend/internal/filestorage/mocks"
	"OnlineShopBackend/internal/usecase/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testAddress = order.OrderAddress{
		Country: "Israel",
		City:    "Haifa",
		Zipcode: "40006",
		Street:  "Daniel 4",
	}

	testUser = order.UserForCart{
		Id:    uuid.New().String(),
		Email: "test@mail.ru",
		Role:  "user",
	}

	testOrderCart = cart.Cart{
		Id:     uuid.New().String(),
		UserId: testUser.Id,
		Items: []cart.CartItem{
			{
				Id:       uuid.NewString(),
				Title:    "test",
				Price:    300,
				Image:    "testurl",
				Quantity: 1,
			},
		},
	}

	testCartUserAddress = order.CartAdressUser{
		Cart:    testOrderCart,
		Address: testAddress,
		User:    testUser,
	}
)

func MockCartUserAddressJson(c *gin.Context, content interface{}, method string) {
	if method == "DELETE" {
		c.Request.Method = "DELETE"
	}
	if method == "PUT" {
		c.Request.Method = "PUT"
	}
	if method == "POST" {
		c.Request.Method = "POST"
	}

	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonbytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// ctx := context.Background()
	// logger := zap.L()
	logger := zap.NewExample()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	orderUsecase := &mocks.OrderUsecaseMock{}
	delivery := NewDelivery(itemUsecase, categoryUsecase, cartUsecase, logger, filestorage, orderUsecase)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockCartUserAddressJson(c, testCartUserAddress, "POST")
	delivery.CreateOrder(c)
	require.Equal(t, 201, w.Code)
}
