package agent

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

type AgentProvider interface {
	ExportMetrics() map[string]map[string]interface{}
	UpdateMetrics()
}

type AgentApp struct {
	remoteWebServer string
	remoteProtocol  string
	getData         AgentProvider
}

func (aa *AgentApp) Init(protocol string, host string, data AgentProvider) {
	aa.remoteProtocol = protocol
	aa.remoteWebServer = host
	aa.getData = data
}

func (aa *AgentApp) UpdateValue() {
	aa.getData.UpdateMetrics()
}

func (aa *AgentApp) SendMetric() {
	exportData := aa.getData.ExportMetrics()
	client := resty.New()
	client.SetHeader("Content-Type", "text/plain")
	client.SetBaseURL(aa.remoteProtocol + aa.remoteWebServer + "/update")

	for k, v := range exportData {
		path := ""
		for m, a := range v {
			path = "/" + k + "/" + m + "/" + fmt.Sprintf("%v", a)
			resp, err := client.R().
				Post(path)
			if err != nil {
				log.Println(err, resp.StatusCode())
			}
			if resp.StatusCode() != 200 {
				log.Println(resp.RawResponse)
			}
		}
	}
}
