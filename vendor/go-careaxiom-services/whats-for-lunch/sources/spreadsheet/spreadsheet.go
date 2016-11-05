package spreadsheet

import (
	"encoding/json"
	"fmt"
	"go-careaxiom-services/whats-for-lunch/menu"
	"net/http"
	"os"
	"strconv"

	"github.com/umar-muneer/go-careaxiom-utilities/authentication"
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
		return
	}
	fmt.Println("retrieving data for old menu")
	offset, _ := strconv.Atoi(os.Getenv("OLD_MENU_SHEET_OFFSET"))
	oldMenu := menu.New(spreadSheetClient, os.Getenv("OLD_MENU_SPREADSHEET_ID"), offset, os.Getenv("OLD_MENU_TITLE"))
	oldMenuEntry, oldMenuEntryErr := oldMenu.GetMenuEntry(req.URL.Query().Get("date"))
	if oldMenuEntryErr != nil {
		fmt.Println("Error while getting old menu entry -> ", oldMenuEntryErr)
		http.Error(res, oldMenuEntryErr.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("retrieving data for new menu")
	newOffset, _ := strconv.Atoi(os.Getenv("NEW_MENU_SHEET_OFFSET"))
	newMenu := menu.New(spreadSheetClient, os.Getenv("NEW_MENU_SPREADSHEET_ID"), newOffset, os.Getenv("NEW_MENU_TITLE"))
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

func postReview(menuType string, date string, score int) error {
	fmt.Println("posting review, Score= ", score, ", Date = ", date, ", Menu Type = ", menuType)
	return nil
}

/*HandleReview review the day's menu*/
func HandleReview(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "post":
	case "POST":
		fmt.Println("posting lunch review")
		date := req.FormValue("date")
		menuType := req.FormValue("menuType")
		score, scoreErr := strconv.Atoi(req.FormValue("score"))
		if menuType != menu.NEWMENUTYPE && menuType != menu.OLDMENUTYPE {
			errorString := "incorrect menu type specified"
			fmt.Println(errorString, " "+menuType)
			http.Error(res, errorString, http.StatusBadRequest)
			return
		}
		if scoreErr != nil {
			fmt.Println("error in score parsing")
			http.Error(res, scoreErr.Error(), http.StatusBadRequest)
			return
		}
		if date == "" {
			fmt.Println("no date found")
			http.Error(res, "no date found", http.StatusBadRequest)
			return
		}
		postError := postReview(menuType, date, score)
		if postError != nil {
			fmt.Println("error while posting review -> ", postError.Error())
			http.Error(res, postError.Error(), http.StatusInternalServerError)
			return
		}
	}
}
