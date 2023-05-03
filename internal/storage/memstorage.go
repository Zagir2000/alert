package storage

type Repository interface {
	AddGaugeValue(name string, value float64)
	AddCounterValue(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllGaugeValues() map[string]float64
	GetAllCounterValues() map[string]int64
}

type MemStorage struct {
	Gaugedata   map[string]float64
	Counterdata map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gaugedata:   make(map[string]float64),
		Counterdata: make(map[string]int64),
	}
}

func (m *MemStorage) AddGaugeValue(name string, value float64) {
	m.Gaugedata[name] = value
}

func (m *MemStorage) AddCounterValue(name string, value int64) {
	m.Counterdata[name] += value
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	value, ok := m.Gaugedata[name]
	return value, ok
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	value, ok := m.Counterdata[name]
	return value, ok
}

func (m *MemStorage) GetAllGaugeValues() map[string]float64 {
	return m.Gaugedata
}

func (m *MemStorage) GetAllCounterValues() map[string]int64 {
	return m.Counterdata
}
