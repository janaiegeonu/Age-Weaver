package cookiesfunc

import (
	"Age-Weaver/storage"
	"encoding/base64"
	"encoding/json"
)

func CreateCookieValue(user storage.User) (string, error) {

	jsonData, err := json.Marshal(user)

	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(jsonData)

	signature := Sign(encoded)

	finalValue := encoded + "." + signature

	return finalValue, nil
}
