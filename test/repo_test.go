package test

import (
	"fmt"
	"github.com/peifengll/go_809_converter/config"
	"github.com/peifengll/go_809_converter/internal/repository"
	"testing"
)

func TestGetCarAndTermBySn(t *testing.T) {
	config.Path = "../conf/global.json"
	config.Load()
	db := config.NewDB()
	repo := repository.NewTerminalRepo(db)
	sn := repo.GetCarAndTermBySn("3D4C3006D6")
	fmt.Printf("%#v", sn)
}
