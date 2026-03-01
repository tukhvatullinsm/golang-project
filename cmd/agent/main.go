package main

import (
	"flag"
	"runtime"
	"time"

	"github.com/tukhvatullinsm/golang-project/internal/handlers/agent"
	"github.com/tukhvatullinsm/golang-project/internal/metrics"
)

var WebServer struct {
	Endpoint string
}

var WebClient struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func init() {
	flag.StringVar(&WebServer.Endpoint, "a", "localhost:8080", "Enter endpoint socket (address:port)")
	flag.DurationVar(&WebClient.ReportInterval, "r", 10, "Report interval for metrics in seconds")
	flag.DurationVar(&WebClient.PollInterval, "p", 2, "Poll interval for metrics in seconds")

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
	// TODO : init runtime memstat object
	runMemStat := runtime.MemStats{}
	// TODO: init agent metrics object
	metricsObj := metrics.MyMetrics{}
	metricsObj.Init(&runMemStat, MetricsName)
	// TODO: init agent handler object
	agentApp := agent.AgentApp{}
	agentApp.Init(scheme, WebServer.Endpoint, &metricsObj)
	// TODO: run main algorithm

	for {
		agentApp.UpdateValue()
		time.Sleep(WebClient.PollInterval * time.Second)
		if int64(metricsObj.PollCount)%5 == 0 {
			agentApp.SendMetric()
		}
	}

}
