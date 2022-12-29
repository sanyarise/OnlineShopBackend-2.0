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
	"OnlineShopBackend/internal/delivery/category"
	"OnlineShopBackend/internal/delivery/item"
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type QuantityOptions struct {
	CategoryName string `form:"categoryName"`
}

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

// CreateItem
//
//	@Summary		Method provides to create store item
//	@Description	Method provides to create store item
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			item	body		item.ShortItem	true	"Data for creating item"
//	@Success		201		{object}	item.ItemId
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/items/create/ [post]
func (delivery *Delivery) CreateItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CreateItem()")
	ctx := context.Background()
	var deliveryItem item.ShortItem
	if err := c.ShouldBindJSON(&deliveryItem); err != nil {
		delivery.logger.Error(fmt.Sprintf("error on bind json from request: %v", err))
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	if deliveryItem.Title == "" || deliveryItem.Description == "" || deliveryItem.Price == 0 {
		delivery.logger.Error(fmt.Errorf("empty item fields in request").Error())
		delivery.SetError(c, http.StatusBadRequest, fmt.Errorf("empty item fields in request"))
		return
	}

	if deliveryItem.Category == "" {
		noCategory, err := delivery.categoryUsecase.GetCategoryByName(ctx, "NoCategory")
		if err != nil {
			delivery.logger.Sugar().Errorf("NoCategory is not exists: %v", err)
			noCategory := models.Category{
				Name:        "NoCategory",
				Description: "Category for items without categories",
			}
			noCategoryId, err := delivery.categoryUsecase.CreateCategory(ctx, &noCategory)
			if err != nil {
				delivery.logger.Error(fmt.Sprintf("error on create no category: %v", err))
				delivery.SetError(c, http.StatusInternalServerError, err)
				return
			}
			deliveryItem.Category = noCategoryId.String()
		} else {
			deliveryItem.Category = noCategory.Id.String()
		}
	}
	categoryId, err := uuid.Parse(deliveryItem.Category)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	modelsItem := models.Item{
		Title:       deliveryItem.Title,
		Description: deliveryItem.Description,
		Price:       deliveryItem.Price,
		Category: models.Category{
			Id: categoryId,
		},
		Vendor: deliveryItem.Vendor,
		Images: deliveryItem.Images,
	}

	id, err := delivery.itemUsecase.CreateItem(ctx, &modelsItem)
	if err != nil {
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, item.ItemId{Value: id.String()})
}

// GetItem - returns item by id
//
//	@Summary		Get item by id
//	@Description	The method allows you to get the product by id.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			itemID	path		string			true	"id of item"
//	@Success		200		{object}	item.OutItem	"Item structure"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/items/{itemID} [get]
func (delivery *Delivery) GetItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetItem()")
	id := c.Param("itemID")
	if id == "" {
		err := fmt.Errorf("empty item in request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(id)

	uid, err := uuid.Parse(id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	modelsItem, err := delivery.itemUsecase.GetItem(ctx, uid)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, item.OutItem{
		Id:          modelsItem.Id.String(),
		Title:       modelsItem.Title,
		Description: modelsItem.Description,
		Category: category.Category{
			Id:          modelsItem.Category.Id.String(),
			Name:        modelsItem.Category.Name,
			Description: modelsItem.Category.Description,
			Image:       modelsItem.Category.Image,
		},
		Price:  modelsItem.Price,
		Vendor: modelsItem.Vendor,
		Images: modelsItem.Images,
	})
}

// UpdateItem - update an item
//
//	@Summary		Method provides to update store item
//	@Description	Method provides to update store item
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			item	body	item.InItem	true	"Data for updating item"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/update [put]
func (delivery *Delivery) UpdateItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UpdateItem()")
	ctx := c.Request.Context()
	var deliveryItem item.InItem
	if err := c.ShouldBindJSON(&deliveryItem); err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(deliveryItem.Id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	categoryUid, err := uuid.Parse(deliveryItem.Category)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	err = delivery.itemUsecase.UpdateItem(ctx, &models.Item{
		Id:          uid,
		Title:       deliveryItem.Title,
		Description: deliveryItem.Description,
		Category: models.Category{
			Id: categoryUid,
		},
		Price:  deliveryItem.Price,
		Vendor: deliveryItem.Vendor,
		Images: deliveryItem.Images,
	})
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ItemsList - returns list of all items
//
//	@Summary		Get list of items
//	@Description	Method provides to get list of items
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int				false	"Quantity of recordings"		default(10)	minimum(0)
//	@Param			offset	query		int				false	"Offset when receiving records"	default(0)	mininum(0)
//	@Success		200		array		item.OutItem	"List of items"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/items/list [get]
func (delivery *Delivery) ItemsList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsList()")
	var options Options
	err := c.Bind(&options)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("options is %v", options))
	ctx := c.Request.Context()
	if options.Limit == 0 {
		quantity, err := delivery.itemUsecase.ItemsQuantity(ctx)
		if err != nil {
			delivery.logger.Error(err.Error())
		}
		if quantity == 0 {
			delivery.logger.Debug("quantity of items is 0")
			c.JSON(http.StatusOK, item.ItemsList{})
			return
		}
		if quantity <= 30 && quantity > 0 {
			options.Limit = quantity
		} else {
			options.Limit = 10
		}
	}
	delivery.logger.Sugar().Debugf("options limit is set in default value: %d", options.Limit)

	list, err := delivery.itemUsecase.ItemsList(ctx, options.Offset, options.Limit)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	items := make([]item.OutItem, len(list))
	for idx, modelsItem := range list {
		items[idx] = item.OutItem{
			Id:          modelsItem.Id.String(),
			Title:       modelsItem.Title,
			Description: modelsItem.Description,
			Category: category.Category{
				Id:          modelsItem.Category.Id.String(),
				Name:        modelsItem.Category.Name,
				Description: modelsItem.Category.Description,
				Image:       modelsItem.Category.Image,
			},
			Price:  modelsItem.Price,
			Vendor: modelsItem.Vendor,
			Images: modelsItem.Images,
		}
	}
	c.JSON(http.StatusOK, items)
}

