package delivery

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (delivery *Delivery) Authorize(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery Authorize()")
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}
		return []byte(delivery.secretKey), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if int64(time.Now().Unix()) > claims["ExpiresAt"].(int64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		interfaceValue := claims["Email"]
		var email string
		switch interfaceValue.(type) {
		case string:
			v := reflect.ValueOf(interfaceValue)
			email = v.String()
		default:
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fmt.Println(email)
		user, err := delivery.userUsecase.GetUserByEmail(c.Request.Context(), email)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if user.Email == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
