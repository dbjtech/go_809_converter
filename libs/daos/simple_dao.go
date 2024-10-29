package daos

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-25 20:07:48
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-29 19:41:34
 * @FilePath: \go_809_converter\libs\daos\simple_dao.go
 * @Description:
 *
 */

import (
	"context"

	"github.com/dbjtech/go_809_converter/libs/database/mysqldb"
	"github.com/dbjtech/go_809_converter/libs/models"
	"gorm.io/gorm"
)

type SimpleDao struct {
	db *gorm.DB
}

func NewSimpleDao() *SimpleDao {
	return &SimpleDao{db: mysqldb.GormDB}
}

func (sd *SimpleDao) GetCarInfo(ctx context.Context, carId string) (car models.CarInfo) {
	sd.db.WithContext(ctx).Raw("select c.car_id,c.vin,c.cnum,c.plate_color from t_car c where car_id = ?", carId).Scan(&car)
	return
}

func (sd *SimpleDao) GetTerminalInfo(ctx context.Context, tid string) (terminal models.TerminalInfo) {
	sd.db.WithContext(ctx).Raw("select t.sn,t.device_mode,t.charge_status,t.acc,"+
		"t.alarm, t.wired_fuel_status, t.dormant_fuel_status from t_terminal_info t where tid = ?", tid).Scan(&terminal)
	return
}

func (sd *SimpleDao) GetCorp(ctx context.Context, cid string) (corp models.Corp) {
	sd.db.WithContext(ctx).Raw("select c.id, c.name, c.mobile, c.name_show "+
		"from t_corp c where tid = ?", cid).Scan(&corp)
	return
}

func (sd *SimpleDao) GetTidBySn(ctx context.Context, sn []string) map[string]string {
	type tmp struct {
		Sn  string `json:"sn"`
		Tid string `json:"tid"`
	}
	result := make(map[string]string)
	var fromDb []tmp
	sd.db.WithContext(ctx).Raw("select sn, tid from t_terminal_info where sn in ?", sn).Scan(&fromDb)
	for _, v := range fromDb {
		result[v.Sn] = v.Tid
	}
	return result
}

func (sd *SimpleDao) GetCarIdByCnum(ctx context.Context, cnum []string) map[string]string {
	type tmp struct {
		Cnum  string `json:"cnum"`
		CarID string `json:"car_id"`
	}
	result := make(map[string]string)
	var fromDb []tmp
	sd.db.WithContext(ctx).Raw("select cnum, car_id from t_car where cnum in ?", cnum).Scan(&fromDb)
	for _, v := range fromDb {
		result[v.Cnum] = v.CarID
	}
	return result
}

func (sd *SimpleDao) GetCarIdByVin(ctx context.Context, vin []string) map[string]string {
	type tmp struct {
		Vin   string `json:"vin"`
		CarID string `json:"car_id"`
	}
	result := make(map[string]string)
	var fromDb []tmp
	sd.db.WithContext(ctx).Raw("select vin, car_id from t_car where vin in ?", vin).Scan(&fromDb)
	for _, v := range fromDb {
		result[v.Vin] = v.CarID
	}
	return result
}
