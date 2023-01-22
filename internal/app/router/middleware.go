package router

import (
	"OnlineShopBackend/internal/delivery/user/jwtauth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Auth-Token, Set-Cookie")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is empty"})
			c.Abort()
			return
		}

		headerSplit := strings.Split(tokenString, " ")
		if len(headerSplit) != 2 || headerSplit[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "header issue"})
			return
		}
		if len(headerSplit[1]) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "empty token"})
			return
		}

		tt := headerSplit[1]

		claims := &jwtauth.Payload{}
		token, err := jwt.ParseWithClaims(tt, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("dsf498uh324seyu2837912sd7*7897"), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid signature"})
				c.Abort()
				return
			}
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": tokenString})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func noOpMiddleware(c *gin.Context) {
	c.Next()
}
