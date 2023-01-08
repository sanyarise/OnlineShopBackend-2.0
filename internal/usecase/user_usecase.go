package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
	"time"

	"go.uber.org/zap"
)

var _ IUserUsecase = &UserUsecase{}

type UserUsecase struct {
	userStore repository.UserStore
	logger    *zap.Logger
}

func NewUserUsecase(userStore repository.UserStore, logger *zap.Logger) IUserUsecase {
	return &UserUsecase{userStore: userStore, logger: logger}
}

type JwtRealisation struct {
	ttl int64
}

type Payload struct {
	jwt.StandardClaims
	Email  string
	Role   string
	UserId uuid.UUID
	Password string
}

type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type SessionLog struct {
	logger zap.Logger
}

type Profile struct {
	Email     string  `json:"email,omitempty"`
	FirstName string  `json:"firstname,omitempty"`
	LastName  string  `json:"lastname,omitempty"`
	Address   Address `json:"address,omitempty"`
	Rights    Rights  `json:"rights,omitempty"`
}

type Address struct {
	Zipcode string `json:"zipcode,omitempty"`
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
	Street  string `json:"street,omitempty"`
}

type Rights struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Name  string    `json:"name,omitempty"`
	Rules []string  `json:"rules,omitempty"`
}

func (usecase *UserUsecase) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	usecase.logger.Debug("Enter in usecase CreateUser()")

	rights, err := usecase.userStore.GetRightsId(ctx, "Customer")
	if err != nil {
		return &models.User{}, err
	}

	usecaseUser := &models.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		Address: models.UserAddress{
			Zipcode: user.Address.Zipcode,
			Country: user.Address.Country,
			City:    user.Address.City,
			Street:  user.Address.Street,
		},
		Rights: models.Rights{
			ID:    rights.ID,
			Name:  rights.Name,
			Rules: rights.Rules,
		},
	}

	id, err := usecase.userStore.Create(ctx, usecaseUser)
	if err != nil {
		return &models.User{}, fmt.Errorf("error on create user: %w", err)
	}
	return id, nil
}

func (usecase *UserUsecase) GetUserByEmail(ctx context.Context, email string, password string) (models.User, error) {
	var user models.User

	user, err := usecase.userStore.GetUserByEmail(ctx, email, password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (usecase *UserUsecase) GetRightsId(ctx context.Context, name string) (*models.Rights, error) {
	//var rights models.Rights

	rights, err := usecase.userStore.GetRightsId(ctx, name)
	if err != nil {
		return nil, err
	}

	return &rights, nil
}

func (usecase *UserUsecase) UpdateUserData(ctx context.Context, user *models.User) (*models.User, error) {
	user, err := usecase.userStore.UpdateUserData(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *UserUsecase) NewJWT(payload Payload) (string, error) {
	key := []byte("dsf498uh324seyu2837912sd7*7897")              //TODO
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload) //jwt.SigningMethodHS256
	return token.SignedString(key)
}

func (usecase *UserUsecase) NewRefreshToken() (string, error) {
	refreshToken := make([]byte, 32)
	t := rand.NewSource(time.Now().Unix())
	r := rand.New(t)

	_, err := r.Read(refreshToken)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", refreshToken), nil
}

func (usecase *UserUsecase) ParseAuthHeader(header string) (*Payload, error) {
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
	userCred := &Payload{
		Email:  cr.Email,
		Role:   cr.Role,
		UserId: cr.UserId,
		Password: cr.Password,
	}
	return userCred, nil
}

func (usecase *UserUsecase) UserIdentity(header string) (*Payload, error) {
	userCr, err := usecase.ParseAuthHeader(header)
	if err != nil {
		return &Payload{}, errors.New("parse error")
	}
	return userCr, nil
}

func (usecase *UserUsecase) CreateSessionJWT(ctx context.Context, user *models.User) (Token, error) {
	payload := Payload{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
		Email:  user.Email,
		Role:   user.Rights.Name,
		UserId: user.ID,
		Password: user.Password,
	}

	accessToken, err := usecase.NewJWT(payload)
	if err != nil {
		return Token{}, fmt.Errorf("unable to create a token")
	}

	refreshToken, err := usecase.NewRefreshToken()
	if err != nil {
		return Token{}, fmt.Errorf("unable to create a refresh token")
	}

	token := Token{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	//if err = usecase.userStore.SaveSession(ctx, refreshToken, payload.ExpiresAt); err != nil {
	//	return Token{}, fmt.Errorf("unable to save session")
	//}

	return token, nil
}
