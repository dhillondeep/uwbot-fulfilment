package handlers

import (
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

	switch intentName {
	case CourseTermAvailability:
		jsonData, _ := context.UWApiClient.Courses.InfoByCatalogNumber(provSubject, provCatalogNum)

		// verify course exists
		if helpers.GetStatusCode(jsonData) != http.StatusOK {
			return responses.CourseNotFound(context.Fields), nil
		}

		var terms []string
		count, err := helpers.IterateContainerData(jsonData, "data.terms_offered", func(path *gabs.Container) error {
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
		jsonData, _ := context.UWApiClient.Courses.InfoByCatalogNumber(provSubject, provCatalogNum)

		// verify course exists
		if helpers.GetStatusCode(jsonData) != http.StatusOK {
			return responses.CourseNotFound(context.Fields), nil
		}

		offered := false
		if _, err := helpers.IterateContainerData(jsonData, "data.terms_offered", func(path *gabs.Container) error {
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

		// verify if course exists
		if helpers.GetStatusCode(jsonData) != http.StatusOK {
			return responses.CourseNotFound(context.Fields), nil
		}

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

		// verify course exists
		if helpers.GetStatusCode(jsonData) != http.StatusOK {
			return responses.CourseNotFound(context.Fields), nil
		}

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
	default:
		return nil, errors.New("handler does not exist for course intent: " + intentName)
	}
}
