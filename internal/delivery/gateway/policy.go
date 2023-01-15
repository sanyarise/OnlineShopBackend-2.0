package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type PolicyGateway interface {
	Ask(*gin.Context) bool
}

type opaRequest struct {
	// Input wraps the OPA request (https://www.openpolicyagent.org/docs/latest/rest-api/#get-a-document-with-input)
	Input *opaInput `json:"input"`
}

type opaInput struct {
	// The token of the requester
	Token string `json:"token"`
	// The path to which the request was made split to an array
	Path []string `json:"path"`
	// The HTTP Method
	Method string `json:"method"`
}

type opaResponse struct {
	Result bool `json:"result"`
}

// PolicyOpaGateway makes policy decision requests to OPA
type PolicyOpaGateway struct {
	endpoint  string
	secretKey string
	logger    *zap.Logger
}

func NewPolicyOpaGateway(endpoint string, secretKey string, logger *zap.Logger) *PolicyOpaGateway {
	return &PolicyOpaGateway{
		endpoint:  endpoint,
		secretKey: secretKey,
		logger:    logger,
	}
}

// Ask requests to OPA with required inputs and returns the decision made by OPA
func (gateway *PolicyOpaGateway) Ask(c *gin.Context) bool {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		gateway.logger.Error(err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return false
	}
	gateway.logger.Sugar().Debugf("tokenString read from request success: %v", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}
		return []byte(gateway.secretKey), nil
	})
	if err != nil {
		gateway.logger.Error(err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return false
	}
	gateway.logger.Sugar().Debugf("Token parse from tokenstring success: %v", token)

	// After splitting, the first element isn't necessary
	// "/finance/salary/alice" -> ["", "finance", "salary", "alice"]
	paths := strings.Split(c.Request.URL.RequestURI(), "/")[1:]

	method := c.Request.Method

	// create input to send to OPA
	input := &opaInput{
		Token:  token.Raw,
		Path:   paths,
		Method: method,
	}
	opaRequest := &opaRequest{
		Input: input,
	}

	gateway.logger.Sugar().Debugf("token: %v, path: %v, method: %v", input.Token, input.Path, input.Method)

	requestBody, err := json.Marshal(opaRequest)
	if err != nil {
		gateway.logger.Sugar().Errorf("error on json marshalling: %v", err)
		return false
	}

	// request OPA
	resp, err := http.Post(gateway.endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		gateway.logger.Sugar().Errorf("PDP request failed, err: %v", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		gateway.logger.Sugar().Errorf("Reading body of response failed, err: %v", err)
		return false
	}

	var opaResponse opaResponse
	err = json.Unmarshal(body, &opaResponse)
	if err != nil {
		gateway.logger.Sugar().Errorf("Unmarshalling response body failed, err: %v", err)
		return false
	}
	gateway.logger.Sugar().Infof("result: %v", opaResponse.Result)
	return opaResponse.Result
}
