package agent

import (
	"log"
	"time"

	"github.com/Zagir2000/alert/internal/metricscollect"
	"github.com/go-resty/resty/v2"
)

type MyAPIError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func sendMetrics(m *metricscollect.RuntimeMetrics) error {

	time.Sleep(time.Duration(reportInterval) * time.Millisecond)
	metrics := m.URLMetrics(flagRunAddr)
	client := resty.New()
	var responseErr MyAPIError
	for _, url := range metrics {
		_, err := client.R().
			SetError(&responseErr).
			SetHeader("Content-Type", "text/plain").
			Post(url)
		return err
	}
	return nil
}
func main() {
	parseFlags()
	Metric := metricscollect.PollIntervalPin(pollInterval)
	Metric.AddValueMetric()
	go Metric.NewСollect()
	for {
		err := sendMetrics(&Metric)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
