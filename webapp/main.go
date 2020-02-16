package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/authenticate", authenticate)
	router.HandleFunc("/signupPage", signUpPage)
	router.HandleFunc("/signup_account", signup)
	router.HandleFunc("/machineRead/{machineID}", dispense)
	router.HandleFunc("/resetPage", resetPage)
	router.HandleFunc("/private/{page}", privatePageGenerator)
	router.HandleFunc("/signout", logout)
	router.HandleFunc("/donate", donate)
	router.HandleFunc("/contribute", contributeData)
	// router.HandleFunc("/update", update)
	router.HandleFunc("/item/{id}", rentItem)

	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/",
		http.FileServer(http.Dir("public/css/"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir("public/js/"))))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/",
		http.FileServer(http.Dir("public/assets/"))))
	router.PathPrefix("/private/css/").Handler(http.StripPrefix("/private/css/",
		http.FileServer(http.Dir("public/css/"))))
	router.PathPrefix("/private/js/").Handler(http.StripPrefix("/private/js/",
		http.FileServer(http.Dir("public/js/"))))
	router.PathPrefix("/private/assets/").Handler(http.StripPrefix("/private/assets/",
		http.FileServer(http.Dir("public/assets/"))))
	router.PathPrefix("/private/img/").Handler(http.StripPrefix("/private/img/",
		http.FileServer(http.Dir("img/"))))

	server := &http.Server{
		Addr:           "",
		Handler:        router,
		ReadTimeout:    time.Duration(10 * int64(time.Second)),
		WriteTimeout:   time.Duration(600 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	// server.ListenAndServeTLS("fullchain.pem", "privkey.pem")
	// server.ListenAndServeTLS("cert.pem", "key.pem")

	server.ListenAndServe()

}
