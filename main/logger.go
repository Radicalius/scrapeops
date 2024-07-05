package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	Attrs map[string]string
}

func NewLogger() *Logger {
	return &Logger{
		Attrs: make(map[string]string),
	}
}

func (l *Logger) With(params ...string) *Logger {
	newLogger := NewLogger()

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
	kv["time"] = time.Now().Format("2006-01-02T15:04:05Z07:00")

	for key, val := range l.Attrs {
		kv[key] = val
	}

	copy(&kv, params...)

	data, _ := json.Marshal(kv)
	fmt.Println(string(data))
}

func copy(attrs *map[string]string, params ...string) {
	for i := 0; i < len(params)/2; i++ {
		(*attrs)[params[i*2]] = params[i*2+1]
	}
}
