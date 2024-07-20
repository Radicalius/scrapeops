package shared

import "encoding/json"

type RawHandlerFunc func([]byte, Context) error
type HandlerFunc[T interface{}] func(T, Context) error
type RawHandlerFuncMap map[string]RawHandlerFunc

type DatabaseConfiguration struct {
	Name   string
	Tables []interface{}
}

type ApiFunc[I any, O any] func(I, Context) (O, error)
type RawApiFunc func([]byte, Context) (*[]byte, error)
type RawApiFuncMap map[string]RawApiFunc

type CronConfig struct {
	Schedule  string
	QueueName string
}

type PluginConfiguration struct {
	Handlers              RawHandlerFuncMap
	DatabaseConfiguration *DatabaseConfiguration
	CronJobs              []CronConfig
	Apis                  RawApiFuncMap
}

var DatabaseName string

func NewPluginConfiguration() *PluginConfiguration {
	return &PluginConfiguration{
		Handlers:              make(RawHandlerFuncMap),
		DatabaseConfiguration: nil,
		CronJobs:              make([]CronConfig, 0),
		Apis:                  make(RawApiFuncMap),
	}
}

func RegisterHandler[T interface{}](pluginConfig *PluginConfiguration, name string, callback HandlerFunc[T]) {
	pluginConfig.Handlers[name] = ToRawHandlerFunc(callback)
}

func ToRawHandlerFunc[T any](callback HandlerFunc[T]) func([]byte, Context) error {
	return func(inp []byte, ctx Context) error {
		var message T
		err := json.Unmarshal(inp, &message)
		if err != nil {
			return err
		}

		return callback(message, ctx)
	}
}

func RegisterDatabase(pluginConfig *PluginConfiguration, name string, tables ...interface{}) {
	DatabaseName = name
	pluginConfig.DatabaseConfiguration = &DatabaseConfiguration{
		Name:   name,
		Tables: tables,
	}
}

func RegisterCron(pluginConfig *PluginConfiguration, schedule string, queueName string) {
	pluginConfig.CronJobs = append(pluginConfig.CronJobs, CronConfig{
		QueueName: queueName,
		Schedule:  schedule,
	})
}

func RegisterApi[I any, O any](pluginConfig *PluginConfiguration, route string, callback ApiFunc[I, O]) {
	pluginConfig.Apis[route] = func(b []byte, ctx Context) (*[]byte, error) {
		var message I
		err := json.Unmarshal(b, &message)
		if err != nil {
			return nil, err
		}

		out, err := callback(message, ctx)
		if err != nil {
			return nil, err
		}

		res, err := json.Marshal(out)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}
}
