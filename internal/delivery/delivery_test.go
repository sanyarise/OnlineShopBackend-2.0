package delivery

import (
	fstor "OnlineShopBackend/internal/filestorage"
	fs "OnlineShopBackend/internal/filestorage/mocks"
	"OnlineShopBackend/internal/handlers/mocks"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
	delivery.Index(c)
	require.Equal(t, 200, w.Code)
}

func TestGetFileList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
	
	filestorage.EXPECT().GetFileList().Return([]fstor.FileInStorageInfo{}, fmt.Errorf("error"))
	delivery.GetFileList(c)
	require.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}
	
	info := new(fstor.FileInStorageInfo)
	info.Name="testName"
	res := []fstor.FileInStorageInfo{*info}
	resBytes, _ := json.Marshal(&res)
	filestorage.EXPECT().GetFileList().Return(res, nil)
	delivery.GetFileList(c)
	require.Equal(t, 200, w.Code)
	require.Equal(t, resBytes, w.Body.Bytes())
}