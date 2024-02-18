package handlers

import (
	"github.com/peifengll/go_809_converter/converter/handlers/senders"
)

var CsCenter *csCenter

type csCenter struct {
	ShowSql       bool
	Interrupted   bool
	VerifyCode    int
	Uwriter       *senders.UpLinkWriter
	LongStopCache map[interface{}]interface{}
}

func InitCeCenter() {
	CsCenter = &csCenter{
		ShowSql:       false,
		Interrupted:   false,
		VerifyCode:    0,
		Uwriter:       nil,
		LongStopCache: make(map[interface{}]interface{}),
	}
}
