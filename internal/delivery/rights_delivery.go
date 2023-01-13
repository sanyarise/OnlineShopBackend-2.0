package delivery

import (
	"OnlineShopBackend/internal/delivery/rights"
	"OnlineShopBackend/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateRights
//
//	@Summary		Method provides to create rights
//	@Description	Method provides to create rights.
//	@Tags			rights
//	@Accept			json
//	@Produce		json
//	@Param			rights	body		rights.ShortRights	true	"Data for creating rights"
//	@Success		201		{object}	rights.RightsId
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/rights/create/ [post]
func (delivery *Delivery) CreateRights(c *gin.Context) {
	delivery.logger.Sugar().Debugf("Enter in delivery CreateRights()")

	var createdRights rights.ShortRights
	if err := c.ShouldBindJSON(&createdRights); err != nil {
		delivery.logger.Error(fmt.Sprintf("error on bind json from request: %v", err))
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	if createdRights.Name == "" {
		err := fmt.Errorf("empty name is not correct")
		if err != nil {
			delivery.logger.Error(err.Error())
			delivery.SetError(c, http.StatusBadRequest, err)
			return
		}
	}
	ctx := c.Request.Context()
	id, err := delivery.rightsUsecase.CreateRights(ctx, &models.Rights{
		Name:  createdRights.Name,
		Rules: createdRights.Rules,
	})
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, rights.RightsId{Value: id.String()})
}

// UpdateItem - update rights
//
//	@Summary		Method provides to update rights
//	@Description	Method provides to update rights
//	@Tags			rights
//	@Accept			json
//	@Produce		json
//	@Param			rights	body	rights.OutRights	true	"Data for updating rights"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/rights/update [put]
func (delivery *Delivery) UpdateRights(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UpdateRights()")

	var updatedRights rights.OutRights
	if err := c.ShouldBindJSON(&updatedRights); err != nil {
		delivery.logger.Error(fmt.Sprintf("error on bind json from request: %v", err))
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	id, err := uuid.Parse(updatedRights.Id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()

	err = delivery.rightsUsecase.UpdateRights(ctx, &models.Rights{
		ID:    id,
		Name:  updatedRights.Name,
		Rules: updatedRights.Rules,
	})
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteRights deleted rights by id
//
//	@Summary		Method provides to delete rights
//	@Description	Method provides to delete rights.
//	@Tags			rights
//	@Accept			json
//	@Produce		json
//	@Param			rightsID	path	string	true	"id of rights"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/rights/delete/{rightsID} [delete]
func (delivery *Delivery) DeleteRights(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteRights()")

	id, err := uuid.Parse(c.Param("rightsID"))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()

	err = delivery.rightsUsecase.DeleteRights(ctx, id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetRights - returns rights by id
//
//	@Summary		Get rights by id
//	@Description	The method allows you to get the rights by id.
//	@Tags			rights
//	@Accept			json
//	@Produce		json
//	@Param			rightsID	path		string				true	"id of rights"
//	@Success		200			{object}	rights.OutRights	"Rights structure"
//	@Failure		400			{object}	ErrorResponse
//	@Failure		403			"Forbidden"
//	@Failure		404			{object}	ErrorResponse	"404 Not Found"
//	@Failure		500			{object}	ErrorResponse
//	@Router			/rights/{rightsID} [get]
func (delivery *Delivery) GetRights(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetRights()")

	id, err := uuid.Parse(c.Param("rightsID"))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()

	right, err := delivery.rightsUsecase.GetRights(ctx, id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, rights.OutRights{
		Id:    right.ID.String(),
		Name:  right.Name,
		Rules: right.Rules,
	})
}

// RightsList - returns list of all rights
//
//	@Summary		Get list of rights
//	@Description	Method provides to get list of rights
//	@Tags			rights
//	@Accept			json
//	@Produce		json
//	@Success		200	array		rights.OutRights	"List of rights"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/rights/list [get]
func (delivery *Delivery) RightsList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery RightsList()")

	ctx := c.Request.Context()

	list, err := delivery.rightsUsecase.RightsList(ctx)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	outList := make([]rights.OutRights, len(list))
	for i, right := range list {
		outList[i] = rights.OutRights{
			Id:    right.ID.String(),
			Name:  right.Name,
			Rules: right.Rules,
		}
	}
	c.JSON(http.StatusOK, outList)
}
