package storage

import (
	"encoding/json"
	"os"
)

func SaveData(value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil
	}
	return os.WriteFile("user.json", data, 0644)
}
