package repository

import (
	"errors"
	"github.com/peifengll/go_809_converter/internal/model"
	"gorm.io/gorm"
	"log"
)

type TerminalRepo struct {
	db *gorm.DB
}

func NewTerminalRepo(db *gorm.DB) *TerminalRepo {
	return &TerminalRepo{db: db}
}

func (r *TerminalRepo) GetTidBy808Sn(sn808 string) string {
	var terminal model.TTerminalInfo
	if err := r.db.Table("t_terminal_info").Select("t_terminal_info.tid").
		Joins("JOIN t_808_sn ON t_terminal_info.sn = t_808_sn.sn").
		Where("t_808_sn.eze_identify_code = ? AND t_terminal_info.unbind = 0", sn808).
		First(&terminal).Error; err != nil {
		return terminal.Tid
	}
	return ""
}
func (r *TerminalRepo) GetTerminalExtendByTid(tid string) *model.TTerminalExtend {
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
func (r *TerminalRepo) GetTerminalByTid(tid string) *model.TTerminalInfo {
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
func (r *TerminalRepo) GetTerminalBySN(sn string) *model.TTerminalInfo {
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
func (r *TerminalRepo) GetWorkStatusBySn(sn string) *model.TTerminalInfo {
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
func (r *TerminalRepo) GetCarAndTermBySn(sn string) {

}
func (r *TerminalRepo) GetCarStatusBySn() {

}
func (r *TerminalRepo) RecordEventHistory() {

}
func (r *TerminalRepo) ExeStatus() {

}
