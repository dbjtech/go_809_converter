package service

import (
	"github.com/peifengll/go_809_converter/internal/model"
	"github.com/peifengll/go_809_converter/internal/repository"
	"log"
)

// TerminalService 结构
type TerminalService struct {
	repo repository.TerminalRepo
}

func (t *TerminalService) GetTerminalByTid(tid string) *model.TTerminalInfo {
	terminal := t.repo.GetTerminalByTid(tid)
	if terminal == nil {
		log.Printf("terminal %s is miss \n", tid)
		return nil
	}
	return terminal
}
