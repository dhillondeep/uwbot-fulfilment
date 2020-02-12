package helpers

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"strings"
)

var termsShortHand = map[string]string{"w": "Winter", "s": "Spring", "f": "Fall"}

// iterates over json array provides gabs container and path. It callbacks the user provided function
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

func GetStatusCode(container *gabs.Container) float64 {
	return container.Path("meta.status").Data().(float64)
}

func ConvertTermShorthandToFull(shorthand string) string {
	return termsShortHand[strings.ToLower(shorthand)]
}

func StringEqualNoCase(first, second string) bool {
	return strings.ToLower(strings.TrimSpace(first)) == strings.ToLower(strings.TrimSpace(second))
}

func StringIsEmpty(text string) bool {
	return strings.TrimSpace(text) == ""
}
