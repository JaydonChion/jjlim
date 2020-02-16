package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {

		generateHTML(w, nil, "loginlayout", "login")

	} else {
		http.Redirect(w, r, "/private/home", http.StatusSeeOther)
	}
}

func signUpPage(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "signuplayout", "signup")

}

func resetPage(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "resetlayout", "reset")

}

func dispense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	mid := vars["machineID"]
	midnum, err := strconv.Atoi(mid)
	if err != nil {
		panic(err)
	}
	check, err := activate(midnum)

	if check {
		w.Write([]byte("dispense"))
	} else {
		w.Write([]byte("no"))
	}

}

func rentItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	userID := vars["id"]
	userIDnum, err := strconv.Atoi(userID)
	itemID := 5
	locationID := 1
	err = rentService(userIDnum, itemID, locationID)

	if err != nil {
		w.Write([]byte("fail"))
	} else {
		w.Write([]byte("success"))
	}

	// cookie := http.Cookie{
	// 	Name:     "error_message",
	// 	Value:    "Thank you for using our service",
	// 	HttpOnly: false,
	// }

	// if err != nil {
	// 	cookie = http.Cookie{
	// 		Name:     "error_message",
	// 		Value:    "Something wrong",
	// 		HttpOnly: false,
	// 	}
	// }

	// http.SetCookie(w, &cookie)
	// http.Redirect(w, r, "/private/homepage", 302)

}

func privatePageGenerator(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	_, err := session(w, r)

	if err != nil {
		generateHTML(w, nil, "loginlayout", "login")
	} else {
		// currentUser, err := getSessionUserData(w, r)
		if err != nil {
			panic(err)
		}

		currentUser, _ := getSessionUserData(w, r)
		currentUser.Name = strings.ToUpper(currentUser.Name)
		if err != nil {

			panic(err)
		}

		cups, err := getCupSaved(currentUser.Id)
		hists, err := getHistory(currentUser.Id)

		if err != nil {
			panic(err)
		}

		cups = (cups * 10)

		temp := make(map[string]interface{})
		temp["user"] = currentUser
		temp["page"] = vars["page"]
		temp["cupSaved"] = cups
		temp["hists"] = hists

		switch pagename := vars["page"]; pagename {

		case "predict":
			pastData, _ := getContribInfo()
			temp["predict"] = pastData
			generateHTML(w, temp, "privatelayout", "privatelayout", "predict")

		case "donate":
			donate, _ := getDonatedFood()
			temp["donate"] = donate
			generateHTML(w, temp, "privatelayout", "privatelayout", "donate")

		case "rewards":

			generateHTML(w, temp, "privatelayout", "privatelayout", "rewards")

		default:

			generateHTML(w, temp, "privatelayout", "privatelayout", "home")
		}

	}

}

func donate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	currentUser, _ := getSessionUserData(w, r)
	location := r.FormValue("location")
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	Mnum := r.FormValue("Mnum")
	halal := r.FormValue("halal")
	amount := r.FormValue("amount")
	descrp := r.FormValue("descrp")

	amountNum, err := strconv.Atoi(amount)
	phoneNum, err := strconv.Atoi(phone)

	if err != nil {
		panic(err)
	}

	err = donateFood(location, name, phoneNum, Mnum, halal, amountNum, descrp, currentUser.Id)

	cookie := http.Cookie{
		Name:     "error_message",
		Value:    "Food donated successfully",
		HttpOnly: false,
	}

	if err != nil {
		cookie = http.Cookie{
			Name:     "error_message",
			Value:    "Fail to donate, please contact us if problem presists",
			HttpOnly: false,
		}
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/private/donate", 302)

}

func getSessionUserData(w http.ResponseWriter, r *http.Request) (currentUser User, err error) {
	c, err := r.Cookie("_cookie")
	if err != nil {
		panic(err)
	} else {
		value := c.Value
		currentUser, err = UserByUUID(value)
	}

	return
}

func contributeData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	tpeople := r.FormValue("tpeople")
	nmale := r.FormValue("nmale")
	nfemale := r.FormValue("nfemale")
	age := r.FormValue("age")
	time := r.FormValue("time")
	foodleft := r.FormValue("foodleft")

	tpeopleNum, err := strconv.Atoi(tpeople)
	nmaleNum, err := strconv.Atoi(nmale)
	nfemaleNum, err := strconv.Atoi(nfemale)
	ageNum, err := strconv.Atoi(age)
	timeNum, err := strconv.Atoi(time)
	foodleftNum, err := strconv.Atoi(foodleft)

	if err != nil {
		panic(err)
	}

	err = contributeInfo(tpeopleNum, nmaleNum, nfemaleNum, ageNum, timeNum, foodleftNum)

	cookie := http.Cookie{
		Name:     "error_message",
		Value:    "Thank you for contributing your data",
		HttpOnly: false,
	}

	if err != nil {
		cookie = http.Cookie{
			Name:     "error_message",
			Value:    "Fail to update datam please try again later",
			HttpOnly: false,
		}
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/private/predict", 302)
}

// func (user *User) emailVerification(w http.ResponseWriter, r *http.Request) (err error) {
// 	from := mail.NewEmail("EIstint", "eistint@gmail.com")
// 	subject := "EIstint account verification"
// 	to := mail.NewEmail("EIstint", user.Email)
// 	plainTextContent := "Confidential account verification information"
// 	// link := 'http://localhost:3000/verify?num="' + strconv.Itoa(user.Id) + "&code=" + user.Uuid[:6]
// 	htmlContent := "Dear " + user.Name + ",<br><br><strong>This email is intended solely for the person named above.</strong> If you are not the intended recipient, please delete it, notify us and do not copy, use, or disclose its contents.<br><br><b>Verification Link: </b><br><a href='https://eistint.com:3000/verify?num=" + strconv.Itoa(user.Id) + "&code=" + user.Uuid[:6] + "'>Click here to verify</a><br><br><b>Or go to this link</b><br>https://eistint.com:3000/verify?num=" + strconv.Itoa(user.Id) + "&code=" + user.Uuid[:6]
// 	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
// 	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
// 	_, err = client.Send(message)

// 	return
// }
