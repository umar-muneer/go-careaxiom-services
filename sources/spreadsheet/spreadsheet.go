package spreadsheet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"whats-for-lunch/authentication"
)

/*Menu representation of a menu for a day*/
type Menu struct {
	New string
	Old string
}

/*GetTomorrowsMenu returns tomorrow's menu*/
func GetTomorrowsMenu(res http.ResponseWriter, req *http.Request) {
	fmt.Println("getting tomorrow's menu")
	_, err := authentication.GetClient()
	if err != nil {
		fmt.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(res).Encode(Menu{New: "Biryani", Old: "Karahi"})
}
