package cookies

import (
	"encoding/base64"
	"encoding/json"
)

func createCookieValue(user User) (string, error) {

	jsonData, err := json.Marshal(user)

	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(jsonData)

	signature := sign(encoded)

	finalValue := encoded + "." + signature

	return finalValue, nil
}
