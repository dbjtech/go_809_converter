package repo

import "C"
import (
	"errors"
	"github.com/peifengll/go_809_converter/libs/constants/terminal"
	"gorm.io/gorm"
	"log"
)

var CarRepo = &carRepo{}

type carRepo struct {
	*gorm.DB
}

func (cr *carRepo) GetCarByCarID(carID string) *TCar {
	// todo  这里暂时还没有去做那个计时器
	t := &TCar{}
	cr.Where("car_id = ?", carID).First(t)
	return t

}

func (cr *carRepo) UpdateFuelCutByCNum(cnum string, optType string, payload int8) *TTerminalInfo {
	terminalj, err := cr.GetSettingsByCNum(cnum)
	if err != nil {
		log.Println(err)
		return nil
	}
	if terminalj == nil {
		return nil
	}
	switch optType {
	case "wired":
		terminalj.WiredFuelExpStatus = payload
		terminalj.WiredFuelExeStatus = int8(terminal.SwitchStatus["downlink"])
	case "dormant":
		terminalj.DormantFuelExpStatus = payload
		terminalj.DormantFuelExeStatus = int8(terminal.SwitchStatus["downlink"])
	}
	// todo 待测试，看更改会不会出问题
	err = cr.Model(&TTerminalInfo{}).Updates(terminalj).Error
	if err != nil {
		return nil
	}
	return terminalj
}

// 用来存连表查询到的数据
type TerminalInfoJ struct {
	Tid                  string
	WiredFuelStatus      int8
	WiredFuelExpStatus   int8
	ConfigChange         int8
	SN                   string
	ID                   int64
	WiredFuelExeStatus   int8
	DormantFuelStatus    int8
	DormantFuelExpStatus int8
	DormantFuelExeStatus int8
	FuelCutLock          int8
}

func (cr *carRepo) GetSettingsByCNum(cnum string) (*TTerminalInfo, error) {
	var terminalj TTerminalInfo
	if err := cr.Table("t_terminal_info").
		Select("t_terminal_info.tid, t_terminal_info.wired_fuel_status, t_terminal_info.wired_fuel_exp_status, "+
			"t_terminal_info.config_change, t_terminal_info.sn, t_terminal_info.id, t_terminal_info.wired_fuel_exe_status, "+
			"t_terminal_info.dormant_fuel_status, t_terminal_info.dormant_fuel_exp_status, "+
			"t_terminal_info.dormant_fuel_exe_status, t_terminal_info.fuel_cut_lock").
		Joins("LEFT OUTER JOIN t_car_terminals ON t_car_terminals.tid = t_terminal_info.tid").
		Joins("LEFT OUTER JOIN t_cars ON t_car_terminals.car_id = t_cars.car_id").
		Where("t_terminal_info.firmware_version = ?", terminal.ResFirmwareType["IV100"]).
		Where("t_cars.cnum = ? OR t_cars.vin = ?", cnum, cnum).
		First(&terminalj).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &terminalj, nil
}
