package main

import (
	"errors"
	"fmt"
)

func ValidatePwd(text string) error {

	if len(text) < 8 {
		return errors.New("Error : password too short")
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
		return errors.New("Error : password need to contain at least one digit")

	}
	if !Has_letters {
		return errors.New("Error : password need to contain at least one letter")

	}
	return nil
}

func main() {
	fmt.Println(ValidatePwd("1fbchgj"))
	fmt.Println(ValidatePwd("7777777777"))
	fmt.Println(ValidatePwd("hhhhhhhhhhkjgsd"))
	fmt.Println(ValidatePwd("1fbchgj5tt"))

}
