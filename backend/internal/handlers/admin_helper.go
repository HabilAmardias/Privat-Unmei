package handlers

import "regexp"

var (
	degreelist = []string{"bachelor", "diploma", "high school", "master", "professor"}
)

func ValidateDegree(degree string) bool {
	for _, item := range degreelist {
		if degree == item {
			return true
		}
	}
	return false
}

func ValidatePhoneNumber(phoneNumber string) bool {

	pattern := `^0\d{9,12}$`

	regex := regexp.MustCompile(pattern)

	return regex.MatchString(phoneNumber)
}
