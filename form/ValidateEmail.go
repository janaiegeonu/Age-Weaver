package form

import (
	"errors"
	"strings"
)

func ValidateEmail(text string) (string, error) {

	if !strings.HasSuffix(text, "@gmail.com") || !strings.HasSuffix(text, "@yahoo.com") || !strings.HasSuffix(text, "@outlook.com") || !strings.HasSuffix(text, "@icloud.com") || !strings.HasSuffix(text, ".org") ||
		!strings.HasSuffix(text, ".net") || !strings.HasSuffix(text, ".edu") || !strings.HasSuffix(text, ".gov") || !strings.HasSuffix(text, ".co") || !strings.HasSuffix(text, ".tech") || !strings.HasSuffix(text, ".info") {
		return "", errors.New("error: invalid email")
	}
	return "", nil
}
