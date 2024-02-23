package service

import (
	"encoding/json"
	"fmt"
	"github.com/peifengll/go_809_converter/libs/constants/terminal"
	"github.com/peifengll/go_809_converter/libs/constants/ucmtiResult"
	"github.com/peifengll/go_809_converter/libs/helpers"
	"log"
	"net/http"
)

type CarService struct {
}

func (cs *CarService) GetCarInfoByCarID(carID string) (map[string]interface{}, error) {
	helpers.CarHelper{}.GetCarInfoByCarID(carID)

	return nil, nil
}

func (cs *CarService) SwitchCarSettings(cnum string, sets string) int {
	var dSettings map[string]interface{}
	err := json.Unmarshal([]byte(sets), &dSettings)
	if err != nil {
		log.Println("Error parsing JSON:", err)
	}

	fuelCut := dSettings["fuel_cut"].(map[string]interface{})
	optType := fuelCut["type"].(string)
	payload := fuelCut["payload"].(string)
	optKey := fmt.Sprintf("%s%s", optType, payload)

	// Simulated example of using DownLinkControl.get(optKey)
	cmd := terminal.DownLinkControl[optKey]
	if cmd == "" {
		log.Println("Invalid command:", sets)
		return ucmtiResult.DENY
	}

	// Simulated example of setting fuel cut by cnum
	// terminal := CarHelper.setFuelCutByCnum(cnum, optType, payload)
	// if terminal == nil {
	// 	return UCMTI_RESULT.DENY
	// }

	// Simulated example of sending request to admin_url
	requestData := map[string]interface{}{
		"sn":     "terminal.sn",
		"cmd":    cmd,
		"remark": "上级平台",
	}
	requestBody, _ := json.Marshal(requestData)
	resp, err := http.Post("http://admin_url/downlink/terminal/control", "application/json", requestBody)
	if err != nil {
		logger.Println("Error sending request:", err)
		return UCMTI_RESULT.FAILURE
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var responseData map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&responseData)
		if err != nil {
			logger.Println("Error decoding response:", err)
			return UCMTI_RESULT.FAILURE
		}

		if responseData["status"].(int) != 200 {
			logger.Println("Control error:", responseData)
			return UCMTI_RESULT.FAILURE
		}
		logger.Println("Control ok")
	} else {
		logger.Println("Control HTTP error status code:", resp.StatusCode)
		return UCMTI_RESULT.FAILURE
	}

	return UCMTI_RESULT.SUCCESS
}

func (cs *CarService) LoadCarSettings(cnum string) (map[string]interface{}, error) {
	// Implement logic to load car settings by cnum
	return nil, nil
}
