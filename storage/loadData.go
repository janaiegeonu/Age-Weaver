package storage

import (
	"encoding/json"
	"os"
)

type User struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var Users []User

func LoadData() error {

	data, err := os.ReadFile("user.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &Users)
}
