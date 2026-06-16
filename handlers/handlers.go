package main

import (
	"Age-Weaver/form"
	"Age-Weaver/storage"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseFiles("templates/login.html", "templates/signin.html"))

func Signin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 path not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Bad Method value", http.StatusMethodNotAllowed)
		return
	}

	type User struct {
		Username string
		Phone    string
		Email    string
		Password string

		UsernameError string
		PhoneError    string
		EmailError    string
		PasswordError string
	}

	UsernameData := r.FormValue("username")
	phoneData := r.FormValue("phone")
	EmailData := r.FormValue("email")
	PasswordData := r.FormValue("password")

	details := User{
		Username: UsernameData,
		Phone:    phoneData,
		Email:    EmailData,
		Password: PasswordData,
	}

	hasError := false
	_, err := form.ValidateName(UsernameData)
	if err != nil {
		details.UsernameError = err.Error()
	}

	_, err = form.ValidateEmail(EmailData)
	if err != nil {
		details.EmailError = err.Error()
	}

	_, err = form.ValidatePhone(phoneData)
	if err != nil {
		details.PhoneError = err.Error()
	}
	_, err = form.ValidatePwd(PasswordData)
	if err != nil {
		details.PasswordError = err.Error()
	}

	if hasError {
		w.WriteHeader(http.StatusUnprocessableEntity)
		templates.ExecuteTemplate(w, "signin.html", details)
		return
	}

	err1 := storage.SaveData(details)
	if err1 != nil {
		http.Error(w, "failed to save user", http.StatusInternalServerError)
		return
	}

}
