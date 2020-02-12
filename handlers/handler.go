package handlers

import (
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"uwbot/helpers"
	"uwbot/models"
)

func HandleWebhook(context *models.ReqContext) (*models.RespContext, error) {
	request := context.DialogflowRequest

	// intents are in form NUM_CATEGORY_NAME
	// we are getting category of the intent
	intentCat := strings.Split(request.QueryResult.Intent.DisplayName, "_")[1]

	// we already have fulfilment text provided to us so we shouldn't do anything
	if !helpers.StringIsEmpty(request.QueryResult.FulfillmentText) {
		return &models.RespContext{
			StatusCode: http.StatusOK,
		}, nil
	}

	switch intentCat {
	case CourseIntent:
		return HandleCourseReq(context)
	case TermIntent:
		return HandleTermReq(context)
	default:
		return nil, errors.New("handler does not exist for intent category: " + intentCat)
	}
}
