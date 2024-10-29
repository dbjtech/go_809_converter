package metrics

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	VersionInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: config.String("serverName", os.Getenv("Name")),
			Name:      "process_info",
			Help:      "程序描述信息",
		},
		[]string{"version", "git_hash", "branch_name", "build_at"},
	)
	LinkHeartBeat = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: config.String("serverName", os.Getenv("Name")),
		Name:      "link_heartbeat",
		Help:      "链路心跳次数",
	}, []string{"name"})
	ConnectCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: config.String("serverName", os.Getenv("Name")),
		Name:      "connect_counter",
		Help:      "链路连断次数",
	}, []string{"name"})
	PacketsDrop = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: config.String("serverName", os.Getenv("Name")),
		Name:      "packet_dropped",
		Help:      "因报错而丢弃的报文",
	}, []string{"topic", "reason"})

	ElapsedTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: config.String("serverName", os.Getenv("Name")),
		Name:      "pkt_elapsed",
		Help:      "报文各阶段处理耗时(单位ms)",
		Buckets:   []float64{10, 20, 50, 100, 150, 300, 500, 1500},
	}, []string{"topic", "type", "action"})

	PacketPoolSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: config.String("serverName", os.Getenv("Name")),
		Name:      "packet_buffer_high",
		Help:      "程序内等待处理的报文个数",
	}, []string{"name"})

	MySQLQuery = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: config.String("serverName", os.Getenv("Name")),
		Name:      "mysql_query",
		Help:      "MySQL查询次数",
	}, []string{"category"})

	CacheSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: config.String("serverName", os.Getenv("Name")),
		Name:      "local_cache",
		Help:      "本地缓存情况，给mysql数据库解压的",
	}, []string{"name"})
)

var Shows = []prometheus.Collector{VersionInfo, ElapsedTime, LinkHeartBeat, ConnectCounter, PacketsDrop,
	PacketPoolSize, MySQLQuery, CacheSize}

func Init809PrometheusClient() {
	prometheus.MustRegister(Shows...)
	SetupVersion("converter")
}

func SetupVersion(serverName string) {
	version := ""
	versionFile, err := os.Open("./version.json")
	if err == nil {
		content, err := io.ReadAll(versionFile)
		if err == nil {
			var payload map[string]string
			err = json.Unmarshal(content, &payload)
			if err == nil {
				version = payload[serverName]
			}
		}
	}
	gitHash, branchName, buildAt := "", "", ""
	gitInfo, err := os.ReadFile("./git_info.txt")
	if err == nil {
		if len(gitInfo) > 0 {
			gis := strings.Split(string(gitInfo), "\n")
			for _, gi := range gis {
				if len(gi) > 0 {
					kv := strings.Split(gi, "=")
					if len(kv) > 1 {
						key := kv[0]
						value := kv[1]
						switch key {
						case "hash":
							gitHash = value
						case "branch":
							branchName = value
						case "build_at":
							buildAt = value
						}
					}
				}
			}
		}
	}
	VersionInfo.WithLabelValues([]string{version, gitHash, branchName, buildAt}...)
}
