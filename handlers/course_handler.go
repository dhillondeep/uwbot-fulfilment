package handlers

import (
	"errors"
	"fmt"
	"github.com/dhillondeep/go-uw-api"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"strings"
)

func HandleCourseReq(qResult *dialogflow.QueryResult, uwClient *uwapi.UWAPI) (*dialogflow.WebhookResponse, error) {
	intentName := qResult.Intent.DisplayName
	fields := qResult.Parameters.Fields

	// fetch course information
	courseInfo := strings.Split(fields["course"].GetStringValue(), " ")
	subject := courseInfo[0]
	catalogNum := courseInfo[1]

	// fetch term information if any
	//termName := fields["term"].GetStringValue()

	respStr := ""
	found := false

	switch intentName {
	case CourseAvailabilityTerm:
		jsonData, _ := uwClient.Courses.InfoByCatalogNumber(subject, catalogNum)
		respStr, found = courseNotFoundError(getStatusCode(jsonData), subject, catalogNum)
		if !found {
			break
		}

		numTermsOffered, _ := jsonData.ArrayCountP("data.terms_offered")

		if numTermsOffered > 0 {
			respStr = "Terms Offered: "
			for i := 0; i < numTermsOffered; i++ {
				termInfo := jsonData.Path(fmt.Sprintf("data.terms_offered.%d", i)).Data().(string)

				respStr += fmt.Sprintf("%s, ", convertTermShorthandToFull(termInfo))
			}

			respStr = trimRespString(respStr)
		}

		break
	case CourseAvailabilityCurrTerm:
		break
	case CourseScheduleCurrTerm:
		break
	case CoursePrerequisites:
		jsonData, _ := uwClient.Courses.PrereqsByCatalogNumber(subject, catalogNum)
		respStr, found = courseNotFoundError(getStatusCode(jsonData), subject, catalogNum)
		if !found {
			break
		}

		prerequisites, ok := jsonData.Path("data.prerequisites").Data().(string)
		if !ok {
			respStr = fmt.Sprintf("I was unable to find prerequisites for %s %s",
				subject, catalogNum)
			break
		}

		respStr = prerequisites
		break
	default:
		return nil, errors.New("handler does not exist for course intent: " + intentName)
	}

	return &dialogflow.WebhookResponse{
		FulfillmentText: respStr,
	}, nil
}
