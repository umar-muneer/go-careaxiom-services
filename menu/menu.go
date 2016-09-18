package menu

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

/*Entry represent a menu for a day*/
type Entry struct {
	MainDish      string
	SecondaryDish string
	Dessert       string
}

func (entry *Entry) newEntry(data []string) *Entry {
	var dessert = ""
	if len(data) > 4 {
		dessert = data[4]
	}

	return &Entry{
		MainDish:      data[2],
		SecondaryDish: data[3],
		Dessert:       dessert,
	}
}

/*SpreadSheetMenu can contain multiple menu entries*/
type SpreadSheetMenu struct {
	Entries     []Entry
	client      *http.Client
	sheetID     string
	sheetOffset int
}

type spreadSheetOutput struct {
	Range  string
	Values [][]string
}

/*New create new menu. can be old or new based on arguments*/
func New(client *http.Client, sheetID string, sheetOffset int) *SpreadSheetMenu {
	return &SpreadSheetMenu{
		client:      client,
		sheetID:     sheetID,
		sheetOffset: sheetOffset,
	}
}

/*GetMenuEntry gets a menu entry for a given date
  date required date in the format DD/MM/YYYY, start of day in Asia/Karachi timezone
*/
func (menu SpreadSheetMenu) GetMenuEntry(date string) (*Entry, error) {
	dayTime, err := time.Parse("01/01/2006", date)
	if err != nil {
		return new(Entry), err
	}
	column := strconv.Itoa(menu.sheetOffset + dayTime.Day())
	cellRange := "A" + string(column) + ":" + "E" + string(column)
	url := os.Getenv("SHEETS_API_URL") + "/" + menu.sheetID + "/values/" + cellRange
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
	menuEntry := &Entry{
		MainDish: spreadSheetOutput.Values[0][2],
	}

	return menuEntry, nil
}
