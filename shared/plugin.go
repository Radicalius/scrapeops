package shared

import "encoding/json"

type RawHandlerFunc func([]byte, Context) error
type HandlerFunc[T interface{}] func(T, Context) error
type RawHandlerFuncMap map[string]RawHandlerFunc

type CronFunc func(Context) error

type DatabaseConfiguration struct {
	Name       string
	Migrations []string
}

type PluginConfiguration struct {
	Handlers              RawHandlerFuncMap
	DatabaseConfiguration *DatabaseConfiguration
	CronJobs              map[string][]CronFunc
}

func NewPluginConfiguration() *PluginConfiguration {
	return &PluginConfiguration{
		Handlers:              make(RawHandlerFuncMap),
		DatabaseConfiguration: nil,
		CronJobs:              make(map[string][]CronFunc),
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

func RegisterDatabase(pluginConfig *PluginConfiguration, name string, migrations []string) {
	pluginConfig.DatabaseConfiguration = &DatabaseConfiguration{
		Name:       name,
		Migrations: migrations,
	}
}

func RegisterCron(pluginConfig *PluginConfiguration, schedule string, callback CronFunc) {
	_, exists := pluginConfig.CronJobs[schedule]
	if !exists {
		pluginConfig.CronJobs[schedule] = make([]CronFunc, 0)
	}

	pluginConfig.CronJobs[schedule] = append(pluginConfig.CronJobs[schedule], callback)
}
