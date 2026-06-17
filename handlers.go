package main

import (
	"Age-Weaver/form"
	"Age-Weaver/storage"
	"fmt"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseFiles("templates/login.html", "templates/signIn.html"))

func DisplaySignUp(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Bad Method value", http.StatusMethodNotAllowed)
		return
	}
	templates.ExecuteTemplate(w, "signIn.html", nil)
}
func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Bad Method value", http.StatusMethodNotAllowed)
		return
	}

	type User struct {
		Username string
		Phone    string
		Email    string
		Password string
	}

	type PageData struct {
		Username string
		Phone    string
		Email    string
		Password string

		UsernameError error
		PhoneError    error
		EmailError    error
		PasswordError error
	}

	UsernameData := r.FormValue("username")
	phoneData := r.FormValue("phone")
	EmailData := r.FormValue("email")
	PasswordData := r.FormValue("password")

	details := storage.User{
		Username: UsernameData,
		Phone:    phoneData,
		Email:    EmailData,
		Password: PasswordData,
	}

	var NameErr, EmailErr, PhoneErr, PasswordErr error
	_, err := form.ValidateName(UsernameData)
	if err != nil {
		NameErr = err
	}

	_, err = form.ValidateEmail(EmailData)
	if err != nil {
		EmailErr = err
	}

	_, err = form.ValidatePhone(phoneData)
	if err != nil {
		PhoneErr = err
	}
	_, err = form.ValidatePwd(PasswordData)
	if err != nil {
		PasswordErr = err
	}

	data := PageData{
		Username: UsernameData,
		Phone:    phoneData,
		Email:    EmailData,
		Password: PasswordData,

		UsernameError: NameErr,
		PhoneError:    PhoneErr,
		EmailError:    EmailErr,
		PasswordError: PasswordErr,
	}

	if NameErr != nil ||
		PhoneErr != nil ||
		EmailErr != nil ||
		PasswordErr != nil {

		templates.ExecuteTemplate(w, "signIn.html", data)
		return
	}

	err1 := storage.SaveData(details)
	if err1 != nil {
		http.Error(w, "failed to save user", http.StatusInternalServerError)
		return
	}

}

func Login(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)

}
func Submit(w http.ResponseWriter, r *http.Request) {

}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := storage.LoadData()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", DisplaySignUp)
	http.HandleFunc("/auth/register", Register)
	http.HandleFunc("/login", Login)
	fmt.Println("running server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
