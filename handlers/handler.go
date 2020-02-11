package handlers

import (
	"errors"
	"github.com/dhillondeep/go-uw-api"
	"github.com/golang/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"io"
	"strings"
)

func handleRequest(req *dialogflow.WebhookRequest, uwClient *uwapi.UWAPI) (*dialogflow.WebhookResponse, error) {
	intentCat := strings.Split(req.QueryResult.Intent.DisplayName, "_")[1]

	// we already have fulfilment text provides so we shouldn't do anything
	if strings.Trim(req.QueryResult.FulfillmentText, " ") != "" {
		return nil, nil
	}

	switch intentCat {
	case CourseIntent:
		return HandleCourseReq(req.QueryResult, uwClient)
	case TermIntent:
		return HandleTermReq(req.QueryResult, uwClient)
	default:
		return nil, errors.New("handler does not exist for intent category: " + intentCat)
	}
}

func HandleWebhook(request io.Reader, uwApiClient *uwapi.UWAPI) (string, error) {
	var err error

	req := dialogflow.WebhookRequest{}
	if err = jsonpb.Unmarshal(request, &req); err != nil {
		logrus.WithError(err).Error("Unable to unmarshal request to jsonpb")
		return "nil", err
	}

	resp, err := handleRequest(&req, uwApiClient)
	if err != nil {
		logrus.WithError(err)
	}

	marshaller := jsonpb.Marshaler{
		Indent: " ",
	}

	respStr, err := marshaller.MarshalToString(resp)
	if err != nil {
		logrus.WithError(err).Error("Unable to marshall request to JSON")
		return "", err
	}

	return respStr, nil
}
