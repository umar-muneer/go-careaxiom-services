package menu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

/*Entry represent a menu for a day*/
type Entry struct {
	Title         string
	MainDish      string
	SecondaryDish string
	Dessert       string
	Score         float64
	ReviewCount   float64
}

const (
	/*NEWMENUTYPE represents our new menu type*/
	NEWMENUTYPE = "new"
	/*OLDMENUTYPE represents our new menu type*/
	OLDMENUTYPE = "old"
)

func newEntry(data []string, title string) *Entry {
	var dessert = ""
	var secondaryDish = ""
	var mainDish = ""
	var score float64
	var reviewCount float64

	if len(data) >= 3 {
		mainDish = data[2]
	}
	if len(data) >= 4 {
		secondaryDish = data[3]
	}
	if len(data) >= 5 {
		dessert = data[4]
	}
	if len(data) >= 6 {
		score, _ = strconv.ParseFloat(data[5], 10)
	}
	if len(data) >= 7 {
		reviewCount, _ = strconv.ParseFloat(data[6], 10)
	}

	return &Entry{
		Title:         title,
		MainDish:      mainDish,
		SecondaryDish: secondaryDish,
		Dessert:       dessert,
		Score:         score,
		ReviewCount:   reviewCount,
	}
}

/*SpreadSheetMenu can contain multiple menu entries*/
type SpreadSheetMenu struct {
	Entries     []Entry
	client      *http.Client
	sheetID     string
	sheetOffset int
	title       string
}

type spreadSheetOutput struct {
	Range  string
	Values [][]string
}

/*New create new menu. can be old or new based on arguments*/
func New(menuType string, client *http.Client) *SpreadSheetMenu {
	var (
		sheetID = ""
		offset  int
		title   = ""
	)
	if menuType == OLDMENUTYPE {
		offset, _ = strconv.Atoi(os.Getenv("OLD_MENU_SHEET_OFFSET"))
		title = os.Getenv("OLD_MENU_TITLE")
		sheetID = os.Getenv("OLD_MENU_SPREADSHEET_ID")
	} else {
		offset, _ = strconv.Atoi(os.Getenv("NEW_MENU_SHEET_OFFSET"))
		title = os.Getenv("NEW_MENU_TITLE")
		sheetID = os.Getenv("NEW_MENU_SPREADSHEET_ID")
	}
	return &SpreadSheetMenu{
		client:      client,
		sheetID:     sheetID,
		sheetOffset: offset,
		title:       title,
	}
}

func (menu SpreadSheetMenu) getMenuEntryRow(date time.Time) string {
	return strconv.Itoa(menu.sheetOffset + date.Day())
}

/*GetMenuEntry gets a menu entry for a given date
  date required date in the format MMM/DD/YYYY, start of day in Asia/Karachi timezone
*/
func (menu SpreadSheetMenu) GetMenuEntry(date string) (*Entry, error) {

	if menu.sheetID == "" {
		fmt.Println("Menu disabled for now, the caterer is probably on leave")
		return &Entry{Title: menu.title}, nil
	}
	dayTime, err := time.Parse("02/01/2006", date)
	fmt.Println("requested day time is", dayTime)
	if err != nil {
		return new(Entry), err
	}
	row := menu.getMenuEntryRow(dayTime)
	cellRange := "A" + row + ":" + "G" + row
	url := os.Getenv("SHEETS_API_URL") + "/" + menu.sheetID + "/values/" + cellRange
	fmt.Println("day of month is", dayTime.Day())
	fmt.Println("url to get menu is -> ", url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	response, responseErr := menu.client.Do(req)
	if responseErr != nil {
		return new(Entry), responseErr
	}
	output, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return new(Entry), readErr
	}
	spreadSheetOutput := new(spreadSheetOutput)
	marshalError := json.Unmarshal(output, spreadSheetOutput)
	if marshalError != nil {
		return new(Entry), marshalError
	}
	menuEntry := newEntry(spreadSheetOutput.Values[0], menu.title)
	return menuEntry, nil
}

type reviewBody struct {
	MajorDimension string      `json:"majorDimension"`
	Range          string      `json:"range"`
	Values         [][]float64 `json:"values"`
}

/*PostReview post a review through this method*/
func (menu SpreadSheetMenu) PostReview(date string, score float64) (float64, error) {
	fmt.Println("posting review, Score= ", score, ", Date = ", date)
	dayTime, err := time.Parse("02/01/2006", date)
	if score < 0 {
		score = 0
	} else if score > 5 {
		score = 5
	}
	if err != nil {
		return 0, err
	}
	entry, _ := menu.GetMenuEntry(date)
	fmt.Printf("review count is: %f", entry.ReviewCount)

	newReviewCount := entry.ReviewCount + 1
	newTotalScore := ((entry.ReviewCount * entry.Score) + score) / newReviewCount

	row := menu.getMenuEntryRow(dayTime)
	cellRange := "F" + row + ":" + "G" + row
	url := os.Getenv("SHEETS_API_URL") + "/" + menu.sheetID + "/values/" + cellRange + "?valueInputOption=RAW"
	fmt.Println("url for posting review score -> " + url)

	reviewData := &reviewBody{
		MajorDimension: "ROWS",
		Values:         [][]float64{{newTotalScore, newReviewCount}},
		Range:          cellRange,
	}
	postBody, marshalErr := json.Marshal(reviewData)
	if marshalErr != nil {
		return 0, marshalErr
	}

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(postBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	_, responseErr := menu.client.Do(req)
	if responseErr != nil {
		return 0, responseErr
	}
	fmt.Printf("final score is %f", newTotalScore)
	fmt.Printf("Total No. of reviewers is %f", newReviewCount)
	return newTotalScore, nil
}
