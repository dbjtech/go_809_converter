package service

import (
	"github.com/peifengll/go_809_converter/internal/helpers"
	"github.com/peifengll/go_809_converter/internal/model"
)

type CorpService struct {
	cor helpers.CorpHelper
}

func (s *CorpService) GetCorpByCid(cid string) *model.TCorp {
	return s.cor.GetCorpInfoByCid(cid)
}
