package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type gauge float64

type counter int64

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *counter `json:"delta,omitempty"`
	Value *gauge   `json:"value,omitempty"`
}
type MemStorage struct {
	gauge      map[string]gauge
	counter    map[string]counter
	dataCache  map[string]any
	lastUpdate time.Time
	lastGet    time.Time
	file       *os.File
	objlog     *zap.SugaredLogger
	buff       *bufio.Writer
	Parameters []string
	ParamsDict map[string]string
}

func New(objlog *zap.SugaredLogger) *MemStorage {
	obj := MemStorage{}
	obj.gauge = make(map[string]gauge)
	obj.counter = make(map[string]counter)
	obj.dataCache = make(map[string]any)
	obj.file = nil
	obj.objlog = objlog
	obj.buff = nil
	obj.Parameters = make([]string, 0)
	obj.ParamsDict = make(map[string]string)
	return &obj
}

func (ms *MemStorage) SetValue(param string, key string, value string) {
	switch param {
	case "gauge":
		num, _ := strconv.ParseFloat(value, 64)
		ms.gauge[key] = gauge(num)
	case "counter":
		num, _ := strconv.ParseInt(value, 10, 64)
		ms.counter[key] += counter(num)
	}
	ms.lastUpdate = time.Now()
}

func (ms *MemStorage) GetValue(param, key string) any {
	switch param {
	case "gauge":
		v, ok := ms.gauge[key]
		if ok {
			return v
		} else {
			return nil
		}
	case "counter":
		v, ok := ms.counter[key]
		if ok {
			return v
		} else {
			return nil
		}
	default:
		return nil
	}
}

func (ms *MemStorage) GetAllValue() map[string]any {
	if ms.lastUpdate.After(ms.lastGet) {
		for k, v := range ms.gauge {
			ms.dataCache[k] = v
		}
		for k, v := range ms.counter {
			ms.dataCache[k] = v
		}
	}
	return ms.dataCache
}

func (ms *MemStorage) Create(filePath string, restore bool, interval int64) error {
	flag := os.O_RDWR | os.O_CREATE | os.O_EXCL
	if interval == 0 {
		flag += os.O_SYNC
	}
	var err error
	ms.file, err = os.OpenFile(filePath, flag, 0666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			ms.objlog.Infow("File already exists, continue to work with that", "filename", filePath)
			flag -= os.O_CREATE | os.O_EXCL
			ms.file, _ = os.OpenFile(filePath, flag, 0666)
		} else {
			return err
		}
	}
	if restore {
		buff := bufio.NewScanner(ms.file)
		for buff.Scan() {
			data := buff.Bytes()
			metricsObj := Metrics{}
			err := json.Unmarshal(data, &metricsObj)
			if err != nil {
				ms.objlog.Infow("Error encode some data from file", "error", err)
			}
			switch metricsObj.MType {
			case "gauge":
				ms.gauge[metricsObj.ID] = *metricsObj.Value
			case "counter":
				ms.counter[metricsObj.ID] = *metricsObj.Delta
			}
			ms.lastUpdate = time.Now()
		}
		if err = buff.Err(); err != nil {
			ms.objlog.Infow("Error restore some data file or end of file", "error", err)
		}

	}
	return nil
}

func (ms *MemStorage) Write() error {
	//var buff []byte
	count := int64(0)
	SaveObjSlice := []Metrics{}
	CurrObjSlice := []Metrics{}
	buff := bufio.NewScanner(ms.file)
	for buff.Scan() {
		data := buff.Bytes()
		metricsObj := Metrics{}
		err := json.Unmarshal(data, &metricsObj)
		if err != nil {
			ms.objlog.Infow("Error encode some data from file", "error", err)
		}
		SaveObjSlice = append(SaveObjSlice, metricsObj)
		ms.ParamsDict[metricsObj.ID] = metricsObj.MType
		count++
	}
	if err := buff.Err(); err != nil {
		ms.objlog.Infow("Error restore some data file or end of file", "error", err)
	}

	for k, v := range ms.gauge {
		CurrObjSlice = append(CurrObjSlice, Metrics{MType: "gauge", ID: k, Value: &v})
	}
	for k, v := range ms.counter {
		CurrObjSlice = append(CurrObjSlice, Metrics{MType: "counter", ID: k, Delta: &v})
	}
	exportData := make([]byte, 0)

	for _, cv := range CurrObjSlice {
		v, ok := ms.ParamsDict[cv.ID]
		if ok && cv.MType == v {
			for _, obj := range SaveObjSlice {
				if obj.MType == cv.MType && obj.ID == cv.ID {
					switch cv.MType {
					case "counter":
						ms.counter[cv.ID] = *cv.Delta
					case "gauge":
						ms.gauge[cv.ID] = *cv.Value
					}
				}
			}
		} else {
			SaveObjSlice = append(SaveObjSlice, cv)
		}
	}

	for _, m := range SaveObjSlice {
		tmp, err := json.Marshal(m)
		if err != nil {
			ms.objlog.Infow("Failed to marshal metrics to JSON", "error", err)
		}
		exportData = append(exportData, tmp...)
		exportData = append(exportData, '\n')
	}

	if len(exportData) != 0 {
		err := ms.file.Truncate(0)
		if err != nil {
			ms.objlog.Infow("Failed to truncate file", "error", err)
		}
		_, err = ms.file.Seek(0, 0)
		if err != nil {
			ms.objlog.Infow("Failed to seek file", "error", err)
		}
		_, err = ms.file.Write(exportData)
		ms.file.Sync()

		if err != nil {
			ms.objlog.Infow("Failed to write metrics to file", "error", err)
			return err
		}
	} else {
		ms.objlog.Infow("Empty data in program")
	}

	return nil
}

func (ms *MemStorage) Close() error {
	return ms.file.Close()
}
