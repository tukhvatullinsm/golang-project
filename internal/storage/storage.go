package storage

import (
	"strconv"
)

type gauge float64

type counter int64

type MemStorage struct {
	gauge   map[string]gauge
	counter map[string][]counter
}

func New() *MemStorage {
	obj := MemStorage{}
	obj.gauge = make(map[string]gauge)
	obj.counter = make(map[string][]counter)
	return &obj
}

func (ms *MemStorage) Set(param, key, value string) {
	switch param {
	case "gauge":
		num, _ := strconv.ParseFloat(value, 64)
		ms.gauge[key] = gauge(num)
	case "counter":
		num, _ := strconv.ParseInt(value, 10, 64)
		ms.counter[key] = append(ms.counter[key], counter(num))
	}
}

func (ms *MemStorage) Get() []string {
	return nil
}
