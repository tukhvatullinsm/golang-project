package storage

import (
	"strconv"
)

type gauge int64

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
	num, _ := strconv.ParseInt(value, 10, 64)
	switch param {
	case "gauge":
		ms.gauge[key] = gauge(num)
	case "counter":
		ms.counter[key] = append(ms.counter[key], counter(num))
	}
}
