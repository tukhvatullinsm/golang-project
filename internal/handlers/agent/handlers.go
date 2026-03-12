package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type AgentProvider interface {
	ExportMetrics() map[string]map[string]interface{}
	UpdateMetrics()
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
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
	backoffSchedule := []time.Duration{
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
	exportData := aa.getData.ExportMetrics()
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Content-Encoding", "gzip")
	for k, v := range exportData {
		expObj := new(Metrics)
		switch k {
		case "counter":
			expObj.MType = "counter"
			expObj.Delta = new(int64)
			for m, a := range v {
				expObj.ID = m
				*expObj.Delta, _ = strconv.ParseInt(fmt.Sprintf("%v", a), 10, 64)
				jsonData, err := json.Marshal(expObj)
				if err != nil {
					log.Printf("Error marshaling JSON: %s", err)
					break
				}
				compressRes, err := Compress(jsonData)
				if err != nil {
					log.Printf("Error compressing JSON: %s", err)
					break
				}
				for _, backoff := range backoffSchedule {
					resp, err := client.R().
						SetBody(compressRes).
						Post(aa.remoteProtocol + aa.remoteWebServer + "/update/")

					if err != nil {
						log.Println(err, resp.StatusCode())
						time.Sleep(backoff)
					} else {
						break
					}
				}
			}
		case "gauge":
			expObj.MType = "gauge"
			expObj.Value = new(float64)
			for m, a := range v {
				expObj.ID = m
				*expObj.Value, _ = strconv.ParseFloat(fmt.Sprintf("%v", a), 64)
				jsonData, err := json.Marshal(expObj)
				if err != nil {
					log.Printf("Error marshaling JSON: %s", err)
					break
				}
				compressRes, err := Compress(jsonData)
				if err != nil {
					log.Printf("Error compressing JSON: %s", err)
					break
				}
				for _, backoff := range backoffSchedule {
					resp, err := client.R().
						SetBody(compressRes).
						Post(aa.remoteProtocol + aa.remoteWebServer + "/update/")

					if err != nil {
						log.Println(err, resp.StatusCode())
						time.Sleep(backoff)
					} else {
						break
					}
				}

			}
		}
	}
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		log.Println("Error creating gzip writer", err)
		return nil, err
	}
	if _, err := gz.Write(data); err != nil {
		log.Println("Error compressing data", err)
	}
	if err := gz.Close(); err != nil {
		log.Println("Error compressing (clear) data", err)
	}
	return b.Bytes(), nil
}
