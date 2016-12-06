package spreadsheet

import (
	"encoding/json"
	"fmt"
	"go-careaxiom-services/whats-for-lunch/menu"
	"net/http"
	"strconv"
	"sync"

	"github.com/umar-muneer/go-careaxiom-utilities/authentication"
)

var apiURL = "https://sheets.googleapis.com/v4/spreadsheets"

var mutex sync.Mutex

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
	oldMenu := menu.New(menu.OLDMENUTYPE, spreadSheetClient)
	oldMenuEntry, oldMenuEntryErr := oldMenu.GetMenuEntry(req.URL.Query().Get("date"))
	if oldMenuEntryErr != nil {
		fmt.Println("Error while getting old menu entry -> ", oldMenuEntryErr)
		http.Error(res, oldMenuEntryErr.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("retrieving data for new menu")
	newMenu := menu.New(menu.NEWMENUTYPE, spreadSheetClient)
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

/*HandleReview review the day's menu*/
func HandleReview(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "post":
	case "POST":
		mutex.Lock()
		fmt.Println("posting lunch review")
		date := req.FormValue("date")
		menuType := req.FormValue("menuType")
		score, scoreErr := strconv.ParseFloat(req.FormValue("score"), 10)
		if menuType != menu.NEWMENUTYPE && menuType != menu.OLDMENUTYPE {
			errorString := "incorrect menu type specified"
			fmt.Println(errorString, " "+menuType)
			http.Error(res, errorString, http.StatusBadRequest)
			mutex.Unlock()
			return
		}
		if scoreErr != nil {
			fmt.Println("error in score parsing")
			http.Error(res, scoreErr.Error(), http.StatusBadRequest)
			mutex.Unlock()
			return
		}
		if date == "" {
			fmt.Println("no date found")
			http.Error(res, "no date found", http.StatusBadRequest)
			mutex.Unlock()
			return
		}
		spreadSheetClient, spreadSheetClientErr := authentication.GetClient()
		if spreadSheetClientErr != nil {
			fmt.Println(spreadSheetClientErr)
			http.Error(res, spreadSheetClientErr.Error(), http.StatusInternalServerError)
			mutex.Unlock()
			return
		}
		selectedMenu := menu.New(menuType, spreadSheetClient)
		currentScore, reviewErr := selectedMenu.PostReview(date, score)
		if reviewErr != nil {
			fmt.Println(reviewErr)
			http.Error(res, reviewErr.Error(), http.StatusInternalServerError)
			mutex.Unlock()
			return
		}
		json.NewEncoder(res).Encode(currentScore)
		mutex.Unlock()
		break
	case "get":
	case "GET":
		date := req.URL.Query().Get("date")
		fmt.Println("retrieving score for date ->", date)
		if date == "" {
			var errString = "incorrect date specified"
			fmt.Println(errString)
			http.Error(res, errString, http.StatusInternalServerError)
			return
		}
		menuType := req.URL.Query().Get("menuType")
		fmt.Println("retrieving score for menu type -> ", menuType)
		if menuType != menu.NEWMENUTYPE && menuType != menu.OLDMENUTYPE {
			var errString = "menu type should either be new or old"
			fmt.Println(errString)
			http.Error(res, errString, http.StatusInternalServerError)
			return
		}
		spreadSheetClient, spreadSheetClientErr := authentication.GetClient()
		if spreadSheetClientErr != nil {
			fmt.Println(spreadSheetClientErr)
			http.Error(res, spreadSheetClientErr.Error(), http.StatusInternalServerError)
			return
		}
		selectedMenu := menu.New(menuType, spreadSheetClient)
		score, err := selectedMenu.GetScore(date)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(score)
		break
	}
}

/*GetScore get score of a menu for a particular date*/
func GetScore(res http.ResponseWriter, req *http.Request) {
	switch req.Method {

	}
}
