package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type User struct {
	Username string
	Email    string
	Phone    string
}

var secretKey = []byte("my-super-secret-key")

func sign(data string) string {
	h := hmac.New(sha256.New, secretKey)

	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

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

func loginHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	// pretend from database/json
	if username == "john" && password == "12345" {

		user := User{
			Username: "john",
			Email:    "john@gmail.com",
			Phone:    "08123456789",
		}

		cookieValue, err := createCookieValue(user)

		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name:     "auth",
			Value:    cookieValue,
			HttpOnly: true,
			Path:     "/",
		}

		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func verifyCookie(value string) (*User, error) {

	parts := strings.Split(value, ".")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cookie format")
	}

	encodedData := parts[0]
	receivedSignature := parts[1]

	expectedSignature := sign(encodedData)

	if expectedSignature != receivedSignature {
		return nil, fmt.Errorf("invalid signature")
	}

	decodedData, err := base64.StdEncoding.DecodeString(encodedData)

	if err != nil {
		return nil, err
	}

	var user User

	err = json.Unmarshal(decodedData, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("auth")

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := verifyCookie(cookie.Value)

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fmt.Fprintf(w,
		"Welcome %s Email: %s Phone: %s",
		user.Username,
		user.Email,
		user.Phone,
	)
}
