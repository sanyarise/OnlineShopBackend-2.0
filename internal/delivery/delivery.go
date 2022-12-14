package delivery

import (
	"OnlineShopBackend/internal/filestorage"
	"OnlineShopBackend/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Delivery struct {
	itemHandlers     handlers.IItemHandlers
	userHandlers     handlers.IUserHandlers
	categoryHandlers handlers.ICategoryHandlers
	logger           *zap.Logger
	filestorage      filestorage.FileStorager
}

// NewDelivery initialize delivery layer
func NewDelivery(itemHandlers handlers.IItemHandlers, categoryHandlers handlers.ICategoryHandlers, userHandlers handlers.IUserHandlers, logger *zap.Logger, fs filestorage.FileStorager) *Delivery {
	logger.Debug("Enter in NewDelivery()")
	return &Delivery{itemHandlers: itemHandlers, categoryHandlers: categoryHandlers, userHandlers: userHandlers, logger: logger, filestorage: fs}
}

// Index is the index handler.
func (delivery *Delivery) Index(c *gin.Context) {
	delivery.logger.Debug("Enter in Index")
	c.String(http.StatusOK, "Hello World!")
}
