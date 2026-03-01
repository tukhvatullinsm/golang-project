package agent

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

type AgentProvider interface {
	ExportMetrics() *map[string]map[string]interface{}
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
	exportData := *aa.getData.ExportMetrics()
	client := resty.New()
	client.SetHeader("Content-Type", "text/plain")
	client.SetBaseURL(aa.remoteProtocol + aa.remoteWebServer + "/update")
	/*
		req, err := http.NewRequest("POST", aa.remoteProtocol+aa.remoteHost+aa.remotePort, nil)
		req.Header.Add("Content-Type", "text/plain")
		if err != nil {
			log.Fatal(err)
		}
	*/
	//client := &http.Client{}
	for k, v := range exportData {
		path := ""
		for m, a := range v {
			path = "/" + k + "/" + m + "/" + fmt.Sprintf("%v", a)
			resp, err := client.R().
				Post(path)
			if err != nil {
				log.Fatal(err, resp.StatusCode())
			}
			if resp.StatusCode() != 200 {
				log.Fatal(resp.RawResponse)
			}
			/*req.URL.Path = path
			req.URL.EscapedPath()
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err, resp.StatusCode)
			}
			resp.Body.Close() */
		}
	}
}
