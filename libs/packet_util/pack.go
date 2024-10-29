package packet_util

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-19 12:33:14
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-19 14:36:13
 * @FilePath: \go_809_converter\libs\packet_util\pack.go
 * @Description:
 *
 */

import "bytes"

func Pack(msg Message) []byte {
	end := byte(0x5d)
	rawBytes := msg.ToBytes()
	rawBytes = bytes.ReplaceAll(rawBytes, []byte{0x5a}, []byte{0x5a, 0x02})
	rawBytes = bytes.ReplaceAll(rawBytes, []byte{0x5b}, []byte{0x5a, 0x01})
	rawBytes = bytes.ReplaceAll(rawBytes, []byte{0x5e}, []byte{0x5e, 0x02})
	rawBytes = bytes.ReplaceAll(rawBytes, []byte{0x5d}, []byte{0x5e, 0x01})
	rawBytes = append([]byte{0x5b}, rawBytes...)
	return append(rawBytes, end)
}
