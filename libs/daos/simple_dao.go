package daos

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

func (sd SimpleDao) GetCarInfo(ctx context.Context, carId string) (car models.CarInfo) {
	sd.db.WithContext(ctx).Raw("select c.car_id,c.vin,c.cnum,c.plate_color from t_car c where car_id = ?", carId).Scan(&car)
	return
}

func (sd SimpleDao) GetTerminalInfo(ctx context.Context, tid string) (terminal models.TerminalInfo) {
	sd.db.WithContext(ctx).Raw("select t.sn,t.device_mode,t.charge_status,t.acc,"+
		"t.alarm, t.wired_fuel_status, t.dormant_fuel_status from t_terminal_info t where tid = ?", tid).Scan(&terminal)
	return
}

func (sd SimpleDao) GetCorp(ctx context.Context, cid string) (corp models.Corp) {
	sd.db.WithContext(ctx).Raw("select c.id, c.name, c.mobile, c.name_show "+
		"from t_corp c where tid = ?", cid).Scan(&corp)
	return
}
