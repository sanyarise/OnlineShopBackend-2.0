package jwtauth

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type Payload struct {
	Email  string `json:"email"`
	Role   string `json:"role"`
	UserId uuid.UUID `json:"userId"`
	jwt.StandardClaims
}

func NewJWT(payload Payload) (string, error) {
	key := []byte("dsf498uh324seyu2837912sd7*7897") //TODO make env
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
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

//func UserIdentity(header string) (*Payload, error) {
//	userCr, err := parseAuthHeader(header)
//	if err != nil {
//		return &Payload{}, err
//	}
//	return userCr, nil
//}

func CreateSessionJWT(ctx context.Context, user *models.User) (Token, error) {
	payload := Payload{
		Email:  user.Email,
		Role:   user.Rights.Name,
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
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

//func checkJWT(tokenString string) error {
//	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("unexpexted signing method: %v", token.Header["alg"])
//		}
//		return []byte("dsf498uh324seyu2837912sd7*7897"), nil //TODO make env
//	})
//
//	if err != nil {
//		return fmt.Errorf("invalid JWT")
//	}
//
//	return nil
//}

//func parseAuthHeader(header string) (*Payload, error) {
//	if header == "" {
//		return &Payload{}, errors.New("empty header")
//	}
//
//	headerSplit := strings.Split(header, " ")
//	if len(headerSplit) != 2 || headerSplit[0] != "Bearer" {
//		return &Payload{}, errors.New("header issue")
//	}
//	if len(headerSplit[1]) == 0 {
//		return &Payload{}, errors.New("empty token")
//	}
//
//	err := checkJWT(headerSplit[1]); if err != nil {
//		return &Payload{}, err
//	}
//
//	parts := strings.Split(headerSplit[1], ".")
//
//	if parts == nil {
//		return &Payload{}, errors.New(parts[2]) //todo
//	}
//
//	email, err := jwt.DecodeSegment(parts[1])
//	if err != nil {
//		return &Payload{}, errors.New("unable to decode")
//	}
//	cr := &Payload{}
//	err = json.Unmarshal(email, &cr)
//	if err != nil {
//		return &Payload{}, errors.New("unable to unmarshall")
//	}
//	payload := &Payload{
//		Email:  cr.Email,
//		Role:   cr.Role,
//		UserId: cr.UserId,
//	}
//	return payload, nil
//}