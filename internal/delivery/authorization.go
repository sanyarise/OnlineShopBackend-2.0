package delivery

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EnforcePolicy is a middleware to check all the requests and ask if the requester has permissions
func (delivery *Delivery) Authorize(c *gin.Context) {
	delivery.logger.Info("Enforcing policy with middleware")

	allow := delivery.policyOpaGateway.Ask(c)

	if allow {
		delivery.logger.Info("Action is allowed, continuing process")
		c.Next()
	} else {
		delivery.logger.Info("Action was not allowed, cancelling process")
		c.AbortWithStatus(http.StatusForbidden)
	}
}

type PolicyGateway interface {
	Ask(*gin.Context) bool
}

type opaRequest struct {
	// Input wraps the OPA request (https://www.openpolicyagent.org/docs/latest/rest-api/#get-a-document-with-input)
	Input *opaInput `json:"input"`
}

type opaInput struct {
	Endpoint string `json:"endpoint"`
	// Users's role
	Role string `json:"role"`
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

func NewPolicyOpaGateway(endpoint string, secretKey string, logger *zap.Logger) PolicyGateway {
	return &PolicyOpaGateway{
		endpoint:  endpoint,
		secretKey: secretKey,
		logger:    logger,
	}
}

// Ask requests to OPA with required inputs and returns the decision made by OPA
func (gateway *PolicyOpaGateway) Ask(c *gin.Context) bool {
	gateway.logger.Debug("Enter in gateway Ask()")
	endpoint, _ := c.Get("endpoint")
	role, _ := c.Get("role")

	// create input to send to OPA
	input := &opaInput{
		Endpoint: endpoint.(string),
		Role:     role.(string),
	}
	opaRequest := &opaRequest{
		Input: input,
	}

	gateway.logger.Sugar().Debugf("endpoint: %v, role: %v", input.Endpoint, input.Role)

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
	gateway.logger.Sugar().Debugf("response: %v", opaResponse)
	gateway.logger.Sugar().Infof("result: %v", opaResponse.Result)
	return opaResponse.Result
}
