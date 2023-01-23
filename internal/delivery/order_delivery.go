package delivery

import (
	"OnlineShopBackend/internal/delivery/cart"
	"OnlineShopBackend/internal/delivery/order"
	"OnlineShopBackend/internal/models"
	"net/http"
	"strings"

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
//	@Param			cartAddressUser	body		order.CartAdressUser	true	"Data for creating order"
//	@Success		201				{object}	order.OrderId			"Order id"
//	@Failure		400				{object}	ErrorResponse
//	@Failure		403				"Forbidden"
//	@Failure		404				{object}	ErrorResponse	"404 Not Found"
//	@Failure		500				{object}	ErrorResponse
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

// GetOrder - get a specific order by id
//
//	@Summary		Get order by id
//	@Description	The method allows you to get the order by id.
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			orderID	path		string		true	"Id of order"
//	@Success		200		{object}	order.Order	"Order structure"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/order/{orderID} [get]
func (d *Delivery) GetOrder(c *gin.Context) {
	d.logger.Sugar().Debug("Enter the delivery GetOrder()")
	ctx := c.Request.Context()
	orderId, err := uuid.Parse(c.Param(("orderID")))
	if err != nil {
		d.logger.Sugar().Errorf("can't parse order id: %s", err)
		d.SetError(c, http.StatusBadRequest, err)
		return
	}
	modelOrder, err := d.orderUsecase.GetOrder(ctx, orderId)
	if err != nil {
		d.logger.Sugar().Errorf("can't get order: %s", err)
		d.SetError(c, http.StatusInternalServerError, err)
		return
	}
	order := order.Order{
		Id:           modelOrder.ID.String(),
		UserId:       modelOrder.User.ID.String(),
		ShipmentTime: modelOrder.ShipmentTime,
		Address:      order.OrderAddress(modelOrder.Address),
		Items:        make([]cart.CartItem, 0, len(modelOrder.Items)),
	}
	for _, item := range modelOrder.Items {
		cartItem := cart.CartItem{
			Id:       item.Id.String(),
			Title:    item.Title,
			Price:    item.Price,
			Image:    firstNotEmpty(item.Images),
			Quantity: item.Quantity,
		}
		order.Items = append(order.Items, cartItem)
	}
	c.JSON(http.StatusOK, order)
}

func firstNotEmpty(arr []string) string {
	for _, item := range arr {
		if item != "" {
			return item
		}
	}
	return ""
}

// GetOrderForUser - get a specific order by UserId
//
//	@Summary		Get all orders by UserId
//	@Description	The method allows you to get all orders by UserId.
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string		true	"Id of the user"
//	@Success		200		{array}		order.Order	"List of orders"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/order/{userID} [get]
func (d *Delivery) GetOrdersForUser(c *gin.Context) {
	d.logger.Sugar().Debug("Enter the delivery GetOrdersForUser()")
	ctx := c.Request.Context()
	userId, err := uuid.Parse(c.Param(("userID")))
	if err != nil {
		d.logger.Sugar().Errorf("can't parse user id: %s", err)
		d.SetError(c, http.StatusBadRequest, err)
		return
	}
	modelOrders, err := d.orderUsecase.GetOrdersForUser(ctx, &models.User{ID: userId})
	if err != nil {
		d.logger.Sugar().Errorf("can't get order: %s", err)
		d.SetError(c, http.StatusInternalServerError, err)
		return
	}
	orders := make([]order.Order, 0, len(modelOrders))
	for _, modelOrder := range modelOrders {
		order := order.Order{
			Id:           modelOrder.ID.String(),
			UserId:       modelOrder.User.ID.String(),
			ShipmentTime: modelOrder.ShipmentTime,
			Address:      order.OrderAddress(modelOrder.Address),
			Items:        make([]cart.CartItem, 0, len(modelOrder.Items)),
		}
		for _, item := range modelOrder.Items {
			cartItem := cart.CartItem{
				Id:       item.Id.String(),
				Title:    item.Title,
				Price:    item.Price,
				Image:    firstNotEmpty(item.Images),
				Quantity: item.Quantity,
			}
			order.Items = append(order.Items, cartItem)
		}
		orders = append(orders, order)
	}
	c.JSON(http.StatusOK, orders)
}

// DeleteOrder - delete a specific order by id
//
//	@Summary		Delete an order by id
//	@Description	The method allows you to delete an order by id.
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			orderID	path	string	true	"Id of the order to delete"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/order/delete/{orderID} [delete]
func (d *Delivery) DeleteOrder(c *gin.Context) {
	d.logger.Debug("Enter in delivery DeleteOrder()")

	ctx := c.Request.Context()

	orderId, err := uuid.Parse(c.Param("orderID"))
	if err != nil {
		d.logger.Sugar().Errorf("Can't parse orderID %s", err)
		d.SetError(c, http.StatusBadRequest, err)
		return
	}

	err = d.orderUsecase.DeleteOrder(ctx, &models.Order{ID: orderId})
	if err != nil {
		d.logger.Sugar().Errorf("Can't delete order with orderID %s", err)
		d.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ChangeAddressOfTheOrder - change address of a specific order by Id
//
//	@Summary		Change address of a  specific order by Id
//	@Description	The method allows you to change address of an order by Id.
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			AddressWithUserAndId	body	order.AddressWithUserAndId	true	"New address with orderID and user structure"
//	@Success		200
//	@Failure		400						{object}	ErrorResponse
//	@Failure		403						"Forbidden"
//	@Failure		404						{object}	ErrorResponse	"404 Not Found"
//	@Failure		500						{object}	ErrorResponse
//	@Router			/order/changeaddress/ [patch]
func (d *Delivery) ChangeAddress(c *gin.Context) {
	d.logger.Sugar().Debug("Enter the delivery ChangeAddress()")
	ctx := c.Request.Context()
	var address order.AddressWithUserAndId
	if err := c.ShouldBindJSON(&address); err != nil {
		d.logger.Sugar().Errorf("can't bind json from request: %s", err)
		d.SetError(c, http.StatusBadRequest, err)
		return
	}
	orderID, err := uuid.Parse(address.OrderId)
	if err != nil {
		d.logger.Sugar().Errorf("can't parse order id: %s", err)
		d.SetError(c, http.StatusBadRequest, err)
		return
	}
	if strings.ToLower(address.User.Role) == "user" {
		d.logger.Sugar().Errorf("the action not allowed: %s", err)
		d.SetError(c, http.StatusForbidden, err)
		return
	}
	userID, err := uuid.Parse(address.User.Id)
	if err != nil {
		d.logger.Sugar().Errorf("can't parse order id: %s", err)
		d.SetError(c, http.StatusBadRequest, err)
		return
	}
	err = d.orderUsecase.ChangeAddress(ctx,
		&models.Order{
			ID:   orderID,
			User: models.User{ID: userID},
		}, models.UserAddress(address.Address))
	if err != nil {
		d.logger.Sugar().Errorf("can't change address for order with id: %s %s", orderID, err)
		d.SetError(c, http.StatusInternalServerError, err)
		return
	}
}

// ChangeStatusOfTheOrder - change status of a specific order by Id
//
//	@Summary		Change status of a specific order by Id
//	@Description	The method allows you to change status of an order by Id.
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			statusWithUserAndId	body	order.StatusWithUserAndId	true	"New status with orderID and User structure"
//	@Success		200
//	@Failure		400					{object}	ErrorResponse
//	@Failure		403					"Forbidden"
//	@Failure		404					{object}	ErrorResponse	"404 Not Found"
//	@Failure		500					{object}	ErrorResponse
//	@Router			/order/changestatus/ [patch]
func (d *Delivery) ChangeStatus(c *gin.Context) {
	d.logger.Sugar().Debug("Enter the delivery ChangeStatus()")
	ctx := c.Request.Context()
	var status order.StatusWithUserAndId
	if err := c.ShouldBindJSON(&status); err != nil {
		d.logger.Sugar().Errorf("can't bind json from request: %s", err)
		d.SetError(c, http.StatusBadRequest, err)
		return
	}
	orderID, err := uuid.Parse(status.OrderId)
	if err != nil {
		d.logger.Sugar().Errorf("can't parse order id: %s", err)
		d.SetError(c, http.StatusBadRequest, err)
		return
	}
	if strings.ToLower(status.User.Role) == "user" {
		d.logger.Sugar().Errorf("the action not allowed: %s", err)
		d.SetError(c, http.StatusForbidden, err)
		return
	}
	userID, err := uuid.Parse(status.User.Id)
	if err != nil {
		d.logger.Sugar().Errorf("can't parse order id: %s", err)
		d.SetError(c, http.StatusBadRequest, err)
		return
	}
	err = d.orderUsecase.ChangeStatus(ctx, &models.Order{
		ID: orderID,
		User: models.User{
			ID: userID,
		},
	}, models.Status(status.Status))
	if err != nil {
		d.logger.Sugar().Errorf("can't change address for order with id: %s %s", orderID, err)
		d.SetError(c, http.StatusInternalServerError, err)
		return
	}
}
