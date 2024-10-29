package exchange

import "github.com/dbjtech/go_809_converter/libs/packet_util"

var ThirdPartyDataQueue = make(chan string, 1000)
var DownLinkVerifyCode uint32
var UpLinkDataQueue = make(chan packet_util.MessageWrapper, 1000)
var ConverterWorker = 0
