package storage

import (
	"encoding/json"
	"os"
)

func SaveCode(user UserCode) error {
	err := LoadCode()
	if err != nil {
		return err
	}

	var updatedCodes []UserCode
	for _, c := range Code {
		if c.Email != user.Email {
			updatedCodes = append(updatedCodes, c)
		}
	}
	Code = updatedCodes

	Code = append(Code, user)

	data, err := json.MarshalIndent(Code, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("code.json", data, 0644)
}
