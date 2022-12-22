package delivery

import (
	"OnlineShopBackend/internal/filestorage"
	"OnlineShopBackend/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title onlineShopBackend
// @version 1.0
// @description Backend service for online store
// @license.name MIT

// @contact.email example@mail.com

// @BasePath /

type Delivery struct {
	itemHandlers     handlers.IItemHandlers
	categoryHandlers handlers.ICategoryHandlers
	logger           *zap.Logger
	filestorage      filestorage.FileStorager
}

// NewDelivery initialize delivery layer
func NewDelivery(itemHandlers handlers.IItemHandlers, categoryHandlers handlers.ICategoryHandlers, logger *zap.Logger, fs filestorage.FileStorager) *Delivery {
	logger.Debug("Enter in NewDelivery()")
	return &Delivery{itemHandlers: itemHandlers, categoryHandlers: categoryHandlers, logger: logger, filestorage: fs}
}

// Index is the index handler.
func (delivery *Delivery) Index(c *gin.Context) {
	delivery.logger.Debug("Enter in Index")
	c.String(http.StatusOK, "Hello World!")
}


func (delivery *Delivery) GetFileList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetFileList()")
	fileInfos, err := delivery.filestorage.GetFileList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, fileInfos)
}
