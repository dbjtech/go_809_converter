package service

import (
	"gorm.io/gorm"
)

// TerminalService 结构
type TerminalService struct {
}

// NewTerminalService 初始化 TerminalService
func NewTerminalService(db *gorm.DB) *TerminalService {
	return &TerminalService{}
}

func (t *TerminalService) GetTerminalByTid() {

}
