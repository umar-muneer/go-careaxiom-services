package authentication

import (
	"fmt"
	"net/http"
	"os"

	"encoding/json"
	"errors"

	"github.com/umar-muneer/go-careaxiom-utilities/filetransfer"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	html = `<html>
            <body>
              <a href="/authenticate/login">Log in with Google</a>
            </body>
          </html>`
	baseURL         = "/authenticate"
	credentialsFile = "lunch.credentials"
)

var oauthConfig = &oauth2.Config{}

/*Parameters authentication parameters to use this package with multiple different projects*/
type Parameters struct {
	LoginPageHTML string
	OAuthScopes   []string
	BaseURL       string
}

var parameters = Parameters{}

/*New set authentication parameters*/
func New(params Parameters) {
	parameters = params
	oauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       parameters.OAuthScopes,
		Endpoint:     google.Endpoint,
	}
}

var oauthState = "why do i need you??"

func saveToken(token *oauth2.Token) error {
	fmt.Println("saving token to file")
	file, err := os.Create(credentialsFile)
	if err != nil {
		errorText := "cannot create token credentials file"
		fmt.Println(errorText)
		return errors.New(errorText)
	}
	defer file.Close()
	s3Folder := os.Getenv("AWS_S3_FOLDER")
	credentialFileName := os.Getenv("CREDENTIALS_FILE_NAME")

	if s3Folder == "" || credentialFileName == "" {
		return errors.New("credential file name or folder key missing")
	}
	s3Writer := &filetransfer.S3IO{
		Bucket: os.Getenv("AWS_S3_BUCKET"),
		Key:    s3Folder + "/" + credentialFileName,
	}
	writeError := json.NewEncoder(s3Writer).Encode(token)
	if writeError != nil {
		return writeError
	}
	return nil
}

func loadToken() (*oauth2.Token, error) {
	fmt.Println("loading token from file")

	s3Folder := os.Getenv("AWS_S3_FOLDER")
	credentialFileName := os.Getenv("CREDENTIALS_FILE_NAME")

	if s3Folder == "" || credentialFileName == "" {
		return nil, errors.New("credential file name or folder key missing")
	}
	s3Reader := &filetransfer.S3IO{
		Bucket: os.Getenv("AWS_S3_BUCKET"),
		// Key:    "whats-for-lunch/lunch.credentials",
		Key: s3Folder + "/" + credentialFileName,
	}
	token := &oauth2.Token{}
	decodeErr := json.NewDecoder(s3Reader).Decode(token)
	return token, decodeErr
}

/*GetClient get an http client to talk to google spreadsheets*/
func GetClient() (*http.Client, error) {
	token, err := loadToken()
	if err != nil {
		fmt.Println("error while loading token", err.Error())
		return nil, errors.New("refresh token not found, cannot proceed further")
	}
	client := oauthConfig.Client(context.Background(), token)
	fmt.Println("RT:", token.RefreshToken, ", AT:", token.AccessToken, ", Expiry:", token.Expiry)
	return client, nil
}

/*LoginHandler main controller method which redirects the browser to a page to authorize the app*/
func LoginHandler(res http.ResponseWriter, req *http.Request) {
	url := oauthConfig.AuthCodeURL(oauthState)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

/*RedirectHandler this is where the google page redirects to send the token information*/
func RedirectHandler(res http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	if state != oauthState {
		fmt.Println("invalid state variable received")
		http.Redirect(res, req, parameters.BaseURL, http.StatusTemporaryRedirect)
		return
	}
	token, err := oauthConfig.Exchange(context.Background(), req.FormValue("code"))
	if err != nil {
		fmt.Println("failed to receive token")
		http.Error(res, "failed to receive token", http.StatusInternalServerError)
		return
	}
	saveErr := saveToken(token)
	if saveErr != nil {
		fmt.Println(saveErr)
		http.Error(res, "cannot save token to file", http.StatusInternalServerError)
	}
}

/*BaseHandler used to launch the entire authentication process*/
func BaseHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("printing base authentication page")
	fmt.Fprintf(res, parameters.LoginPageHTML)
}
