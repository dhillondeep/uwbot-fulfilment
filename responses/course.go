package responses

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"strings"
	"uwbot/helpers"
	"uwbot/models"
)

func genericCourseCarouselCardResp(title, description string, fields *models.Fields) *models.RespContext {
	return FbCarouselCard(&models.FbCarouselItem{
		Title:    title,
		Subtitle: strings.TrimSpace(description),
		Buttons: []models.FbButton{
			{
				Type:  "web_url",
				Url:   fmt.Sprintf(helpers.UWFlowCourseUrl, fields.Subject, fields.CatalogNum),
				Title: "More Info",
			},
		},
	})
}

func createCourseSectionSchedule(class *gabs.Container) string {
	startDate := class.Path("date.start_date").Data()
	endDate := class.Path("date.end_date").Data()

	startTime := class.Path("date.start_time").Data()
	endTime := class.Path("date.end_time").Data()
	weekdays := class.Path("date.weekdays").Data()

	building := class.Path("location.building").Data()
	room := class.Path("location.room").Data()

	var schStrClass string

	if startDate != nil {
		schStrClass += fmt.Sprintf("Start Date: %s\n", startDate.(string))
	}

	if endDate != nil {
		schStrClass += fmt.Sprintf("End Date: %s\n", endDate.(string))
	}

	if startTime != nil && endTime != nil && weekdays != nil {
		schStrClass += fmt.Sprintf("Timmings: %s - %s %s\n",
			startTime.(string), endTime.(string), weekdays.(string))
	}

	if building != nil && room != nil {
		schStrClass += fmt.Sprintf("Building: %s %s\n", building.(string), room.(string))
	}

	return schStrClass
}

func CourseNotFound(fields *models.Fields) *models.RespContext {
	return TextResponsef("Sorry, I am unable to find %s %s", fields.Subject, fields.CatalogNum)
}

func CourseOfferingNotFound(fields *models.Fields) *models.RespContext {
	return TextResponsef("There are no offerings for %s %s", fields.Subject, fields.CatalogNum)
}

func CoursePrereqNotFound(fields *models.Fields) *models.RespContext {
	return TextResponsef("There are no prerequisites for %s %s", fields.Subject, fields.CatalogNum)
}

func CourseAvailableInTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s is offered in %s!", fields.Subject, fields.CatalogNum, fields.Term)
}

func CourseNotAvailableInTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("Unfortunately %s %s is not offered in %s", fields.Subject, fields.CatalogNum, fields.Term)
}

func NoCourseSectionAvailable(fields *models.Fields) *models.RespContext {
	return TextResponsef("There are currently no sections for %s %s", fields.Subject, fields.CatalogNum)
}

func CourseNotOfferedNextTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s is not offered next term", fields.Subject, fields.CatalogNum)
}

func CourseOfferedNextTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s is offered next term!", fields.Subject, fields.CatalogNum)
}

func CourseNotOfferedPrevTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s was not offered last term", fields.Subject, fields.CatalogNum)
}

func CourseOfferedPrevTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s was offered last term!", fields.Subject, fields.CatalogNum)
}

func CourseNotOfferedCurrTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s is not offered this term", fields.Subject, fields.CatalogNum)
}

func CourseOfferedCurrTerm(fields *models.Fields) *models.RespContext {
	return TextResponsef("%s %s is being offered this term!", fields.Subject, fields.CatalogNum)
}

func CoursePrerequisites(prerequisites string) *models.RespContext {
	cleanPrerequisites := strings.Trim(strings.Replace(prerequisites, "Prereq: ", "", 1), ".")
	return TextResponse(cleanPrerequisites)
}

func TermsWhenCourseAvailable(terms []string, fields *models.Fields) *models.RespContext {
	return genericCourseCarouselCardResp(
		fmt.Sprintf("Terms when %s %s is offered", fields.Subject, fields.CatalogNum),
		strings.Join(terms, "\n"), fields)
}

func SectionsAvailableForCourse(sections []string, fields *models.Fields) *models.RespContext {
	return genericCourseCarouselCardResp(
		fmt.Sprintf("Sections available for %s %s", fields.Subject, fields.CatalogNum),
		strings.Join(sections, "\n"), fields)
}

func SectionInformationItem(path *gabs.Container, fields *models.Fields, section string) *models.FbCarouselItem {
	sectionInfo := createCourseSectionSchedule(path)
	return &models.FbCarouselItem{
		Title:    fmt.Sprintf("%s %s %s", fields.Subject, fields.CatalogNum, section),
		Subtitle: strings.TrimSpace(sectionInfo),
		Buttons: []models.FbButton{
			{
				Type:  "web_url",
				Url:   fmt.Sprintf(helpers.UWFlowCourseUrl, fields.Subject, fields.CatalogNum),
				Title: "More Info",
			},
		},
	}
}

func SectionEnrollmentInfoItem(enrollmentInfo string, fields *models.Fields, section string) *models.FbCarouselItem {
	return &models.FbCarouselItem{
		Title:    fmt.Sprintf("%s %s %s", fields.Subject, fields.CatalogNum, section),
		Subtitle: strings.TrimSpace(enrollmentInfo),
		Buttons: []models.FbButton{
			{
				Type:  "web_url",
				Url:   fmt.Sprintf(helpers.UWFlowCourseUrl, fields.Subject, fields.CatalogNum),
				Title: "More Info",
			},
		},
	}
}

func CourseInformation(title, desc string, fields *models.Fields) *models.RespContext {
	return genericCourseCarouselCardResp(title, desc, fields)
}
