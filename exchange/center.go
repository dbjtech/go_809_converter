package exchange

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-18 20:52:48
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2025-11-14 20:56:26
 * @FilePath: \go_809_converter\exchange\center.go
 * @Description:
 *
 */

import (
	"sync"

	"github.com/dbjtech/go_809_converter/libs/packet_util"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type VerifyCode struct {
	CodeMap map[string]uint32
	Lock    sync.Mutex
}

func (vc *VerifyCode) Get(cvtName string) uint32 {
	vc.Lock.Lock()
	defer vc.Lock.Unlock()
	code, _ := vc.CodeMap[cvtName]
	return code
}

func (vc *VerifyCode) Set(cvtName string, code uint32) {
	vc.Lock.Lock()
	defer vc.Lock.Unlock()
	vc.CodeMap[cvtName] = code
}

func NewVerifyCode() *VerifyCode {
	return &VerifyCode{
		CodeMap: make(map[string]uint32),
	}
}

var ThirdPartyDataQueuePool = make(map[string]chan string)
var DownLinkVerifyCode = NewVerifyCode()
var UpLinkDataQueuePool = make(map[string]chan packet_util.MessageWrapper)
var ConverterWorker = 0
var TaskMarker = cmap.New[int64]()
var JtwConverterUpLinkDataQueuePool = make(map[string]chan packet_util.MessageWrapper)
