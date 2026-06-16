package form

import (
	"errors"
)

func ValidatePwd(text string) (string, error) {

	if len(text) < 8 {
		return "", errors.New("error : password too short")
	}
	Has_digits := false
	Has_letters := false
	for _, char := range text {
		if char >= '0' && char <= '9' {
			Has_digits = true
		}
		if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
			Has_letters = true
		}
	}
	if !Has_digits {
		return "", errors.New("error : password need to contain at least one digit")

	}
	if !Has_letters {
		return "", errors.New("error : password need to contain at least one letter")

	}
	return "", nil
}
