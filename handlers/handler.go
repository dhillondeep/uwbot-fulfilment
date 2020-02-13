package handlers

import (
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"uwbot/helpers"
	"uwbot/models"
)

func fetchAndCreateFields(context *models.ReqContext) {
	dialogflowFields := context.DialogflowRequest.QueryResult.Parameters.Fields

	context.Fields = &models.Fields{}

	// fetch course event if any
	helpers.DoIfFieldsContains(dialogflowFields, "course", func(s string) {
		context.Fields.Subject = helpers.CourseSubjectReg.FindString(s)
		context.Fields.CatalogNum = helpers.CourseCatalogReg.FindString(s)
	})

	// fetch term event if any
	helpers.DoIfFieldsContains(dialogflowFields, "term", func(s string) {
		context.Fields.Term = s
	})

	// fetch section event if any
	helpers.DoIfFieldsContains(dialogflowFields, "section", func(s string) {
		context.Fields.Section = s
	})
}

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

	// parse and store dialogflow fields in context
	fetchAndCreateFields(context)

	switch intentCat {
	case CourseIntent:
		return HandleCourseReq(context)
	default:
		return nil, errors.New("handler does not exist for intent category: " + intentCat)
	}
}
