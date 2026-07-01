package main

import (
	"Age-Weaver/cookiesfunc"
	"Age-Weaver/form"
	"Age-Weaver/storage"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mail.v2"
)

type Userdata struct {
	Username string
	Email    string
	Phone    string
}

func GeneratePin() string {

	bytes := make([]byte, 3)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
func HashPin(pin string) string {
	hash := sha256.New()
	hash.Write([]byte(pin))
	return hex.EncodeToString(hash.Sum(nil))
}

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.ParseFiles("templates/login.html", "templates/signIn.html", "templates/dashboard.html", "templates/forgottenPwd.html", "templates/code.html", "templates/reset.html")
	if err != nil {
		http.Error(w, "Template Parsing Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		http.Error(w, "Template Execution Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
func DisplaySignUp(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Bad Method value", http.StatusMethodNotAllowed)
		return
	}
	renderTemplate(w, "signIn.html", nil)
}

type PageData struct {
	Username string
	Phone    string
	Email    string
	Password string

	UsernameError string
	PhoneError    string
	EmailError    string
	PasswordError string
}

func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Bad Method value", http.StatusMethodNotAllowed)
		return
	}

	UsernameData := r.FormValue("username")
	phoneData := r.FormValue("phone")
	EmailData := r.FormValue("email")
	PasswordData := r.FormValue("pwd")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(PasswordData), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error securing password", http.StatusInternalServerError)
		return
	}

	details := storage.User{
		Username: UsernameData,
		Phone:    phoneData,
		Email:    EmailData,
		Password: string(hashedPassword),
	}

	var nameErrStr, emailErrStr, phoneErrStr, passwordErrStr string
	var hasError bool

	_, err = form.ValidateName(UsernameData)
	if err != nil {
		nameErrStr = err.Error()
		hasError = true
	}

	_, err = form.ValidateEmail(EmailData)
	if err != nil {
		emailErrStr = err.Error()
		hasError = true
	}

	_, err = form.ValidatePhone(phoneData)
	if err != nil {
		phoneErrStr = err.Error()
		hasError = true
	}

	_, err = form.ValidatePwd(PasswordData)
	if err != nil {
		passwordErrStr = err.Error()
		hasError = true
	}

	data := PageData{
		Username: UsernameData,
		Phone:    phoneData,
		Email:    EmailData,
		Password: PasswordData,

		UsernameError: nameErrStr,
		PhoneError:    phoneErrStr,
		EmailError:    emailErrStr,
		PasswordError: passwordErrStr,
	}

	if hasError {
		renderTemplate(w, "signIn.html", data)
		return
	}

	err1 := storage.SaveData(details)
	if err1 != nil {
		http.Error(w, "failed to save user", http.StatusInternalServerError)
		return
	}

	User_Credentials := storage.User{
		Username: details.Username,
		Email:    details.Email,
		Phone:    details.Phone,
	}

	cookieValue, err := cookiesfunc.CreateCookieValue(User_Credentials)
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
}

type Loginerror struct {
	Error string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		renderTemplate(w, "login.html", nil)

	}
	if r.Method == http.MethodPost {

		storage.LoadData()

		email := r.FormValue("email")
		password := r.FormValue("pwd")
		var matchedUser storage.User

		LogSuccessful := false
		for _, user := range storage.Users {
			if user.Email == email {
				// VERIFY HASHED PASSWORD HERE
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
				if err == nil {
					matchedUser = user
					LogSuccessful = true
					break
				}
			}
		}

		if LogSuccessful {

			User_Credentials := storage.User{
				Username: matchedUser.Username,
				Email:    matchedUser.Email,
				Phone:    matchedUser.Phone,
			}

			cookieValue, err := cookiesfunc.CreateCookieValue((User_Credentials))

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

		} else {

			errormessage := "Account not found please"

			data := Loginerror{
				Error: errormessage,
			}
			renderTemplate(w, "login.html", data)

		}

	}

}

type Emaildata struct {
	Username string
	Code     string
}
type Sent struct {
	Sentdata string
}

func ForgottenPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "forgottenPwd.html", nil)
		return
	}

	type recoveryerror struct {
		Error string
	}
	email := r.FormValue("email")

	var receiver string
	var username string

	Valid_Email := false
	for _, user := range storage.Users {
		if user.Email == email {
			Valid_Email = true
			receiver = user.Email
			username = user.Username
			break
		}
	}

	if Valid_Email {
		code := GeneratePin()

		hashCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error securing code", http.StatusInternalServerError)
			return
		}

		details := storage.UserCode{
			Email:    receiver,
			Codedata: string(hashCode),
		}
		if err := storage.SaveCode(details); err != nil {
			http.Error(w, "Server error saving code", http.StatusInternalServerError)
			return
		}

		data := Emaildata{
			Username: username,
			Code:     code,
		}

		tmpl, err := template.ParseFiles("templates/emailcode.html")
		if err != nil {
			fmt.Println("Error parsing template:", err)
			http.Error(w, "Internal server error", 500)
			return
		}

		var bodyBuffer bytes.Buffer
		if err := tmpl.Execute(&bodyBuffer, data); err != nil {
			fmt.Println("Error executing template:", err)
			return
		}

		m := mail.NewMessage()
		m.SetHeader("From", "egeonujanai@gmail.com")
		m.SetHeader("To", receiver)
		m.SetHeader("Subject", "Age-Weaver (Confirmation Code)")
		m.SetBody("text/html", bodyBuffer.String())
		d := mail.NewDialer("smtp.gmail.com", 587, "egeonujanai@gmail.com", "giwg nrsr xqui xfrq")

		err = d.DialAndSend(m)
		if err != nil {
			http.Error(w, "email not sent ", 500)
			return
		}

		RedirectData := fmt.Sprintf("/code?email=%s", receiver)
		http.Redirect(w, r, RedirectData, http.StatusSeeOther)
		return

	} else {
		errordata := recoveryerror{
			Error: "Email credentials not found",
		}
		renderTemplate(w, "forgottenPwd.html", errordata)
		return
	}
}

type Confirmation struct {
	Email string
	Pin   string
	Error string
}
type InvalidPin struct {
	Error string
	Email string
	Pin   string
}

type User_code struct {
	Confirmation_Code string
}

func VerifyPin(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if r.Method == http.MethodGet {
		data := Confirmation{
			Email: email,
			Error: "",
		}
		renderTemplate(w, "code.html", data)
		return
	}

	if r.Method == http.MethodPost {
		userEnteredCode := r.FormValue("code")

		err := storage.LoadCode()
		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}

		Codesent := false
		for _, codeStorage := range storage.Code {
			if codeStorage.Email == email {
				err := bcrypt.CompareHashAndPassword([]byte(codeStorage.Codedata), []byte(userEnteredCode))
				if err == nil {
					Codesent = true
					break
				}
			}
		}

		if Codesent {

			http.Redirect(w, r, fmt.Sprintf("/reset?email=%s", email), http.StatusSeeOther)
			return
		} else {
			errordata := InvalidPin{
				Error: "Invalid Code",
				Email: email,
			}
			renderTemplate(w, "code.html", errordata)
			return
		}
	}
}

func ResetPWD(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		renderTemplate(w, "reset.html", nil)
		return
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("auth")

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := cookiesfunc.VerifyCookie(cookie.Value)

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "dashboard.html", user)

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
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/forgot", ForgottenPassword)
	http.HandleFunc("/code", VerifyPin)
	http.HandleFunc("/reset", ResetPWD)
	fmt.Println("running server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
