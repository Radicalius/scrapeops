package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
)

type Point struct {
	Timestamp int64
	Data      float64
}

type Metrics struct {
	Data     map[string][]Point
	Counters map[string]float64
	Hists    map[string][]float64
}

func NewMetrics() *Metrics {
	return &Metrics{
		Data:     make(map[string][]Point),
		Counters: make(map[string]float64),
		Hists:    make(map[string][]float64),
	}
}

func (m *Metrics) IncrementCounter(counterName string) {
	originalVal, exists := m.Counters[counterName]
	if !exists {
		originalVal = 0
	}

	m.Counters[counterName] = originalVal + 1
}

func (m *Metrics) Observe(histName string, value float64) {
	appendToMapArray(&m.Hists, histName, value)
}

func (m *Metrics) Flush() error {
	now := time.Now().Unix()

	for counter := range m.Counters {
		appendToMapArray(&m.Data, counter, Point{
			Timestamp: now,
			Data:      m.Counters[counter],
		})
		m.Counters[counter] = 0
	}

	for hist := range m.Hists {
		p50, err := stats.Percentile(m.Hists[hist], 50)
		if err != nil {
			continue
		}
		appendToMapArray(&m.Data, fmt.Sprintf("%s_p50", hist), Point{
			Timestamp: now,
			Data:      p50,
		})

		p75, err := stats.Percentile(m.Hists[hist], 75)
		if err != nil {
			continue
		}
		appendToMapArray(&m.Data, fmt.Sprintf("%s_p75", hist), Point{
			Timestamp: now,
			Data:      p75,
		})

		p90, err := stats.Percentile(m.Hists[hist], 90)
		if err != nil {
			continue
		}
		appendToMapArray(&m.Data, fmt.Sprintf("%s_p90", hist), Point{
			Timestamp: now,
			Data:      p90,
		})

		p99, err := stats.Percentile(m.Hists[hist], 99)
		if err != nil {
			continue
		}
		appendToMapArray(&m.Data, fmt.Sprintf("%s_p99", hist), Point{
			Timestamp: now,
			Data:      p99,
		})

		m.Hists[hist] = make([]float64, 0)
	}

	return nil
}

func appendToMapArray[T interface{}](ma *map[string][]T, key string, value T) {
	values, exists := (*ma)[key]
	if !exists {
		values = make([]T, 0)
	}

	(*ma)[key] = append(values, value)
}

func (m *Metrics) InitMetricsApis() {
	http.HandleFunc("/api/grafana/metric", func(w http.ResponseWriter, r *http.Request) {
		metricName := r.URL.Query().Get("name")
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")

		metric, exists := m.Data[metricName]
		if !exists {
			w.WriteHeader(400)
			return
		}

		startTime, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			startTime = 0
		}

		endTime, err := strconv.ParseInt(end, 10, 64)
		if err != nil {
			endTime = 100000000000000
		}

		res := make([]Point, 0)
		for _, point := range metric {
			if point.Timestamp > startTime/1000 && point.Timestamp < endTime/1000 {
				res = append(res, point)
			}
		}

		data, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(500)
		}

		w.Header().Add("content-type", "application/json")
		w.Write(data)
	})
}
