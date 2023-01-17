/*
 * Backend for Online Shop
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package router

import (
	"OnlineShopBackend/internal/delivery"
	"OnlineShopBackend/internal/delivery/swagger/docs"

	"net/http"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

type Router struct {
	*gin.Engine
	delivery *delivery.Delivery
	logger   *zap.Logger
}

// NewRouter returns a new router.
func NewRouter(delivery *delivery.Delivery, logger *zap.Logger) *Router {
	logger.Debug("Enter in NewRouter()")
	gin := gin.Default()
	gin.Use(cors.Default())
	gin.Use(ginzap.RecoveryWithZap(logger, true))
	gin.Static("/files", "./static/files")
	docs.SwaggerInfo.BasePath = "/"
	gin.Group("/docs").Any("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router := &Router{
		delivery: delivery,
		logger:   logger,
	}

	routes := Routes{
		{
			"Index",
			http.MethodGet,
			"/",
			delivery.Index,
		},
		{
			"GetFileList",
			http.MethodGet,
			"/images/list",
			delivery.GetFileList,
		},
		{
			"CreateCategory",
			http.MethodPost,
			"/categories/create",
			delivery.CreateCategory,
		},
		{
			"GetCategory",
			http.MethodGet,
			"/categories/:categoryID",
			delivery.GetCategory,
		},
		{
			"GetCategoryList",
			http.MethodGet,
			"/categories/list",
			delivery.GetCategoryList,
		},
		{
			"UpdateCategory",
			http.MethodPut,
			"/categories/:categoryID",
			delivery.UpdateCategory,
		},
		{
			"UploadCategoryImage",
			http.MethodPost,
			"/categories/image/upload/:categoryID",
			delivery.UploadCategoryImage,
		},
		{
			"DeleteCategoryImage",
			http.MethodDelete,
			"/categories/image/delete", //?id=25f32441-587a-452d-af8c-b3876ae29d45&name=20221209194557.jpeg
			delivery.DeleteCategoryImage,
		},
		{
			"DeleteCategory",
			http.MethodDelete,
			"/categories/delete/:categoryID",
			delivery.DeleteCategory,
		},
		{
			"CreateItem",
			http.MethodPost,
			"/items/create",
			delivery.CreateItem,
		},
		{
			"GetItem",
			http.MethodGet,
			"/items/:itemID",
			delivery.GetItem,
		},
		{
			"GetItemsByCategory",
			http.MethodGet,
			"/items/", //?param=categoryName&offset=20&limit=10&sort_type=name&sort_order=asc (sort_type == name or price, sort_order == asc or desc)
			delivery.GetItemsByCategory,
		},
		{
			"UpdateItem",
			http.MethodPut,
			"/items/update",
			delivery.UpdateItem,
		},
		{
			"UploadItemImage",
			http.MethodPost,
			"/items/image/upload/:itemID",
			delivery.UploadItemImage,
		},
		{
			"DeleteItemImage",
			http.MethodDelete,
			"/items/image/delete", //?id=25f32441-587a-452d-af8c-b3876ae29d45&name=20221209194557.jpeg
			delivery.DeleteItemImage,
		},
		{
			"ItemsQuantity",
			http.MethodGet,
			"/items/quantity", //?categoryName={categoryName} - if category name is empty returns quantity of all items
			delivery.ItemsQuantity,
		},
		{
			"ItemsList",
			http.MethodGet,
			"/items/list", //?offset=20&limit=10&sort_type=name&sort_order=asc (sort_type == name or price, sort_order == asc or desc)
			delivery.ItemsList,
		},
		{
			"SearchLine",
			http.MethodGet,
			"/items/search/", //?param=searchRequest&offset=20&limit=10&sort_type=name&sort_order=asc (sort_type == name or price, sort_order == asc or desc)
			delivery.SearchLine,
		},
		{
			"DeleteItem",
			http.MethodDelete,
			"/items/delete/:itemID",
			delivery.DeleteItem,
		},
		{
			"GetCart",
			http.MethodGet,
			"/cart/:cartID",
			delivery.GetCart,
		},
		{
			"CreateCart",
			http.MethodPost,
			"/cart/create/:userID",
			delivery.CreateCart,
		},
		{
			"AddItemToCart",
			http.MethodPut,
			"/cart/addItem",
			delivery.AddItemToCart,
		},
		{
			"DeleteItemFromCart",
			http.MethodDelete,
			"/cart/delete/:cartID/:itemID",
			delivery.DeleteItemFromCart,
		},
		{
			"DeleteCart",
			http.MethodDelete,
			"/cart/delete/:cartID",
			delivery.DeleteCart,
		},
		{
			"CreateUser",
			http.MethodPost,
			"/user/create",
			delivery.CreateUser,
		},

		{
			"LoginUser",
			http.MethodPost,
			"/user/login",
			delivery.LoginUser,
		},

		{
			"LogoutUser",
			http.MethodPost,
			"/user/logout",
			delivery.LogoutUser,
		},
	}

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			gin.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			gin.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			gin.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			gin.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			gin.DELETE(route.Pattern, route.HandlerFunc)
		}
	}
	router.Engine = gin
	return router
}
