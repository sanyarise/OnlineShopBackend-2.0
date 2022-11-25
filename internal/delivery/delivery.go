package delivery

import (
	"OnlineShopBackend/internal/filestorage"
	"OnlineShopBackend/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Delivery struct {
	handlers    *handlers.Handlers
	logger      *zap.Logger
	filestorage filestorage.FileStorager
}

// NewDelivery initialize delivery layer
func NewDelivery(handlers *handlers.Handlers, logger *zap.Logger, fs filestorage.FileStorager) *Delivery {
	logger.Debug("Enter in NewDelivery()")
	return &Delivery{handlers: handlers, logger: logger, filestorage: fs}
}

// Index is the index handler.
func (delivery *Delivery) Index(c *gin.Context) {
	delivery.logger.Debug("Enter in Index")
	c.String(http.StatusOK, "Hello World!")
}
