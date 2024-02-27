package converters

import (
	"fmt"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/internal/model"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
)

type LocationConverter struct {
	*BaseConverter
}

func (c *LocationConverter) Convert(item string) ([]byte, error) {
	locations := gjson.Get(item, "res.location")
	length := int(locations.Get("#").Int())
	var carInfo *model.TCar = nil

	for i := 0; i < length; i++ {
		location := locations.Get(strconv.Itoa(i))
		log.Printf("convert %d/%d\n", i+1, length)
		termStatus := po.NormalStatus()
		alarm := 0
		timestamp := location.Get("timestamp").Int()
		if timestamp == 0 {
			continue
		}
		if carInfo == nil {
			carId := location.Get("car_id").String()

		}

	}
	return nil, fmt.Errorf("Convert method is not implemented")
}
