package jwtauth

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type Payload struct {
	jwt.StandardClaims
	Email    string
	Role     string
	UserId   uuid.UUID
}

func NewJWT(payload Payload) (string, error) {
	key := []byte("dsf498uh324seyu2837912sd7*7897")              //TODO
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload) //jwtauth.SigningMethodHS256
	return token.SignedString(key)
}

func NewRefreshToken() (string, error) {
	refreshToken := make([]byte, 32)
	t := rand.NewSource(time.Now().Unix())
	r := rand.New(t)

	_, err := r.Read(refreshToken)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", refreshToken), nil
}

func ParseAuthHeader(header string) (*Payload, error) {
	if header == "" {
		return &Payload{}, errors.New("empty header")
	}
	headerSplit := strings.Split(header, " ")
	if len(headerSplit) != 2 || headerSplit[0] != "Bearer" {
		return &Payload{}, errors.New("header issue")
	}
	if len(headerSplit[1]) == 0 {
		return &Payload{}, errors.New("empty token")
	}

	parts := strings.Split(headerSplit[1], ".")

	if parts == nil {
		return &Payload{}, errors.New(parts[2]) //todo
	}

	email, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return &Payload{}, errors.New("unable to decode")
	}
	cr := &Payload{}
	err = json.Unmarshal(email, &cr)
	if err != nil {
		return &Payload{}, errors.New("unable to unmarshall")
	}
	payload := &Payload{
		Email:    cr.Email,
		Role:     cr.Role,
		UserId:   cr.UserId,
	}
	return payload, nil
}

func UserIdentity(header string) (*Payload, error) {
	userCr, err := ParseAuthHeader(header)
	if err != nil {
		return &Payload{}, errors.New("parse error")
	}
	return userCr, nil
}

func CreateSessionJWT(ctx context.Context, user *models.User) (Token, error) {
	payload := Payload{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
		Email:    user.Email,
		Role:     user.Rights.Name,
		UserId:   user.ID,
	}

	accessToken, err := NewJWT(payload)
	if err != nil {
		return Token{}, fmt.Errorf("unable to create a token")
	}

	refreshToken, err := NewRefreshToken()
	if err != nil {
		return Token{}, fmt.Errorf("unable to create a refresh token")
	}

	token := Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	//if err = usecase.userStore.SaveSession(ctx, refreshToken, payload.ExpiresAt); err != nil {
	//	return Token{}, fmt.Errorf("unable to save session")
	//}

	return token, nil
}

