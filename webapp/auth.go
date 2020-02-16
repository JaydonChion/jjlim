package main

import (
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func signup(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		panic("Cannot parse form")
	}
	phoneNum, err := strconv.Atoi(request.PostFormValue("phone"))

	timeStampString := request.PostFormValue("dobDate")
	layOut := "2006-01-02"
	timeStamp, err := time.Parse(layOut, timeStampString)

	if err != nil {
		panic(err)
	}

	user := User{
		Name:     request.PostFormValue("name"),
		Email:    request.PostFormValue("email"),
		Password: ToNullString(request.PostFormValue("password")),
		Gender:   ToNullString(request.PostFormValue("gender")),
		Phone:    phoneNum,
		Mnum:     request.PostFormValue("Mnum"),
		Dob:      timeStamp}

	if err := user.Create(); err != nil {
		cookie := http.Cookie{
			Name:     "error_message",
			Value:    "Unable to create account",
			HttpOnly: false,
		}
		http.SetCookie(writer, &cookie)

	} else {
		// err := user.emailVerification(writer, request)

		// if err != nil {
		cookie := http.Cookie{
			Name:     "error_message",
			Value:    "Account created successfully",
			HttpOnly: false,
		}
		http.SetCookie(writer, &cookie)

		// } else {
		// 	cookie := http.Cookie{
		// 		Name:     "error_message",
		// 		Value:    "Account created successfully, please check your email for verification",
		// 		HttpOnly: false,
		// 	}
		// 	http.SetCookie(writer, &cookie)

		// }

	}

	http.Redirect(writer, request, "/", 302)
}

// POST /authenticate
// Authenticate the user given the email and password
func authenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	user, err := UserByEmail(request.PostFormValue("email"))

	if err != nil {

		cookie := http.Cookie{
			Name:     "error_message",
			Value:    "No existing account, please sign up a new account",
			HttpOnly: false,
		}
		http.SetCookie(writer, &cookie)
		http.Redirect(writer, request, "/", 302)
		// panic(err)
	} else {

		// if !user.Verification_status {

		// 	err = user.emailVerification(writer, request)

		// 	var message string
		// 	if err != nil {

		// 		message = "Unable to log in. if problem persist, please contact for help"
		// 	} else {

		// 		message = "Please check your registered email address for email verification"
		// 	}

		// 	cookie := http.Cookie{
		// 		Name:     "error_message",
		// 		Value:    message,
		// 		HttpOnly: false,
		// 	}

		// 	http.SetCookie(writer, &cookie)
		// 	http.Redirect(writer, request, "/", 302)
		// }

		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(request.PostFormValue("password")))
		if passErr == nil {

			session, err := user.CreateSession()
			if err != nil {
				panic(err)
			}
			cookie := http.Cookie{
				Name:     "_cookie",
				Value:    session.Uuid,
				HttpOnly: true,
			}

			// intro, err := getFBStatus(user.Id)
			// if err != nil {
			// 	panic(err)
			// }
			http.SetCookie(writer, &cookie)
			// if intro {
			// 	http.Redirect(writer, request, "/private/home", 302)

			// } else {
			// 	http.Redirect(writer, request, "/introduction", 302)

			// }
			http.Redirect(writer, request, "/private/home", 302)

		} else {

			cookie := http.Cookie{
				Name:     "error_message",
				Value:    "wrong email/password",
				HttpOnly: false,
			}
			http.SetCookie(writer, &cookie)
			http.Redirect(writer, request, "/", 302)
		}
	}

}

// GET /logout
// Logs the user out
func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		// fmt.Println("failed to get cookie")
		session := Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	http.Redirect(writer, request, "/", 302)
}
