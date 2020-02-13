package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dhillondeep/go-uw-api"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"io"
	"os"
	"strings"
	"uwbot/handlers"
	"uwbot/models"
)

var uwApiClient uwapi.UWAPI

func createReqContext(reader io.Reader) (*models.ReqContext, error) {
	req := &dialogflow.WebhookRequest{}
	if err := jsonpb.Unmarshal(reader, req); err != nil {
		return nil, errors.Wrap(err, "Unable to unmarshal request to jsonpb")
	}

	return &models.ReqContext{
		UWApiClient:       &uwApiClient,
		DialogflowRequest: req,
	}, nil
}

// handles webhook requests from dialogflow using gin
func ginHandler(c *gin.Context) {
	reqContext, err := createReqContext(c.Request.Body)
	if err != nil {
		logrus.Error(err)
	}

	respContext, err := handlers.HandleWebhook(reqContext)
	if err != nil {
		logrus.Error(err)
	}

	// send response back
	c.JSON(respContext.StatusCode, respContext.Resp)
}

// handles webhook requests from dialogflow using AWS Lambda
func lambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqContext, err := createReqContext(strings.NewReader(request.Body))
	if err != nil {
		logrus.Error(err)
	}

	respContext, err := handlers.HandleWebhook(reqContext)
	if err != nil {
		logrus.Error(err)
	}

	// marshall resp context to dialogflow compatible json
	bodyStr, err := jsoniter.MarshalToString(respContext.Resp)
	if err != nil {
		logrus.Error(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: respContext.StatusCode,
		Body:       bodyStr,
	}, nil
}

func main() {
	// create a client for UWaterloo API
	uwApiClient = uwapi.Create(os.Getenv("UW_API_KEY"))

	// if AWS_LAMBDA_USE env is set, run the application serverless on AWS lambda
	// else use gin web server to create a webhook
	if _, exists := os.LookupEnv("AWS_LAMBDA_USE"); exists {
		lambda.Start(lambdaHandler)
	} else {
		r := gin.Default()
		r.POST("/webhook", ginHandler)

		if err := r.Run(); err != nil {
			logrus.WithError(err).Fatal("Couldn't start webhook server")
		}
	}
}
