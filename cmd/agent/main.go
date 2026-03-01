package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/caarlos0/env"
	"github.com/tukhvatullinsm/golang-project/internal/handlers/agent"
	"github.com/tukhvatullinsm/golang-project/internal/metrics"
)

type Config struct {
	Endpoint       string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

func init() {

}

const (
	//pollInterval   = 2 * time.Second
	//reportInterval = 10 * time.Second
	scheme = "http://"
)

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
	// TODO: Init Client Configuration
	os.Setenv("ADDRESS", "localhost:443")
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Endpoint == "" {
		flag.StringVar(&cfg.Endpoint, "a", "localhost:8080", "Enter endpoint socket (address:port)")
	}
	if cfg.ReportInterval == 0 {
		flag.DurationVar(&cfg.ReportInterval, "r", 10, "Report interval for metrics in seconds")
	}
	if cfg.PollInterval == 0 {
		flag.DurationVar(&cfg.PollInterval, "p", 2, "Poll interval for metrics in seconds")
	}
	flag.Parse()
	// TODO : init runtime memstat object
	runMemStat := runtime.MemStats{}
	// TODO: init agent metrics object
	metricsObj := metrics.MyMetrics{}
	metricsObj.Init(&runMemStat, MetricsName)
	// TODO: init agent handler object
	agentApp := agent.AgentApp{}
	agentApp.Init(scheme, cfg.Endpoint, &metricsObj)
	// TODO: run main algorithm

	for {
		agentApp.UpdateValue()
		time.Sleep(cfg.PollInterval * time.Second)
		if int64(metricsObj.PollCount)%5 == 0 {
			agentApp.SendMetric()
		}
	}

}
