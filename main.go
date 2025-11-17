package main

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-17 20:18:14
 * @LastEditors: yangtongbing 1280758415@qq.com
 * @LastEditTime: 2025-02-18 09:57:41
 * @FilePath: main.go
 * @Description:
 *
 */

import (
	"context"
	"errors"

	"github.com/dbjtech/go_809_converter/converter"
	"github.com/dbjtech/go_809_converter/libs/database/mysqldb"
	"github.com/dbjtech/go_809_converter/libs/packet_util"

	//"github.com/dbjtech/go_809_converter/libs/database/mysqldb" // 临时注释掉
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gookit/config/v2"

	"github.com/dbjtech/go_809_converter/exchange"
	"github.com/dbjtech/go_809_converter/libs"
	"github.com/dbjtech/go_809_converter/metrics"
	"github.com/dbjtech/go_809_converter/receivers"
	"github.com/dbjtech/go_809_converter/senders"
	"github.com/dbjtech/go_809_converter/traces"
	"github.com/linketech/microg/v4"
	"github.com/linketech/microg/v4/aly/slsprovider"
	flag "github.com/spf13/pflag"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/resource"

	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
)

var normalTcp bool
var jtwTcp bool

func parseCommand() {
	flag.StringVarP(&libs.ConfigType, "type", "t", "toml",
		"config file type, choose one from [json, yaml, toml(default)]. will ignore at setting config-file")
	flag.StringVarP(&libs.Environment, "env", "e", "develop", "program environment")
	flag.StringVarP(&libs.ConfigFile, "config-file", "c", "", "config file path")
	flag.IntVarP(&exchange.ConverterWorker, "worker", "w", 3, "converter worker numbers")
	flag.BoolVarP(&normalTcp, "normal-tcp", "n", false, "normal tcp connect mode")
	flag.BoolVarP(&jtwTcp, "jtw-tcp", "j", false, "Jiao Tong Wei tcp connect mode")
	flag.Parse()
}

func setup() *slsprovider.Config {
	libs.NewConfig()
	_ = config.Set("normalTcp", normalTcp, true)
	_ = config.Set("jtwTcp", jtwTcp, true)
	metrics.Init809PrometheusClient()
	otelConf := initTracer("809_converter")
	//redisdb.InitDefaultRedis()
	mysqldb.InitDefaultGormDB()
	return otelConf
}

func main() {
	parseCommand()
	otelConf := setup()
	if otelConf != nil {
		defer slsprovider.Shutdown(otelConf)
	}
	engine := gin.New()
	if traces.TraceSwitch().WebTraceable {
		engine.Use(otelgin.Middleware("809_converter"))
	}
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "x-requested-with")
	corsConfig.AllowAllOrigins = true
	engine.Use(cors.New(corsConfig))
	ctx, cancel := context.WithCancel(context.Background())
	initDataQueue()
	var wg sync.WaitGroup
	// 启动第三方数据消费服务
	for i := 0; i < exchange.ConverterWorker; i++ {
		senders.TransformThirdPartyData(ctx)
	}
	time.Sleep(time.Second)
	// 启动第三方数据接收服务
	receivers.StartThirdPartyReceiver(ctx, &wg)
	if normalTcp {
		// 上级平台的下行链路，上级平台连接本服务
		go receivers.StartDownlink(ctx, &wg)
		time.Sleep(time.Second)
		// 上级平台的上行链路，本服务连接上级平台
		go senders.StartUpLink(ctx, &wg)
	}

	// if jtwTcp { // 交委的暂时不改造
	// 	// 交委转换服务下行
	// 	go receivers.StartJtwConverterDownLink(ctx, &wg)
	// 	// 交委转换服务上行
	// 	go senders.StartJtwConverterUpLink(ctx, &wg)
	// }

	converter.SetRoute(engine)
	addr := ":" + config.String(libs.Environment+".converter.consolePort", "13031")
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			microg.I("listen err: %s\n", err)
		}
	}()
	microg.I("web server listen at %s", addr)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	microg.I("Shutdown Server ...")

	//if err := srv.Shutdown(ctx); err != nil {
	//	microg.E("Server Shutdown: %v", err)
	//	return
	//}
	cancel()
	wg.Wait()
	microg.I("Server exited")
	if err := srv.Shutdown(ctx); err != nil {
		microg.E("Web Server Shutdown: %v", err)
	}
	microg.W("Web Server exiting")
}

func initTracer(serviceName string) *slsprovider.Config {
	if !traces.TraceSwitch().WebTraceable {
		return nil
	}
	cfg, err := slsprovider.NewConfig(
		slsprovider.WithServiceName(serviceName),
		slsprovider.WithResource(resource.Empty()),
	)
	if err != nil {
		panic(err)
	}
	if err := slsprovider.Start(cfg); err != nil {
		panic(err)
	}
	return cfg
}

func initDataQueue() {
	// 获取 libs.Environment + ".converter" 下的所有子元素，并判断是否开启，和是否开启扩展协议
	configConverter := config.SubDataMap(libs.Environment + ".converter")
	for key, _ := range configConverter {
		exchange.ThirdPartyDataQueuePool[key] = make(chan string, 1000)
		exchange.UpLinkDataQueuePool[key] = make(chan packet_util.MessageWrapper, 1000)
		exchange.JtwConverterUpLinkDataQueuePool[key] = make(chan packet_util.MessageWrapper, 1000)
	}
}
