package service

import (
	"context"
	"errors"
	"sync"

	"github.com/dbjtech/go_809_converter/libs/daos"
	"github.com/dbjtech/go_809_converter/libs/models"
	"github.com/dbjtech/go_809_converter/metrics"
)

type SimpleService struct {
	dao *daos.SimpleDao
}

var CarTermService *SimpleService

var once sync.Once

func initService() *SimpleService {
	if CarTermService != nil {
		return CarTermService
	}
	once.Do(func() {
		CarTermService = &SimpleService{
			dao: daos.NewSimpleDao(),
		}
	})
	return CarTermService
}

func (s *SimpleService) GetCarInfoByCarID(ctx context.Context, carID string) (models.Car, error) {
	car := s.dao.GetCarInfo(ctx, carID)
	if car.CarID == "" {
		return models.Car{}, errors.New("car not found")
	}
	return models.Car{
		Vin:        car.Vin,
		Cnum:       car.Cnum,
		PlateColor: car.PlateColor,
	}, nil
}

func (s *SimpleService) GetTerminalInfoByTid(ctx context.Context, tid string) (models.Terminal, error) {
	terminal := s.dao.GetTerminalInfo(ctx, tid)
	if terminal.Sn == "" {
		return models.Terminal{}, errors.New("terminal not found")
	}
	return models.Terminal{
		Sn:                terminal.Sn,
		WiredFuelStatus:   uint8(terminal.WiredFuelStatus),
		DormantFuelStatus: uint8(terminal.DormantFuelStatus),
		Acc:               uint8(terminal.Acc),
		ChargeStatus:      uint8(terminal.ChargeStatus),
		DeviceMode:        uint8(terminal.DeviceMode),
		Alarm:             uint8(terminal.Alarm),
	}, nil
}

func (s *SimpleService) GetCorpNameByCid(ctx context.Context, cid string) (corpName string, err error) {
	corp := s.dao.GetCorp(ctx, cid)
	if corp.Id == 0 {
		return "", errors.New("corp not found")
	}
	return corp.Name, nil
}

func GetCarInfoByCarID(ctx context.Context, carID string) (models.Car, error) {
	if CarTermService == nil {
		initService()
	}
	metrics.MySQLQuery.WithLabelValues("car_info").Inc()
	return CarTermService.GetCarInfoByCarID(ctx, carID)
}

func GetTerminalInfoByTid(ctx context.Context, tid string) (models.Terminal, error) {
	if CarTermService == nil {
		initService()
	}
	metrics.MySQLQuery.WithLabelValues("terminal").Inc()
	return CarTermService.GetTerminalInfoByTid(ctx, tid)
}

func GetCorpNameByCid(ctx context.Context, cid string) (corpName string, err error) {
	if CarTermService == nil {
		initService()
	}
	metrics.MySQLQuery.WithLabelValues("corp").Inc()
	return CarTermService.GetCorpNameByCid(ctx, cid)
}
