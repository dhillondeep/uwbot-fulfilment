package handlers

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"net/http"
	"strings"
)

func getStatusCode(container *gabs.Container) float64 {
	return container.Path("meta.status").Data().(float64)
}

func courseNotFoundError(statusCode float64, subject, catalogNum string) (string, bool) {
	if statusCode != http.StatusOK {
		return fmt.Sprintf("Sorry, I was unable to find this information "+
			"about %s %s at University of Waterloo",
			subject, catalogNum), false
	}

	return "", true
}

func convertTermShorthandToFull(shorthand string) string {
	switch shorthand {
	case "W":
		return "Winter"
	case "S":
		return "Spring"
	case "F":
		return "Fall"
	}

	return ""
}

func trimRespString(respStr string) string {
	return strings.Trim(strings.Trim(respStr, " "), ",")
}
