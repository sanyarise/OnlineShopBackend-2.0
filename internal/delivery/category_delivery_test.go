package delivery

import (
	"OnlineShopBackend/internal/delivery/category"
	fs "OnlineShopBackend/internal/filestorage/mocks"
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/mocks"
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
	testCategoryNoId = &models.Category{
		Name:        "testName",
		Description: "testDescription",
	}
	testShortCategory = category.ShortCategory{
		Name:        "testName",
		Description: "testDescription",
	}
	testCategoryId = category.CategoryId{
		Value: testId.String(),
	}
	testCategoryWithId = category.Category{
		Id:          testId.String(),
		Name:        "testName",
		Description: "testDescription",
	}
	testModelsCategoryWithId = &models.Category{
		Id:          testId,
		Name:        "testName",
		Description: "testDescription",
	}
	testModelsCategoryWithId2 = &models.Category{
		Id:          testId,
		Name:        "testName",
		Description: "testDescription",
	}
	testEmptyCategory = category.ShortCategory{
		Name:        "",
		Description: "",
	}
	testEmptyModelsCategory = models.Category{
		Id:          uuid.Nil,
		Name:        "",
		Description: "",
	}
	testWrong = WrongStruct{
		Name:        5,
		Description: 6,
	}
	testList    = []models.Category{*testModelsCategoryWithId}
	testOutList = category.CategoriesList{
		List: []category.Category{
			testCategoryWithId,
		},
	}
	testCategoryWithImage = models.Category{
		Id:          testId,
		Name:        "testName",
		Description: "testDescription",
		Image:       "testImagePath",
	}
	testNoCategory = models.Category{
		Name:        "NoCategory",
		Description: "Category for items from deleting categories",
	}
	testNoCategoryWithId = models.Category{
		Id:          testId,
		Name:        "NoCategory",
		Description: "Category for items from deleting categories",
	}
	testCategoryWithImage2 = &models.Category{
		Id:          testId,
		Name:        "testName",
		Description: "testDescription",
		Image:       "testImagePath",
	}
	testModelsItemNoCat = models.Item{
		Id:          testId,
		Title:       "testTitle",
		Description: "testDescription",
		Category:    testNoCategoryWithId,
		Price:       10,
		Vendor:      "testVendor",
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
	itemHandlers := mocks.NewMockIItemUsecase(ctrl)
	categoryHandlers := mocks.NewMockICategoryUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemHandlers, categoryHandlers, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	MockCatJson(c, testShortCategory, post)
	bytesRes, _ := json.Marshal(&testCategoryId)
	categoryHandlers.EXPECT().CreateCategory(ctx, testCategoryNoId).Return(testId, nil)
	delivery.CreateCategory(c)
	require.Equal(t, 201, w.Code)
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
	MockCatJson(c, testEmptyCategory, post)
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
	MockJson(c, testShortCategory, post)
	categoryHandlers.EXPECT().CreateCategory(ctx, testCategoryNoId).Return(uuid.Nil, fmt.Errorf("error"))
	delivery.CreateCategory(c)
	require.Equal(t, 500, w.Code)
}

func TestUpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemUsecase(ctrl)
	categoryHandlers := mocks.NewMockICategoryUsecase(ctrl)
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
	categoryHandlers.EXPECT().UpdateCategory(ctx, testModelsCategoryWithId).Return(nil)
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
	categoryHandlers.EXPECT().UpdateCategory(ctx, testModelsCategoryWithId).Return(fmt.Errorf("error"))
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
			Value: testId.String() + "l",
		},
	}
	MockCatJson(c, testCategoryWithId, put)
	delivery.UpdateCategory(c)
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
	itemHandlers := mocks.NewMockIItemUsecase(ctrl)
	categoryHandlers := mocks.NewMockICategoryUsecase(ctrl)
	filestorage := fs.NewMockFileStorager(ctrl)
	delivery := NewDelivery(itemHandlers, categoryHandlers, logger, filestorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	categoryHandlers.EXPECT().GetCategoryList(ctx).Return([]models.Category{}, fmt.Errorf("error"))
	delivery.GetCategoryList(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	testBytes, _ := json.Marshal(&testOutList.List)
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
	itemHandlers := mocks.NewMockIItemUsecase(ctrl)
	categoryHandlers := mocks.NewMockICategoryUsecase(ctrl)
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, fmt.Errorf("error"))
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
			Value: testId.String() + "n",
		},
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
	testBytes, _ := json.Marshal(&testCategoryWithId)
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	delivery.GetCategory(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, testBytes, w.Body.Bytes())
}

func TestUploadCategoryImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemUsecase(ctrl)
	categoryHandlers := mocks.NewMockICategoryUsecase(ctrl)
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
			Value: testId.String() + "l",
		},
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, fmt.Errorf("error"))
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	categoryHandlers.EXPECT().UpdateCategory(ctx, &testCategoryWithImage).Return(fmt.Errorf("error"))
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	categoryHandlers.EXPECT().UpdateCategory(ctx, &testCategoryWithImage).Return(nil)
	delivery.UploadCategoryImage(c)
	require.Equal(t, 201, w.Code)
}

func TestDeleteCategoryImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemUsecase(ctrl)
	categoryHandlers := mocks.NewMockICategoryUsecase(ctrl)
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
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testName", testId.String()+"l"))
	filestorage.EXPECT().DeleteCategoryImage(testId.String()+"l", "testName").Return(nil)
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, fmt.Errorf("error"))
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testImagePath", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testImagePath").Return(nil)
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(&testCategoryWithImage, nil)
	categoryHandlers.EXPECT().UpdateCategory(ctx, testModelsCategoryWithId2).Return(fmt.Errorf("error"))
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.URL, _ = url.Parse(fmt.Sprintf("?id=%s&name=testImagePath", testId.String()))
	filestorage.EXPECT().DeleteCategoryImage(testId.String(), "testImagePath").Return(nil)
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(&testCategoryWithImage, nil)
	categoryHandlers.EXPECT().UpdateCategory(ctx, testModelsCategoryWithId2).Return(nil)
	delivery.DeleteCategoryImage(c)
	require.Equal(t, 200, w.Code)
}

func TestDeleteCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	itemHandlers := mocks.NewMockIItemUsecase(ctrl)
	categoryHandlers := mocks.NewMockICategoryUsecase(ctrl)
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
			Value: testId.String(),
		},
	}
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(&testEmptyModelsCategory, fmt.Errorf("error"))
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	itemHandlers.EXPECT().ItemsQuantityInCategory(ctx, testModelsCategoryWithId.Name).Return(0, nil)
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	itemHandlers.EXPECT().ItemsQuantityInCategory(ctx, testModelsCategoryWithId.Name).Return(-1, fmt.Errorf("error"))
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemHandlers.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(0, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(fmt.Errorf("error"))
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemHandlers.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(0, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(nil)
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemHandlers.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(1, nil)
	itemHandlers.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage2.Name, 0, 1).Return([]models.Item{*testModelsItemWithId}, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(nil)
	categoryHandlers.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(&testNoCategoryWithId, nil)
	itemHandlers.EXPECT().UpdateItem(ctx, &testModelsItemNoCat).Return(fmt.Errorf("error"))
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemHandlers.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(1, nil)
	itemHandlers.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage2.Name, 0, 1).Return(nil, fmt.Errorf("error"))
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemHandlers.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(1, nil)
	itemHandlers.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage2.Name, 0, 1).Return([]models.Item{*testModelsItemWithId}, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(nil)
	categoryHandlers.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(nil, fmt.Errorf("error"))
	categoryHandlers.EXPECT().CreateCategory(ctx, &testNoCategory).Return(uuid.Nil, fmt.Errorf("error"))
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testCategoryWithImage2, nil)
	itemHandlers.EXPECT().ItemsQuantityInCategory(ctx, testCategoryWithImage2.Name).Return(1, nil)
	itemHandlers.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage2.Name, 0, 1).Return([]models.Item{*testModelsItemWithId}, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(nil)
	categoryHandlers.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(nil, fmt.Errorf("error"))
	categoryHandlers.EXPECT().CreateCategory(ctx, &testNoCategory).Return(testId, nil)
	itemHandlers.EXPECT().UpdateItem(ctx, &testModelsItemNoCat).Return(fmt.Errorf("error"))
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(&testNoCategoryWithId, nil)
	delivery.DeleteCategory(c)
	require.Equal(t, 400, w.Code)

	/*w = httptest.NewRecorder()
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(nil)
	itemHandlers.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	itemHandlers.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage.Name, 0, 1).Return([]models.Item{testModelsItemWithImage}, nil)
	categoryHandlers.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(&testEmptyModelsCategory, fmt.Errorf("error"))
	categoryHandlers.EXPECT().CreateCategory(ctx, &testNoCategory).Return(testId, nil)
	itemHandlers.EXPECT().UpdateItem(ctx, &testModelsItemNoCat).Return(fmt.Errorf("error"))
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
	categoryHandlers.EXPECT().GetCategory(ctx, testId).Return(testModelsCategoryWithId, nil)
	categoryHandlers.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	filestorage.EXPECT().DeleteCategoryImageById(testId.String()).Return(nil)
	itemHandlers.EXPECT().ItemsQuantity(ctx).Return(1, nil)
	itemHandlers.EXPECT().GetItemsByCategory(ctx, testCategoryWithImage.Name, 0, 1).Return([]models.Item{testModelsItemWithImage}, nil)
	categoryHandlers.EXPECT().GetCategoryByName(ctx, "NoCategory").Return(&testNoCategoryWithId, nil)
	itemHandlers.EXPECT().UpdateItem(ctx, &testModelsItemNoCat).Return(nil)
	delivery.DeleteCategory(c)
	require.Equal(t, 200, w.Code)*/
}
