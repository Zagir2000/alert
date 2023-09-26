package metricscollect

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/Zagir2000/alert/internal/agent/hash"
	"github.com/Zagir2000/alert/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/johncgriffin/overflow"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const (
	counterMetric   string = "counter"     //counter metric
	gaugeMetric     string = "gauge"       //gauge metric
	randomValueName string = "RandomValue" // random value
	pollCountName   string = "PollCount"
)

// Константы для хедера.
const (
	contentType  string = "application/json"
	compressType string = "gzip"
)

// Структура для сбора и отправки метрик.
type RuntimeMetrics struct {
	RuntimeMemstats map[string]float64
	PollCount       int64
	RandomValue     float64
	pollInterval    time.Duration
}

// Структура для ошибки при отправке на сервер.
type SendMetricsError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Инициализицаия структуры для сбора и отправки метрик.
func IntervalPin(pollIntervalFlag int) RuntimeMetrics {
	return RuntimeMetrics{pollInterval: time.Duration(pollIntervalFlag), RuntimeMemstats: make(map[string]float64), PollCount: 0, RandomValue: 0}
}

// Добавление метрик в стрктуру RuntimeMetrics для будущей их отправки.
func (m *RuntimeMetrics) AddValueMetric() error {

	mapstats := make(map[string]float64)
	var RtMetrics runtime.MemStats
	runtime.ReadMemStats(&RtMetrics)
	mapstats["Alloc"] = float64(RtMetrics.Alloc)
	mapstats["BuckHashSys"] = float64(RtMetrics.BuckHashSys)
	mapstats["Frees"] = float64(RtMetrics.Frees)
	mapstats["GCCPUFraction"] = float64(RtMetrics.GCCPUFraction)
	mapstats["GCSys"] = float64(RtMetrics.GCSys)
	mapstats["HeapAlloc"] = float64(RtMetrics.HeapAlloc)
	mapstats["HeapIdle"] = float64(RtMetrics.HeapIdle)
	mapstats["HeapInuse"] = float64(RtMetrics.HeapInuse)
	mapstats["HeapObjects"] = float64(RtMetrics.HeapObjects)
	mapstats["HeapReleased"] = float64(RtMetrics.HeapReleased)
	mapstats["HeapSys"] = float64(RtMetrics.HeapSys)
	mapstats["LastGC"] = float64(RtMetrics.LastGC)
	mapstats["Lookups"] = float64(RtMetrics.Lookups)
	mapstats["MCacheInuse"] = float64(RtMetrics.MCacheInuse)
	mapstats["MCacheSys"] = float64(RtMetrics.MCacheSys)
	mapstats["MSpanInuse"] = float64(RtMetrics.MSpanInuse)
	mapstats["MSpanSys"] = float64(RtMetrics.MSpanSys)
	mapstats["Mallocs"] = float64(RtMetrics.Mallocs)
	mapstats["NextGC"] = float64(RtMetrics.NextGC)
	mapstats["NumForcedGC"] = float64(RtMetrics.NumForcedGC)
	mapstats["NumGC"] = float64(RtMetrics.NumGC)
	mapstats["OtherSys"] = float64(RtMetrics.OtherSys)
	mapstats["PauseTotalNs"] = float64(RtMetrics.PauseTotalNs)
	mapstats["StackInuse"] = float64(RtMetrics.StackInuse)
	mapstats["StackSys"] = float64(RtMetrics.StackSys)
	mapstats["Sys"] = float64(RtMetrics.Sys)
	mapstats["TotalAlloc"] = float64(RtMetrics.TotalAlloc)
	m.RandomValue = rand.Float64() * 10000
	if m.PollCount < 0 {
		return errors.New("counter is negative number")
	}

	checkCounterInOverflow, ok := overflow.Add64(m.PollCount, 1)

	if !ok {
		m.PollCount = 0
		return errors.New("counter is overflow")
	}
	m.PollCount = checkCounterInOverflow
	m.RuntimeMemstats = mapstats
	time.Sleep(m.pollInterval * time.Second)
	return nil
}

