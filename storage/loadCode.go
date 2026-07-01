package storage

import (
	"encoding/json"
	"os"
)

type UserCode struct {
	Email    string `json:"email"`
	Codedata string `json:"codedata"`
}

var Code []UserCode

func LoadCode() error {
	data, err := os.ReadFile("code.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &Code)
}
