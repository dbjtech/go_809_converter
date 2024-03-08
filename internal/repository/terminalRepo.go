package repository

import (
	"errors"
	"github.com/peifengll/go_809_converter/internal/model"
	"github.com/peifengll/go_809_converter/libs/constants/terminal"
	"gorm.io/gorm"
	"log"
)

type TerminalRepo interface {
	GetTidBy808Sn(sn808 string) string
	GetTerminalExtendByTid(tid string) *model.TTerminalExtend
	GetTerminalByTid(tid string) *model.TTerminalInfo
	GetTerminalBySN(sn string) *model.TTerminalInfo
	GetWorkStatusBySn(sn string) *model.TTerminalInfo
	GetCarAndTermBySn(sn string) *model.CarAndTerminal
	GetCarStatusBySn(sn string) *model.CarStatus
	RecordEventHistory(sn, name, remark, locateErrorInfo string, eventType int16) error
	ExeStatus(code int) (int, bool)
}

type terminalRepo struct {
	db *gorm.DB
}

func NewTerminalRepo(db *gorm.DB) TerminalRepo {
	return &terminalRepo{db: db}
}

func (r *terminalRepo) GetTidBy808Sn(sn808 string) string {
	var terminal model.TTerminalInfo
	if err := r.db.Table("t_terminal_info").Select("t_terminal_info.tid").
		Joins("JOIN t_808_sn ON t_terminal_info.sn = t_808_sn.sn").
		Where("t_808_sn.eze_identify_code = ? AND t_terminal_info.unbind = 0", sn808).
		First(&terminal).Error; err != nil {
		return terminal.Tid
	}
	return ""
}
func (r *terminalRepo) GetTerminalExtendByTid(tid string) *model.TTerminalExtend {
	var terminalExtend model.TTerminalExtend
	if err := r.db.Select("tid,producer_id, terminal_version, terminal_id, vehicle_identification_number").Where("tid = ?", tid).First(&terminalExtend).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			log.Println(err)
		}
	}
	return &terminalExtend
}
func (r *terminalRepo) GetTerminalByTid(tid string) *model.TTerminalInfo {
	terminal := model.TTerminalInfo{}
	err := r.db.Model(&terminal).
		Where("tid = ?", tid).
		Select("tid", "wired_fuel_status", "wired_fuel_exe_status",
			"sn", "dormant_fuel_status", "acc", "dormant_fuel_exe_status",
			"config_change", "alarm", "charge_status", "device_mode").
		First(&terminal).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
		}
		return nil
	}
	return &terminal
}
func (r *terminalRepo) GetTerminalBySN(sn string) *model.TTerminalInfo {
	terminal := model.TTerminalInfo{}
	err := r.db.Model(&terminal).
		Where("sn = ?", sn).
		Select("id", "tid", "iccid", "cid", "firmware_version",
			"firmware_c", "firmware_l", "activate_time", "imsi").
		First(&terminal).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
		}
		return nil
	}

	return &terminal
}
func (r *terminalRepo) GetWorkStatusBySn(sn string) *model.TTerminalInfo {
	t := model.TTerminalInfo{}
	if err := r.db.Model(&t).Where("sn = ?", sn).
		Select("id,tid,status,gsm,temp,gps,last_pkt_time").
		First(&t).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
		}
		return nil
	}
	return &t
}
func (r *terminalRepo) GetCarAndTermBySn(sn string) *model.CarAndTerminal {
	terminalAndCar := model.CarAndTerminal{}

	err := r.db.Table("t_terminal_info as T").
		Select("T.id, T.cid, T.tid, T.status, C.car_id,T.gps_lid, T.login, C.vin, C.cnum, C.id as t_car_id,T.gps, T.last_pkt_time,T.last_position_type").
		Joins("JOIN t_car_terminal ct ON t.tid = ct.tid").
		Joins("JOIN t_car C ON ct.car_id = c.car_id").
		Where("t.sn = ?", sn).
		First(&terminalAndCar).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
		}
		return nil
	}
	return &terminalAndCar
}

func (r *terminalRepo) GetCarStatusBySn(sn string) *model.CarStatus {
	carStatus := model.CarStatus{}
	err := r.db.Table("t_terminal_info as t").
		Select("t.id,t.tid,t.status,c.car_id,t.gps_lid,c.vin,c.cnum,c.id as t_car_id,t.gps,t.last_pkt_time,t.last_position_type").
		Joins("JOIN t_car_terminal ct ON t.tid = ct.tid").
		Joins("JOIN t_car c ON ct.car_id = c.car_id").
		Where("t.sn = ?", sn).
		First(&carStatus).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
		}
		return nil
	}
	return &carStatus
}

func (r *terminalRepo) RecordEventHistory(sn, name, remark, locateErrorInfo string, eventType int16) error {
	panic("not  implemented")
}

func (r *terminalRepo) ExeStatus(code int) (int, bool) {
	if code == terminal.TerminalExe.SUCCESS {
		return terminal.SwitchStatus["worked"], true
	}
	return terminal.SwitchStatus["sent"], false
}
