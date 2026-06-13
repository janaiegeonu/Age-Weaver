package main

import (
	"errors"
	"unicode"
)

func ValidatePhone(text string) error {

	if len(text) < 10 {
		return errors.New("Error: invalid phone number length")
	}
	Not_Numerics := false
	for _, char := range text {
		if !unicode.IsDigit(char) {
			Not_Numerics = true
		}
		if char == '+' || char == '-' || char == '(' || char == ')' || char == ' ' {
			Not_Numerics = false

		}
	}
	if Not_Numerics {
		return errors.New("Error: invalid phone number")

	}

	has_symbols := false

	for _, char := range text {
		if char == '+' || char == '(' || char == ')' || char == '-' {
			has_symbols = true
		}
	}
	if has_symbols {

		if text[0] != '+' && text[0] != '(' {
			return errors.New("Error: invalid phone number")
		}
	}

	if text[0] == '-' && text[len(text)-1] == '-' && text[len(text)-1] == '(' && text[len(text)-1] == ')' {
		return errors.New("Error: invalid phone number")

	}
	count_hypen := 0
	count_plus := 0
	count_bracket := 0

	for _, char := range text {
		if char == '-' {
			count_hypen++
		}
	}

	if count_hypen > 3 {
		return errors.New("Error: invalid phone number")

	}

	for _, char := range text {
		if char == '+' {
			count_plus++
		}
	}

	if count_plus > 1 {
		return errors.New("Error: invalid phone number")

	}

	for _, char := range text {
		if char == '(' || char == ')' {
			count_bracket++
		}
	}

	if count_bracket > 2 {
		return errors.New("Error: invalid phone number")

	}

	return nil

}
