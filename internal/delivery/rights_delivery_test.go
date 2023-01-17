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
	Id   int
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
		Id:   5,
		Name: 10,
	}
	testRights = rights.OutRights{
		Id:   testId.String(),
		Name: "test",
	}
	testUpdatedRights = rights.OutRights{
		Id:   testId.String(),
		Name: "Updatedtest",
	}
	testModelUpdateRights = &models.Rights{
		ID:   testId,
		Name: "Updatedtest",
	}
	testModelRights = &models.Rights{
		ID:   testId,
		Name: "test",
	}
	testModelsRightsList = []models.Rights{
		*testModelRights,
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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testWrongRights, put)
	delivery.UpdateRights(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testUpdatedRights, put)
	rightsUsecase.EXPECT().UpdateRights(ctx, testModelUpdateRights).Return(err)
	delivery.UpdateRights(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockRightsJson(c, testUpdatedRights, put)
	rightsUsecase.EXPECT().UpdateRights(ctx, testModelUpdateRights).Return(nil)
	delivery.UpdateRights(c)
	require.Equal(t, 200, w.Code)
}

func TestDeleteRights(t *testing.T) {
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

	c.Params = []gin.Param{
		{
			Key:   "rightsID",
			Value: testId.String() + "k",
		},
	}

	delivery.DeleteRights(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = []gin.Param{
		{
			Key:   "rightsID",
			Value: testId.String(),
		},
	}

	rightsUsecase.EXPECT().DeleteRights(ctx, testId).Return(err)
	delivery.DeleteRights(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = []gin.Param{
		{
			Key:   "rightsID",
			Value: testId.String(),
		},
	}

	rightsUsecase.EXPECT().DeleteRights(ctx, testId).Return(nil)
	delivery.DeleteRights(c)
	require.Equal(t, 200, w.Code)
}

func TestGetRights(t *testing.T) {
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

	c.Params = []gin.Param{
		{
			Key:   "rightsID",
			Value: testId.String() + "k",
		},
	}
	delivery.GetRights(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = []gin.Param{
		{
			Key:   "rightsID",
			Value: testId.String(),
		},
	}
	rightsUsecase.EXPECT().GetRights(ctx, testId).Return(nil, err)
	delivery.GetRights(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = []gin.Param{
		{
			Key:   "rightsID",
			Value: testId.String(),
		},
	}
	rightsUsecase.EXPECT().GetRights(ctx, testId).Return(testModelRights, nil)
	delivery.GetRights(c)
	require.Equal(t, 200, w.Code)
}

func TestRightsList(t *testing.T) {
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

	rightsUsecase.EXPECT().RightsList(ctx).Return(nil, err)
	delivery.RightsList(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	rightsUsecase.EXPECT().RightsList(ctx).Return(testModelsRightsList, nil)
	delivery.RightsList(c)
	require.Equal(t, 200, w.Code)
}
