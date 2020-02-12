package handlers

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"uwbot/helpers"
	"uwbot/models"
	"uwbot/responses"
)

const uwflowCourseUrl = "https://uwflow.com/course/%s%s"

func HandleCourseReq(context *models.ReqContext) (*models.RespContext, error) {
	intentName := context.DialogflowRequest.QueryResult.Intent.DisplayName
	fields := context.DialogflowRequest.QueryResult.Parameters.Fields

	// fetch course information if any
	courseInfo := strings.Split(fields["course"].GetStringValue(), " ")
	subject := courseInfo[0]
	catalogNum := courseInfo[1]

	// fetch term information if any
	termName := fields["term"].GetStringValue()

	// fetch section information if any
	sectionGiven := fields["section"].GetStringValue()

	switch intentName {
	case CourseTermAvailability:
		jsonData, _ := context.UWApiClient.Courses.InfoByCatalogNumber(subject, catalogNum)

		// verify course exists
		if helpers.GetStatusCode(jsonData) != http.StatusOK {
			return responses.CourseNotFoundError(subject, catalogNum), nil
		}

		var termsStr string
		count, err := helpers.IterateContainerData(jsonData, "data.terms_offered", func(path *gabs.Container) error {
			termsStr += fmt.Sprintf("%s\n", helpers.ConvertTermShorthandToFull(path.Data().(string)))
			return nil
		})
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return responses.FbCarouselCard(models.FbCarouselItem{
				Title:    fmt.Sprintf("Terms when %s %s is offered", subject, catalogNum),
				Subtitle: strings.TrimSpace(termsStr),
				Buttons: []models.FbButton{
					{
						Type:  "web_url",
						Url:   fmt.Sprintf(uwflowCourseUrl, subject, catalogNum),
						Title: "More Info",
					},
				},
			}), nil
		} else {
			return responses.CourseOfferingError(subject, catalogNum), nil
		}
	case CourseAvailabilityGivenTerm:
		jsonData, _ := context.UWApiClient.Courses.InfoByCatalogNumber(subject, catalogNum)

		// verify course exists
		if helpers.GetStatusCode(jsonData) != http.StatusOK {
			return responses.CourseNotFoundError(subject, catalogNum), nil
		}

		offered := false
		if _, err := helpers.IterateContainerData(jsonData, "data.terms_offered", func(path *gabs.Container) error {
			termInfo := path.Data().(string)
			if helpers.StringEqualNoCase(termInfo, termName) {
				offered = true
			}
			return nil
		}); err != nil {
			return nil, err
		}

		if offered {
			return responses.TextResponsef("%s %s is available in %s term!", subject, catalogNum, termName), nil
		} else {
			return responses.TextResponsef("%s %s is not available in %s term!", subject, catalogNum, termName), nil
		}
	case CourseSections:
		jsonData, _ := context.UWApiClient.Courses.ScheduleByCatalogNumber(subject, catalogNum)

		// verify course exists
		if helpers.GetStatusCode(jsonData) != http.StatusOK {
			return responses.CourseNotFoundError(subject, catalogNum), nil
		}

		var sectionsStr string
		count, err := helpers.IterateContainerData(jsonData, "data", func(path *gabs.Container) error {
			section := path.Path("section").Data().(string)
			sectionsStr += fmt.Sprintf("%s\n", section)
			return nil
		})
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return responses.FbCarouselCard(models.FbCarouselItem{
				Title:    fmt.Sprintf("Sections available for %s %s", subject, catalogNum),
				Subtitle: strings.TrimSpace(sectionsStr),
				Buttons: []models.FbButton{
					{
						Type:  "web_url",
						Url:   fmt.Sprintf(uwflowCourseUrl, subject, catalogNum),
						Title: "More Info",
					},
				},
			}), nil
		} else {
			return responses.TextResponsef("No sections are available for %s %s", subject, catalogNum), nil
		}
	case CoursePrerequisites:
		jsonData, _ := context.UWApiClient.Courses.PrereqsByCatalogNumber(subject, catalogNum)

		// verify course exists
		if helpers.GetStatusCode(jsonData) != http.StatusOK {
			return responses.CourseNotFoundError(subject, catalogNum), nil
		}

		prerequisites, ok := jsonData.Path("data.prerequisites").Data().(string)
		if !ok {
			return responses.CoursePrerequisitesNotFound(subject, catalogNum), nil
		}

		return responses.TextResponse(strings.Trim(
			strings.Replace(prerequisites, "Prereq: ", "", 1), ".")), nil
	case CourseSectionSchedule:
		jsonData, _ := context.UWApiClient.Courses.ScheduleByCatalogNumber(subject, catalogNum)

		// verify course exists
		if helpers.GetStatusCode(jsonData) != http.StatusOK {
			return responses.CourseNotFoundError(subject, catalogNum), nil
		}

		var items []models.FbCarouselItem

		// iterate over each section
		if _, err := helpers.IterateContainerData(jsonData, "data", func(path *gabs.Container) error {
			section := path.Path("section").Data().(string)

			if sectionGiven == "" || helpers.StringEqualNoCase(sectionGiven, section) {
				// iterate over each class for that section
				if _, err := helpers.IterateContainerData(path, "classes", func(path *gabs.Container) error {
					infoStr := helpers.CreateCourseSectionSchedule(path)
					items = append(items, models.FbCarouselItem{
						Title:    section,
						Subtitle: strings.TrimSpace(infoStr),
						Buttons: []models.FbButton{
							{
								Type:  "web_url",
								Url:   fmt.Sprintf(uwflowCourseUrl, subject, catalogNum),
								Title: "More Info",
							},
						},
					})

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
			return responses.CourseSectionInfoNotFound(subject, catalogNum, sectionGiven), nil
		}
	default:
		return nil, errors.New("handler does not exist for course intent: " + intentName)
	}
}
