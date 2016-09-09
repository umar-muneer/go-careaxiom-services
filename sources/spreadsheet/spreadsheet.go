package spreadsheet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"whats-for-lunch/authentication"
)

var apiURL = "https://sheets.googleapis.com/v4/spreadsheets"

/*Output representation of a menu for a day*/
type Output struct {
	New     string
	Old     string
	Dessert string
}

/*GetTomorrowsMenu returns tomorrow's menu*/
func GetTomorrowsMenu(res http.ResponseWriter, req *http.Request) {
	fmt.Println("getting tomorrow's menu")
	spreadSheetClient, err := authentication.GetClient()
	if err != nil {
		fmt.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	var oldURL = apiURL + "/" + os.Getenv("OLD_MENU_SPREADSHEET_ID") + "/values" + "/A6:E36"
	fmt.Println("spreadsheet url is", oldURL)
	spreadSheetResponse, spreadSheetErr := spreadSheetClient.Get(oldURL)
	if spreadSheetErr != nil {
		fmt.Println("error reading from spreadsheet", spreadSheetErr.Error())
		http.Error(res, spreadSheetErr.Error(), http.StatusInternalServerError)
	}
	output, _ := ioutil.ReadAll(spreadSheetResponse.Body)
	fmt.Println(string(output[:]), "yayayaya")
	json.NewEncoder(res).Encode(Output{New: "Biryani", Old: "Karahi", Dessert: "Custard"})
}
