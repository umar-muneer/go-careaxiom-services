package main

import (
	"fmt"
	"net/http"

	"go-careaxiom-utilities/authentication"
	"log"
	"os"
	"whats-for-lunch/sources/spreadsheet"
)

func main() {
	fmt.Println("registering routes")
	http.HandleFunc("/whats-for-lunch/authenticate", authentication.BaseHandler)
	http.HandleFunc("/whats-for-lunch/authenticate/login", authentication.LoginHandler)
	http.HandleFunc("/whats-for-lunch/authenticate/redirect", authentication.RedirectHandler)
	http.HandleFunc("/whats-for-lunch", spreadsheet.GetMenu)
	fmt.Println("starting server here")
	var port = ":"
	if os.Getenv("PORT") != "" {
		port += os.Getenv("PORT")
	} else {
		port += "8081"
	}
	log.Fatal(http.ListenAndServe(port, nil))
}
