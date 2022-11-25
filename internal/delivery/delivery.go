package delivery

import (
	"OnlineShopBackend/internal/filestorage"
	"OnlineShopBackend/internal/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Delivery struct {
	h  *handlers.Handlers
	l  *zap.Logger
	fs filestorage.FileStorager
}

func NewDelivery(handlers *handlers.Handlers, logger *zap.Logger, fs filestorage.FileStorager) *Delivery {
	log.Println("Enter in NewDelivery()")
	return &Delivery{h: handlers, l: logger, fs: fs}
}

// Index is the index handler.
func (d *Delivery) Index(c *gin.Context) {
	d.l.Debug("Enter in Index")
	c.String(http.StatusOK, "Hello World!")
}
