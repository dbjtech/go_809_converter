package exchange

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-18 20:52:48
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-29 19:26:00
 * @FilePath: \go_809_converter\exchange\center.go
 * @Description:
 *
 */

import (
	"github.com/dbjtech/go_809_converter/libs/packet_util"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var ThirdPartyDataQueue = make(chan string, 1000)
var DownLinkVerifyCode uint32
var UpLinkDataQueue = make(chan packet_util.MessageWrapper, 1000)
var ConverterWorker = 0
var TaskMarker = cmap.New[int64]()
