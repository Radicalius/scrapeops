package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type LogCollector struct {
	Logs []map[string]string
}

func (lc *LogCollector) Collect(log map[string]string) {
	lc.Logs = append(lc.Logs, log)
}

func (lc *LogCollector) Find(startTime int64, endTime int64, params map[string]string) []map[string]string {
	res := make([]map[string]string, 0)
	for _, log := range lc.Logs {
		for i := range params {
			if log[i] != params[i] {
				continue
			}
		}

		logTime, _ := strconv.ParseInt(log["time"], 10, 64)
		if logTime > startTime && logTime < endTime {
			res = append(res, log)
		}
	}

	return res
}

func InitLogCollector() *LogCollector {
	res := &LogCollector{
		Logs: make([]map[string]string, 0),
	}

	http.HandleFunc("/api/grafana/logs", func(w http.ResponseWriter, r *http.Request) {
		params := make(map[string]string)
		for value := range r.URL.Query() {
			if value != "start" && value != "end" {
				params[value] = r.URL.Query()[value][0]
			}
		}

		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")

		startTime, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			startTime = 0
		}

		endTime, err := strconv.ParseInt(end, 10, 64)
		if err != nil {
			endTime = 100000000000000
		}

		data := res.Find(startTime, endTime, params)
		dataBytes, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Error marshaling data: %s", err.Error())
			w.WriteHeader(500)
			return
		}

		w.Header().Add("content-type", "application/json")
		w.Write(dataBytes)
	})

	return res
}

type Logger struct {
	Attrs     map[string]string
	collector *LogCollector
}

func NewLogger(collector *LogCollector) *Logger {
	return &Logger{
		Attrs:     make(map[string]string),
		collector: collector,
	}
}

func (l *Logger) With(params ...string) *Logger {
	newLogger := NewLogger(l.collector)

	for key, val := range l.Attrs {
		newLogger.Attrs[key] = val
	}

	copy(&newLogger.Attrs, params...)

	return newLogger
}

func (l *Logger) Fatal(message string, params ...string) {
	l.Log("fatal", message, params...)
	os.Exit(1)
}

func (l *Logger) Error(message string, params ...string) {
	l.Log("error", message, params...)
}

func (l *Logger) Warn(message string, params ...string) {
	l.Log("warning", message, params...)
}

func (l *Logger) Info(message string, params ...string) {
	l.Log("info", message, params...)
}

func (l *Logger) Log(level string, message string, params ...string) {
	kv := make(map[string]string)
	kv["message"] = message
	kv["level"] = level
	kv["time"] = fmt.Sprintf("%d", time.Now().Unix())
	kv["formattedTime"] = time.Now().Format("2006-01-02T15:04:05Z07:00")

	for key, val := range l.Attrs {
		kv[key] = val
	}

	copy(&kv, params...)

	l.collector.Collect(kv)

	data, _ := json.Marshal(kv)
	fmt.Println(string(data))
}

func copy(attrs *map[string]string, params ...string) {
	for i := 0; i < len(params)/2; i++ {
		(*attrs)[params[i*2]] = params[i*2+1]
	}
}

var baseLogger = NewLogger(InitLogCollector())
