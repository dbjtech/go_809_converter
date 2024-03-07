package repository

import (
	"github.com/peifengll/go_809_converter/internal/model"
	"gorm.io/gorm"
	"log"
)

type CorpRepo interface {
	GetCorpByCid(cid string) *model.TCorp
}

func NewCorpRepo(db *gorm.DB) CorpRepo {
	return &corpRepo{db: db}
}

type corpRepo struct {
	db *gorm.DB
}

func (r *corpRepo) GetCorpByCid(cid string) *model.TCorp {
	crop := model.TCorp{}
	err := r.db.Model(&model.TCorp{}).Where("cid = ?", cid).First(&crop).Error
	if err != nil {
		log.Println(err)
		return nil
	}
	return &crop

}
