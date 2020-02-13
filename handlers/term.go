package handlers

import (
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"net/http"
	"uwbot/helpers"
	"uwbot/models"
	"uwbot/responses"
)

func HandleTermReq(context *models.ReqContext) (*models.RespContext, error) {
	intentName := context.DialogflowRequest.QueryResult.Intent.DisplayName
	provSubject := context.Fields.Subject
	provCatalogNum := context.Fields.CatalogNum

	termsListed, _ := context.UWApiClient.Terms.List()
	nextTerm := termsListed.Path("data.next_term").String()
	prevTerm := termsListed.Path("data.previous_term").String()
	currTerm := termsListed.Path("data.current_term").String()

	switch intentName {
	case CourseAvailNextTerm:
		nextTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(nextTerm, provSubject, provCatalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(nextTermInfo) == http.StatusNoContent {
			return responses.CourseNotOfferedNextTerm(context.Fields), nil
		}

		return responses.CourseOfferedNextTerm(context.Fields), nil
	case CourseAvailPrevTerm:
		prevTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(prevTerm, provSubject, provCatalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(prevTermInfo) == http.StatusNoContent {
			return responses.CourseNotOfferedPrevTerm(context.Fields), nil
		}

		return responses.CourseOfferedPrevTerm(context.Fields), nil
	case CourseAvailCurrTerm:
		currTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(currTerm, provSubject, provCatalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(currTermInfo) == http.StatusNoContent {
			return responses.CourseNotOfferedCurrTerm(context.Fields), nil
		}

		return responses.CourseOfferedCurrTerm(context.Fields), nil
	case CourseEnrolmentInfo:
		sectionGiven := context.Fields.Section
		currTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(currTerm, provSubject, provCatalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(currTermInfo) == http.StatusNoContent {
			return responses.CourseNotOfferedCurrTerm(context.Fields), nil
		}

		var items []*models.FbCarouselItem
		count, err := helpers.IterateContainerData(currTermInfo, "data", func(path *gabs.Container) error {
			courseSection := path.Path("section").Data().(string)

			if helpers.StringIsEmpty(sectionGiven) || helpers.StringEqualNoCase(sectionGiven, courseSection) {
				enrollmentCap := path.Path("enrollment_capacity").String()
				enrollmentTotal := path.Path("enrollment_total").String()
				waitingCap := path.Path("waiting_capacity").String()
				waitingTotal := path.Path("waiting_total").String()

				subTextStr := fmt.Sprintf("Enrollment Capacity: %s\nEnrollment Total: %s\nWaiting: %s/%s\n",
					enrollmentCap, enrollmentTotal, waitingCap, waitingTotal)

				items = append(items, responses.SectionEnrollmentInfoItem(subTextStr, context.Fields))
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return responses.FbCarousel(items), nil
		} else {
			return responses.NoCourseSectionAvailable(context.Fields), nil
		}
	default:
		return nil, errors.New("handler does not exist for term intent: " + intentName)
	}
}
