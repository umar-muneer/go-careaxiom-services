package menu

import (
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
}

func newEntry(data []string, title string) *Entry {
	var dessert = ""
	var secondaryDish = ""
	var mainDish = ""

	if len(data) > 2 {
		mainDish = data[2]
	}
	if len(data) > 3 {
		secondaryDish = data[3]
	}
	if len(data) > 4 {
		dessert = data[4]
	}

	return &Entry{
		Title:         title,
		MainDish:      mainDish,
		SecondaryDish: secondaryDish,
		Dessert:       dessert,
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
func New(client *http.Client, sheetID string, sheetOffset int, title string) *SpreadSheetMenu {
	return &SpreadSheetMenu{
		client:      client,
		sheetID:     sheetID,
		sheetOffset: sheetOffset,
		title:       title,
	}
}

/*GetMenuEntry gets a menu entry for a given date
  date required date in the format MMM/DD/YYYY, start of day in Asia/Karachi timezone
*/
func (menu SpreadSheetMenu) GetMenuEntry(date string) (*Entry, error) {
	dayTime, err := time.Parse("02/01/2006", date)
	fmt.Println("requested day time is", dayTime)
	if err != nil {
		return new(Entry), err
	}
	column := strconv.Itoa(menu.sheetOffset + dayTime.Day())
	cellRange := "A" + string(column) + ":" + "E" + string(column)
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
