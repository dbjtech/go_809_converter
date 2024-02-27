package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const PROJECT_ID = "QJCG"

var RedisKeyHelper = &redisKeyHelper{}

type redisKeyHelper struct{}

func (r *redisKeyHelper) GetPushStatusID(oid, tag string) string {
	return fmt.Sprintf("qjcg:%s:%s", oid, tag)
}

func (r *redisKeyHelper) GetPushKey(oid, t string) string {
	secret := "7c2d6047c7ad95f79cdb985e26a92141"
	s := oid + t + secret
	hasher := md5.New()
	hasher.Write([]byte(s))
	key := hex.EncodeToString(hasher.Sum(nil))
	return key
}

func (r *redisKeyHelper) GetPushRegistID(oid string) string {
	return fmt.Sprintf("%s_wspush_jqcg_registered:%s", PROJECT_ID, oid)
}

func (r *redisKeyHelper) GetTerminalCacheKey(tid string) string {
	return fmt.Sprintf("%s_terminal_info_cache:%s", PROJECT_ID, tid)
}

func (r *redisKeyHelper) GetCarCacheKey(carID string) string {
	return fmt.Sprintf("%s_car_info_cache:%s", PROJECT_ID, carID)
}

func (r *redisKeyHelper) GetTerminalLastPvtKey(tid string) string {
	return fmt.Sprintf("%s_terminal_last_pvt_key:%s", PROJECT_ID, tid)
}

func (r *redisKeyHelper) GetCarLastPvtKey(tid string) string {
	return fmt.Sprintf("%s_car_last_pvt_key:%s", PROJECT_ID, tid)
}

func (r *redisKeyHelper) GetCorpCacheHashKey(cid string) string {
	return fmt.Sprintf("%s_corp_info:%s", PROJECT_ID, cid)
}

func (r *redisKeyHelper) GetCorpSimpleKey(cid string) string {
	return fmt.Sprintf("%s_corp:%s", PROJECT_ID, cid)
}

func (r *redisKeyHelper) GetPushableMobileKey(mobile string) string {
	return fmt.Sprintf("%s_pushable_mobile:%s", PROJECT_ID, mobile)
}

func (r *redisKeyHelper) GetThirdPartyInterfaceYH() string {
	return fmt.Sprintf("%s_third_party_interface_yh", PROJECT_ID)
}

func (r *redisKeyHelper) GetGWSerialKey() string {
	return fmt.Sprintf("%s_qurey_gw:serial", PROJECT_ID)
}

func (r *redisKeyHelper) GetWarningSettingsKey(cid string) string {
	return fmt.Sprintf("%s:warn_option:%s", PROJECT_ID, cid)
}

func (r *redisKeyHelper) GetEventLongStopPushedKey(tid string) string {
	return fmt.Sprintf("event:long_stop_pushed:%s", tid)
}
