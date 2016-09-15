package menu

import (
	"fmt"
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
func (menu Menu) GetMenuEntry(date string) error {
	dayTime, err := time.Parse("01/1/2006", date)
	if err != nil {
		return err
	}
	column := strconv.Itoa(menu.sheetOffset + dayTime.Day())
	cellRange := "A" + string(column) + ":" + "E" + string(column)
	url := os.Getenv("SHEETS_API_URL") + "/" + menu.sheetID + "/values/" + cellRange
	response, responseErr := menu.client.Get(url)
	if responseErr != nil {
		return responseErr
	}
	output, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return readErr
	}
	fmt.Println(string(output[:]))
	return nil
}
