package leaves

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/umar-muneer/go-careaxiom-utilities/authentication"
)

type spreadSheetOutput struct {
	Range  string
	Values [][]string
}

type leaveStatus struct {
	EmailID      string
	EmployeeName string
	Total        float64
	Taken        float64
	Earned       float64
	Balance      float64
}

func createLeaveStatusMap(values [][]string) (result map[string]*leaveStatus) {
	result = make(map[string]*leaveStatus, 0)
	for i := 0; i < len(values); i++ {
		row := values[i]
		var total float64
		var taken float64
		var earned float64
		var balance float64
		if len(row) <= 2 {
			continue
		}
		fmt.Println("parsing row ", row)
		if len(row) >= 6 {
			total, _ = strconv.ParseFloat(row[5], 10)
		}
		if len(row) >= 8 {
			taken, _ = strconv.ParseFloat(row[7], 10)
		}
		if len(row) >= 9 {
			earned, _ = strconv.ParseFloat(row[8], 10)
		}
		if len(row) >= 10 {
			balance, _ = strconv.ParseFloat(row[9], 10)
		}
		result[row[0]] = &leaveStatus{
			EmailID:      row[0],
			EmployeeName: row[1],
			Total:        total,
			Taken:        taken,
			Balance:      balance,
			Earned:       earned,
		}
	}
	return result
}
func getLeavesStatus(employeeID string) (*leaveStatus, error) {
	fmt.Println("calculating leaves balance for " + string(employeeID))
	spreadSheetClient, spreadSheetClientError := authentication.GetClient()
	if spreadSheetClientError != nil {
		return nil, spreadSheetClientError
	}

	url := os.Getenv("SHEETS_API_URL") + "/" +
		os.Getenv("LEAVES_BALANCE_SPREADSHEET_ID") + "/values/" +
		os.Getenv("LEAVES_BALANCE_SHEET_NAME") + "!A5:J100"
	fmt.Println("url to retrieve leaves balance is -> ", url)

	spreadSheetRequest, _ := http.NewRequest("GET", url, nil)
	spreadSheetRequest.Header.Add("Accept", "application/json")

	response, responseErr := spreadSheetClient.Do(spreadSheetRequest)
	if responseErr != nil {
		return nil, responseErr
	}
	output, outputErr := ioutil.ReadAll(response.Body)
	if outputErr != nil {
		return nil, outputErr
	}

	spreadSheetOutput := new(spreadSheetOutput)
	marshalError := json.Unmarshal(output, spreadSheetOutput)
	if marshalError != nil {
		return nil, marshalError
	}

	leaveStatusMap := createLeaveStatusMap(spreadSheetOutput.Values)

	if _, present := leaveStatusMap[employeeID]; present == false {
		fmt.Println("Employee ID ", employeeID, " not found")
		return nil, nil
	}
	fmt.Println("result retrieved for employeeID "+employeeID, leaveStatusMap[employeeID])
	return leaveStatusMap[employeeID], nil
}

/*GetLeavesStatus get leaves balance of an employee*/
func GetLeavesStatus(res http.ResponseWriter, req *http.Request) {
	var employeeID = req.URL.Query().Get("employeeID")

	leaveStatus, leaveStatusErr := getLeavesStatus(employeeID)

	if leaveStatusErr != nil {
		fmt.Println("error while calculating leave balance for " + employeeID)
		http.Error(res, leaveStatusErr.Error(), http.StatusInternalServerError)
		return
	}
	if leaveStatus == nil {
		http.Error(res, "Employee ID "+employeeID+" not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(res).Encode(leaveStatus)
}

func deleteCache() {
	fmt.Println("Deleting Leaves Cache")
}

/*HandleCache method to handle leaves cache operation*/
func HandleCache(res http.ResponseWriter, req *http.Request) {
	if req.Method == "DELETE" {
		deleteCache()
	}
}
