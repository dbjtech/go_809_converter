package converters

import (
	"fmt"
	"github.com/peifengll/go_809_converter/converter/handlers/po"
	"github.com/peifengll/go_809_converter/libs/constants/businessType"
	"github.com/peifengll/go_809_converter/libs/pack"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
)

type BaseConverter struct {
	dB      *gorm.DB
	redis   *redis.Client
	TraceID string
}

func NewBaseConverter(db *gorm.DB) *BaseConverter {
	return &BaseConverter{
		dB: db,
	}
}

func (c *BaseConverter) Convert(item string) []byte {
	log.Println(fmt.Errorf("Convert method is not implemented"))
	return nil
}

func (c *BaseConverter) Handle(item string) []byte {
	return c.Convert(item)
}

func (c *BaseConverter) SetTraceID(traceID string) {
	c.TraceID = traceID
}

func (c *BaseConverter) GetTraceID() string {
	return c.TraceID
}

func (c *BaseConverter) BuildUpWarnExtends(warnCode int, cnum string, color byte, sn string) []byte {
	btype := businessType.UP_WARN_MSG_EXTENDS
	var data string
	if sn != "" {
		data = fmt.Sprintf(`{"src": "DBJ", "warn_code": %d, "sn": "%s"}`, warnCode, sn)
	} else {
		data = fmt.Sprintf(`{"src": "DBJ", "warn_code": %d}`, warnCode)
	}
	upDict := po.UpWarnMsgExtends{
		VehicleNo:    cnum,
		VehicleColor: color,
		DataType:     uint16(btype),
		DataLength:   0,
		Data:         data,
	}
	return pack.BuildMessageP(btype, upDict.Encode(), 0)
}
