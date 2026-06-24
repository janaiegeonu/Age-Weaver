package main

import (
	"Age-Weaver/form"
	"Age-Weaver/storage"
	"fmt"
	"net/http"
	"text/template"
)

type Userdata struct {
	Username string
	Email    string
	Phone    string
}

var secretKey = []byte("my-super-secret-key")
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

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "login.html", nil)

	}
	if r.Method == http.MethodPost {

		storage.LoadData()

		email := r.FormValue("email")
		password := r.FormValue("pwd")
		var matchedUser storage.User

		LogSuccessful := false
		for _, user := range storage.Users {
			if user.Email == email &&
				user.Password == password {
				matchedUser = user
				LogSuccessful = true
				break
			}
		}

		if LogSuccessful {

			User_Credentials := Userdata{
				Username: matchedUser.Username,
				Email:    matchedUser.Email,
				Phone:    matchedUser.Phone,
			}

			cookieValue, err := createCookieValue(User(User_Credentials))

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

}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := storage.LoadData()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", DisplaySignUp)
	http.HandleFunc("/auth/register", Register)
	http.HandleFunc("/login", loginHandler)
	fmt.Println("running server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
