package helpers

import (
	"github.com/peifengll/go_809_converter/internal/model"
	"github.com/peifengll/go_809_converter/internal/repository"
)

type CarHelper struct {
	carRepo repository.CarRepoInterface
}

func (ch *CarHelper) GetCarInfoByCarID(carID string) *model.TCar {
	car := ch.carRepo.GetCarInfoByCarID(carID)
	return car
}

func (ch *CarHelper) SetFuelCutByCNum(cnum string, optType string, payload int8) *model.TTerminalInfo {
	return ch.carRepo.UpdateFuelCutByCNum(cnum, optType, payload)
}

func (ch *CarHelper) GetSettingsByCNum(cnum string) *model.TTerminalInfo {
	return ch.carRepo.GetSettingsByCNum(cnum)
}
