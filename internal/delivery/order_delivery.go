package delivery

import (
	"OnlineShopBackend/internal/delivery/order"
	"OnlineShopBackend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Create order - create an order out of cart and user
//
//	@Summary		Create order
//	@Description	The method allows you to create an order out of cart and user info
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Success		201		{object}	order.OrderId	"Order id"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/order/create/ [post]
func (d *Delivery) CreateOrder(c *gin.Context) {
	d.logger.Debug("Eneter in delivery CreateOrder")
	ctx := c.Request.Context()
	var cart order.CartAdressUser
	if err := c.ShouldBindJSON(&cart); err != nil {
		d.logger.Sugar().Errorf("can't bind json from request: %s", err)
		d.SetError(c, http.StatusBadRequest, err)
		return
	}
	id, err := uuid.Parse(cart.User.Id)
	if err != nil {
		d.logger.Sugar().Errorf("can't parse user id: %s", err)
		d.SetError(c, http.StatusInternalServerError, err)
		return
	}
	user := models.User{
		ID:    id,
		Email: cart.User.Email,
	}
	id, err = uuid.Parse(cart.Cart.Id)
	if err != nil {
		d.logger.Sugar().Errorf("can't parse cart id: %s", err)
		d.SetError(c, http.StatusInternalServerError, err)
		return
	}
	cartModel := models.Cart{
		Id:     id,
		UserId: user.ID,
		Items:  make([]models.ItemWithQuantity, 0, len(cart.Cart.Items)),
	}
	for _, item := range cart.Cart.Items {
		id, err = uuid.Parse(item.Id)
		if err != nil {
			d.logger.Sugar().Errorf("can't parse item id: %s", err)
			d.SetError(c, http.StatusInternalServerError, err)
			return
		}
		itemM := models.ItemWithQuantity{
			Item: models.Item{
				Id:    id,
				Title: item.Title,
				Price: item.Price,
			},
			Quantity: item.Quantity,
		}
		cartModel.Items = append(cartModel.Items, itemM)
	}

	addressMdl := models.UserAddress{
		Country: cart.Address.Country,
		City:    cart.Address.City,
		Zipcode: cart.Address.Zipcode,
		Street:  cart.Address.Street,
	}

	ordr, err := d.orderUsecase.PlaceOrder(ctx, &cartModel, user, addressMdl)
	if err != nil {
		d.logger.Sugar().Errorf("can't create order: %s", err)
		d.SetError(c, http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusCreated, order.OrderId{Value: ordr.ID.String()})
}
