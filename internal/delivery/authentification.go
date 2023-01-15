package delivery

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (delivery *Delivery) Authentificate(c *gin.Context) {
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
			delivery.logger.Error("timeout of token is over")
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
		if user.Password != claims["Password"].(string) {
			delivery.logger.Sugar().Errorf("password from jwt not ident to user password")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		endpoint := delivery.getEndpointName(c.Request.Method, c.Request.RequestURI)
		c.Set("userId", user.ID.String())
		c.Set("role", user.Rights.Name)
		if !strings.Contains(endpoint, "unknown") {
			c.Set("endpoint", endpoint)
		}
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func (delivery *Delivery) getEndpointName(method string, uri string) string {
	delivery.logger.Sugar().Debugf("Enter in delivery getEndpointName() with args: method: %s, uri: %s", method, uri)

	switch {
	case method == "GET":
		switch {
		case strings.Contains(uri, "images"):
			return "imagesList"
		case strings.Contains(uri, "list/rights"):
			return "rightsList"
		case strings.Contains(uri, "user/list"):
			return "userList"
		case strings.Contains(uri, "user") && len(uri) > 10:
			return "userById"
		default:
			return fmt.Sprintf("unknown uri: %v", uri)
		}
	case method == "POST":
		switch {
		case strings.Contains(uri, "categories/create"):
			return "createCategory"
		case strings.Contains(uri, "categories/image/upload"):
			return "upCatImg"
		case strings.Contains(uri, "items/create"):
			return "createItem"
		case strings.Contains(uri, "items/image/upload"):
			return "upItImg"
		case strings.Contains(uri, "rights/create"):
			return "createRights"
		default:
			return fmt.Sprintf("unknown uri: %v", uri)
		}
	case method == "PUT":
		switch {
		case strings.Contains(uri, "categories"):
			return "updateCategory"
		case strings.Contains(uri, "items/update"):
			return "updateItem"
		case strings.Contains(uri, "rights/update"):
			return "updateRights"
		case strings.Contains(uri, "user/profile/edit"):
			return "updateUserProfile"
		case strings.Contains(uri, "user/update/rights"):
			return "updateUserRights"
		case strings.Contains(uri, "user/update/password"):
			return "updateUserPassword"
		default:
			return fmt.Sprintf("unknown uri: %v", uri)
		}
	case method == "DELETE":
		switch {
		case strings.Contains(uri, "categories/image/"):
			return "catImgDel"
		case strings.Contains(uri, "categories/delete"):
			return "categoryDelete"
		case strings.Contains(uri, "items/image"):
			return "itImgDel"
		case strings.Contains(uri, "items/delete"):
			return "itemDelete"
		case strings.Contains(uri, "rights/delete"):
			return "rightsDelete"
		default:
			return fmt.Sprintf("unknown uri: %v", uri)
		}
	default:
		return fmt.Sprintf("unknown method: %v", method)
	}
}
