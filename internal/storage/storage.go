package storage

type gauge int64

type counter int64

type MemStorage struct {
	Gauge   map[string]gauge
	Counter map[string][]counter
}

func (ms *MemStorage) NewMemStorage() {
	ms.Gauge = make(map[string]gauge)
	ms.Counter = make(map[string][]counter)
}
