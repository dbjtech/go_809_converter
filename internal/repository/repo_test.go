package repository

import (
	"fmt"
	"github.com/peifengll/go_809_converter/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
)

func TestLocationRepo_SaveNewLocation(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "peifeng", "127.0.0.1", 3306, "test",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	db = db.Debug()
	//db.AutoMigrate(&model.TLocation{})
	c := model.TLocation{
		CarId: "666666",
	}
	repo := NewLocationRepo(db)
	repo.SaveNewLocation(&c)
	if c.ID == 0 {
		log.Fatal()
	}
	fmt.Println(c)
}
