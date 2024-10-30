package main

/*
 * @Author: SimingLiu siming.liu@linketech.cn
 * @Date: 2024-10-17 20:18:14
 * @LastEditors: SimingLiu siming.liu@linketech.cn
 * @LastEditTime: 2024-10-28 19:31:32
 * @FilePath: \go_809_converter\main.go
 * @Description:
 *
 */

import (
	"context"
	"errors"
	"github.com/gin-contrib/cors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dbjtech/go_809_converter/converter"
	"github.com/dbjtech/go_809_converter/libs/database/mysqldb"
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

func parseCommand() {
	flag.StringVarP(&libs.ConfigType, "type", "t", "toml",
		"config file type, choose one from [json, yaml, toml(default)]. will ignore at setting config-file")
	flag.StringVarP(&libs.Environment, "env", "e", "develop", "program environment")
	flag.StringVarP(&libs.ConfigFile, "config-file", "c", "", "config file path")
	flag.IntVarP(&exchange.ConverterWorker, "worker", "w", 3, "converter worker numbers")
	flag.Parse()
}

func setup() *slsprovider.Config {
	libs.NewConfig()
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
	var wg sync.WaitGroup
	go receivers.StartDownlink(ctx, &wg)
	go receivers.StartThirdPartyReceiver(ctx, &wg)
	time.Sleep(time.Second)
	go senders.StartUpLink(ctx, &wg)
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
