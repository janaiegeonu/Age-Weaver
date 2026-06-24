package cookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func sign(data string) string {
	h := hmac.New(sha256.New, secretKey)

	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}
