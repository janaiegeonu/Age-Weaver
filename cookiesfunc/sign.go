package cookiesfunc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

var secretKey = []byte("my-super-secret-key")

func Sign(data string) string {
	h := hmac.New(sha256.New, secretKey)

	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}
