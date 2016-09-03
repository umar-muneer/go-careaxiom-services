package authentication

import (
	"fmt"
	"net/http"
	"os"

	"encoding/json"
	"errors"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	HTML = `<html>
            <body>
              <a href="/whats-for-lunch/authenticate/login">Log in with Google</a>
            </body>
          </html>`
	BASE_URL         = "/whats-for-lunch/authenticate"
	CREDENTIALS_FILE = "lunch.credentials"
)

var oauthConfig = &oauth2.Config{
	RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets"},
	Endpoint:     google.Endpoint,
}
var oauthState = "why do i need you??"

func saveToken(token *oauth2.Token) error {
	fmt.Println("saving token to file")
	file, err := os.Create(CREDENTIALS_FILE)
	if err != nil {
		errorText := "cannot create token credentials file"
		fmt.Println(errorText)
		return errors.New(errorText)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(token)
	return nil
}

func loadToken() (*oauth2.Token, error) {
	fmt.Println("loading token from file")
	file, err := os.Open(CREDENTIALS_FILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	token := &oauth2.Token{}
	decodeErr := json.NewDecoder(file).Decode(token)

	return token, decodeErr
}

func GetClient() (*http.Client, error) {
	token, err := loadToken()
	if err != nil {
		return nil, errors.New("refresh token not found, cannot proceed further")
	}
	client := oauthConfig.Client(context.Background(), token)
	return client, nil
}

func LoginHandler(res http.ResponseWriter, req *http.Request) {
	url := oauthConfig.AuthCodeURL(oauthState)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
func RedirectHandler(res http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	if state != oauthState {
		fmt.Println("invalid state variable received")
		http.Redirect(res, req, BASE_URL, http.StatusTemporaryRedirect)
		return
	}
	token, err := oauthConfig.Exchange(context.Background(), req.FormValue("code"))
	if err != nil {
		fmt.Println("failed to recive token")
		http.Error(res, "failed to receive token", http.StatusInternalServerError)
		return
	}
	saveErr := saveToken(token)
	if saveErr != nil {
		http.Error(res, "cannot save token to file", http.StatusInternalServerError)
	}
}
func BaseHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("printing base authentication page")
	fmt.Fprintf(res, HTML)
}
