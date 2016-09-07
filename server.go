package main

import (
	"fmt"
	"net/http"

	"whats-for-lunch/sources/spreadsheet"

	"whats-for-lunch/authentication"
)

func main() {
	fmt.Println("registering routes")
	http.HandleFunc("/whats-for-lunch/authenticate", authentication.BaseHandler)
	http.HandleFunc("/whats-for-lunch/authenticate/login", authentication.LoginHandler)
	http.HandleFunc("/whats-for-lunch/authenticate/redirect", authentication.RedirectHandler)
	http.HandleFunc("/whats-for-lunch/tomorrow", spreadsheet.GetTomorrowsMenu)
	fmt.Println("starting server here")
	http.ListenAndServe("0.0.0.0:8081", nil)
}
