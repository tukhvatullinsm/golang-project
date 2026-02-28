package main

import (
	"runtime"
	"time"

	"github.com/tukhvatullinsm/golang-project/internal/handlers/agent"
	"github.com/tukhvatullinsm/golang-project/internal/metrics"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	host           = "localhost"
	scheme         = "http://"
	port           = "8080"
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
	agentApp.Init(scheme, host, port, &metricsObj)
	// TODO: run main algorithm

	for {
		agentApp.UpdateValue()
		time.Sleep(pollInterval)
		if int64(metricsObj.PollCount)%5 == 0 {
			agentApp.SendMetric()
		}
	}

}
