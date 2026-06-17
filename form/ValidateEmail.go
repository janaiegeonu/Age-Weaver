package form

import (
	"errors"
	"strings"
)

func ValidateEmail(text string) (string, error) {
	ValidSuffixes := []string{
		"@gmail.com", "@yahoo.com", "@outlook.com", "@icloud.com",
		".org", ".net", ".edu", ".gov", ".co", ".tech", ".info",
	}

	hasValidSuffix := false
	for _, suffix := range ValidSuffixes {
		if strings.HasSuffix(text, suffix) {
			hasValidSuffix = true
			break
		}
	}

	if !hasValidSuffix {
		return "", errors.New("error: invalid email")
	}

	return "", nil
}
