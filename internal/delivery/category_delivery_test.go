package delivery

import (
	fs "OnlineShopBackend/internal/filestorage/mocks"
	"OnlineShopBackend/internal/handlers"
	"OnlineShopBackend/internal/handlers/mocks"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type WrongStruct struct {
	Name        int   `json:"name"`
	Description int32 `json:"description"`
}

var (
	testCategoryNoId = handlers.Category{
		Name: "testName",
	}
	testCategoryWithId = handlers.Category{
		Id:   testId.String(),
		Name: "testName",
	}
	testEmptyCategory = handlers.Category{}
	testWrong         = WrongStruct{
		Name:        5,
		Description: 6,
	}
	testList              = []handlers.Category{testCategoryWithId}
	testCategoryWithImage = handlers.Category{
		Id:    testId.String(),
		Name:  "testName",
		Image: "testImagePath",
	}
	testNoCategory = handlers.Category{
		Name:        "NoCategory",
		Description: "Category for items from deleting categories",
	}
	testNoCategoryWithId = handlers.Category{
		Id:          testId.String(),
		Name:        "NoCategory",
		Description: "Category for items from deleting categories",
	}
	testHandlersItemNoCat = handlers.Item{
		Id:    testId.String(),
		Title: "testTitle",
		Category: testNoCategoryWithId,
		Images: []string{"testName"},
	}
)

func MockCatJson(c *gin.Context, content interface{}, method string) {
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

func MockCatFile(c *gin.Context, fileType string, file []byte) {
	c.Request.Method = "POST"
	if fileType == "jpeg" {
		c.Request.Header.Set("Content-Type", "image/jpeg")
	} else {
		c.Request.Header.Set("Content-Type", "image/png")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(file))
}

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemHandlers(ctrl)
	categoryHandlers := mocks.NewMockICategoryHandlers(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemHandlers, categoryHandlers, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockCatJson(c, testCategoryNoId, post)
	res := gin.H{"success": testId.String()}
	bytesRes, _ := json.Marshal(&res)
	categoryHandlers.EXPECT().CreateCategory(ctx, testCategoryNoId).Return(testId, nil)
	delivery.CreateCategory(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, bytesRes, w.Body.Bytes())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockCatJson(c, testWrong, post)
	c.Request.Header.Set("Content-Type", "application/text")
	delivery.CreateCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testEmptyCategory, post)
	delivery.CreateCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockJson(c, testCategoryNoId, post)
	categoryHandlers.EXPECT().CreateCategory(ctx, testCategoryNoId).Return(uuid.Nil, fmt.Errorf("error"))
	delivery.CreateCategory(c)
	require.Equal(t, 500, w.Code)
}

func TestUpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemHandlers(ctrl)
	categoryHandlers := mocks.NewMockICategoryHandlers(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemHandlers, categoryHandlers, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatJson(c, testCategoryWithId, put)
	categoryHandlers.EXPECT().UpdateCategory(ctx, testCategoryWithId).Return(nil)
	delivery.UpdateCategory(c)
	require.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatJson(c, testCategoryWithId, put)
	categoryHandlers.EXPECT().UpdateCategory(ctx, testCategoryWithId).Return(fmt.Errorf("error"))
	delivery.UpdateCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatJson(c, testWrong, put)
	delivery.UpdateCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.UpdateCategory(c)
	require.Equal(t, 400, w.Code)
}

func TestGetCategoryList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemHandlers(ctrl)
	categoryHandlers := mocks.NewMockICategoryHandlers(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemHandlers, categoryHandlers, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	categoryHandlers.EXPECT().GetCategoryList(ctx).Return([]handlers.Category{}, fmt.Errorf("error"))
	delivery.GetCategoryList(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	testBytes, _ := json.Marshal(&testList)
	categoryHandlers.EXPECT().GetCategoryList(ctx).Return(testList, nil)
	delivery.GetCategoryList(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, testBytes, w.Body.Bytes())
}

func TestGetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemHandlers(ctrl)
	categoryHandlers := mocks.NewMockICategoryHandlers(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemHandlers, categoryHandlers, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.GetCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testEmptyCategory, fmt.Errorf("error"))
	delivery.GetCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	testBytes, _ := json.Marshal(&testCategoryWithId)
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithId, nil)
	delivery.GetCategory(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, testBytes, w.Body.Bytes())
}

func TestUploadCategoryImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemHandlers(ctrl)
	categoryHandlers := mocks.NewMockICategoryHandlers(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemHandlers, categoryHandlers, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.UploadCategoryImage(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}

	delivery.UploadCategoryImage(c)
	require.Equal(t, 415, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatFile(c, "jpeg", testFile)
	filestorage.EXPECT().PutCategoryImage(testId.String(), carbon.Now().ToShortDateTimeString()+".jpeg", testFile).Return("", fmt.Errorf("error"))
	delivery.UploadCategoryImage(c)
	require.Equal(t, 507, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatFile(c, "png", testFile)
	filestorage.EXPECT().PutCategoryImage(testId.String(), carbon.Now().ToShortDateTimeString()+".png", testFile).Return("testImagePath", nil)
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testEmptyCategory, fmt.Errorf("error"))
	delivery.UploadCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatFile(c, "png", testFile)
	filestorage.EXPECT().PutCategoryImage(testId.String(), carbon.Now().ToShortDateTimeString()+".png", testFile).Return("testImagePath", nil)
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithId, nil)
	categoryHandlers.EXPECT().UpdateCategory(ctx, testCategoryWithImage).Return(fmt.Errorf("error"))
	delivery.UploadCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	MockCatFile(c, "png", testFile)
	filestorage.EXPECT().PutCategoryImage(testId.String(), carbon.Now().ToShortDateTimeString()+".png", testFile).Return("testImagePath", nil)
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithId, nil)
	categoryHandlers.EXPECT().UpdateCategory(ctx, testCategoryWithImage).Return(nil)
	delivery.UploadCategoryImage(c)
	require.Equal(t, 201, w.Code)
}

func TestDeleteCategoryImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemHandlers(ctrl)
	categoryHandlers := mocks.NewMockICategoryHandlers(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemHandlers, categoryHandlers, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	MockCatJson(c, testWrong, post)
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testName").Return(fmt.Errorf("error"))
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testName").Return(nil)
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testEmptyCategory, fmt.Errorf("error"))
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testImagePath", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testImagePath").Return(nil)
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithImage, nil)
	categoryHandlers.EXPECT().UpdateCategory(ctx, testCategoryWithId).Return(fmt.Errorf("error"))
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testImagePath", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testImagePath").Return(nil)
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithImage, nil)
	categoryHandlers.EXPECT().UpdateCategory(ctx, testCategoryWithId).Return(nil)
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 200, w.Code)
}

func TestDeleteCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemHandlers(ctrl)
	categoryHandlers := mocks.NewMockICategoryHandlers(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemHandlers, categoryHandlers, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	delivery.DeleteCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String() + "1",
		},
	}
	delivery.DeleteCategory(c)
	require.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testEmptyCategory, fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithId, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithImage, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(fmt.Errorf("error"))
	itemHandlers.EXPECT().ItemsQuantity(ctx).Return(-1, fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithImage, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(fmt.Errorf("error"))
	itemHandlers.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	itemHandlers.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage.Name, 0, 1).Return([]handlers.Item{}, fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithImage, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(fmt.Errorf("error"))
	itemHandlers.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	itemHandlers.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage.Name, 0, 1).Return([]handlers.Item{testHandlersItemWithImage}, nil)
	categoryHandlers.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(testEmptyCategory, fmt.Errorf("error"))
	categoryHandlers.EXPECT().CreateCategory(ctx, testNoCategory).Return(uuid.Nil, fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithImage, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(fmt.Errorf("error"))
	itemHandlers.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	itemHandlers.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage.Name, 0, 1).Return([]handlers.Item{testHandlersItemWithImage}, nil)
	categoryHandlers.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(testEmptyCategory, fmt.Errorf("error"))
	categoryHandlers.EXPECT().CreateCategory(ctx, testNoCategory).Return(testId, nil)
	itemHandlers.EXPECT().UpdateItem(ctx, testHandlersItemNoCat).Return(fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "categoryID",
			Value: testId.String(),
		},
	}
	categoryHandlers.EXPECT().GetCategory(ctx, testId.String()).Return(testCategoryWithImage, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(fmt.Errorf("error"))
	itemHandlers.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	itemHandlers.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage.Name, 0, 1).Return([]handlers.Item{testHandlersItemWithImage}, nil)
	categoryHandlers.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(testNoCategoryWithId, nil)
	itemHandlers.EXPECT().UpdateItem(ctx, testHandlersItemNoCat).Return(fmt.Errorf("error"))
	delivery.DeleteCategory(c)
	require.Equal(t, 200, w.Code)
}
