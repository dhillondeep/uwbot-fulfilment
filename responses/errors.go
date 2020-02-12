package responses

import (
	"uwbot/models"
)

/********* Course Errors *********/

func CourseNotFoundError(subject, catalogNum string) *models.RespContext {
	respStr := "Sorry, I was unable to find this information about %s %s at University of Waterloo"
	return TextResponsef(respStr, subject, catalogNum)
}

func CourseOfferingError(subject, catalogNum string) *models.RespContext {
	respStr := "Sorry, I was unable to find any offerings for %s %s"
	return TextResponsef(respStr, subject, catalogNum)
}

func CoursePrerequisitesNotFound(subject, catalogNum string) *models.RespContext {
	respStr := "Sorry, I was unable to find information about prerequisites for %s %s"
	return TextResponsef(respStr, subject, catalogNum)
}

func CourseSectionInfoNotFound(subject, catalogNum, section string) *models.RespContext {
	respStr := "Sorry, I was unable to find information about %s %s %s"
	return TextResponsef(respStr, subject, catalogNum, section)
}

/********* Term Errors *********/

func CourseSectionsNotFound(subject, catalogNum string) *models.RespContext {
	respStr := "Sorry, I wasn't able to find any sections for %s %s"
	return TextResponsef(respStr, subject, catalogNum)
}
