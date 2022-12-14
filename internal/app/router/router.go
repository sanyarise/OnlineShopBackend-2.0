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
	"fmt"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	router   *gin.Engine
	delivery *delivery.Delivery
	logger   *zap.Logger
}

// NewRouter returns a new router.
func NewRouter(delivery *delivery.Delivery, logger *zap.Logger) *Router {
	logger.Debug("Enter in NewRouter()")
	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/files", "./storage/files")
	routes := Routes{
		{
			"Index",
			http.MethodGet,
			"/",
			delivery.Index,
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
			"/category/image/delete", //?id=25f32441-587a-452d-af8c-b3876ae29d45&name=20221209194557.jpeg
			delivery.DeleteCategoryImage,
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
			"/items/", //?param=categoryName&offset=20&limit=10
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
			"/items/quantity",
			delivery.ItemsQuantity,
		},
		{
			"ItemsList",
			http.MethodGet,
			"/items/list", //?offset=20&limit=10
			delivery.ItemsList,
		},
		{
			"SearchLine",
			http.MethodGet,
			"/items/search/", //?param=searchRequest&offset=20&limit=10
			delivery.SearchLine,
		},
		{
			"GetCart",
			http.MethodGet,
			"/cart/:userID",
			delivery.GetCart,
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
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}
	return &Router{router: router, delivery: delivery, logger: logger}
}

func (router *Router) Run(port string) error {
	router.logger.Debug(fmt.Sprintf("Enter in router Run(), port: %s", port))
	return router.router.Run(port)
}
