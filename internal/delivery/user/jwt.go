package user

import (
	"github.com/golang-jwt/jwt"
)

func NewJWT(payload Payload) (string, error) {
	key := []byte("dsf498uh324seyu2837912sd7*(*7897") //TODO
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	return token.SignedString(key)
}
