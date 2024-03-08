package service

import (
	"encoding/json"
	"fmt"
	"github.com/peifengll/go_809_converter/internal/helpers"
	"github.com/peifengll/go_809_converter/internal/model"
	"github.com/peifengll/go_809_converter/libs/constants/terminal"
	"github.com/peifengll/go_809_converter/libs/constants/ucmtiResult"
	"log"
	"net/http"
)

var CarServiceObj *CarService

type CarService struct {
	ch *helpers.CarHelper
}

func (cs *CarService) GetCarInfoByCarID(carID string) *model.TCar {
	car := cs.ch.GetCarInfoByCarID(carID)
	return car
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
		log.Println("Error sending request:", err)
		return ucmtiResult.FAILURE
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var responseData map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&responseData)
		if err != nil {
			log.Println("Error decoding response:", err)
			return ucmtiResult.FAILURE
		}

		if responseData["status"].(int) != 200 {
			log.Println("Control error:", responseData)
			return ucmtiResult.FAILURE
		}
		log.Println("Control ok")
	} else {
		log.Println("Control HTTP error status code:", resp.StatusCode)
		return ucmtiResult.FAILURE
	}

	return ucmtiResult.SUCCESS
}

func (cs *CarService) LoadCarSettings(cnum string) map[string]interface{} {
	panic("not implemented")
	// Implement logic to load car settings by cnum
	return nil
}
