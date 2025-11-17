package receivers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	"github.com/dbjtech/go_809_converter/exchange"
	"github.com/dbjtech/go_809_converter/libs"
	"github.com/dbjtech/go_809_converter/libs/constants"
	"github.com/dbjtech/go_809_converter/metrics"
	"github.com/gookit/config/v2"

	"github.com/linketech/microg/v4"
)

type receiveBuffer struct {
	header    []byte
	buf       []byte
	size      int
	headMatch uint8
	left      uint8
}

// add 添加一个字符, 缓存起来，遇到头部信息，才开始添加到缓存
func (r *receiveBuffer) add(char byte) {
	if r.size == 0 { // 等匹配上头部信息之后才真正添加内容
		// 如果 char 是换行符或者空格，则直接跳过
		if char == '\n' || char == ' ' {
			return
		}
		if char == r.header[r.headMatch] {
			r.headMatch += 1
			if r.headMatch == uint8(len(r.header)) { // 头部信息全匹配
				r.headMatch = 0
				for i, v := range r.header {
					r.buf[i] = v
				}
				r.size = len(r.header)
				r.left = 1
			}
		} else {
			r.headMatch = 0
		}
	} else {
		r.buf[r.size] = char
		r.size++
		if char == '{' {
			r.left += 1
		} else if char == '}' {
			r.left -= 1
		}
	}
}

// matched 判断是否已经匹配到完整的数据包
func (r *receiveBuffer) matched() bool {
	return r.size > 0 && r.left == 0
}

// reset 重置状态
func (r *receiveBuffer) reset() {
	r.size = 0
	r.left = 0
}

// flush 提取数据
func (r *receiveBuffer) flush() string {
	cache := string(r.buf[:r.size])
	r.size = 0
	r.left = 0
	return cache
}

func newReceiveBuffer() *receiveBuffer {
	return &receiveBuffer{
		header: []byte(`{"res":`),
		buf:    make([]byte, 51200),
	}
}

func StartThirdPartyReceiver(ctx context.Context, wg *sync.WaitGroup) {
	configConverter := config.SubDataMap(libs.Environment + ".converter")
	for key, v := range configConverter {
		m := v.(map[string]any)
		enable, ok := m["enable"].(bool)
		if !ok || !enable {
			microg.W("3rd Party Receiver %s is disabled", key)
			continue
		}
		go func(key string) {
			thirdpartPortKey := libs.Environment + ".converter." + key + ".thirdpartPort"
			port := config.Int(thirdpartPortKey, 11223)
			addr := fmt.Sprintf(":%d", port)
			listener, err := net.Listen("tcp", addr)
			innerCtx := context.WithValue(ctx, constants.TracerKeyCvtName, key)
			if err != nil {
				microg.W(innerCtx, "Failed to listen on %s: %v", addr, err)
			}
			defer listener.Close()
			microg.I(innerCtx, "3rd Party Server is listening on %s for %s", addr, key)
			for {
				// 接收连接
				conn, err := listener.Accept()
				if err != nil {
					microg.E(innerCtx, "Failed to accept connection: %v", err)
					continue
				}
				// 连接保持
				go handleConnection(innerCtx, wg, conn)
			}
		}(key)
	}
}

func handleConnection(ctx context.Context, wg *sync.WaitGroup, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			microg.E(err)
		}
	}()
	defer conn.Close()
	microg.I(ctx, "New connection from %s", conn.RemoteAddr().String())
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	buffer := newReceiveBuffer()
	tempBuffer := make([]byte, 1024)
	cvtName := ctx.Value(constants.TracerKeyCvtName).(string)
	onLineKey := fmt.Sprintf("%s_3rd_party_On", cvtName)
	offLineKey := fmt.Sprintf("%s_3rd_party_Off", cvtName)
	metrics.ConnectCounter.WithLabelValues(onLineKey).Inc()
	defer func() {
		metrics.ConnectCounter.WithLabelValues(offLineKey).Inc()
	}()
	thirdPartyDataQueue := exchange.ThirdPartyDataQueuePool[cvtName]
	for {
		err := conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		if err != nil {
			microg.E(ctx, "Failed to set read deadline: %v", err)
			return
		}
		select {
		case <-ctx.Done():
			microg.W(ctx, "exit third party receiver %s", conn.RemoteAddr().String())
			return
		default:
			n, err := conn.Read(tempBuffer)
			if n == 0 {
				if err != nil {
					if !errors.Is(err, os.ErrDeadlineExceeded) {
						microg.W(ctx, "Connection error: %s", err.Error())
						return
					}
				}
				time.Sleep(time.Millisecond * 10)
				continue
			}
			tempCache := tempBuffer[:n]
			for _, char := range tempCache {
				buffer.add(char)
				if buffer.matched() {
					raw := buffer.flush()
					if len(raw) > 0 {
						if len(thirdPartyDataQueue) < cap(thirdPartyDataQueue) {
							thirdPartyDataQueue <- raw
						} else {
							traceID := gjson.Get(raw, "trace_id").String()
							if traceID != "" {
								zapField := zap.String("trace_id", traceID)
								microg.W(ctx, zapField, "Third party entrance data queue is full")
							}
						}
					}
				}
			}
		}
	}
}