// ItemsQuantity returns quantity of all items
//
//	@Summary		Get quantity of items
//	@Description	Method provides to get quantity of items
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	item.ItemsQuantity	"Quantity of items"
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/quantity [get]
func (delivery *Delivery) ItemsQuantity(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsQuantity()")
	var options QuantityOptions
	err := c.Bind(&options)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("options is %v", options))
	ctx := c.Request.Context()
	if options.CategoryName == "" {
		quantity, err := delivery.itemUsecase.ItemsQuantity(ctx)
		if err != nil {
			delivery.logger.Error(err.Error())
			delivery.SetError(c, http.StatusInternalServerError, err)
			return
		}
		itemsQuantity := item.ItemsQuantity{Quantity: quantity}
		c.JSON(http.StatusOK, itemsQuantity)
	} else {
		quantity, err := delivery.itemUsecase.ItemsQuantityInCategory(ctx, options.CategoryName)
		if err != nil {
			delivery.logger.Error(err.Error())
			delivery.SetError(c, http.StatusInternalServerError, err)
			return
		}
		itemsQuantity := item.ItemsQuantity{Quantity: quantity}
		c.JSON(http.StatusOK, itemsQuantity)
	}
}

// SearchLine - returns list of items with parameters
//
//	@Summary		Get list of items by search parameters
//	@Description	Method provides to get list of items by search parameters
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			param	query		string			false	"Search param"
//	@Param			limit	query		int				false	"Quantity of recordings"		default(10)	minimum(0)
//	@Param			offset	query		int				false	"Offset when receiving records"	default(0)	mininum(0)
//	@Success		200		array		item.OutItem	"List of items"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/items/search [get]
func (delivery *Delivery) SearchLine(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery SearchLine()")
	var options SearchOptions
	err := c.Bind(&options)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("options is %v", options))
	if options.Param == "" {
		err = fmt.Errorf("empty search request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	if options.Limit == 0 {
		options.Limit = 10
	}

	delivery.logger.Sugar().Debugf("options limit is set in default value: %d", options.Limit)
	ctx := c.Request.Context()
	list, err := delivery.itemUsecase.SearchLine(ctx, options.Param, options.Offset, options.Limit)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	items := make([]item.OutItem, len(list))
	for idx, modelsItem := range list {
		items[idx] = item.OutItem{
			Id:          modelsItem.Id.String(),
			Title:       modelsItem.Title,
			Description: modelsItem.Description,
			Category: category.Category{
				Id:          modelsItem.Category.Id.String(),
				Name:        modelsItem.Category.Name,
				Description: modelsItem.Category.Description,
				Image:       modelsItem.Category.Image,
			},
			Price:  modelsItem.Price,
			Vendor: modelsItem.Vendor,
			Images: modelsItem.Images,
		}
	}
	c.JSON(http.StatusOK, items)
}

