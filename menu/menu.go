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
	MainDish     string
	SecondarDish string
	Dessert      string
}

/*Menu can contain multiple menu entries*/
type Menu struct {
	Entries     []Entry
	client      *http.Client
	sheetID     string
	sheetOffset int
}

/*SpreadSheetOutput contains mapped json output*/
type SpreadSheetOutput struct {
	Range  string
	Values [][]string
}

/*New create new menu. can be old or new based on arguments*/
func New(client *http.Client, sheetID string, sheetOffset int) *Menu {
	return &Menu{
		client:      client,
		sheetID:     sheetID,
		sheetOffset: sheetOffset,
	}
}

/*GetMenuEntry gets a menu entry for a given date
  date required date in the format DD/MM/YYYY, start of day in Asia/Karachi timezone
*/
func (menu Menu) GetMenuEntry(date string) (*SpreadSheetOutput, error) {
	dayTime, err := time.Parse("01/01/2006", date)
	if err != nil {
		return new(SpreadSheetOutput), err
	}
	column := strconv.Itoa(menu.sheetOffset + dayTime.Day())
	cellRange := "A" + string(column) + ":" + "E" + string(column)
	url := os.Getenv("SHEETS_API_URL") + "/" + menu.sheetID + "/values/" + cellRange
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	response, responseErr := menu.client.Do(req)
	if responseErr != nil {
		return new(SpreadSheetOutput), responseErr
	}
	output, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return new(SpreadSheetOutput), readErr
	}
	spreadSheetOutput := new(SpreadSheetOutput)
	marshalError := json.Unmarshal(output, spreadSheetOutput)
	if marshalError != nil {
		return new(SpreadSheetOutput), marshalError
	}
	return spreadSheetOutput, nil
}
