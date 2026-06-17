package storage

import (
	"encoding/json"
	"os"
)

func SaveData(user User) error {

	err := LoadData()
	if err != nil {
		return err
	}

	Users = append(Users, user)

	data, err := json.MarshalIndent(Users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("user.json", data, 0644)
}
