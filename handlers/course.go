package handlers

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/pkg/errors"
	"net/http"
	"uwbot/helpers"
	"uwbot/models"
	"uwbot/responses"
)

func HandleCourseReq(context *models.ReqContext) (*models.RespContext, error) {
	intentName := context.DialogflowRequest.QueryResult.Intent.DisplayName
	provSubject := context.Fields.Subject
	provCatalogNum := context.Fields.CatalogNum
	provTerm := context.Fields.Term

	courseData, _ := context.UWApiClient.Courses.InfoByCatalogNumber(provSubject, provCatalogNum)

	// verify course exists
	if helpers.GetStatusCode(courseData) != http.StatusOK {
		return responses.CourseNotFound(context.Fields), nil
	}

	switch intentName {
	case CourseTermAvailability:
		var terms []string
		count, err := helpers.IterateContainerData(courseData, "data.terms_offered", func(path *gabs.Container) error {
			terms = append(terms, helpers.ConvertTermShorthandToFull(path.Data().(string)))
			return nil
		})
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return responses.TermsWhenCourseAvailable(terms, context.Fields), nil
		} else {
			return responses.CourseOfferingNotFound(context.Fields), nil
		}
	case CourseAvailabilityGivenTerm:
		offered := false
		if _, err := helpers.IterateContainerData(courseData, "data.terms_offered", func(path *gabs.Container) error {
			termInfo := path.Data().(string)
			if helpers.StringEqualNoCase(termInfo, provTerm) {
				offered = true
			}
			return nil
		}); err != nil {
			return nil, err
		}

		if offered {
			return responses.CourseAvailableInTerm(context.Fields), nil
		} else {
			return responses.CourseNotAvailableInTerm(context.Fields), nil
		}
	case CourseSections:
		jsonData, _ := context.UWApiClient.Courses.ScheduleByCatalogNumber(provSubject, provCatalogNum)

		var sections []string
		count, err := helpers.IterateContainerData(jsonData, "data", func(path *gabs.Container) error {
			sections = append(sections, path.Path("section").Data().(string))
			return nil
		})
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return responses.SectionsAvailableForCourse(sections, context.Fields), nil
		} else {
			return responses.NoCourseSectionAvailable(context.Fields), nil
		}
	case CoursePrerequisites:
		jsonData, _ := context.UWApiClient.Courses.PrereqsByCatalogNumber(provSubject, provCatalogNum)

		// verify if prerequisite exist
		if helpers.GetStatusCode(jsonData) != http.StatusOK {
			return responses.CoursePrereqNotFound(context.Fields), nil
		}

		prerequisites := jsonData.Path("data.prerequisites").Data().(string)
		return responses.TextResponse(helpers.CleanPrereqString(prerequisites)), nil
	case CourseSectionsInformation:
		provSection := context.Fields.Section
		jsonData, _ := context.UWApiClient.Courses.ScheduleByCatalogNumber(provSubject, provCatalogNum)

		var items []*models.FbCarouselItem

		// iterate over each section
		if _, err := helpers.IterateContainerData(jsonData, "data", func(path *gabs.Container) error {
			section := path.Path("section").Data().(string)

			if helpers.StringIsEmpty(provSection) || helpers.StringEqualNoCase(provSection, section) {
				// iterate over each class for that section
				if _, err := helpers.IterateContainerData(path, "classes", func(path *gabs.Container) error {
					sectionInfo := helpers.CreateCourseSectionSchedule(path)
					items = append(items, responses.SectionInformationItem(sectionInfo, context.Fields, section))
					return nil
				}); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return nil, err
		}

		if len(items) > 0 {
			return responses.FbCarousel(items), nil
		} else {
			return responses.NoCourseSectionAvailable(context.Fields), nil
		}
	case CourseAvailNextTerm:
		termsListed, _ := context.UWApiClient.Terms.List()
		nextTerm := termsListed.Path("data.next_term").String()

		nextTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(nextTerm, provSubject, provCatalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(nextTermInfo) == http.StatusNoContent {
			return responses.CourseNotOfferedNextTerm(context.Fields), nil
		}

		return responses.CourseOfferedNextTerm(context.Fields), nil
	case CourseAvailPrevTerm:
		termsListed, _ := context.UWApiClient.Terms.List()
		prevTerm := termsListed.Path("data.previous_term").String()

		prevTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(prevTerm, provSubject, provCatalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(prevTermInfo) == http.StatusNoContent {
			return responses.CourseNotOfferedPrevTerm(context.Fields), nil
		}

		return responses.CourseOfferedPrevTerm(context.Fields), nil
	case CourseAvailCurrTerm:
		termsListed, _ := context.UWApiClient.Terms.List()
		currTerm := termsListed.Path("data.current_term").String()

		currTermInfo, _ := context.UWApiClient.Terms.ClassScheduleByTerm(currTerm, provSubject, provCatalogNum)

		// verify course info exists for next term
		if helpers.GetStatusCode(currTermInfo) == http.StatusNoContent {
			return responses.CourseNotOfferedCurrTerm(context.Fields), nil
		}

		return responses.CourseOfferedCurrTerm(context.Fields), nil
	case CourseEnrolmentInfo:
		termsListed, _ := context.UWApiClient.Terms.List()
		currTerm := termsListed.Path("data.current_term").String()

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

				items = append(items, responses.SectionEnrollmentInfoItem(subTextStr, context.Fields, courseSection))
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
	case CourseInformation:
		title := fmt.Sprintf("%s %s: %s", context.Fields.Subject,
			context.Fields.CatalogNum, courseData.Path("data.title").Data().(string))
		desc := courseData.Path("data.description").Data().(string)

		return responses.CourseInformation(title, desc, context.Fields), nil

	default:
		return nil, errors.New("handler does not exist for course intent: " + intentName)
	}
}
