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
		delivery.logger.Error(err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	delivery.logger.Sugar().Debugf("tokenString read from request success: %v", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}
		return []byte(delivery.secretKey), nil
	})
	if err != nil {
		delivery.logger.Error(err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	delivery.logger.Sugar().Debugf("Token parse from tokenstring success: %v", token)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		delivery.logger.Sugar().Debugf("claims from token read success: %v", claims)
		interfaceValue := claims["Email"]
		delivery.logger.Sugar().Debugf("interface value email : %v", interfaceValue)
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
		c.Set("userId", user.ID.String())
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
