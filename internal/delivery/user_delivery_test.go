package delivery

import (
	fs "OnlineShopBackend/internal/filestorage/mocks"
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase"
	"OnlineShopBackend/internal/usecase/mocks"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	//"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testModelNewUser = &models.User{
	Firstname: "Test",
	Password:  "password",
	Email:     "test@gmail.com",
}

//var testNewUserId = uuid.New()

var testTokenSet = &usecase.Token{
	AccessToken:  "",
	RefreshToken: "",
}



func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	userUsecase := mocks.NewMockIUserUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemUsecase, userUsecase, categoryUsecase, cartUsecase, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	userUsecase.EXPECT().CreateUser(ctx, testModelNewUser).Return(testModelNewUser, nil)
	delivery.CreateUser(c)
	require.Equal(t, 201, w.Code)



}


