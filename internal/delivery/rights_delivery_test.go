package delivery

import (
	auth "OnlineShopBackend/internal/delivery/mocks"
	"OnlineShopBackend/internal/delivery/rights"
	fs "OnlineShopBackend/internal/filestorage/mocks"
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/mocks"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type WrongShortRights struct {
	Name int
}
type WrongRights struct {
	Id int
	Name int
}

var (
	testShortRights = rights.ShortRights{
		Name: "test",
	}
	testModelRightsNoId = models.Rights{
		Name: "test",
	}
	testSREmptyName = rights.ShortRights{
		Name: "",
	}
	testWrongShortRights = WrongShortRights{
		Name: 5,
	}
	testWrongRights = WrongRights{
		Id: 5,
		Name: 10,
	}
	testRights = rights.OutRights{
		Id: testId.String(),
		Name: "test",
	}
	testModelRights = &models.Rights{
		ID: testId,
		Name: "test",
	}
)

func MockRightsJson(c *gin.Context, content interface{}, method string) {
	if method == "POST" {
		c.Request.Method = "POST"
	}
	if method == "PUT" {
		c.Request.Method = "PUT"
	}

	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestCreateRights(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	rightsUsecase := mocks.NewMockIRightsUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	authorization := auth.NewMockPolicyGateway(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, rightsUsecase, logger, filestorage, authorization, "")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testWrongShortRights, post)
	delivery.CreateRights(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testSREmptyName, post)
	delivery.CreateRights(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testShortRights, post)
	rightsUsecase.EXPECT().CreateRights(ctx, &testModelRightsNoId).Return(uuid.Nil, err)
	delivery.CreateRights(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testShortRights, post)
	rightsUsecase.EXPECT().CreateRights(ctx, &testModelRightsNoId).Return(testId, nil)
	delivery.CreateRights(c)
	require.Equal(t, 201, w.Code)
}

func TestUpdateRights(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemUsecase := mocks.NewMockIItemUsecase(ctrl)
	categoryUsecase := mocks.NewMockICategoryUsecase(ctrl)
	cartUsecase := mocks.NewMockICartUsecase(ctrl)
	rightsUsecase := mocks.NewMockIRightsUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	authorization := auth.NewMockPolicyGateway(ctrl)
	delivery := NewDelivery(itemUsecase, nil, categoryUsecase, cartUsecase, rightsUsecase, logger, filestorage, authorization, "")


}
