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
	//gin.Use(cors.Default())
	gin.Use(CORSMiddleware())
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:3000"}
	//config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	//config.AllowHeaders = []string{"Authorization"}
	//gin.Use(cors.New(config))

	//gin.Use(cors.New(cors.Config{
	//	AllowOrigins: []string{"https://accounts.google.com", "https://accounts.google.com/o/oauth2/auth?", "http://localhost:8000", "http://localhost:3000", "http://localhost:8000/user/login/google", "*"}, //,
	//	AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "*"},
	//	AllowHeaders:     []string{"Origin", "Authorization", "Credentials", "*"},
	//	ExposeHeaders:    []string{"Content-Length", "*"},
	//	AllowCredentials: true,
	//}))

	gin.Use(ginzap.RecoveryWithZap(logger, true))
	gin.Static("/files", "./static/files")
	docs.SwaggerInfo.BasePath = "/"
	gin.Group("/docs").Any("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router := &Router{
		delivery: delivery,
		logger:   logger,
	}

	routes := Routes{
		/*{
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
		},*/
		/*{
			"CreateCategory",
			http.MethodPost,
			"/categories/create",
			delivery.CreateCategory,
		},*/
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
		/*{
			"UpdateCategory",
			http.MethodPut,
			"/categories/:categoryID",
			delivery.UpdateCategory,
		},*/
		/*{
			"UploadCategoryImage",
			http.MethodPost,
			"/categories/image/upload/:categoryID",
			delivery.UploadCategoryImage,
		},*/
		/*{
			"DeleteCategoryImage",
			http.MethodDelete,
			"/categories/image/delete", //?id=25f32441-587a-452d-af8c-b3876ae29d45&name=20221209194557.jpeg
			delivery.DeleteCategoryImage,
		},*/
		/*{
			"DeleteCategory",
			http.MethodDelete,
			"/categories/delete/:categoryID",
			delivery.DeleteCategory,
		},*/
		/*{
			"CreateItem",
			http.MethodPost,
			"/items/create",
			delivery.CreateItem,
		},*/
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
		/*{
			"UpdateItem",
			http.MethodPut,
			"/items/update",
			delivery.UpdateItem,
		},*/
		/*{
			"UploadItemImage",
			http.MethodPost,
			"/items/image/upload/:itemID",
			delivery.UploadItemImage,
		},*/
		/*{
			"DeleteItemImage",
			http.MethodDelete,
			"/items/image/delete", //?id=25f32441-587a-452d-af8c-b3876ae29d45&name=20221209194557.jpeg
			delivery.DeleteItemImage,
		},*/
		{
			"ItemsQuantity",
			http.MethodGet,
			"/items/quantity", //?categoryName={categoryName} - if category name is empty returns quantity of all items
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
		/*{
			"DeleteItem",
			http.MethodDelete,
			"/items/delete/:itemID",
			delivery.DeleteItem,
		},*/
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
			"/cart/deleteItem",
			delivery.DeleteItemFromCart,
		},
		{
			"DeleteCart",
			http.MethodDelete,
			"/cart/delete/:cartID",
			delivery.DeleteCart,
		},
		/*{
			"CreateRights",
			http.MethodPost,
			"/rights/create",
			delivery.CreateRights,
		},
		{
			"UpdateRights",
			http.MethodPut,
			"/rights/update",
			delivery.UpdateRights,
		},
		{
			"DeleteRights",
			http.MethodDelete,
			"/rights/delete/:rightsID",
			delivery.DeleteRights,
		},*/
		{
			"GetRights",
			http.MethodGet,
			"/rights/:rightsID",
			delivery.GetRights,
		},
		/*{
			"RightsList",
			http.MethodGet,
			"/rights/list",
			delivery.RightsList,
		},*/
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
			http.MethodGet,
			"/user/logout",
			delivery.LogoutUser,
		},
		{
			"LoginUserGoogle",
			http.MethodGet,
			"/user/login/google",
			delivery.LoginUserGoogle,
		},

		{
			"LoginUserYandex",
			http.MethodGet,
			"/user/login/yandex",
			delivery.LoginUserYandex,
		},

		{
			"callbackGoogle",
			http.MethodGet,
			"/user/callbackGoogle",
			delivery.CallbackGoogle,
		},

		{
			"callbackYandex",
			http.MethodPost,
			"/user/callbackYandex",
			delivery.CallbackYandex,
		},
		/*{
			"userProfile",
			http.MethodGet,
			"/user/profile",
			delivery.UserProfile,
		},*/
		/*{
			"userProfileUpdate",
			http.MethodPut,
			"/user/profile/edit",
			delivery.UserProfileUpdate,
		},*/
		{
			"tokenUpdate",
			http.MethodPost,
			"/user/token/update",
			delivery.TokenUpdate,
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
	gin.GET("/", delivery.Authentificate, delivery.Authorize, delivery.Index)
	gin.GET("/images/list", delivery.Authentificate, delivery.Authorize, delivery.GetFileList)
	gin.POST("/categories/create", delivery.Authentificate, delivery.Authorize, delivery.CreateCategory)
	gin.PUT("/categories/:categoryID", delivery.Authentificate, delivery.Authorize, delivery.UpdateCategory)
	gin.POST("/categories/image/upload/:categoryID", delivery.Authentificate, delivery.Authorize, delivery.UploadCategoryImage)
	gin.DELETE("/categories/image/delete", delivery.Authentificate, delivery.Authorize, delivery.DeleteCategoryImage) //?id=25f32441-587a-452d-af8c-b3876ae29d45&name=20221209194557.jpeg
	gin.DELETE("/categories/delete/:categoryID", delivery.Authentificate, delivery.Authorize, delivery.DeleteCategory)
	gin.POST("/items/create", delivery.Authentificate, delivery.Authorize, delivery.CreateItem)
	gin.PUT("/items/update", delivery.Authentificate, delivery.Authorize, delivery.UpdateItem)
	gin.POST("/items/image/upload/:itemID", delivery.Authentificate, delivery.Authorize, delivery.UploadItemImage)
	gin.DELETE("/items/image/delete", delivery.Authentificate, delivery.Authorize, delivery.DeleteItemImage) //?id=25f32441-587a-452d-af8c-b3876ae29d45&name=20221209194557.jpeg
	gin.DELETE("/items/delete/:itemID", delivery.Authentificate, delivery.Authorize, delivery.DeleteItem)
	gin.POST("/rights/create", delivery.Authentificate, delivery.Authorize, delivery.CreateRights)
	gin.PUT("/rights/update", delivery.Authentificate, delivery.Authorize, delivery.UpdateRights)
	gin.DELETE("/rights/delete/:rightsID", delivery.Authentificate, delivery.Authorize, delivery.DeleteRights)
	gin.GET("/rights/list", delivery.Authentificate, delivery.Authorize, delivery.GetRights)
	gin.GET("/user/profile", delivery.Authentificate, delivery.Authorize, delivery.UserProfile)
	gin.PUT("/user/profile/edit", delivery.Authentificate, delivery.Authorize, delivery.UserProfileUpdate)
	gin.GET("/user/:userID", delivery.Authentificate, delivery.Authorize, delivery.GetUserById)
	gin.GET("/user/list", delivery.Authentificate, delivery.Authorize, delivery.GetUsersList)
	gin.PUT("/user/update/rights", delivery.Authentificate, delivery.Authorize, delivery.ChangeUserRole)
	gin.PUT("/user/update/password", delivery.Authentificate, delivery.Authorize, delivery.ChangeUserPassword)

	router.Engine = gin
	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		//c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Auth-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
