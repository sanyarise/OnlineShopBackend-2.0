package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var _ IUserUsecase = &UserUsecase{}

type UserUsecase struct {
	userStore repository.UserStore
	logger    *zap.Logger
}

func NewUserUsecase(userStore repository.UserStore, logger *zap.Logger) IUserUsecase {
	logger.Debug("Enter in usecase NewUserUsecase()")
	return &UserUsecase{userStore: userStore, logger: logger}
}

type Payload struct {
	jwt.StandardClaims
	Email    string
	Role     string
	UserId   uuid.UUID
	Password string
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (usecase *UserUsecase) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase CreateUser() with args: ctx, user: %v", user)

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

func (usecase *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetUserByEmail() with args: ctx, email: %s", email)
	var user *models.User

	user, err := usecase.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (usecase *UserUsecase) GetRightsId(ctx context.Context, name string) (*models.Rights, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase user GetRightsId() with args: ctx, name: %s", name)
	rights, err := usecase.userStore.GetRightsId(ctx, name)
	if err != nil {
		return nil, err
	}
	return &rights, nil
}

func (usecase *UserUsecase) UpdateUserData(ctx context.Context, user *models.User) (*models.User, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateUserData() with args: ctx, user: %v", user)
	user, err := usecase.userStore.UpdateUserData(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *UserUsecase) NewJWT(payload Payload, key string) (string, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase user NewJWT() with args: payload: %v, key: %s", payload, key)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload) //jwt.SigningMethodHS256
	return token.SignedString(key)
}

func (usecase *UserUsecase) NewRefreshToken() (string, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase user NewRefreshToken()")
	refreshToken := make([]byte, 32)
	t := rand.NewSource(time.Now().Unix())
	r := rand.New(t)

	_, err := r.Read(refreshToken)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", refreshToken), nil
}

func (usecase *UserUsecase) UserIdentity(header string) (*Payload, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase UserIdentity() with args: header: %s", header)
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
		Email:    cr.Email,
		Role:     cr.Role,
		UserId:   cr.UserId,
		Password: cr.Password,
	}
	return userCred, nil
}

func (usecase *UserUsecase) CreateSessionJWT(ctx context.Context, user *models.User, key string) (Token, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase user CreateSessionJWT() with args: ctx, user: %v, key: %s", user, key)
	payload := Payload{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
		Email:    user.Email,
		Role:     user.Rights.Name,
		UserId:   user.ID,
		Password: user.Password,
	}

	accessToken, err := usecase.NewJWT(payload, key)
	if err != nil {
		return Token{}, fmt.Errorf("unable to create a token")
	}

	refreshToken, err := usecase.NewRefreshToken()
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
