package cookiesfunc

import (
	"Age-Weaver/storage"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func VerifyCookie(value string) (*storage.User, error) {

	parts := strings.Split(value, ".")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cookie format")
	}

	encodedData := parts[0]
	receivedSignature := parts[1]

	expectedSignature := Sign(encodedData)

	if expectedSignature != receivedSignature {
		return nil, fmt.Errorf("invalid signature")
	}

	decodedData, err := base64.StdEncoding.DecodeString(encodedData)

	if err != nil {
		return nil, err
	}

	var user storage.User

	err = json.Unmarshal(decodedData, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
