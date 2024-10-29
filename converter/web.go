package converter

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-27 16:57:55
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-29 20:19:25
 * @FilePath: \go_809_converter\converter\web.go
 * @Description:
 *
 */

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dbjtech/go_809_converter/exchange"
	"github.com/dbjtech/go_809_converter/libs/cache"
	"github.com/dbjtech/go_809_converter/libs/daos"
	"github.com/dbjtech/go_809_converter/metrics"

	"github.com/dbjtech/go_809_converter/libs/database/mysqldb"
	"github.com/gin-gonic/gin"
	"github.com/gookit/config/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetRoute(r *gin.Engine) {
	baseSet(r, "dbj-809-converter service")
	r.StaticFS("/static", http.Dir("converter/static")) //外挂静态文件目录，方便升级前端页面
	r.GET("/metrics", getMetrics(promhttp.Handler()))
	r.PUT("/cache/manager", removeCache)
	r.POST("/cache/manager", showCache)
	r.POST("/push/time", getPushTime)
}
func baseSet(r *gin.Engine, name string) {
	r.GET("/", welcome(name))
	r.GET("/ping", ping)
	r.GET("/config", getConfig)
}

func welcome(s string) gin.HandlerFunc {
	wl := fmt.Sprintf("This is %v Backend Server", s)
	return func(c *gin.Context) {
		c.String(http.StatusOK, wl)
	}
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func getConfig(c *gin.Context) {
	// ${env}.main.eventer.windowSize
	// mysql.pool_size
	// redis.pool_size
	env := config.String("env")
	windowSizeKey := env + ".main.eventer.windowSize"
	mysqlPoolSizeKey := "mysql_db.pool_size"
	windowSize := c.Query(windowSizeKey)
	mysqlPoolSize := c.Query(mysqlPoolSizeKey)
	if windowSize != "" {
		ws, _ := strconv.Atoi(windowSize)
		if ws > 0 {
			_ = config.Set(windowSizeKey, ws, true)
		}
	}
	if mysqlPoolSize != "" {
		db, _ := mysqldb.GormDB.DB()
		psz, _ := strconv.Atoi(mysqlPoolSize)
		if psz > 0 {
			_ = config.Set(mysqlPoolSizeKey, psz, false)
			db.SetMaxOpenConns(psz)
		}
	}

	configData := config.Data()
	c.PureJSON(http.StatusOK, configData)
}

func getMetrics(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		extraSetting()
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func extraSetting() {
	metrics.PacketPoolSize.WithLabelValues("3rd_part_input").Set(float64(len(exchange.ThirdPartyDataQueue)))
	metrics.PacketPoolSize.WithLabelValues("uplink_4_send").Set(float64(len(exchange.UpLinkDataQueue)))
	metrics.CacheSize.WithLabelValues("all_count").Set(float64(cache.Manager.Count()))
}

func removeCache(c *gin.Context) {
	type cacheQueryBody struct {
		CacheFrom string   `json:"cacheFrom"`
		CacheList []string `json:"cacheList"`
	}
	var cacheQuery cacheQueryBody
	err := c.BindJSON(&cacheQuery)
	ctx := c.Request.Context()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	if len(cacheQuery.CacheFrom) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "cache list is empty",
		})
		return
	}
	result := make([]string, len(cacheQuery.CacheList))
	switch cacheQuery.CacheFrom {
	case "vin":
		carIDs := daos.NewSimpleDao().GetCarIdByVin(ctx, cacheQuery.CacheList)
		for i, v := range cacheQuery.CacheList {
			carID, ok := carIDs[v]
			if !ok {
				result[i] = "车辆不存在"
				continue
			}
			cache.Manager.Remove(carID[:16])
			result[i] = "清除成功"
		}
	case "cnum":
		carIDs := daos.NewSimpleDao().GetCarIdByCnum(ctx, cacheQuery.CacheList)
		for i, v := range cacheQuery.CacheList {
			carID, ok := carIDs[v]
			if !ok {
				result[i] = "车辆不存在"
				continue
			}
			cache.Manager.Remove(carID[:16])
			result[i] = "清除成功"
		}
	case "sn":
		tids := daos.NewSimpleDao().GetTidBySn(ctx, cacheQuery.CacheList)
		for i, v := range cacheQuery.CacheList {
			tid, ok := tids[v]
			if !ok {
				result[i] = "设备不存在"
				continue
			}
			cache.Manager.Remove(tid[:16])
			result[i] = "清除成功"
		}
	case "id":
		for i, v := range cacheQuery.CacheList {
			if len(v) > 16 {
				v = v[:16]
			}
			cache.Manager.Remove(v)
			result[i] = "清除成功"
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "操作成功",
	})
}

func showCache(c *gin.Context) {
	type cacheQueryBody struct {
		CacheFrom string   `json:"cacheFrom"`
		CacheList []string `json:"cacheList"`
	}
	var cacheQuery cacheQueryBody
	err := c.BindJSON(&cacheQuery)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	if cacheQuery.CacheFrom == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "未知的缓存来源",
		})
		return
	}
	result := make([]any, len(cacheQuery.CacheList))
	ctx := c.Request.Context()
	switch cacheQuery.CacheFrom {
	case "vin":
		carIDs := daos.NewSimpleDao().GetCarIdByVin(ctx, cacheQuery.CacheList)
		for i, v := range cacheQuery.CacheList {
			carID, ok := carIDs[v]
			if !ok {
				result[i] = "车辆不存在"
				continue
			}
			data := cache.Manager.Get(carID[:16])
			result[i] = data
		}
	case "cnum":
		carIDs := daos.NewSimpleDao().GetCarIdByCnum(ctx, cacheQuery.CacheList)
		for i, v := range cacheQuery.CacheList {
			carID, ok := carIDs[v]
			if !ok {
				result[i] = "车辆不存在"
				continue
			}
			data := cache.Manager.Get(carID[:16])
			result[i] = data
		}
	case "sn":
		tids := daos.NewSimpleDao().GetTidBySn(ctx, cacheQuery.CacheList)
		for i, v := range cacheQuery.CacheList {
			tid, ok := tids[v]
			if !ok {
				result[i] = "设备不存在"
				continue
			}
			data := cache.Manager.Get(tid[:16])
			result[i] = data
		}
	case "id":
		for i, v := range cacheQuery.CacheList {
			if len(v) > 16 {
				v = v[:16]
			}
			data := cache.Manager.Get(v)
			result[i] = data
		}
	case "all":
		info := cache.Manager.ComputeAvailable()
		result = append(result, fmt.Sprintf("删除缓存:%d个，现存缓存:%d个", info["removed"], info["cached"]))
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "操作成功",
	})
}

func getPushTime(c *gin.Context) {
	type pushTimeBody struct {
		Items []string `json:"items"`
	}
	var body pushTimeBody
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	result := make([]int64, len(body.Items))
	for i, v := range body.Items {
		result[i], _ = exchange.TaskMarker.Get(v)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "操作成功",
	})
}
