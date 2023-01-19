package delivery

import (
	"OnlineShopBackend/internal/delivery/order"
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	testAddress = order.OrderAddress{
		Country: "Israel",
		City:    "Haifa",
		Zipcode: "3372707",
		Street:  "Daniel 4",
	}

	testUser = order.UserForCart{
		Id:    uuid.New().String(),
		Email: "test@mail.ru",
	}

	testCartUserAddress = order.CartAdressUser{
		Cart:    testCart,
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

	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}
