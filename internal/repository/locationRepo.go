package repository

import (
	"github.com/peifengll/go_809_converter/internal/model"
	"gorm.io/gorm"
	"log"
)

type LocationRepo interface {
	SaveNewLocation(data *model.TLocation)
}

func NewLocationRepo(db *gorm.DB) LocationRepo {
	return &locationRepo{db: db}
}

type locationRepo struct {
	db *gorm.DB
}

func (l *locationRepo) GetLocationByLat(lat int64) string {
	k := model.TLocation{}
	err := l.db.Model(&model.TLocation{}).Select("address").
		Where(" id<=200000 and lat = ?", lat).
		First(&k).Error
	if err != nil {
		log.Println(err)
	}
	return k.Address
}

func (l *locationRepo) SaveNewLocation(data *model.TLocation) {
	l.db.Model(&model.TLocation{}).Create(data)
}
