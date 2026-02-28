package agent

import (
	"fmt"
	"log"
	"net/http"
)

type AgentProvider interface {
	ExportMetrics() *map[string]map[string]interface{}
	UpdateMetrics()
}

type AgentApp struct {
	remoteHost     string
	remotePort     string
	remoteProtocol string
	getData        AgentProvider
}

func (aa *AgentApp) Init(protocol string, host string, port string, data AgentProvider) {
	aa.remotePort = port
	aa.remoteProtocol = protocol
	aa.remoteHost = host
	aa.getData = data
}

func (aa *AgentApp) UpdateValue() {
	aa.getData.UpdateMetrics()
}

func (aa *AgentApp) SendMetric() {
	exportData := *aa.getData.ExportMetrics()
	req, err := http.NewRequest("POST", aa.remoteProtocol+aa.remoteHost+aa.remotePort, nil)
	req.Header.Add("Content-Type", "text/plain")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	for k, v := range exportData {
		path := ""
		for m, a := range v {
			path = "/update/" + k + "/" + m + "/" + fmt.Sprintf("%v", a)
			req.URL.Path = path
			req.URL.EscapedPath()
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err, resp.StatusCode)
			}
			resp.Body.Close()
		}
	}
}
