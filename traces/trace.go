package traces

import (
	"os"
	"strconv"
	"sync"
)

type traceSwitch struct {
	WebTraceable   bool
	RedisTraceable bool
	MySQLTraceable bool
	KafkaTraceable bool
}

var ts *traceSwitch
var initTraceSwitch sync.Once

func TraceSwitch() *traceSwitch {
	if ts == nil {
		initTraceSwitch.Do(func() {
			nts := newTraceSwitch()
			ts = &nts
		})
	}
	return ts
}

func newTraceSwitch() traceSwitch {
	return traceSwitch{
		WebTraceable:   webTrace(),
		RedisTraceable: redisTrace(),
		MySQLTraceable: mysqlTrace(),
		KafkaTraceable: kafkaTrace(),
	}
}

func webTrace() bool {
	ts := os.Getenv("WEB_TRACE")
	switchOn := false
	if ts != "" {
		switchOn, _ = strconv.ParseBool(ts)
	}
	return switchOn
}

func redisTrace() bool {
	ts := os.Getenv("REDIS_TRACE")
	switchOn := false
	if ts != "" {
		switchOn, _ = strconv.ParseBool(ts)
	}
	return switchOn
}

func mysqlTrace() bool {
	ts := os.Getenv("MYSQL_TRACE")
	switchOn := false
	if ts != "" {
		switchOn, _ = strconv.ParseBool(ts)
	}
	return switchOn
}

func kafkaTrace() bool {
	ts := os.Getenv("KAFKA_TRACE")
	switchOn := false
	if ts != "" {
		switchOn, _ = strconv.ParseBool(ts)
	}
	return switchOn
}
