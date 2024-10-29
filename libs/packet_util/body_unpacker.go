package packet_util

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-18 20:52:48
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-27 10:21:51
 * @FilePath: \go_809_converter\libs\packet_util\body_unpacker.go
 * @Description:
 *
 */

import (
	"context"

	"github.com/dbjtech/go_809_converter/libs/constants"
	"github.com/linketech/microg/v4"
)

type bodyUnpacker = func(ctx context.Context, body []byte) MessageWithBody
type MessageWithBody interface {
	FromBytes(body []byte) error
	ToBytes() []byte
	String() string
}

var unpackPool = map[uint16]bodyUnpacker{
	constants.UP_CONNECT_REQ:                upLoginUnpacker,
	constants.UP_CONNECT_RSP:                upLoginRespUnpacker,
	constants.UP_LINKTEST_REQ:               emptyUnpacker,
	constants.UP_LINKTEST_RSP:               emptyUnpacker,
	constants.DOWN_LINKTEST_REQ:             emptyUnpacker,
	constants.DOWN_LINKTEST_RSP:             emptyUnpacker,
	constants.DOWN_CONNECT_REQ:              downLinkLoginUnpacker,
	constants.DOWN_CONNECT_RSP:              downLinkLoginRespUnpacker,
	constants.UP_EXG_MSG:                    upExgMsgUnpacker,
	constants.UP_EXG_MSG_REGISTER:           upExgMsgRegisterUnpacker,
	constants.UP_EXG_MSG_REAL_LOCATION:      realLocationUnpacker,
	constants.UP_EXG_MSG_TERMINAL_INFO:      carExtraInfoUnpacker,
	constants.UP_BASE_MSG:                   upBaseMsgUnpacker,
	constants.UP_BASE_MSG_VEHICLE_ADDED_ACK: upBaseMsgVehicleAddedUnpacker,
	constants.UP_WARN_MSG:                   upWarnMsgUnpacker,
	constants.UP_WARN_MSG_EXTENDS:           upWarnMsgExtendsUnpacker,
	//constants.DOWN_BASE_MSG_VEHICLE_ADDED: down_base_msg_vehicle_added_unpacker,
	//constants.DOWN_CTRL_MSG_TEXT_INFO: down_ctrl_msg_text_info_unpacker,
	//constants.DOWN_CTRL_MSG: down_ctrl_msg_unpacker,
	//constants.UP_CTRL_MSG_TEXT_INFO_ACK: up_ctrl_msg_text_info_unpacker,
	//constants.UP_CTRL_MSG: up_ctrl_msg_unpacker,
}

func upExgMsgUnpacker(ctx context.Context, body []byte) MessageWithBody {
	upExgMsg := newUpExgMsg()
	err := upExgMsg.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return upExgMsg
}

func upBaseMsgUnpacker(ctx context.Context, body []byte) MessageWithBody {
	upBaseMsg := newUpBaseMsg()
	err := upBaseMsg.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return upBaseMsg
}

func realLocationUnpacker(ctx context.Context, body []byte) MessageWithBody {
	upExgMsgRealLocation := newUpExgMsgRealLocation()
	err := upExgMsgRealLocation.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return upExgMsgRealLocation
}

func carExtraInfoUnpacker(ctx context.Context, body []byte) MessageWithBody {
	carExtraInfo := newCarExtraInfo()
	err := carExtraInfo.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return carExtraInfo
}

func downLinkLoginUnpacker(ctx context.Context, body []byte) MessageWithBody {
	downConnectReq := newDownConnectReq()
	err := downConnectReq.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return downConnectReq
}

func downLinkLoginRespUnpacker(ctx context.Context, body []byte) MessageWithBody {
	downConnectRsp := newDownConnectRsp()
	err := downConnectRsp.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return downConnectRsp
}

func upLoginUnpacker(ctx context.Context, body []byte) MessageWithBody {
	upConnectReq := newUpConnectReq()
	err := upConnectReq.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return upConnectReq
}

func upLoginRespUnpacker(ctx context.Context, body []byte) MessageWithBody {
	upConnectRsp := newUpConnectResp()
	err := upConnectRsp.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return upConnectRsp
}

func upExgMsgRegisterUnpacker(ctx context.Context, body []byte) MessageWithBody {
	upExgMsgRegister := newUpExgMsgRegister()
	err := upExgMsgRegister.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return upExgMsgRegister
}

func emptyUnpacker(ctx context.Context, body []byte) MessageWithBody {
	emptyBody := &EmptyBody{}
	return emptyBody
}

func upBaseMsgVehicleAddedUnpacker(ctx context.Context, body []byte) MessageWithBody {
	upBaseMsgVehicleAdded := newUpBaseMsgVehicleAdded()
	err := upBaseMsgVehicleAdded.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return upBaseMsgVehicleAdded
}

func upWarnMsgUnpacker(ctx context.Context, body []byte) MessageWithBody {
	upWarnMsg := newUpWarnMsg()
	err := upWarnMsg.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return upWarnMsg
}

func upWarnMsgExtendsUnpacker(ctx context.Context, body []byte) MessageWithBody {
	upWarnMsgExtends := newUpWarnMsgExtends()
	err := upWarnMsgExtends.FromBytes(body)
	if err != nil {
		microg.E(ctx, err)
		return nil
	}
	return upWarnMsgExtends
}
