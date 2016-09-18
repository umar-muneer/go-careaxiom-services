package spreadsheet

import (
	"encoding/json"
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
	oldMenuEntry, oldMenuEntryErr := oldMenu.GetMenuEntry(req.URL.Query().Get("date"))
	if oldMenuEntryErr != nil {
		fmt.Println("Error while getting old menu entry -> ", oldMenuEntryErr)
		http.Error(res, oldMenuEntryErr.Error(), http.StatusInternalServerError)
		return
	}
	newOffset, _ := strconv.Atoi(os.Getenv("NEW_MENU_SHEET_OFFSET"))
	newMenu := menu.New(spreadSheetClient, os.Getenv("NEW_MENU_SPREADSHEET_ID"), newOffset)
	newMenuEntry, newMenuEntryErr := newMenu.GetMenuEntry(req.URL.Query().Get("date"))
	if newMenuEntryErr != nil {
		fmt.Println("Error while getting new menu entry", newMenuEntryErr)
		http.Error(res, newMenuEntryErr.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(Output{
		Old: *oldMenuEntry,
		New: *newMenuEntry,
	})
}
