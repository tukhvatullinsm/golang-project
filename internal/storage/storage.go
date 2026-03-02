package storage

import (
	"strconv"
	"time"
)

type gauge float64

type counter int64

type MemStorage struct {
	gauge      map[string]gauge
	counter    map[string]counter
	dataCache  map[string]any
	lastUpdate time.Time
	lastGet    time.Time
}

func New() *MemStorage {
	obj := MemStorage{}
	obj.gauge = make(map[string]gauge)
	obj.counter = make(map[string]counter)
	obj.dataCache = make(map[string]any)
	return &obj
}

func (ms *MemStorage) SetValue(param, key, value string) {
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

func (ms *MemStorage) GetAllValue() *map[string]any {
	if ms.lastUpdate.After(ms.lastGet) {
		for k, v := range ms.gauge {
			ms.dataCache[k] = v
		}
		for k, v := range ms.counter {
			ms.dataCache[k] = v
		}
	}
	return &ms.dataCache
}
