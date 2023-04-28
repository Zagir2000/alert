package storage

type gauge float64
type counter int64

type MemStorage struct {
	gaugedata   map[string]gauge
	counterdata map[string]counter
}
