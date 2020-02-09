package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// handles webhook requests from dialogflow
func handleWebhook(c *gin.Context) {
	var err error

	request := dialogflow.WebhookRequest{}
	if err = jsonpb.Unmarshal(c.Request.Body, &request); err != nil {
		logrus.WithError(err).Error("Unable to unmarshal request to jsonpb")
		c.Status(http.StatusBadRequest)
		return
	}

	fmt.Println(request.GetQueryResult().GetOutputContexts())
}

func main() {
	var err error

	r := gin.Default()
	r.POST("/webhook", handleWebhook)

	if err = r.Run(); err != nil {
		logrus.WithError(err).Fatal("Couldn't start server")
	}
}
