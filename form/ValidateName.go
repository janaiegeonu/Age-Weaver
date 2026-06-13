package main

import (
	"errors"
	"fmt"
	"unicode"
)

func ValidateName(text string) error {

	space := false

	for _, char := range text {
		if char == ' ' {
			space = true
		}
	}
	if !space {
		return errors.New("Error: please enter your fullname")
	}

	Has_numeric := false

	for _, char := range text {
		if char >= '0' && char <= '9' {
			Has_numeric = true
		}
	}

	Has_symbols := false

	for _, char := range text {
		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			Has_symbols = true
		}
		if char == '-' || char == '\'' {
			Has_symbols = false
		}
	}

	if Has_symbols || Has_numeric {
		return errors.New("Error: alphabetic format only ")
	}


	
	return nil

}

func main() {

	fmt.Println(ValidateName("janai"))
	fmt.Println(ValidateName("janai egeonu"))
	fmt.Println(ValidateName("janai egeonu123"))
	fmt.Println(ValidateName("egeonu janai??"))
	fmt.Println(ValidateName("D'janai jay"))

}
