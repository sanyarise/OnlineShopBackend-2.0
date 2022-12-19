/*
 * Backend for Online Shop
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package delivery

import (
	"OnlineShopBackend/internal/handlers"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"go.uber.org/zap"
)

type Options struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type SearchOptions struct {
	Param  string `form:"param"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

type ImageOptions struct {
	Id   string `form:"id"`
	Name string `form:"name"`
}

type ItemsQuantity struct {
	Quantity int `json:"quantity"`
}

type DeliveryItem struct {
	Id          string   `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       int32    `json:"price,omitempty"`
	Category    string   `json:"category,omitempty"`
	Vendor      string   `json:"vendor,omitempty"`
	Images      []string `json:"image,omitempty"`
}

// CreateItem - create a new item
func (delivery *Delivery) CreateItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CreateItem()")
	ctx := c.Request.Context()
	var deliveryItem DeliveryItem
	if err := c.ShouldBindJSON(&deliveryItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if deliveryItem.Title == "" && deliveryItem.Description == "" && deliveryItem.Category == "" && deliveryItem.Price == 0 && deliveryItem.Vendor == "" {
		c.JSON(http.StatusBadRequest, "empty item is not correct")
		return
	}
	item := handlers.Item{
		Id:          deliveryItem.Id,
		Title:       deliveryItem.Title,
		Description: deliveryItem.Description,
		Price:       deliveryItem.Price,
		Category: handlers.Category{
			Id: deliveryItem.Category,
		},
		Vendor: deliveryItem.Vendor,
		Images: deliveryItem.Images,
	}
	id, err := delivery.itemHandlers.CreateItem(ctx, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": id.String()})
}

// GetItem - returns item on id
func (delivery *Delivery) GetItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetItem()")
	id := c.Param("itemID")
	delivery.logger.Debug(id)
	ctx := c.Request.Context()
	item, err := delivery.itemHandlers.GetItem(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// UpdateItem - update an item
func (delivery *Delivery) UpdateItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UpdateItem()")
	ctx := c.Request.Context()
	var deliveryItem handlers.Item
	if err := c.ShouldBindJSON(&deliveryItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := delivery.itemHandlers.UpdateItem(ctx, deliveryItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ItemsList - returns list of all items
func (delivery *Delivery) ItemsList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsList()")
	var options Options
	err := c.Bind(&options)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery.logger.Debug(fmt.Sprintf("options is %v", options))

	if options.Limit == 0 {
		quantity, err := delivery.itemHandlers.ItemsQuantity(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if quantity < 100 {
			options.Limit = quantity
		} else {
			options.Limit = 20
		}
	}

	delivery.logger.Debug("options limit is set in default value")

	list, err := delivery.itemHandlers.ItemsList(c.Request.Context(), options.Offset, options.Limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ItemsQuantity returns quantity of all items
func (delivery *Delivery) ItemsQuantity(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsQuantity()")
	ctx := c.Request.Context()
	quantity, err := delivery.itemHandlers.ItemsQuantity(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	itemsQuantity := ItemsQuantity{Quantity: quantity}
	c.JSON(http.StatusOK, itemsQuantity)
}

// SearchLine - returns list of items with parameters
func (delivery *Delivery) SearchLine(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery SearchLine()")

	var options SearchOptions
	err := c.Bind(&options)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if options.Param == "" {
		delivery.logger.Sugar().Error("empty search request")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	delivery.logger.Debug(fmt.Sprintf("options is %v", options))

	if options.Limit == 0 {
		options.Limit = 10
	}

	delivery.logger.Debug("options limit is set in default value: 10")
	list, err := delivery.itemHandlers.SearchLine(c.Request.Context(), options.Param, options.Offset, options.Limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetItemsByCategory returns list of items in category
func (delivery *Delivery) GetItemsByCategory(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetItemsByCategory()")
	var options SearchOptions
	err := c.Bind(&options)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if options.Param == "" {
		delivery.logger.Sugar().Error("empty category name")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	delivery.logger.Debug(fmt.Sprintf("options is %v", options))

	if options.Limit == 0 {
		options.Limit = 10
		delivery.logger.Debug("options limit is set in default value")
	}



	items, err := delivery.itemHandlers.GetItemsByCategory(c.Request.Context(), options.Param, options.Offset, options.Limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// UploadItemImage - upload an image
func (delivery *Delivery) UploadItemImage(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UploadItemImage()")
	ctx := c.Request.Context()
	id := c.Param("itemID")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty item id"})
		return
	}
	var name string
	contentType := c.ContentType()

	if contentType == "image/jpeg" {
		name = carbon.Now().ToShortDateTimeString() + ".jpeg"
	} else if contentType == "image/png" {
		name = carbon.Now().ToShortDateTimeString() + ".png"
	} else {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{})
		return
	}

	file, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{})
		return
	}

	delivery.logger.Info("Read id", zap.String("id", id))
	delivery.logger.Info("File len=", zap.Int32("len", int32(len(file))))
	path, err := delivery.filestorage.PutItemImage(id, name, file)
	if err != nil {
		c.JSON(http.StatusInsufficientStorage, gin.H{})
		return
	}

	item, err := delivery.itemHandlers.GetItem(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item.Images = append(item.Images, path)

	err = delivery.itemHandlers.UpdateItem(ctx, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "upload image success"})
}

// DeleteItemImage delete an item image
func (delivery *Delivery) DeleteItemImage(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteItemImage()")
	var imageOptions ImageOptions
	err := c.Bind(&imageOptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery.logger.Debug(fmt.Sprintf("image options is %v", imageOptions))

	if imageOptions.Id == "" || imageOptions.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("empty item id or file name")})
		return
	}
	err = delivery.filestorage.DeleteItemImage(imageOptions.Id, imageOptions.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	ctx := c.Request.Context()
	item, err := delivery.itemHandlers.GetItem(ctx, imageOptions.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for idx, imagePath := range item.Images {
		if strings.Contains(imagePath, imageOptions.Name) {
			item.Images = append(item.Images[:idx], item.Images[idx+1:]...)
			break
		}
	}
	err = delivery.itemHandlers.UpdateItem(ctx, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "delete image success"})
}
