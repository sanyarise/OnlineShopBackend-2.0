/*
 * Backend for Online Shop
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package app

import (
	"OnlineShopBackend/internal/delivery"

	"net/http"

	"github.com/gin-gonic/gin"
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
	router *gin.Engine
	del    *delivery.Delivery
}

// NewRouter returns a new router.
func NewRouter(del *delivery.Delivery) *Router {
	router := gin.Default()
	return &Router{router: router, del: del}
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

func (r *Router) Routes() {
	routes := Routes{
		{
			"Index",
			http.MethodGet,
			"/",
			Index,
		},

		{
			"CreateCategory",
			http.MethodPost,
			"/categories/:category",
			delivery.CreateCategory,
		},

		{
			"CreateItem",
			http.MethodPost,
			"/items",
			r.del.CreateItem,
		},

		{
			"GetItem",
			http.MethodGet,
			"/items/:itemID",
			delivery.GetItem,
		},

		{
			"UpdateItem",
			http.MethodPut,
			"/items/:itemID",
			r.del.UpdateItem,
		},

		{
			"UploadFile",
			http.MethodPost,
			"/items/:itemID/upload",
			delivery.UploadFile,
		},

		{
			"GetCart",
			http.MethodGet,
			"/cart/:userID",
			delivery.GetCart,
		},

		{
			"GetCategoryList",
			http.MethodGet,
			"/items/categories/:category",
			delivery.GetCategoryList,
		},

		{
			"ItemsList",
			http.MethodGet,
			"/items",
			delivery.ItemsList,
		},

		{
			"SearchLine",
			http.MethodGet,
			"/search/:searchRequest",
			delivery.SearchLine,
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
			r.router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			r.router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			r.router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			r.router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			r.router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}
}

func (r *Router) Run(port string) error {
	return r.router.Run(port)
}
