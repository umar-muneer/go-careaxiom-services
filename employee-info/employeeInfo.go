package employeeInfo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/umar-muneer/go-careaxiom-utilities/authentication"
)

type employeeInfo struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	BirthDate   string `json:"birthDate"`
	JoiningDate string `json:"joiningDate"`
}

func (info employeeInfo) String() string {
	return info.Name
}

type spreadSheetOutput struct {
	Range  string
	Values [][]string
}

func createEmployeeInfo(data []string) (result *employeeInfo) {
	info := new(employeeInfo)
	if len(data) >= 1 {
		info.Name = data[0]
	}
	if len(data) >= 2 {
		info.Email = data[1]
	}
	if len(data) >= 3 {
		info.BirthDate = data[2]
	}
	if len(data) >= 4 {
		info.JoiningDate = data[3]
	}
	return info
}
func createBirthdaysMap(data [][]string) (result map[string][]*employeeInfo) {
	fmt.Println("creating birthdays map")
	result = map[string][]*employeeInfo{}
	for i := 0; i < len(data); i++ {
		info := createEmployeeInfo(data[i])
		employees, _ := result[info.BirthDate]
		employees = append(employees, info)
		result[info.BirthDate] = employees
	}
	return result
}
func createAnniversariesMap(data [][]string) (result map[string]*employeeInfo) {
	return nil
}

func getBirthdaysAndAnniversariesFromSpreadsheet() ([][]string, error) {
	spreadSheetClient, err := authentication.GetClient()
	if err != nil {
		return nil, err
	}
	url := os.Getenv("SHEETS_API_URL") + "/" +
		os.Getenv("BIRTHDAY_ANNIVERSARIES_SPREADSHEET_ID") + "/values/" +
		os.Getenv("BIRTHDAY_ANNIVERSARIES_SHEET_NAME") + "!A1:D200"

	spreadSheetRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	spreadSheetRequest.Header.Add("Accept", "application/json")
	response, err := spreadSheetClient.Do(spreadSheetRequest)
	if err != nil {
		return nil, err
	}

	output, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	spreadSheetOutput := new(spreadSheetOutput)
	marshalErr := json.Unmarshal(output, spreadSheetOutput)
	if marshalErr != nil {
		return nil, marshalErr
	}
	return spreadSheetOutput.Values, nil
}

/*GetEmployeesWithBirthdays get employess with birthdays according to input */
func GetEmployeesWithBirthdays(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "get":
	case "GET":
		date := req.URL.Query().Get("date")
		if date == "" {
			fmt.Println("no date passed")
			http.Error(res, "no date passed", http.StatusBadRequest)
			return
		}
		output, err := getBirthdaysAndAnniversariesFromSpreadsheet()

		if err != nil {
			fmt.Println("error while reading birthday info from sheet")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		birthdays := createBirthdaysMap(output)
		fmt.Println("Birthday Employees On", date, " are-> ", birthdays[date])
		json.NewEncoder(res).Encode(birthdays[date])
		break
	}
}

/*GetEmployeesWithWorkAnniversaries get employees with anniversaries according to input */
func GetEmployeesWithWorkAnniversaries(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "get":
	case "GET":
		date := req.URL.Query().Get("date")
		if date == "" {
			fmt.Println("no date passed")
			http.Error(res, "no date passed", http.StatusBadRequest)
			return
		}
		fmt.Println("which employees have anniversaries today")
		break
	}
}
