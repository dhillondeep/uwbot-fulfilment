package helpers

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"regexp"
	"strings"
)

// used to provide more info on cards
const UWFlowCourseUrl = "https://uwflow.com/course/%s%s"

var CourseSubjectReg = regexp.MustCompile(`[A-Za-z]+`)       // regex to get course subject substring
var CourseCatalogReg = regexp.MustCompile(`[0-9]+[A-Za-z]*`) // regex to get course catalog number substring

var termsShortHand = map[string]string{"w": "Winter", "s": "Spring", "f": "Fall"}

// IterateContainerData iterates over container provided using the path provided and for each
// element, calls callback function providing with element container path
func IterateContainerData(data *gabs.Container, path string, callback func(path *gabs.Container) error) (int, error) {
	numItems, err := data.ArrayCountP(path)
	if err != nil {
		return -1, err
	}

	for i := 0; i < numItems; i++ {
		if err := callback(data.Path(fmt.Sprintf("%s.%d", path, i))); err != nil {
			return numItems, err
		}
	}

	return numItems, nil
}

// GetStatusCode looks into UW API response and finds status code
func GetStatusCode(container *gabs.Container) float64 {
	return container.Path("meta.status").Data().(float64)
}

// ConvertTermShorthandToFull converts UW API shot hand notation for terms to full form
func ConvertTermShorthandToFull(shorthand string) string {
	return termsShortHand[strings.ToLower(shorthand)]
}

// StringEqualNoCase compares two strings without case consideration
func StringEqualNoCase(first, second string) bool {
	return strings.ToLower(strings.TrimSpace(first)) == strings.ToLower(strings.TrimSpace(second))
}

// StringIsEmpty returns true or false based on if string is empty
func StringIsEmpty(text string) bool {
	return strings.TrimSpace(text) == ""
}

// DoIfFieldsContains check if key provided exists inside Dialogflow fields and if it does,
// it calls the provided function with value retrieved from the fields map
func DoIfFieldsContains(fields map[string]*structpb.Value, key string, function func(string)) {
	if val, ok := fields[key]; ok {
		function(strings.TrimSpace(val.GetStringValue()))
	}
}
