package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/caarlos0/env"
	"github.com/tukhvatullinsm/golang-project/internal/handlers/agent"
	"github.com/tukhvatullinsm/golang-project/internal/metrics"
)

type Config struct {
	Endpoint       string `env:"ADDRESS"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
}

var MetricsName = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Endpoint == "" {
		flag.StringVar(&cfg.Endpoint, "a", "localhost:8080", "Enter endpoint socket (address:port)")
	}
	if cfg.ReportInterval == 0 {
		flag.Int64Var(&cfg.ReportInterval, "r", 10, "Report interval for metrics in seconds")
	}
	if cfg.PollInterval == 0 {
		flag.Int64Var(&cfg.PollInterval, "p", 2, "Poll interval for metrics in seconds")
	}
	flag.Parse()

	runMemStat := runtime.MemStats{}
	metricsObj := metrics.MyMetrics{}
	metricsObj.Init(&runMemStat, MetricsName)

	agentApp := agent.AgentApp{}
	agentApp.Init("http://", cfg.Endpoint, &metricsObj)

	for {
		agentApp.UpdateValue()
		time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
		if int64(metricsObj.PollCount)%5 == 0 {
			agentApp.SendMetric()
		}
	}

}
