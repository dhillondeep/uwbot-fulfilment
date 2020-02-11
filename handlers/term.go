package handlers

import (
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/dhillondeep/go-uw-api"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"net/http"
	"strings"
	"warrior_bot/helpers"
	"warrior_bot/models"
	"warrior_bot/responses"
)

func HandleTermReq(qResult *dialogflow.QueryResult, uwClient *uwapi.UWAPI) (*dialogflow.WebhookResponse, error) {
	intentName := qResult.Intent.DisplayName
	fields := qResult.Parameters.Fields

	// fetch course information if any
	courseInfo := strings.Split(fields["course"].GetStringValue(), " ")
	subject := courseInfo[0]
	catalogNum := courseInfo[1]

	termsListed, _ := uwClient.Terms.List()
	nextTerm := termsListed.Path("data.next_term").String()
	prevTerm := termsListed.Path("data.previous_term").String()
	currTerm := termsListed.Path("data.current_term").String()

	switch intentName {
	case CourseAvailNextTerm:
		nextTermInfo, _ := uwClient.Terms.ClassScheduleByTerm(nextTerm, subject, catalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(nextTermInfo) == http.StatusNoContent {
			return responses.TextResponsef("%s %s is not offered next term", subject, catalogNum)
		}

		return responses.TextResponsef("yes %s %s is offered next term", subject, catalogNum)
	case CourseAvailPrevTerm:
		prevTermInfo, _ := uwClient.Terms.ClassScheduleByTerm(prevTerm, subject, catalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(prevTermInfo) == http.StatusNoContent {
			return responses.TextResponsef("%s %s was not offered previous term", subject, catalogNum)
		}

		return responses.TextResponsef("yes %s %s was offered previous term", subject, catalogNum)
	case CourseAvailCurrTerm:
		currTermInfo, _ := uwClient.Terms.ClassScheduleByTerm(currTerm, subject, catalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(currTermInfo) == http.StatusNoContent {
			return responses.TextResponsef("%s %s is not offered this term", subject, catalogNum)
		}

		return responses.TextResponsef("yes %s %s is offered this term", subject, catalogNum)
	case CourseEnrolmentInfo:
		currTermInfo, _ := uwClient.Terms.ClassScheduleByTerm(currTerm, subject, catalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(currTermInfo) == http.StatusNoContent {
			return responses.TextResponsef("%s %s is not offered this term", subject, catalogNum)
		}

		var items []models.FbCarouselItem
		count, err := helpers.IterateContainerData(currTermInfo, "data", func(path *gabs.Container) error {
			courseSection := path.Path("section").Data().(string)
			enrollmentCap := path.Path("enrollment_capacity").String()
			enrollmentTotal := path.Path("enrollment_total").String()
			waitingCap := path.Path("waiting_capacity").String()
			waitingTotal := path.Path("waiting_total").String()

			subTextStr := fmt.Sprintf(
				"Enrollment Capacity: %s\nEnrollment Total: %s\nWaiting: %s/%s\n",
				enrollmentCap, enrollmentTotal, waitingCap, waitingTotal)

			items = append(items, models.FbCarouselItem{
				Title:    fmt.Sprintf("%s %s %s", subject, catalogNum, courseSection),
				Subtitle: strings.TrimSpace(subTextStr),
				Buttons: []models.FbButton{
					{
						Type:  "web_url",
						Url:   fmt.Sprintf(uwflowCourseUrl, subject, catalogNum),
						Title: "More Info",
					},
				},
			})

			return nil
		})
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return responses.FbCarousel(items)
		} else {
			return responses.CourseSectionsNotFound(subject, catalogNum)
		}
	default:
		return nil, errors.New("handler does not exist for term intent: " + intentName)
	}
}
