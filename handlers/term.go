package handlers

import (
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"net/http"
	"strings"
	"uwbot/helpers"
	"uwbot/models"
	"uwbot/responses"
)

func HandleTermReq(context *models.ReqContext) (*models.RespContext, error) {
	intentName := context.DialogflowRequest.QueryResult.Intent.DisplayName
	fields := context.DialogflowRequest.QueryResult.Parameters.Fields

	// fetch course information if any
	courseInfo := strings.Split(fields["course"].GetStringValue(), " ")
	subject := courseInfo[0]
	catalogNum := courseInfo[1]

	// fetch section information if any
	sectionGiven := fields["section"].GetStringValue()

	termsListed, _ := context.UWApiClient.Terms.List()
	nextTerm := termsListed.Path("data.next_term").String()
	prevTerm := termsListed.Path("data.previous_term").String()
	currTerm := termsListed.Path("data.current_term").String()

	switch intentName {
	case CourseAvailNextTerm:
		nextTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(nextTerm, subject, catalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(nextTermInfo) == http.StatusNoContent {
			return responses.TextResponsef("%s %s is not offered next term", subject, catalogNum), nil
		}

		return responses.TextResponsef("yes %s %s is offered next term", subject, catalogNum), nil
	case CourseAvailPrevTerm:
		prevTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(prevTerm, subject, catalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(prevTermInfo) == http.StatusNoContent {
			return responses.TextResponsef("%s %s was not offered previous term", subject, catalogNum), nil
		}

		return responses.TextResponsef("yes %s %s was offered previous term", subject, catalogNum), nil
	case CourseAvailCurrTerm:
		currTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(currTerm, subject, catalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(currTermInfo) == http.StatusNoContent {
			return responses.TextResponsef("%s %s is not offered this term", subject, catalogNum), nil
		}

		return responses.TextResponsef("yes %s %s is offered this term", subject, catalogNum), nil
	case CourseEnrolmentInfo:
		currTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(currTerm, subject, catalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(currTermInfo) == http.StatusNoContent {
			return responses.TextResponsef("%s %s is not offered this term", subject, catalogNum), nil
		}

		var items []models.FbCarouselItem
		count, err := helpers.IterateContainerData(currTermInfo, "data", func(path *gabs.Container) error {
			courseSection := path.Path("section").Data().(string)

			if sectionGiven == "" || helpers.StringEqualNoCase(sectionGiven, courseSection) {
				enrollmentCap := path.Path("enrollment_capacity").String()
				enrollmentTotal := path.Path("enrollment_total").String()
				waitingCap := path.Path("waiting_capacity").String()
				waitingTotal := path.Path("waiting_total").String()

				subTextStr := fmt.Sprintf(
					"Enrollment Capacity: %s\nEnrollment Total: %s\nWaiting: %s/%s\n",
					enrollmentCap, enrollmentTotal, waitingCap, waitingTotal)

				items = append(items, models.FbCarouselItem{
					Title:    fmt.Sprintf("%s", courseSection),
					Subtitle: strings.TrimSpace(subTextStr),
					Buttons: []models.FbButton{
						{
							Type:  "web_url",
							Url:   fmt.Sprintf(uwflowCourseUrl, subject, catalogNum),
							Title: "More Info",
						},
					},
				})
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return responses.FbCarousel(items), nil
		} else {
			return responses.NoCourseSecFound(subject, catalogNum), nil
		}
	default:
		return nil, errors.New("handler does not exist for term intent: " + intentName)
	}
}
