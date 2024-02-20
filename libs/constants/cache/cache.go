package cache

import "strings"

var CACHE = struct {
	TERMINAL_INFO   string
	LOCATION_INFO   string
	CAR_INFO        string
	CORP_INFO       string
	DEFAULT         string
	bb809_converter string
}{
	TERMINAL_INFO:   "terminal",
	LOCATION_INFO:   "location",
	CAR_INFO:        "car",
	CORP_INFO:       "corp",
	DEFAULT:         "default",
	bb809_converter: "cache_change_chanel",
}

func GetType(key string) string {
	if key == "" {
		return ""
	}
	if strings.Contains(key, "terminal") {
		return CACHE.TERMINAL_INFO
	}
	if strings.Contains(key, "car") {
		return CACHE.CAR_INFO
	}

	return CACHE.DEFAULT
}

var CacheOperate = struct {
	DEL string
	NTS string
}{
	DEL: "delete",
	NTS: "notice",
}

var TIME = struct {
	ONEMINUTE int
	AQUARTER  int
	ONEDAY    int
}{
	ONEMINUTE: 60,
	AQUARTER:  900,
	ONEDAY:    86400,
}
