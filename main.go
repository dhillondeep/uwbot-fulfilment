package main

import (
	"net/http"
	"os"
	"warrior_bot/handlers"

	"github.com/dhillondeep/go-uw-api"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

var uwApiClient uwapi.UWAPI

// handles webhook requests from dialogflow
func handleWebhook(c *gin.Context) {
	var err error

	req := dialogflow.WebhookRequest{}
	if err = jsonpb.Unmarshal(c.Request.Body, &req); err != nil {
		logrus.WithError(err).Error("Unable to unmarshal request to jsonpb")
		c.Status(http.StatusBadRequest)
		return
	}

	resp, err := handlers.HandleRequest(&req, &uwApiClient)
	if err != nil {
		logrus.WithError(err)
	}

	marshaller := jsonpb.Marshaler{
		Indent: " ",
	}
	respStr, err := marshaller.MarshalToString(resp)
	if err != nil {
		logrus.WithError(err).Error("Unable to marshall request to JSON")
		c.Status(http.StatusBadRequest)
		return
	}

	// send response back
	c.String(http.StatusOK, respStr)
}

func main() {
	var err error

	r := gin.Default()
	r.POST("/webhook", handleWebhook)

	// create a client for Uwaterloo API
	uwApiClient = uwapi.Create(os.Getenv("UW_API_KEY"))

	if err = r.Run(); err != nil {
		logrus.WithError(err).Fatal("Couldn't start server")
	}
}
