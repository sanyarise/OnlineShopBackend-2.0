package delivery

import (
	"OnlineShopBackend/internal/delivery/rights"
	"OnlineShopBackend/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	c.JSON(http.StatusOK, rights.RightsList{List: outList})
}
