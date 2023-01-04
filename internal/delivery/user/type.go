package user

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtRealisation struct {
	ttl int64
}

type Payload struct {
	jwt.StandardClaims
	Name uuid.UUID
}

type Credentials struct {
	Email string `json:"email"`
	Password string `json:"password"`

}

