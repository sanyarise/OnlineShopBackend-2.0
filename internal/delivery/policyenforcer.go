package delivery

import (
	"OnlineShopBackend/internal/delivery/gateway"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EnforcePolicy is a middleware to check all the requests and ask if the requester has permissions
func (delivery *Delivery) EnforcePolicy(c *gin.Context, gateway gateway.PolicyGateway) {
	delivery.logger.Info("Enforcing policy with middleware")

	allow := gateway.Ask(c)

	if allow {
		delivery.logger.Info("Action is allowed, continuing process")
		c.Next()
	} else {
		delivery.logger.Info("Action was not allowed, cancelling process")
		c.String(http.StatusForbidden, "Action not allowed")
	}
}