// GetItemsByCategory returns list of items in category
//
//	@Summary		Get list of items by category name
//	@Description	Method provides to get list of items by category name
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			param	query		string			false	"Category name"
//	@Param			limit	query		int				false	"Quantity of recordings"		default(10)	minimum(0)
//	@Param			offset	query		int				false	"Offset when receiving records"	default(0)	mininum(0)
//	@Success		200		array		item.OutItem	"List of items"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/items [get]
func (delivery *Delivery) GetItemsByCategory(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetItemsByCategory()")
	var options SearchOptions
	err := c.Bind(&options)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("options is %v", options))
	if options.Param == "" {
		err = fmt.Errorf("empty search request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	if options.Limit == 0 {
		options.Limit = 10
	}
	delivery.logger.Sugar().Debugf("options limit is set in default value: %d", options.Limit)

	ctx := c.Request.Context()
	list, err := delivery.itemUsecase.GetItemsByCategory(ctx, options.Param, options.Offset, options.Limit)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	items := make([]item.OutItem, len(list))
	for idx, modelsItem := range list {
		items[idx] = item.OutItem{
			Id:          modelsItem.Id.String(),
			Title:       modelsItem.Title,
			Description: modelsItem.Description,
			Category: category.Category{
				Id:          modelsItem.Category.Id.String(),
				Name:        modelsItem.Category.Name,
				Description: modelsItem.Category.Description,
				Image:       modelsItem.Category.Image,
			},
			Price:  modelsItem.Price,
			Vendor: modelsItem.Vendor,
			Images: modelsItem.Images,
		}
	}
	c.JSON(http.StatusOK, items)
}

// UploadItemImage - upload an image
//
//	@Summary		Upload an image of item
//	@Description	Method provides to upload an image of item
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"id of item"
//	@Param			image	formData	file	true	"picture of item"
//	@Success		201
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		415	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Failure		507	{object}	ErrorResponse
//	@Router			/items/image/upload/:itemID [post]
func (delivery *Delivery) UploadItemImage(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UploadItemImage()")
	ctx := c.Request.Context()
	id := c.Param("itemID")
	if id == "" {
		err := fmt.Errorf("empty search request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	var name string
	contentType := c.ContentType()

	if contentType == "image/jpeg" {
		name = carbon.Now().ToShortDateTimeString() + ".jpeg"
	} else if contentType == "image/png" {
		name = carbon.Now().ToShortDateTimeString() + ".png"
	} else {
		err := fmt.Errorf("unsupported media type: %s", contentType)
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusUnsupportedMediaType, err)
		return
	}

	file, err := io.ReadAll(c.Request.Body)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusUnsupportedMediaType, err)
		return
	}

	delivery.logger.Info("Read id", zap.String("id", id))
	delivery.logger.Info("File len=", zap.Int32("len", int32(len(file))))
	path, err := delivery.filestorage.PutItemImage(id, name, file)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInsufficientStorage, err)
		return
	}

	item, err := delivery.itemUsecase.GetItem(ctx, uid)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	item.Images = append(item.Images, path)

	err = delivery.itemUsecase.UpdateItem(ctx, item)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// DeleteItemImage delete an item image
//
//	@Summary		Delete an item image by item id
//	@Description	The method allows you to delete an item image by item id.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			id		query	string	true	"Item id"
//	@Param			name	query	string	true	"Image name"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/image/delete [delete]
func (delivery *Delivery) DeleteItemImage(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteItemImage()")
	var imageOptions ImageOptions
	err := c.Bind(&imageOptions)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	delivery.logger.Debug(fmt.Sprintf("image options is %v", imageOptions))

	if imageOptions.Id == "" || imageOptions.Name == "" {
		err := fmt.Errorf("empty image options in request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(imageOptions.Id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
	}
	err = delivery.filestorage.DeleteItemImage(imageOptions.Id, imageOptions.Name)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return

	}
	ctx := c.Request.Context()
	item, err := delivery.itemUsecase.GetItem(ctx, uid)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	for idx, imagePath := range item.Images {
		if strings.Contains(imagePath, imageOptions.Name) {
			item.Images = append(item.Images[:idx], item.Images[idx+1:]...)
			break
		}
	}
	err = delivery.itemUsecase.UpdateItem(ctx, item)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteItem deleted item by id
//
//	@Summary		Method provides to delete item
//	@Description	Method provides to delete item.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			itemID	path	string	true	"id of item"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/delete/{itemID} [delete]
func (delivery *Delivery) DeleteItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteItem()")
	id := c.Param("itemID")
	if id == "" {
		err := fmt.Errorf("empty item id in request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	deletedItem, err := delivery.itemUsecase.GetItem(ctx, uid)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("deletedItem: %v", deletedItem))

	err = delivery.itemUsecase.DeleteItem(ctx, uid)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	err = delivery.itemUsecase.UpdateItemsInCategoryCash(ctx, deletedItem, "delete")
	if err != nil {
		delivery.logger.Sugar().Errorf("error on update cash in category items list: %v", err)
	}

	if len(deletedItem.Images) > 0 {
		err = delivery.filestorage.DeleteItemImagesFolderById(id)
		if err != nil {
			delivery.logger.Error(err.Error())
		}
	}
	delivery.logger.Sugar().Infof("Item with id: %s deleted success", id)
	c.JSON(http.StatusOK, gin.H{})
}