// Добавление метрик с помощью пакета gopsutil и сбор дополнительных метрик.
func (m *RuntimeMetrics) AddVaueMetricGopsutil() error {
	cpuStat, err := cpu.Times(true)
	if err != nil {
		return err
	}
	mapstats := make(map[string]float64)
	mem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	mapstats["TotalMemory"] = float64(mem.Total)
	mapstats["FreeMemory"] = float64(mem.Free)
	for _, k := range cpuStat {
		mapstats["CPUutilization1"] = k.Idle
	}

	m.RuntimeMemstats = mapstats
	return nil
}

// Функция для отправления метрик на сервер пачками.
func (m *RuntimeMetrics) jsonMetricsToBatch() []byte {
	metrics := make([]models.Metrics, 0, len(m.RuntimeMemstats)+2)
	for k, v := range m.RuntimeMemstats {
		valueGauge := v
		jsonGauge := &models.Metrics{
			ID:    k,
			MType: gaugeMetric,
			Value: &valueGauge,
		}
		metrics = append(metrics, *jsonGauge)
	}
	URLRandomGauge := &models.Metrics{
		ID:    randomValueName,
		MType: gaugeMetric,
		Value: &m.RandomValue,
	}
	metrics = append(metrics, *URLRandomGauge)
	URLCount := &models.Metrics{
		ID:    pollCountName,
		MType: counterMetric,
		Delta: &m.PollCount,
	}
	metrics = append(metrics, *URLCount)
	out, err := json.Marshal(metrics)
	if err != nil {
		log.Fatal(err)
	}
	return out
}

// Функция для отправление на сервер по одной метрике за запрос.
func SendMetrics(res []byte, hash, hostpath string) error {
	url := strings.Join([]string{"http:/", hostpath, "updates/"}, "/")

	responseErr := &SendMetricsError{}
	client := resty.New()
	_, err := client.R().
		SetError(responseErr).
		SetHeader("Content-Encoding", compressType).
		SetHeader("Content-Type", contentType).
		SetHeader("HashSHA256", hash).
		SetBody(res).
		Post(url)
	if err != nil {
		if errors.Is(err, syscall.ECONNRESET) && errors.Is(err, syscall.ECONNREFUSED) {
			_, err := client.R().
				SetHeader("Content-Type", contentType).
				SetHeader("Content-Encoding", compressType).
				SetHeaderVerbatim("HashSHA256", hash).
				SetBody(res).
				Post(url)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

// Функция которая добавляет и сжимает метрики для последующей отправки в канал.
func (m *RuntimeMetrics) NewСollect(ctx context.Context, cancel context.CancelFunc, jobs chan []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

			err := m.AddValueMetric()
			if err != nil {
				log.Println("Error in collect metrics 1:", err)
			}
			out := m.jsonMetricsToBatch()
			res, err := gzipCompress(out)
			if err != nil {
				log.Println("Error in comress metrics:", err)
			}
			jobs <- res

		}

	}
}

// Функция которая добавляет и сжимает метрики для последующей отправки в канал.
func (m *RuntimeMetrics) NewСollectMetricGopsutil(ctx context.Context, cancel context.CancelFunc, jobs chan []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(1) * time.Second):
			err := m.AddVaueMetricGopsutil()
			if err != nil {
				log.Println("Error in collect metrics 1", err)
			}
			out := m.jsonMetricsToBatch()
			res, err := gzipCompress(out)
			if err != nil {
				log.Println("Error in comress metrics:", err)
			}
			jobs <- res
		}

	}
}

// Функция для сжатия данных.
func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// Функция для отправки метрик на сервер.
func SendMetricsGor(jobs <-chan []byte, runAddr, secretKey string) error {
	for j := range jobs {

		hash := hash.CrateHash(secretKey, j, sha256.New)
		err := SendMetrics(j, hash, runAddr)
		if err != nil {
			log.Println("Error in send metrics:", err)
			return err
		}
	}
	return nil
}
