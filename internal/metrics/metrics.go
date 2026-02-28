package metrics

import (
	"math/rand/v2"
	"reflect"
	"runtime"
)

type gauge float64
type counter float64

type MyMetrics struct {
	PollCount   counter
	RandomValue gauge
	RunMetrics  map[string]gauge
	SysRuntime  *runtime.MemStats
}

func (mm *MyMetrics) Init(rt *runtime.MemStats, metrics []string) {
	mm.PollCount = 0
	mm.RandomValue = 0.0
	mm.SysRuntime = rt
	mm.RunMetrics = make(map[string]gauge, len(metrics))
	for _, metric := range metrics {
		mm.RunMetrics[metric] = gauge(0.0)
	}
}

func (mm *MyMetrics) UpdateMetrics() {
	mm.PollCount++
	mm.RandomValue = gauge(rand.Float64())
	runtime.ReadMemStats(mm.SysRuntime)
	t := reflect.TypeOf(mm.SysRuntime).Elem()
	v := reflect.ValueOf(mm.SysRuntime).Elem()
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		if _, ok := mm.RunMetrics[field.Name]; ok {
			switch v.FieldByName(field.Name).Type().Name() {
			case "uint64":
				mm.RunMetrics[field.Name] = gauge(v.FieldByName(field.Name).Uint())
			case "float64":
				mm.RunMetrics[field.Name] = gauge(v.FieldByName(field.Name).Float())
			}
		}
	}
}

func (mm *MyMetrics) ExportMetrics() *map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})
	result["counter"] = make(map[string]interface{})
	result["gauge"] = make(map[string]interface{})
	result["counter"]["PollCount"] = mm.PollCount
	result["gauge"]["RandomValue"] = mm.RandomValue
	for k, v := range mm.RunMetrics {
		result["gauge"][k] = v
	}

	return &result
}

/*
func getStructFieldAttr(v interface{}) map[string]string {
	t := reflect.TypeOf(v).Elem()
	numField := t.NumField()
	structAttr := make(map[string]string, numField)
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Map {
			maps.Copy(structAttr, getStructFieldAttr(field))
		}
		structAttr[field.Name] = field.Type.Name()
	}
	return structAttr
}
*/
