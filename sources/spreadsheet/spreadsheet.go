package spreadsheet

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"whats-for-lunch/authentication"
	"whats-for-lunch/menu"
)

var apiURL = "https://sheets.googleapis.com/v4/spreadsheets"

/*Output representation of a menu for a day*/
type Output struct {
	New menu.Entry
	Old menu.Entry
}

/*GetMenu returns tomorrow's menu*/
func GetMenu(res http.ResponseWriter, req *http.Request) {
	fmt.Println("getting menu for date")
	spreadSheetClient, err := authentication.GetClient()
	if err != nil {
		fmt.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	offset, _ := strconv.Atoi(os.Getenv("OLD_MENU_SHEET_OFFSET"))
	oldMenu := menu.New(spreadSheetClient, os.Getenv("OLD_MENU_SPREADSHEET_ID"), offset)
	menuEntryErr := oldMenu.GetMenuEntry(req.URL.Query().Get("date"))
	if menuEntryErr != nil {
		fmt.Println(menuEntryErr)
		http.Error(res, menuEntryErr.Error(), http.StatusInternalServerError)
	}
	// var oldURL = apiURL + "/" + os.Getenv("OLD_MENU_SPREADSHEET_ID") + "/values" + "/A6:E36"
	// fmt.Println("spreadsheet url is", oldURL)
	// spreadSheetResponse, spreadSheetErr := spreadSheetClient.Get(oldURL)
	// if spreadSheetErr != nil {
	// 	fmt.Println("error reading from spreadsheet", spreadSheetErr.Error())
	// 	http.Error(res, spreadSheetErr.Error(), http.StatusInternalServerError)
	// }
	// output, _ := ioutil.ReadAll(spreadSheetResponse.Body)
	// fmt.Println(string(output[:]), "yayayaya")
	//json.NewEncoder(res).Encode(Output{New: "Biryani", Old: "Karahi"})
}
