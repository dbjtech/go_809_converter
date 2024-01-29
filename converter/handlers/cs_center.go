package handlers

var CsCenter *csCenter

type csCenter struct {
	ShowSql       bool
	Interrupted   bool
	Verify_code   int
	Uwriter       interface{}
	LongStopCache map[interface{}]interface{}
}

func InitCeCenter() {
	CsCenter = &csCenter{
		ShowSql:       false,
		Interrupted:   false,
		Verify_code:   0,
		Uwriter:       nil,
		LongStopCache: make(map[interface{}]interface{}),
	}
}
