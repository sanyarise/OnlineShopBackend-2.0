package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Id     uuid.UUID   `json:"id"`
	Error  string      `json:"message,omitempty"`
	Errors []string    `json:"errors,omitempty"`
	Info   interface{} `json:"info,omitempty"`
}

func (delivery *Delivery) SetError(c *gin.Context, statusCode int, errs ...error) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(statusCode, gin.H{})
		}
	}()
	var response = ErrorResponse{
		Id: uuid.New(),
	}

	if len(errs) == 0 {
		return
	}

	if len(errs) > 0 {
		response.Error = errs[0].Error()

		if len(errs) > 1 {
			for _, err := range errs {
				response.Errors = append(response.Errors, c.Error(err).Error())
			}
		}
	}
	c.JSON(statusCode, response)

	fields := getContextFields(c)
	if statusCode >= 400 && statusCode < 500 {
		delivery.logger.Warn(errs[len(errs)-1].Error(), fields...)
	} else if statusCode >= 500 {
		delivery.logger.Error(errs[len(errs)-1].Error(), fields...)
	}
}

func getContextFields(c *gin.Context) []zap.Field {
	var fields = []zap.Field{zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()),
	}
	return fields
}
