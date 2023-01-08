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
	"OnlineShopBackend/internal/delivery/user"
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase"
	"github.com/dghubble/sessions"
	"net/http"
	//"github.com/dghubble/sessions"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"

	"golang.org/x/oauth2"
	og2 "golang.org/x/oauth2/google"
	yandex "golang.org/x/oauth2/yandex"
	//"golang.org/x/oauth2/yandex"
)

// CreateUser create a new user
//
//		@Summary		Create a new user
//		@Description	Method provides to create a user
//		@Tags			user
//		@Accept			json
//		@Produce		json
//	 	@Param			user	body	user.Credentials	true	"User data" //TODO
//		@Success		200	{object} user.Token
//		@Failure		400	"Bad Request"
//		@Failure		404	{object}	ErrorResponse	"404 Not Found"
//		@Failure		500	{object}	ErrorResponse
//		@Router			/user/create [post]
func (delivery *Delivery) CreateUser(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CreateUser()")
	ctx := c.Request.Context()
	var newUser *models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check user in database
	if existedUser, err := delivery.userUsecase.GetUserByEmail(ctx, newUser.Email, user.GeneratePasswordHash(newUser.Password)); err == nil && existedUser.ID != uuid.Nil {
		c.JSON(http.StatusContinue, gin.H{"error": err.Error()})
		return
	}

	// Password validation check
	if err := user.ValidationCheck(*newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashPassword := user.GeneratePasswordHash(newUser.Password)
	newUser.Password = hashPassword

	// Create a user
	createdUser, err := delivery.userUsecase.CreateUser(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	delivery.logger.Info("success: user was created")

	token, err := delivery.userUsecase.CreateSessionJWT(c.Request.Context(), createdUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//user.IssueSession(delivery.logger, createdUser.ID.String())

	c.JSON(http.StatusOK, token)
}

// LoginUser login user
//
//		@Summary		Login user
//		@Description	Method provides to login a user
//		@Tags			user
//		@Accept			json
//		@Produce		json
//	 @Param			user	body	user.Credentials	true	"User data"
//		@Success		200	{object} user.Token //TODO
//		@Failure		404	"Bad Request"
//		@Failure		404	{object}	ErrorResponse	"404 Not Found"
//		@Failure		500	{object}	ErrorResponse
//		@Router			/user/login [post]
func (delivery *Delivery) LoginUser(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LoginUser()")
	var userCredentials usecase.Credentials
	if err := c.ShouldBindJSON(&userCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userExist, err := delivery.userUsecase.GetUserByEmail(c.Request.Context(), userCredentials.Email, user.GeneratePasswordHash(userCredentials.Password))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if userExist.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error56": err.Error()})
		return
	}

	token, err := delivery.userUsecase.CreateSessionJWT(c.Request.Context(), &userExist)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

// UserProfile user profile
//
//	@Summary		User profile
//	@Description	Method provides to get profile info
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security	    ApiKeyAuth || firebase
//	@Success		200	{object} user.Profile //TODO
//	@Failure		404	"Bad Request"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/user/profile [get]
func (delivery *Delivery) UserProfile(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UserProfile()")
	header := c.GetHeader("Authorization")
	userCr, err := delivery.userUsecase.UserIdentity(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	userData, err := delivery.userUsecase.GetUserByEmail(c.Request.Context(), userCr.Email, userCr.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	userProfile := usecase.Profile{
		Email:     userData.Email,
		FirstName: userData.Firstname,
		LastName:  userData.Lastname,
		Address: usecase.Address{
			Zipcode: userData.Address.Zipcode,
			Country: userData.Address.Country,
			City:    userData.Address.City,
			Street:  userData.Address.Street,
		},
		Rights: usecase.Rights{
			ID:    userData.Rights.ID,
			Name:  userData.Rights.Name,
			Rules: userData.Rights.Rules,
		},
	}
	 c.JSON(http.StatusCreated, userProfile)
}

func (delivery *Delivery) UserProfileUpdate(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UserProfileUpdate()")
	header := c.GetHeader("Authorization")
	userCr, err := delivery.userUsecase.UserIdentity(header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var newInfoUser *models.User
	if err = c.ShouldBindJSON(&newInfoUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser := models.User{
		ID:        userCr.UserId,
		Firstname: newInfoUser.Firstname,
		Lastname:  newInfoUser.Lastname,
		Address:   models.UserAddress{
			Zipcode: newInfoUser.Address.Zipcode,
			Country: newInfoUser.Address.Country,
			City:    newInfoUser.Address.City,
			Street:  newInfoUser.Address.Street,
		},
	}

	userUpdated, err := delivery.userUsecase.UpdateUserData(c.Request.Context(), &updatedUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, userUpdated)
}

func (delivery *Delivery) TokenUpdate(c *gin.Context)  {

}

// LoginUserGoogle -
func (delivery *Delivery) LoginUserGoogle(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LogoutUser()")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
	//c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	oauth2Config := &oauth2.Config{
		ClientID:     "435235643575-7g5u2gfhfrhgm3e2mtv682ev5ch54k7q.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8000/user/callbackGoogle",
		//RedirectURL: "http://localhost:3000",
		Endpoint: og2.Endpoint,
		Scopes:   []string{"profile", "email"},
	}

	url := oauth2Config.AuthCodeURL("random")
	//http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
	c.Redirect(http.StatusTemporaryRedirect, url)
	//stateConfig := gologin.DefaultCookieConfig
	//google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueSession(), nil)).ServeHTTP(c.Writer, c.Request)
}

// LoginUserYandex -
func (delivery *Delivery) LoginUserYandex(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LoginUserYandex()")
	//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
	//c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	oauth2ConfigYandex := &oauth2.Config{
		ClientID:     "c0ba17e8f61d47fdb0a978131d1c2d48",
		ClientSecret: "2e41418551724d918919ac09d4f6a1eb",
		Endpoint:     yandex.Endpoint,
		RedirectURL:  "http://localhost:8000/user/callbackYandex",
		Scopes:       []string{"email"},
	}
	c.Redirect(http.StatusTemporaryRedirect, oauth2ConfigYandex.AuthCodeURL("random")) //

	delivery.logger.Debug(yandex.Endpoint.TokenURL)
	c.JSON(http.StatusOK, yandex.Endpoint)
}

// CallbackGoogle -
func (delivery *Delivery) CallbackGoogle(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CallbackGoogle()")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
	//c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	//oauth2Config := &oauth2.Config{
	//	ClientID:     "614400740650-ioroeqq2rvn45k5tv5rc8noa7058m1l9.apps.googleusercontent.com",
	//	ClientSecret: "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",
	//	RedirectURL:  "http://localhost:8000/user/callbackGoogle",
	//	Endpoint:     og2.Endpoint,
	//	Scopes:       []string{"profile", "email"},
	//}
	//stateConfig := gologin.DefaultCookieConfig
	//google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueSession(), nil))
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000")
}

// CallbackYandex -
func (delivery *Delivery) CallbackYandex(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CallbackYandex()")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
	//c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	//c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	//stateConfig := gologin.DefaultCookieConfig
	//gologinoauth2.StateHandler(stateConfig, )
	//c.JSON(http.StatusOK, gin.H{})

}

const (
	sessionName    = "example-google-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "key"
	sessionUserID  = "userId"
)

var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

// LogoutUser logout
//
//	@Summary		Logout
//	@Description	Method provides to log out
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		404	"Bad Request"
//	@Router			/user/logout [get]
func (delivery *Delivery) LogoutUser(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery LogoutUser()")
	sessionStore.Destroy(c.Writer, sessionUserID)
	c.SetCookie(sessionName + "-tmp", "", 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"you have been successfully logged out": nil})
}

