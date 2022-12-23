package delivery

import (
	"OnlineShopBackend/internal/filestorage"
	"OnlineShopBackend/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//	@title			Online Shop Backend Service
//	@version		1.0
//	@description	Backend service for online store
//	@license.name	MIT

//	@contact.url	https://github.com/GBteammates/OnlineShopBackend

//	@BasePath	/

type Delivery struct {
	itemHandlers     handlers.IItemHandlers
	categoryHandlers handlers.ICategoryHandlers
	logger           *zap.Logger
	filestorage      filestorage.FileStorager
}

type FilesInfo struct {
	Name       string `json:"name" example:"20221213125935.jpeg"`
	Path       string `json:"path" example:"storage\\files\\categories\\d0d3df2d-f6c8-4956-9d76-998ee1ec8a39\\20221213125935.jpeg"`
	CreateDate string `json:"created_date" example:"2022-12-13 12:46:16.0964549 +0300 MSK"`
	ModifyDate string `json:"modify_date" example:"2022-12-13 12:46:16.0964549 +0300 MSK"`
}

type FileListResponse struct {
	Files []FilesInfo `json:"files"`
}

// NewDelivery initialize delivery layer
func NewDelivery(itemHandlers handlers.IItemHandlers, categoryHandlers handlers.ICategoryHandlers, logger *zap.Logger, fs filestorage.FileStorager) *Delivery {
	logger.Debug("Enter in NewDelivery()")
	return &Delivery{itemHandlers: itemHandlers, categoryHandlers: categoryHandlers, logger: logger, filestorage: fs}
}

// Index is the index handler.
func (delivery *Delivery) Index(c *gin.Context) {
	delivery.logger.Debug("Enter in Index")
	c.String(http.StatusOK, "Hello World!")
}

// GetFileList returns list of files
//
//	@Summary		Get list of files
//	@Description	Method provides to get list of files.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	FileListResponse	"List of files"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/images/list [get]
func (delivery *Delivery) GetFileList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetFileList()")
	fileInfos, err := delivery.filestorage.GetFileList()
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	var files FileListResponse
	files.Files = make([]FilesInfo, len(fileInfos))
	for i, info := range fileInfos {
		files.Files[i] = FilesInfo{
			Name:       info.Name,
			Path:       info.Path,
			CreateDate: info.CreateDate,
			ModifyDate: info.ModifyDate,
		}
	}

	c.JSON(http.StatusOK, files)
}
