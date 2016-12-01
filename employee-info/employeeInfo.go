package employeeInfo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/umar-muneer/go-careaxiom-utilities/authentication"
)

type employeeInfo struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	birthDate   time.Time
	joiningDate time.Time
}

func (info employeeInfo) String() string {
	return info.Name
}
func getHashKey(date time.Time) string {
	return fmt.Sprintf("%d/%d", date.Day(), date.Month())
}

type spreadSheetOutput struct {
	Range  string
	Values [][]string
}

func createEmployeeInfo(data []string) (result *employeeInfo, err error) {
	var (
		birthDate   time.Time
		joiningDate time.Time
	)
	info := new(employeeInfo)
	if len(data) >= 1 {
		info.Name = data[0]
	}
	if len(data) >= 2 {
		info.Email = data[1]
	}
	if len(data) >= 3 {
		birthDate, err = time.Parse("02/01/2006", data[2])
	}
	if len(data) >= 4 {
		joiningDate, err = time.Parse("02/01/2006", data[3])
	}
	if err != nil {
		return nil, err
	}
	info.birthDate = birthDate
	info.joiningDate = joiningDate
	return info, nil
}
func createBirthdaysMap(data [][]string) (result map[string][]*employeeInfo, err error) {
	fmt.Println("creating birthdays map")
	result = map[string][]*employeeInfo{}
	for i := 0; i < len(data); i++ {
		info, err := createEmployeeInfo(data[i])
		if err != nil {
			return nil, err
		}
		key := getHashKey(info.birthDate)
		employees, _ := result[key]
		employees = append(employees, info)
		result[key] = employees
	}
	return result, nil
}
func createAnniversariesMap(data [][]string) (result map[string][]*employeeInfo, err error) {
	fmt.Println("creating anniversaries map")
	result = map[string][]*employeeInfo{}
	for i := 0; i < len(data); i++ {
		info, err := createEmployeeInfo(data[i])
		if err != nil {
			return nil, err
		}
		key := getHashKey(info.joiningDate)
		employees, _ := result[key]
		employees = append(employees, info)
		result[key] = employees
	}
	return result, nil
}

func getBirthdaysAndAnniversariesFromSpreadsheet() ([][]string, error) {
	fmt.Println("retrieving data from birthdays and anniversaries sheet")
	spreadSheetClient, err := authentication.GetClient()
	if err != nil {
		return nil, err
	}
	url := os.Getenv("SHEETS_API_URL") + "/" +
		os.Getenv("BIRTHDAY_ANNIVERSARIES_SPREADSHEET_ID") + "/values/" +
		os.Getenv("BIRTHDAY_ANNIVERSARIES_SHEET_NAME") + "!A2:D200"

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
		date, err := time.Parse("02/01/2006", req.URL.Query().Get("date"))
		if err != nil {
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
		birthdays, err := createBirthdaysMap(output)
		if err != nil {
			fmt.Println("error while creating birthday map", err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Employees with birthdays on", date, " are-> ", birthdays[getHashKey(date)])
		json.NewEncoder(res).Encode(birthdays[getHashKey(date)])
		break
	}
}

/*GetEmployeesWithWorkAnniversaries get employees with anniversaries according to input */
func GetEmployeesWithWorkAnniversaries(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "get":
	case "GET":
		date, err := time.Parse("02/01/2006", req.URL.Query().Get("date"))
		if err != nil {
			fmt.Println("no date passed")
			http.Error(res, "no date passed", http.StatusBadRequest)
			return
		}
		output, err := getBirthdaysAndAnniversariesFromSpreadsheet()
		if err != nil {
			fmt.Println("error while reading anniversary info from sheet", err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		anniversaries, err := createAnniversariesMap(output)
		if err != nil {
			fmt.Println("error while creating anniversary map", err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		key := getHashKey(date)
		fmt.Println("Employees with anniversaries on", date, " are-> ", anniversaries[key])
		json.NewEncoder(res).Encode(anniversaries[key])
		break
	}
}
