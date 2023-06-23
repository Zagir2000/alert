package metricscollect

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	counterMetric   string = "counter"
	gaugeMetric     string = "gauge"
	randomValueName string = "RandomValue"
	pollCountName   string = "PollCount"
	contentType     string = "application/json"
	compressType    string = "gzip"
)

type RuntimeMetrics struct {
	RuntimeMemstats map[string]float64
	PollCount       int64
	RandomValue     float64
	pollInterval    time.Duration
}

type SendMetricsError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func IntervalPin(pollIntervalFlag int) RuntimeMetrics {
	return RuntimeMetrics{pollInterval: time.Duration(pollIntervalFlag), RuntimeMemstats: make(map[string]float64), PollCount: 0, RandomValue: 0}
}

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

func (m *RuntimeMetrics) AddVaueMetricGopsutil() error {
	cp, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}
	mapstats := make(map[string]float64)
	mem, err := mem.VirtualMemory()
	mapstats["TotalMemory"] = float64(mem.Total)
	mapstats["FreeMemory"] = float64(mem.Free)
	mapstats["CPUutilization1"] = cp[0]

	m.RuntimeMemstats = mapstats
	return nil
}

func (m *RuntimeMetrics) jsonMetricsToBatch() []byte {
	var metrics []models.Metrics
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

func (m *RuntimeMetrics) SendMetrics(res []byte, hash, hostpath string) error {
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

func (m *RuntimeMetrics) NewСollect(ctx context.Context, cancel context.CancelFunc, jobs chan []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

			err := m.AddValueMetric()
			if err != nil {
				log.Println("Error in collect metrics:", err)
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

func (m *RuntimeMetrics) NewСollectMetricGopsutil(ctx context.Context, cancel context.CancelFunc, jobs chan []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := m.AddVaueMetricGopsutil()
			if err != nil {
				log.Println("Error in collect metrics", err)
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

func (m *RuntimeMetrics) SendMetricsGor(ctx context.Context, cancel context.CancelFunc, jobs <-chan []byte, runAddr, secretKey string) {
	for j := range jobs {
		select {
		case <-ctx.Done():
			fmt.Println(j)
			return
		default:
			hash := hash.CrateHash(secretKey, j)
			err := m.SendMetrics(j, hash, runAddr)
			if err != nil {
				log.Println("Error in send metrics:", err)
				cancel()
			}
		}

	}

}
