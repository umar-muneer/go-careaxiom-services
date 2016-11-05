package main

import (
	"fmt"
	"net/http"

	"go-careaxiom-services/whats-for-lunch/sources/spreadsheet"
	"log"
	"os"

	"go-careaxiom-services/leaves"

	"github.com/umar-muneer/go-careaxiom-utilities/authentication"
)

func main() {
	fmt.Println("registering routes for all sub apps")
	authentication.New(authentication.Parameters{
		OAuthScopes: []string{"https://www.googleapis.com/auth/spreadsheets"},
		LoginPageHTML: `<html>
            					<body>
              					<a href="/authenticate/login">Authenticate Careaxiom Services API</a>
            					</body>
          					</html>`,
		BaseURL: "/authenticate",
	})
	http.HandleFunc("/authenticate", authentication.BaseHandler)
	http.HandleFunc("/authenticate/login", authentication.LoginHandler)
	http.HandleFunc("/authenticate/redirect", authentication.RedirectHandler)
	http.HandleFunc("/whats-for-lunch", spreadsheet.GetMenu)
	http.HandleFunc("/lunch/review", spreadsheet.HandleReview)
	http.HandleFunc("/leaves/status", leaves.GetLeavesStatus)
	http.HandleFunc("/leaves/cache", leaves.HandleCache)
	fmt.Println("starting server here")
	var port = ":"
	if os.Getenv("PORT") != "" {
		port += os.Getenv("PORT")
	} else {
		port += "8081"
	}
	fmt.Println("binding server on port ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
