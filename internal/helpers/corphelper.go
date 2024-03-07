package helpers

import (
	"context"
	"encoding/json"
	"github.com/peifengll/go_809_converter/internal/model"
	"github.com/peifengll/go_809_converter/internal/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CorpHelper struct {
	corpRepo repository.CorpRepo
	db       *gorm.DB
}

func NewCorpHelper(corpRepo repository.CorpRepo, db *gorm.DB) *CorpHelper {
	return &CorpHelper{
		corpRepo: corpRepo,
		db:       db,
	}
}

func (h *CorpHelper) GetOrCreateCorpCache(cid string, rebuild bool, redisClient *redis.Client) {
	panic("not implemented")
	ctx := context.Background()
	name := RedisKeyHelper.GetCorpCacheHashKey(cid)
	data, err := redisClient.HGetAll(ctx, name).Result()
	var t_op model.TOperator
	if err != nil || rebuild || len(data) == 0 {
		err := h.db.Raw("SELECT `mobile`, `type`, `status`, `auth_type` FROM t_operator WHERE `cid`=? AND status=1", cid).First(&t_op).Error
		if err != nil {
			return
		}
		corpInfo := h.GetCorpInfoByCid(cid)
		if err != nil {
			return
		}
		cacheData := map[string]any{
			"users":     "",
			"corp_info": corpInfo,
		}
		cacheBytes, err := json.Marshal(cacheData)
		if err != nil {
		}
		pipe := redisClient.TxPipeline()
		pipe.HMSet(ctx, name, map[string]interface{}{"data": string(cacheBytes)})
		pipe.Expire(ctx, name, 3600*24*7)
		_, err = pipe.Exec(ctx)
		if err != nil {
			return
		}
		return
	}

}

func (h *CorpHelper) GetCorpInfoByCid(cid string) *model.TCorp {
	//corp_key := RedisKeyHelper.GetCorpSimpleKey(cid)
	var crop *model.TCorp
	crop = h.corpRepo.GetCorpByCid(cid)
	return crop
}
