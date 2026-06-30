package storage

import (
	"encoding/json"
	"os"
)

func SaveCode(user UserCode) error {

	err := LoadData()
	if err != nil {
		return err
	}

	Code = append(Code, user)

	data, err := json.MarshalIndent(Code, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("code.json", data, 0644)
}
