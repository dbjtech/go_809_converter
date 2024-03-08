package converters

import (
	"github.com/peifengll/go_809_converter/internal/helpers"
	"github.com/peifengll/go_809_converter/internal/service"
	"github.com/peifengll/go_809_converter/libs/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ConverterInterface interface {
	Convert(item string) []byte
	Handle(item string) []byte
	SetTraceID(traceID string)
	GetTraceID() string
	BuildUpWarnExtends(warnCode int, cnum string, color byte, sn string) []byte
}

func NewRequestConverters(db *gorm.DB, rds *redis.Client,
	carService *service.CarService, corpHelper *helpers.CorpHelper,
	carIdWhitelist *utils.CarIdWhitelist,
	terminalService *service.TerminalService,
) map[string]ConverterInterface {
	base := baseConverter{
		dB:    db,
		redis: rds,
	}
	carEx := carExtraInfoConverter{
		baseConverter: &base,
		carService:    carService,
	}
	carInfo := carInfoConverter{
		baseConverter: &base,
		corpHelper:    corpHelper,
	}
	carRe := carRegisterConverter{&base}
	locationC := locationConverter{
		baseConverter:   &base,
		carIdWhitelist:  carIdWhitelist,
		carService:      carService,
		terminalService: terminalService,
	}
	online := onlineOfflineConverter{
		baseConverter:   &base,
		carService:      carService,
		terminalService: terminalService,
	}
	return map[string]ConverterInterface{
		"S13":  &locationC,
		"S99":  &carRe,
		"S10":  &online,
		"S106": &carEx,
		"S991": &carInfo,
	}
}
